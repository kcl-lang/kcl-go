// Copyright 2023 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package hello_plugin

import (
	"testing"

	"kcl-lang.io/kcl-go/pkg/plugin"
)

func TestPlugin_add(t *testing.T) {
	result_json := plugin.Invoke("kcl_plugin.hello.add", []interface{}{111, 22}, nil)
	if result_json != "133" {
		t.Fatal(result_json)
	}
}
