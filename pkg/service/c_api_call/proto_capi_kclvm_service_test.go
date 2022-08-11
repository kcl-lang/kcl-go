//go:build cgo && kclvm_service_capi
// +build cgo,kclvm_service_capi

package capicall

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func TestPing(t *testing.T) {
	client := PROTOCAPI_NewKclvmServiceClient()
	out, err := client.Ping(nil)
	assert.Nil(t, err)
	out, err = client.Ping(&gpyrpc.Ping_Args{Value: "hello"})
	assert.Nil(t, err)
	assert.Equal(t, "hello", out.Value)
}

func TestExecProgram(t *testing.T) {

	workdir, _ := filepath.Abs(EXEC_DATA_PATH)
	args := &gpyrpc.ExecProgram_Args{
		WorkDir:       workdir,
		KFilenameList: []string{"hello.k"},
		Args: []*gpyrpc.CmdArgSpec{
			{Name: "__kcl_test_run", Value: "___test_schema_@@@__"},
			{Name: "__kcl_test_debug", Value: "true"},
		},
		Overrides:         []*gpyrpc.CmdOverrideSpec{},
		DisableYamlResult: false,
	}
	client := PROTOCAPI_NewKclvmServiceClient()
	out, err := client.ExecProgram(args)
	assert.Nil(t, err)
	assert.NotEmpty(t, out.JsonResult)
	assert.NotEmpty(t, out.YamlResult)
}
