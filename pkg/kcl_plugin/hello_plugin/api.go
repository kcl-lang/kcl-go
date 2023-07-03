// Copyright 2023 The KCL Authors. All rights reserved.

package hello_plugin

import (
	"fmt"

	"kcl-lang.io/kcl-go/pkg/kcl_plugin"
)

func init() {
	kcl_plugin.RegisterPlugin(kcl_plugin.Plugin{
		Name:      "hello",
		ResetFunc: _reset_plugin,
		MethodMap: _MethodMap,
	})
}

var _MethodMap = map[string]kcl_plugin.MethodSpec{
	// func set_global_int(v int64)
	"set_global_int": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			set_global_int(args.IntArg(0))
			return nil, nil
		},
	},

	// func get_global_int() int64
	"get_global_int": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			v := get_global_int()
			return &kcl_plugin.MethodResult{V: float64(v)}, nil
		},
	},

	// func say_hello(msg string)
	"say_hello": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			msg := fmt.Sprint(args.StrArg(0))
			say_hello(msg)

			return nil, nil
		},
	},

	// func add(a, b int64) int64
	"add": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			v := add(args.IntArg(0), args.IntArg(1))
			return &kcl_plugin.MethodResult{V: float64(v)}, nil
		},
	},

	// func tolower(s string) string
	"tolower": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			v := tolower(args.StrArg(0))
			return &kcl_plugin.MethodResult{V: v}, nil
		},
	},

	// func update_dict(d map[string]interface{}, key, value string) map[string]interface{}
	"update_dict": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			d := args.Arg(0).(map[string]interface{})
			v := update_dict(d, args.StrArg(1), args.StrArg(2))
			return &kcl_plugin.MethodResult{V: v}, nil
		},
	},

	// func list_append(list []interface{}, values ...interface{}) []interface{}
	"list_append": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			d := args.Arg(0).([]interface{})
			v := list_append(d, args.Args[1:]...)
			return &kcl_plugin.MethodResult{V: v}, nil
		},
	},

	// func foo(a, b interface{}, x []interface{}, values map[string]interface{}) interface{} {
	"foo": {
		Type: &kcl_plugin.MethodType{},
		Body: func(args *kcl_plugin.MethodArgs) (*kcl_plugin.MethodResult, error) {
			values := args.KwArgs

			a := args.Arg(0)
			b := args.Arg(1)
			x := args.Args[2:]

			v := foo(a, b, nil, x, values)

			return &kcl_plugin.MethodResult{V: v}, nil
		},
	},
}
