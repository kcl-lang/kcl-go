// Copyright 2021 The KCL Authors. All rights reserved.

package parser

import (
	"encoding/json"
)

func _JSONString(v interface{}) string {
	if s, ok := v.(string); ok {
		v = []byte(s)
	}
	if x, ok := v.([]byte); ok {
		var m map[string]interface{}
		if err := json.Unmarshal(x, &m); err != nil {
			return string(x)
		}
		result, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			return string(x)
		}
		return string(result)
	}
	x, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}
	return string(x)
}
