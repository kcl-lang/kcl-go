// Copyright The KCL Authors. All rights reserved.

// Package kcl defines the top-level interface for KCL.
package kcl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/chai2010/jsonv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type KclType = gpyrpc.KclType

type KCLResultList[T KCLResultType] struct {
	// When the list is empty, the result is the raw result.
	list            []KCLResult[T]
	raw_json_result string
	raw_yaml_result string
}

// GetResult returns the result
func (k *KCLResultList[T]) GetResult() (T, error) {
	var result T
	if len(k.list) == 0 {
		return result, fmt.Errorf("result is nil")
	}

	return k.list[0].GetResult(), nil
}

func (p *KCLResultList[T]) Len() int {
	return len(p.list)
}

func (p *KCLResultList[T]) Get(i int) *KCLResult[T] {
	if i == 0 {
		return p.First()
	}

	if i == p.Len()-1 {
		return p.Tail()
	}

	if i >= 0 && i < p.Len() {
		return &p.list[i]
	}
	return nil
}

func (p *KCLResultList[T]) First() *KCLResult[T] {
	if p.Len() > 0 {
		return &p.list[0]
	}
	return nil
}

func (p *KCLResultList[T]) Tail() *KCLResult[T] {
	if p.Len() > 0 {
		return &p.list[len(p.list)-1]
	}
	return nil
}

func (p *KCLResultList[T]) Slice() []KCLResult[T] {
	return p.list
}

func (p *KCLResultList[T]) GetRawJsonResult() string {
	return p.raw_json_result
}

func (p *KCLResultList[T]) GetRawYamlResult() string {
	return p.raw_yaml_result
}

type KCLResultType interface {
	string | bool | map[string]any | int | float64 | []interface{} | any
}

// KCLResult denotes the result for the Run API.
type KCLResult[T KCLResultType] struct {
	result T
}

// NewResult constructs a KCLResult using the value
func NewResult[T KCLResultType](value T) KCLResult[T] {
	return KCLResult[T]{
		result: value,
	}
}

// GetResult returns the result
func (k *KCLResult[T]) GetResult() T {
	return k.result
}

func (m *KCLResult[T]) Get(key string, target ...interface{}) interface{} {
	ss := strings.Split(key, ".")
	if len(ss) == 0 {
		return nil
	}
	var subKeys []interface{}
	for _, x := range ss[1:] {
		subKeys = append(subKeys, x)
	}

	rv := jsonv.Get(m.result, ss[0], subKeys...)
	if len(target) == 0 {
		return rv
	}

	if m, ok := rv.(map[string]interface{}); ok {
		if err := mapstructure.Decode(m, target[0]); err == nil {
			return target[0]
		}
	}

	return rv
}

func (m *KCLResult[T]) GetValue(key string, target ...interface{}) (value KCLResultType, err error) {
	ss := strings.Split(key, ".")
	if len(ss) == 0 {
		return nil, nil
	}

	var subKeys []interface{}
	for _, x := range ss[1:] {
		subKeys = append(subKeys, x)
	}

	rv, err := jsonv.GetValue(m.result, ss[0], subKeys...)
	if err != nil {
		return nil, err
	}
	if len(target) == 0 {
		return rv, nil
	}

	switch rv := rv.(type) {
	case map[string]interface{}:
		if err := mapstructure.Decode(rv, target[0]); err == nil {
			return target[0], nil
		}

		return rv, err
	case string:
		if pv, ok2 := target[0].(*string); ok2 {
			*pv = rv
			return rv, nil
		}

		return "", fmt.Errorf("target expect *string type: got = %T", target[0])
	case int:
		if pv, ok2 := target[0].(*int); ok2 {
			*pv = rv
			return rv, nil
		}

		return "", fmt.Errorf("target expect *int type: got = %T", target[0])
	case float64:
		switch pTarget0 := target[0].(type) {
		case *int:
			*pTarget0 = int(rv)
			return rv, nil
		case *float64:
			*pTarget0 = rv
			return rv, nil
		default:
			return "", fmt.Errorf("%s expect *float64 or *int type: got = %T", key, target[0])
		}
	default:
		return rv, fmt.Errorf("unknown type: got = %T", target[0])
	}
}

func (m *KCLResult[T]) YAMLString() string {
	out, _ := yaml.Marshal(m.result)
	return string(out)
}

func (m *KCLResult[T]) JSONString() string {
	var prefix = ""
	var indent = "    "
	x, _ := json.MarshalIndent(m.result, prefix, indent)
	return string(x)
}

func MustRun(path string, opts ...Option) *KCLResultList[KCLResultType] {
	v, err := Run(path, opts...)
	if err != nil {
		panic(err)
	}

	return v
}

func Run(path string, opts ...Option) (*KCLResultList[KCLResultType], error) {
	return run[KCLResultType]([]string{path}, opts...)
}

// RunWithOpts is the same as Run, but it does not require a path as input.
// Note: you need to specify the path in options by method WithKFilenameList()
// or the workdir in method WorkDir(),
// or it will return an error.
func RunWithOpts(opts ...Option) (*KCLResultList[KCLResultType], error) {
	return run[KCLResultType]([]string{}, opts...)
}

func RunFiles(paths []string, opts ...Option) (*KCLResultList[KCLResultType], error) {
	return run[KCLResultType](paths, opts...)
}

func GetSchemaType(file, code, schemaName string) ([]*gpyrpc.KclType, error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.GetSchemaType(&gpyrpc.GetSchemaType_Args{
		File:       file,
		Code:       code,
		SchemaName: schemaName,
	})
	if err != nil {
		return nil, err
	}

	return resp.SchemaTypeList, nil
}

func GetFullSchemaType(pathList []string, schemaName string, opts ...Option) ([]*gpyrpc.KclType, error) {
	opts = append(opts, *NewOption().Merge(WithKFilenames(pathList...)))
	args, err := ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.GetFullSchemaType(&gpyrpc.GetFullSchemaType_Args{
		ExecArgs:   args.ExecProgram_Args,
		SchemaName: schemaName,
	})

	if err != nil {
		return nil, err
	}

	return resp.SchemaTypeList, nil
}

func GetSchemaTypeMapping(file, code, schemaName string) (map[string]*gpyrpc.KclType, error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.GetSchemaTypeMapping(&gpyrpc.GetSchemaTypeMapping_Args{
		File:       file,
		Code:       code,
		SchemaName: schemaName,
	})
	if err != nil {
		return nil, err
	}

	return resp.SchemaTypeMapping, nil
}

func ExecResultToKCLResult[T KCLResultType](o *Option, resp *gpyrpc.ExecProgram_Result, logger io.Writer, hooks Hooks) (*KCLResultList[T], error) {
	for _, hook := range hooks {
		hook.Do(o, resp)
	}
	if logger != nil && resp.LogMessage != "" {
		_, err := logger.Write([]byte(resp.LogMessage))
		if err != nil {
			return nil, err
		}
	}
	if resp.ErrMessage != "" {
		return nil, errors.New(resp.ErrMessage)
	}

	var result KCLResultList[T]
	if strings.TrimSpace(resp.JsonResult) == "" {
		return &result, nil
	}

	documents, err := splitDocuments(resp.YamlResult)
	if err != nil {
		return &result, nil
	}

	for _, d := range documents {
		var m T
		if err := yaml.Unmarshal([]byte(d), &m); err != nil {
			return nil, err
		}
		result.list = append(result.list, KCLResult[T]{
			result: m,
		})
	}
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	for _, m := range result.list {
		encoder.Encode(m.result)
	}
	result.raw_json_result = resp.JsonResult
	result.raw_yaml_result = resp.YamlResult
	return &result, nil
}

func runWithHooks[T KCLResultType](pathList []string, hooks Hooks, opts ...Option) (*KCLResultList[T], error) {
	args, err := ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.ExecProgram(args.ExecProgram_Args)
	if err != nil {
		return nil, err
	}
	return ExecResultToKCLResult[T](&args, resp, args.GetLogger(), hooks)
}

func run[T KCLResultType](pathList []string, opts ...Option) (*KCLResultList[T], error) {
	return runWithHooks[T](pathList, DefaultHooks, opts...)
}

// splitDocuments returns a slice of all documents contained in a YAML string. Multiple documents can be divided by the
// YAML document separator (---). It allows for white space and comments to be after the separator on the same line,
// but will return an error if anything else is on the line.
func splitDocuments(s string) ([]string, error) {
	docs := make([]string, 0)
	if len(s) == 0 {
		return docs, nil
	}

	// The YAML document separator is any line that starts with ---
	yamlSeparatorRegexp := regexp.MustCompile(`\n---.*\n`)

	// Find all separators, check them for invalid content, and append each document to docs
	separatorLocations := yamlSeparatorRegexp.FindAllStringIndex(s, -1)
	prev := 0
	for i := range separatorLocations {
		loc := separatorLocations[i]
		separator := s[loc[0]:loc[1]]
		// If the next non-whitespace character on the line following the separator is not a comment, return an error
		trimmedContentAfterSeparator := strings.TrimSpace(separator[4:])
		if len(trimmedContentAfterSeparator) > 0 && trimmedContentAfterSeparator[0] != '#' {
			return nil, fmt.Errorf("invalid document separator: %s", strings.TrimSpace(separator))
		}

		docs = append(docs, s[prev:loc[0]])
		prev = loc[1]
	}

	docs = append(docs, s[prev:])

	return docs, nil
}
