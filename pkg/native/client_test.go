//go:build !windows
// +build !windows

package native

import (
	"runtime"
	"testing"

	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const code = `
import kcl_plugin.hello

name = "kcl"
three = hello.add(1,2)
`

func TestExecProgramWithPlugin(t *testing.T) {
	// TODO: windows support
	if runtime.GOOS != "windows" {
		client := NewNativeServiceClient()
		result, err := client.ExecProgram(&gpyrpc.ExecProgram_Args{
			KFilenameList: []string{"main.k"},
			KCodeList:     []string{code},
		})
		if err != nil {
			t.Fatal(err)
		}
		if result.ErrMessage != "" {
			t.Fatal("error message must be empty")
		}
	}
}
