// Copyright The KCL Authors. All rights reserved.

// Package kcl defines the top-level interface for KCL.
package kcl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/chai2010/jsonv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"kcl-lang.io/kcl-go/pkg/loader"
	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type KclType = gpyrpc.KclType

type KCLResultList struct {
	// When the list is empty, the result is the raw result.
	list            []KCLResult
	raw_json_result string
	raw_yaml_result string
}

// ToString returns the result as string.
func (p *KCLResultList) ToString() (string, error) {
	if len(p.list) == 0 {
		return "", fmt.Errorf("result is nil")
	}
	var resS string
	err := p.ToType(&resS)
	if err != nil {
		return "", err
	}
	return resS, nil
}

// ToBool returns the result as bool.
func (p *KCLResultList) ToBool() (*bool, error) {
	if len(p.list) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	var resB bool
	err := p.ToType(&resB)
	if err != nil {
		return nil, err
	}
	return &resB, nil
}

// ToMap returns the result as map[string]interface{}.
func (p *KCLResultList) ToMap() (map[string]interface{}, error) {
	if len(p.list) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	var resMap map[string]interface{}
	err := p.ToType(&resMap)
	if err != nil {
		return nil, err
	}
	return resMap, nil
}

// ToInt returns the result as int.
func (p *KCLResultList) ToInt() (*int, error) {
	if len(p.list) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	var resI int
	err := p.ToType(&resI)
	if err != nil {
		return nil, err
	}
	return &resI, nil
}

// ToFloat64 returns the result as float64.
func (p *KCLResultList) ToFloat64() (*float64, error) {
	if len(p.list) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	var resF float64
	err := p.ToType(&resF)
	if err != nil {
		return nil, err
	}
	return &resF, nil
}

// ToList returns the result as []interface{}.
func (p *KCLResultList) ToList() ([]interface{}, error) {
	if len(p.list) == 0 {
		return nil, fmt.Errorf("result is nil")
	}
	var resList []interface{}
	err := p.ToType(&resList)
	if err != nil {
		return nil, err
	}
	return resList, nil
}

// ToType returns the result as target type.
func (p *KCLResultList) ToType(target interface{}) error {
	if len(p.list) == 0 {
		return fmt.Errorf("result is nil")
	}

	srcVal := reflect.ValueOf(p.list[0].result)
	targetVal := reflect.ValueOf(target)

	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return fmt.Errorf("failed to convert result to %T", target)
	}

	if srcVal.Type() != targetVal.Elem().Type() {
		return fmt.Errorf("failed to convert result to %T: type mismatch", target)
	}

	targetVal.Elem().Set(srcVal)

	return nil
}

func (p *KCLResultList) Len() int {
	return len(p.list)
}

func (p *KCLResultList) Get(i int) *KCLResult {
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
func (p *KCLResultList) First() *KCLResult {
	if p.Len() > 0 {
		return &p.list[0]
	}
	return nil
}
func (p *KCLResultList) Tail() *KCLResult {
	if p.Len() > 0 {
		return &p.list[len(p.list)-1]
	}
	return nil
}

func (p *KCLResultList) Slice() []KCLResult {
	return p.list
}

func (p *KCLResultList) GetRawJsonResult() string {
	return p.raw_json_result
}
func (p *KCLResultList) GetRawYamlResult() string {
	return p.raw_yaml_result
}

// KCLResult denotes the result for the Run API.
type KCLResult struct {
	result any
}

// NewResult constructs a KCLResult using the value
func NewResult(value any) KCLResult {
	return KCLResult{
		result: value,
	}
}

// ToString returns the result as string.
func (p *KCLResult) ToString() (string, error) {
	var resS string
	err := p.ToType(&resS)
	if err != nil {
		return "", err
	}
	return resS, nil
}

// ToBool returns the result as bool.
func (p *KCLResult) ToBool() (*bool, error) {
	var resB bool
	err := p.ToType(&resB)
	if err != nil {
		return nil, err
	}
	return &resB, nil
}

// ToMap returns the result as map[string]interface{}.
func (p *KCLResult) ToMap() (map[string]interface{}, error) {
	var resMap map[string]interface{}
	err := p.ToType(&resMap)
	if err != nil {
		return nil, err
	}
	return resMap, nil
}

// ToInt returns the result as int.
func (p *KCLResult) ToInt() (*int, error) {
	var resI int
	err := p.ToType(&resI)
	if err != nil {
		return nil, err
	}
	return &resI, nil
}

// ToFloat64 returns the result as float64.
func (p *KCLResult) ToFloat64() (*float64, error) {
	var resF float64
	err := p.ToType(&resF)
	if err != nil {
		return nil, err
	}
	return &resF, nil
}

// ToList returns the result as []interface{}.
func (p *KCLResult) ToList() ([]interface{}, error) {
	var resList []interface{}
	err := p.ToType(&resList)
	if err != nil {
		return nil, err
	}
	return resList, nil
}

// ToType returns the result as target type.
func (p *KCLResult) ToType(target interface{}) error {
	srcVal := reflect.ValueOf(p.result)
	targetVal := reflect.ValueOf(target)

	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return fmt.Errorf("failed to convert result to %T", target)
	}

	if srcVal.Type() != targetVal.Elem().Type() {
		return fmt.Errorf("failed to convert result to %T: type mismatch", target)
	}

	targetVal.Elem().Set(srcVal)

	return nil
}

func (m *KCLResult) Get(key string, target ...interface{}) interface{} {
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

func (m *KCLResult) GetValue(key string, target ...interface{}) (value interface{}, err error) {
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

func (m *KCLResult) YAMLString() string {
	out, _ := yaml.Marshal(m.result)
	return string(out)
}

func (m *KCLResult) JSONString() string {
	var prefix = ""
	var indent = "    "
	x, _ := json.MarshalIndent(m.result, prefix, indent)
	return string(x)
}

func MustRun(path string, opts ...Option) *KCLResultList {
	v, err := Run(path, opts...)
	if err != nil {
		panic(err)
	}

	return v
}

func Run(path string, opts ...Option) (*KCLResultList, error) {
	return run([]string{path}, opts...)
}

// RunWithOpts is the same as Run, but it does not require a path as input.
// Note: you need to specify the path in options by method WithKFilenameList()
// or the workdir in method WorkDir(),
// or it will return an error.
func RunWithOpts(opts ...Option) (*KCLResultList, error) {
	return run([]string{}, opts...)
}

func RunFiles(paths []string, opts ...Option) (*KCLResultList, error) {
	return run(paths, opts...)
}

func GetSchemaType(filename string, src any, schemaName string) ([]*gpyrpc.KclType, error) {
	mapping, err := GetSchemaTypeMapping(filename, src, schemaName)
	if err != nil {
		return nil, err
	}
	return getValues(mapping), nil
}

func GetFullSchemaType(pathList []string, schemaName string, opts ...Option) ([]*gpyrpc.KclType, error) {
	mapping, err := GetFullSchemaTypeMapping(pathList, schemaName, opts...)
	if err != nil {
		return nil, err
	}
	return getValues(mapping), nil
}

func GetFullSchemaTypeMapping(pathList []string, schemaName string, opts ...Option) (map[string]*gpyrpc.KclType, error) {
	opts = append(opts, *NewOption().Merge(WithKFilenames(pathList...)))
	args, err := ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.GetSchemaTypeMapping(&gpyrpc.GetSchemaTypeMapping_Args{
		ExecArgs:   args.ExecProgram_Args,
		SchemaName: schemaName,
	})

	if err != nil {
		return nil, err
	}

	return resp.SchemaTypeMapping, nil
}

func GetSchemaTypeMapping(filename string, src any, schemaName string) (map[string]*gpyrpc.KclType, error) {
	source, err := loader.ReadSource(filename, src)
	if err != nil {
		return nil, err
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.GetSchemaTypeMapping(&gpyrpc.GetSchemaTypeMapping_Args{
		ExecArgs: &gpyrpc.ExecProgram_Args{
			KFilenameList: []string{filename},
			KCodeList:     []string{string(source)},
		},
		SchemaName: schemaName,
	})
	if err != nil {
		return nil, err
	}
	return resp.SchemaTypeMapping, nil
}

func getValues(myMap map[string]*gpyrpc.KclType) []*gpyrpc.KclType {
	var values []*gpyrpc.KclType
	for _, value := range myMap {
		values = append(values, value)
	}
	return values
}

func ExecResultToKCLResult(o *Option, resp *gpyrpc.ExecProgram_Result, logger io.Writer, hooks Hooks) (*KCLResultList, error) {
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

	var result KCLResultList
	if strings.TrimSpace(resp.JsonResult) == "" {
		return &result, nil
	}

	documents, err := SplitDocuments(resp.YamlResult)
	if err != nil {
		return &result, nil
	}

	for _, d := range documents {
		var m any
		if err := yaml.Unmarshal([]byte(d), &m); err != nil {
			return nil, err
		}
		result.list = append(result.list, KCLResult{
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

func runWithHooks(pathList []string, hooks Hooks, opts ...Option) (*KCLResultList, error) {
	args, err := ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.ExecProgram(args.ExecProgram_Args)
	if err != nil {
		return nil, err
	}
	return ExecResultToKCLResult(&args, resp, args.GetLogger(), hooks)
}

func run(pathList []string, opts ...Option) (*KCLResultList, error) {
	return runWithHooks(pathList, DefaultHooks, opts...)
}

// SplitDocuments returns a slice of all documents contained in a YAML string. Multiple documents can be divided by the
// YAML document separator (---). It allows for white space and comments to be after the separator on the same line,
// but will return an error if anything else is on the line.
func SplitDocuments(s string) ([]string, error) {
	docs := make([]string, 0)
	if len(s) > 0 {
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
	}
	return docs, nil
}
