// Copyright 2022 The KCL Authors. All rights reserved.

package hello_plugin

import (
	"testing"

	"kusionstack.io/kclvm-go/pkg/kcl_plugin"
)

func TestPlugin_global_int(t *testing.T) {
	if !kcl_plugin.CgoEnabled {
		t.Skip("cgo disabled")
	}

	kcl_plugin.Invoke("kcl_plugin.hello.set_global_int", []interface{}{123}, nil)
	result_json := kcl_plugin.Invoke("kcl_plugin.hello.get_global_int", nil, nil)
	if result_json != "123" {
		t.Fatal(result_json)
	}

	kcl_plugin.ResetPlugin()

	result_json = kcl_plugin.Invoke("kcl_plugin.hello.get_global_int", nil, nil)
	if result_json != "0" {
		t.Fatal(result_json)
	}

	kcl_plugin.Invoke("kcl_plugin.hello.set_global_int", []interface{}{1024}, nil)
	result_json = kcl_plugin.Invoke("kcl_plugin.hello.get_global_int", nil, nil)
	if result_json != "1024" {
		t.Fatal(result_json)
	}
}

func TestPlugin_add(t *testing.T) {
	if !kcl_plugin.CgoEnabled {
		t.Skip("cgo disabled")
	}

	result_json := kcl_plugin.Invoke("kcl_plugin.hello.add", []interface{}{111, 22}, nil)
	if result_json != "133" {
		t.Fatal(result_json)
	}
}

func TestPlugin_tolower(t *testing.T) {
	if !kcl_plugin.CgoEnabled {
		t.Skip("cgo disabled")
	}
	result_json := kcl_plugin.Invoke("kcl_plugin.hello.tolower", []interface{}{"KCL"}, nil)
	if result_json != `"kcl"` {
		t.Fatal(result_json)
	}
}

func TestPlugin_update_dict(t *testing.T) {
	if !kcl_plugin.CgoEnabled {
		t.Skip("cgo disabled")
	}
	dict := map[string]interface{}{
		"name": 123,
	}

	result_json := kcl_plugin.Invoke("kcl_plugin.hello.update_dict", []interface{}{dict, "name", "KusionStack"}, nil)
	if result_json != `{"name":"KusionStack"}` {
		t.Fatal(result_json)
	}
}

func TestPlugin_list_append(t *testing.T) {
	if !kcl_plugin.CgoEnabled {
		t.Skip("cgo disabled")
	}
	list := []interface{}{"abc"}
	result_json := kcl_plugin.Invoke("kcl_plugin.hello.list_append", []interface{}{list, "name", 123}, nil)
	if result_json != `["abc","name",123]` {
		t.Fatal(result_json)
	}
}
