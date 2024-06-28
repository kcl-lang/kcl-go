// Copyright 2023 The KCL Authors. All rights reserved.

package plugin

import (
	"strings"
	"testing"
)

func init() {
	RegisterPlugin(Plugin{
		Name:      "strings",
		ResetFunc: func() {},
		MethodMap: map[string]MethodSpec{
			"join": {
				Type: &MethodType{},
				Body: func(args *MethodArgs) (*MethodResult, error) {
					var ss []string
					for i := range args.Args {
						ss = append(ss, args.StrArg(i))
					}
					return &MethodResult{strings.Join(ss, ".")}, nil
				},
			},
			"panic": {
				Type: &MethodType{},
				Body: func(args *MethodArgs) (*MethodResult, error) {
					panic(args.Args)
				},
			},
		},
	})
}

func TestPlugin_strings_join(t *testing.T) {
	if !CgoEnabled {
		t.Skip("cgo disabled")
	}
	result_json := Invoke("kcl_plugin.strings.join", []interface{}{"KCL", "KCL", 123}, nil)
	if result_json != `"KCL.KCL.123"` {
		t.Fatal(result_json)
	}
}

func TestPlugin_strings_panic(t *testing.T) {
	if !CgoEnabled {
		t.Skip("cgo disabled")
	}
	result_json := Invoke("kcl_plugin.strings.panic", []interface{}{"KCL", "KCL", 123}, nil)
	if result_json != `{"__kcl_PanicInfo__":"[KCL KCL 123]"}` {
		t.Fatal(result_json)
	}
}
