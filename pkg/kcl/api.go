// Copyright 2021 The KCL Authors. All rights reserved.

// Package kcl defines the top-level interface for the Kusion Configuration Language (KCL).
package kcl

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chai2010/jsonv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

type KclType = gpyrpc.KclType

type KCLResultList struct {
	list            []KCLResult
	raw_json_result string
	raw_yaml_result string
	escaped_time    string
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

func (p *KCLResultList) GetPyEscapedTime() string {
	return p.escaped_time
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

	var result KCLResultList
	if strings.TrimSpace(resp.JsonResult) == "" {
		return &result, nil
	}

	var mList []map[string]interface{}
	if err := json.Unmarshal([]byte(resp.JsonResult), &mList); err != nil {
		return nil, err
	}
	if len(mList) == 0 {
		return nil, fmt.Errorf("kcl.Run: invalid result: %s", resp.JsonResult)
	}

	for _, m := range mList {
		if len(m) != 0 {
			result.list = append(result.list, m)
		}
	}

	result.raw_json_result = resp.JsonResult
	result.raw_yaml_result = resp.YamlResult
	result.escaped_time = resp.EscapedTime
	return &result, nil
}
