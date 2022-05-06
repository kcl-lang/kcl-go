// Copyright 2022 The KCL Authors. All rights reserved.

package ast

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func readSource(filename string, src interface{}) (data []byte, err error) {
	if src == nil {
		src, err = os.ReadFile(filename)
		if err != nil {
			return
		}
	}

	switch src := src.(type) {
	case []byte:
		return src, nil
	case string:
		return []byte(src), nil
	case io.Reader:
		d, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported src type: %T", src)
	}
}

func LoadJson(filename string, src interface{}) (map[string]interface{}, error) {
	data, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func json_decodeMap(v interface{}) (map[string]interface{}, error) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}

	var data []byte
	switch v := v.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		d, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data = d
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func JSONString(v interface{}) string {
	return json_String(v)
}

func JSONMap(x interface{}) map[string]interface{} {
	m, err := json_decodeMap(x)
	if err == nil {
		return m
	}
	fmt.Println("err:", err)
	return nil
}

func json_String(v interface{}) string {
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
