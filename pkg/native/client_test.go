//go:build !windows
// +build !windows

package native

import (
	"path"
	"runtime"
	"testing"

	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const code = `
import kcl_plugin.hello

name = "kcl"
sum = hello.add(option("a"), option("b"))
`

func TestExecProgramWithPlugin(t *testing.T) {
	// TODO: windows support
	if runtime.GOOS != "windows" {
		client := NewNativeServiceClient()
		result, err := client.ExecProgram(&gpyrpc.ExecProgram_Args{
			KFilenameList: []string{"main.k"},
			KCodeList:     []string{code},
			Args: []*gpyrpc.CmdArgSpec{
				{
					Name:  "a",
					Value: "1",
				},
				{
					Name:  "b",
					Value: "2",
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if result.ErrMessage != "" {
			t.Fatal("error message must be empty")
		}
	}
}

func TestExecArtifactWithPlugin(t *testing.T) {
	// TODO: windows support
	if runtime.GOOS != "windows" {
		output := path.Join(t.TempDir(), "example")
		client := NewNativeServiceClient()
		// BuildProgram
		buildResult, err := client.BuildProgram(&gpyrpc.BuildProgram_Args{
			ExecArgs: &gpyrpc.ExecProgram_Args{
				KFilenameList: []string{"main.k"},
				KCodeList:     []string{code},
			},
			Output: output,
		})
		if err != nil {
			t.Fatal(err)
		}
		// ExecArtifact
		execResult, err := client.ExecArtifact(&gpyrpc.ExecArtifact_Args{
			ExecArgs: &gpyrpc.ExecProgram_Args{
				Args: []*gpyrpc.CmdArgSpec{
					{
						Name:  "a",
						Value: "1",
					},
					{
						Name:  "b",
						Value: "2",
					},
				},
			},
			Path: buildResult.Path,
		})
		if err != nil {
			t.Fatal(err)
		}
		if execResult.ErrMessage != "" {
			t.Fatal("error message must be empty")
		}
	}
}
