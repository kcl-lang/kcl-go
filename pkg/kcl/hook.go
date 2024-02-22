package kcl

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type Hook interface {
	Do(o *Option, r *gpyrpc.ExecProgram_Result) error
}

type Hooks []Hook

type typeAttributeHook struct{}

func (t *typeAttributeHook) Do(o *Option, r *gpyrpc.ExecProgram_Result) error {
	// Deal the `_type` attribute
	if o != nil && r != nil && !o.fullTypePath && o.IncludeSchemaTypePath {
		return resultTypeAttributeHook(r)
	}
	return nil
}

func resultTypeAttributeHook(r *gpyrpc.ExecProgram_Result) error {
	// Modify the _type fields
	var data []map[string]interface{}
	// Unmarshal the JSON string into a Node
	err := json.Unmarshal([]byte(r.JsonResult), &data)
	if err != nil {
		return nil
	}
	// Modify the _type fields
	modifyTypeList(data)
	// Marshal the modified Node back to YAML
	yamlOutput, err := yaml.Marshal(&data)
	if err != nil {
		return nil
	}
	// Marshal the modified Node back to JSON
	jsonOutput, err := json.Marshal(&data)
	if err != nil {
		return nil
	}
	r.JsonResult = string(jsonOutput)
	r.YamlResult = string(yamlOutput)
	return nil
}

func modifyTypeList(dataList []map[string]interface{}) {
	for _, data := range dataList {
		modifyType(data)
	}
}

func modifyType(data map[string]interface{}) {
	for key, value := range data {
		if key == "_type" {
			if v, ok := data[key].(string); ok {
				parts := strings.Split(v, ".")
				data[key] = parts[len(parts)-1]
			}
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			modifyType(nestedMap)
		}
	}
}
