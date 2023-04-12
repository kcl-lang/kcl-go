// Copyright 2023 The KCL Authors. All rights reserved.

//go:build !cgo
// +build !cgo

package kcl_plugin

const CgoEnabled = false

func Invoke(method string, args []interface{}, kwargs map[string]interface{}) (result_json string) {
	panic("unsupport")
}
