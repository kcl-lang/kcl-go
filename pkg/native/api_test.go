//go:build cgo
// +build cgo

package native

import (
	"fmt"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/plugin"
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin"
)

const codeWithPlugin = `
import kcl_plugin.my_plugin

value1 = my_plugin.config_append({key1 = "value1"}, "key2", "value2")
value2 = my_plugin.list_append([1, 2, 3], 4)
`

func TestNativeRun(t *testing.T) {
	yaml := MustRun("main.k", kcl.WithCode(code), kcl.WithOptions("a=1", "b=2")).GetRawYamlResult()
	fmt.Println(yaml)
}

func TestNativeRunWithPlugin(t *testing.T) {
	plugin.RegisterPlugin(plugin.Plugin{
		Name: "my_plugin",
		MethodMap: map[string]plugin.MethodSpec{
			"config_append": {
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					config := args.MapArg(0)
					k := args.StrArg(1)
					v := args.StrArg(2)
					config[k] = v
					return &plugin.MethodResult{V: config}, nil
				},
			},
			"list_append": {
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					values := args.ListArg(0)
					v := args.Arg(1)
					values = append(values, v)
					return &plugin.MethodResult{V: values}, nil
				},
			},
		},
	})

	yaml := MustRun("main.k", kcl.WithCode(codeWithPlugin)).GetRawYamlResult()
	assert2.Equal(t, yaml, "value1:\n  key1: value1\n  key2: value2\nvalue2:\n- 1\n- 2\n- 3\n- 4")
}
