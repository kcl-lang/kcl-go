// Copyright The KCL Authors. All rights reserved.

// Package kcl defines the top-level interface for KCL.
package kcl

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/chai2010/jsonv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type KclType = gpyrpc.KclType

type KCLResultList struct {
	// When the list is empty, the result is the raw result.
	list []KCLResult
	// The result is the raw result whose type is not []KCLResult.
	result          interface{}
	raw_json_result string
	raw_yaml_result string
}

// ToString returns the result as string.
func (p *KCLResultList) ToString() (string, error) {
	if p == nil || p.result == nil {
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
	if p == nil || p.result == nil {
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
	if p == nil || p.result == nil {
		return nil, fmt.Errorf("result is nil")
	}
	var resMap map[string]interface{}
	err := p.ToType(&resMap)
	if err != nil {
		return nil, err
	}
	return resMap, nil
}

// ToFloat64 returns the result as float64.
func (p *KCLResultList) ToFloat64() (*float64, error) {
	if p == nil || p.result == nil {
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
	if p == nil || p.result == nil {
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
	if p == nil || p.result == nil {
		return fmt.Errorf("result is nil")
	}

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

func (p *KCLResultList) Len() int {
	return len(p.list)
}

func (p *KCLResultList) Get(i int) KCLResult {
	if i == 0 {
		return p.First()
	}
	if i == p.Len()-1 {
		return p.Tail()
	}
	if i >= 0 && i < p.Len() {
		return p.list[i]
	}
	return nil
}
func (p *KCLResultList) First() KCLResult {
	if p.Len() > 0 {
		return p.list[0]
	}
	return nil
}
func (p *KCLResultList) Tail() KCLResult {
	if p.Len() > 0 {
		return p.list[len(p.list)-1]
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

type KCLResult map[string]interface{}

func (m KCLResult) Get(key string, target ...interface{}) interface{} {
	ss := strings.Split(key, ".")
	if len(ss) == 0 {
		return nil
	}

	var subKeys []interface{}
	for _, x := range ss[1:] {
		subKeys = append(subKeys, x)
	}

	rv := jsonv.Get((map[string]interface{})(m), ss[0], subKeys...)
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

func (m KCLResult) GetValue(key string, target ...interface{}) (value interface{}, err error) {
	ss := strings.Split(key, ".")
	if len(ss) == 0 {
		return nil, nil
	}

	var subKeys []interface{}
	for _, x := range ss[1:] {
		subKeys = append(subKeys, x)
	}

	rv, err := jsonv.GetValue((map[string]interface{})(m), ss[0], subKeys...)
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
		} else {
			return rv, err
		}
	case string:
		if pv, ok2 := target[0].(*string); ok2 {
			*pv = rv
			return rv, nil
		} else {
			return "", fmt.Errorf("target expect *string type: got = %T", target[0])
		}
	case int:
		if pv, ok2 := target[0].(*int); ok2 {
			*pv = rv
			return rv, nil
		} else {
			return "", fmt.Errorf("target expect *int type: got = %T", target[0])
		}
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

func (m KCLResult) YAMLString() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

func (m KCLResult) JSONString() string {
	var prefix = ""
	var indent = "    "
	x, _ := json.MarshalIndent(m, prefix, indent)
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

func run(pathList []string, opts ...Option) (*KCLResultList, error) {
	args, err := ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.ExecProgram(args.ExecProgram_Args)
	if err != nil {
		return nil, err
	}
	// Output log message
	logger := args.GetLogger()
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

	var mList []map[string]interface{}

	if err := json.Unmarshal([]byte(resp.JsonResult), &mList); err != nil {
		err = nil
		if err := json.Unmarshal([]byte(resp.JsonResult), &result.result); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
	}
	result.list = make([]KCLResult, 0, len(mList))
	for _, m := range mList {
		if len(m) != 0 {
			result.list = append(result.list, m)
		}
	}

	result.raw_json_result = resp.JsonResult
	result.raw_yaml_result = resp.YamlResult
	return &result, nil
}
