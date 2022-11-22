//go:build cgo
// +build cgo

package capicall

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func TestPing(t *testing.T) {
	client := PROTOCAPI_NewKclvmServiceClient()
	err := client.Ping(nil, nil)
	assert.Nil(t, err)
	out := &gpyrpc.Ping_Result{}
	err = client.Ping(&gpyrpc.Ping_Args{Value: "hello"}, out)
	assert.Nil(t, err)
	assert.Equal(t, "hello", out.Value)
}

func TestExecProgram(t *testing.T) {
	workdir, err := filepath.Abs("./exec_data/correct_data")
	if err != nil {
		t.Fatal(err)
	}
	args := &gpyrpc.ExecProgram_Args{
		WorkDir:       workdir,
		KFilenameList: []string{"complex.k"},
		Args: []*gpyrpc.CmdArgSpec{
			{Name: "__kcl_test_run", Value: "___test_schema_@@@__"},
			{Name: "__kcl_test_debug", Value: "true"},
		},
		Overrides:         []*gpyrpc.CmdOverrideSpec{},
		DisableYamlResult: false,
	}
	client := PROTOCAPI_NewKclvmServiceClient()
	err = client.ExecProgram(args, &gpyrpc.ExecProgram_Result{})
	if err != nil {
		t.Fatal(err)
	}
}
