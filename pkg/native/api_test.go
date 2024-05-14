//go:build cgo
// +build cgo

package native

import (
	"fmt"
	"testing"

	"kcl-lang.io/kcl-go/pkg/kcl"
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin"
)

func TestNativeRun(t *testing.T) {
	yaml := MustRun("main.k", kcl.WithCode(code), kcl.WithOptions("a=1", "b=2")).GetRawYamlResult()
	fmt.Println(yaml)
}
