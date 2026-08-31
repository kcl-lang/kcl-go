//go:build !rpc && cgo
// +build !rpc,cgo

package kcl

import (
	"fmt"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"kcl-lang.io/kcl-go/pkg/plugin"
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin"
)

const code = `
import kcl_plugin.hello

name = "kcl"
sum = hello.add(option("a"), option("b"))
`
const codeWithPlugin = `
import kcl_plugin.my_plugin

value1 = my_plugin.config_append({key1 = "value1"}, "key2", "value2")
value2 = my_plugin.list_append([1, 2, 3], 4)
`

func TestNativeRun(t *testing.T) {
	yaml := MustRun("main.k", WithCode(code), WithOptions("a=1", "b=2")).GetRawYamlResult()
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

	yaml := MustRun("main.k", WithCode(codeWithPlugin)).GetRawYamlResult()
	assert2.Equal(t, yaml, "value1:\n  key1: value1\n  key2: value2\nvalue2:\n- 1\n- 2\n- 3\n- 4")
}

// TestXAMLStringMocked exercises the XAMLString pipeline without running
// the KCL runtime, using NewResult to fabricate a result that already
// carries the marker shape PR-3 will emit. Until PR-3 lands, the marker
// cannot come from the runtime, so we exercise the renderer end-to-end
// against hand-built Go values.
func TestXAMLStringMocked(t *testing.T) {
	res := NewResult(map[string]any{
		"TextView": map[string]any{
			"__kcl_info_meta__": []any{"android:id", "android:text"},
			"android:id":        "@+id/userNameTextView",
			"android:text":      "用户名",
		},
	})
	got := res.XAMLString()
	assert2.True(t, strings.Contains(got, `android:id="`), "expected android:id attribute, got %q", got)
	assert2.True(t, strings.Contains(got, `android:text="`), "expected android:text attribute, got %q", got)
	assert2.False(t, strings.Contains(got, "__kcl_info_meta__"), "marker key must not be emitted, got %q", got)
	assert2.False(t, strings.Contains(got, "<android:id>"), "attr key must not be emitted as child, got %q", got)
}

// TestXAMLStringEmpty verifies that XAMLString returns the empty string
// for nil/empty results instead of panicking.
func TestXAMLStringEmpty(t *testing.T) {
	var r *KCLResult
	assert2.Equal(t, "", r.XAMLString())
}
