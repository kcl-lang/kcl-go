//go:build cgo
// +build cgo

package native

/*
#include <stdlib.h>
#include <stdint.h>
typedef struct kclvm_service kclvm_service;
*/
import "C"
import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"kcl-lang.io/kcl-go/pkg/3rdparty/dlopen"
	"kcl-lang.io/kcl-go/pkg/plugin"
	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type validator interface {
	Validate() error
}

var libInit sync.Once

var lib *dlopen.LibHandle

type NativeServiceClient struct {
	client *C.kclvm_service
}

var _ service.KclvmService = (*NativeServiceClient)(nil)

func NewNativeServiceClient() *NativeServiceClient {
	libInit.Do(func() {
		lib = loadServiceNativeLib()
	})
	c := new(NativeServiceClient)
	c.client = NewKclvmService(C.uint64_t(plugin.GetInvokeJsonProxyPtr()))
	runtime.SetFinalizer(c, func(x *NativeServiceClient) {
		DeleteKclvmService(x.client)
		x.client = nil
	})
	return c
}

func cApiCall[I interface {
	*TI
	proto.Message
}, O interface {
	*TO
	protoreflect.ProtoMessage
}, TI any, TO any](c *NativeServiceClient, callName string, in I) (O, error) {
	if callName == "" {
		return nil, errors.New("kclvm service c api call: empty method name")
	}

	if in == nil {
		in = new(TI)
	}

	if x, ok := proto.Message(in).(validator); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}
	inBytes, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	cCallName := C.CString(callName)

	defer C.free(unsafe.Pointer(cCallName))

	cIn := C.CString(string(inBytes))

	defer C.free(unsafe.Pointer(cIn))

	cOut := KclvmServiceCall(c.client, cCallName, cIn)

	defer KclvmServiceFreeString(cOut)

	msg := C.GoString(cOut)

	if strings.HasPrefix(msg, "ERROR:") {
		return nil, errors.New(strings.TrimPrefix(msg, "ERROR:"))
	}

	var out O = new(TO)
	err = proto.Unmarshal([]byte(msg), out)
	if err != nil {
		return nil, errors.New(msg)
	}

	return out, nil
}

func (c *NativeServiceClient) Ping(in *gpyrpc.Ping_Args) (*gpyrpc.Ping_Result, error) {
	return cApiCall[*gpyrpc.Ping_Args, *gpyrpc.Ping_Result](c, "KclvmService.Ping", in)
}

func (c *NativeServiceClient) ExecProgram(in *gpyrpc.ExecProgram_Args) (*gpyrpc.ExecProgram_Result, error) {
	return cApiCall[*gpyrpc.ExecProgram_Args, *gpyrpc.ExecProgram_Result](c, "KclvmService.ExecProgram", in)
}

func (c *NativeServiceClient) BuildProgram(in *gpyrpc.BuildProgram_Args) (*gpyrpc.BuildProgram_Result, error) {
	return cApiCall[*gpyrpc.BuildProgram_Args, *gpyrpc.BuildProgram_Result](c, "KclvmService.BuildProgram", in)
}

func (c *NativeServiceClient) ExecArtifact(in *gpyrpc.ExecArtifact_Args) (*gpyrpc.ExecProgram_Result, error) {
	return cApiCall[*gpyrpc.ExecArtifact_Args, *gpyrpc.ExecProgram_Result](c, "KclvmService.ExecArtifact", in)
}

func (c *NativeServiceClient) ParseFile(in *gpyrpc.ParseFile_Args) (*gpyrpc.ParseFile_Result, error) {
	return cApiCall[*gpyrpc.ParseFile_Args, *gpyrpc.ParseFile_Result](c, "KclvmService.ParseFile", in)
}

func (c *NativeServiceClient) ParseProgram(in *gpyrpc.ParseProgram_Args) (*gpyrpc.ParseProgram_Result, error) {
	return cApiCall[*gpyrpc.ParseProgram_Args, *gpyrpc.ParseProgram_Result](c, "KclvmService.ParseProgram", in)
}

func (c *NativeServiceClient) ListOptions(in *gpyrpc.ParseProgram_Args) (*gpyrpc.ListOptions_Result, error) {
	return cApiCall[*gpyrpc.ParseProgram_Args, *gpyrpc.ListOptions_Result](c, "KclvmService.ListOptions", in)
}

func (c *NativeServiceClient) LoadPackage(in *gpyrpc.LoadPackage_Args) (*gpyrpc.LoadPackage_Result, error) {
	return cApiCall[*gpyrpc.LoadPackage_Args, *gpyrpc.LoadPackage_Result](c, "KclvmService.LoadPackage", in)
}

func (c *NativeServiceClient) FormatCode(in *gpyrpc.FormatCode_Args) (*gpyrpc.FormatCode_Result, error) {
	return cApiCall[*gpyrpc.FormatCode_Args, *gpyrpc.FormatCode_Result](c, "KclvmService.FormatCode", in)
}

func (c *NativeServiceClient) FormatPath(in *gpyrpc.FormatPath_Args) (*gpyrpc.FormatPath_Result, error) {
	return cApiCall[*gpyrpc.FormatPath_Args, *gpyrpc.FormatPath_Result](c, "KclvmService.FormatPath", in)
}

func (c *NativeServiceClient) LintPath(in *gpyrpc.LintPath_Args) (*gpyrpc.LintPath_Result, error) {
	return cApiCall[*gpyrpc.LintPath_Args, *gpyrpc.LintPath_Result](c, "KclvmService.LintPath", in)
}

func (c *NativeServiceClient) OverrideFile(in *gpyrpc.OverrideFile_Args) (*gpyrpc.OverrideFile_Result, error) {
	return cApiCall[*gpyrpc.OverrideFile_Args, *gpyrpc.OverrideFile_Result](c, "KclvmService.OverrideFile", in)
}

func (c *NativeServiceClient) GetSchemaType(in *gpyrpc.GetSchemaType_Args) (*gpyrpc.GetSchemaType_Result, error) {
	return cApiCall[*gpyrpc.GetSchemaType_Args, *gpyrpc.GetSchemaType_Result](c, "KclvmService.GetSchemaType", in)
}

func (c *NativeServiceClient) GetSchemaTypeMapping(in *gpyrpc.GetSchemaTypeMapping_Args) (*gpyrpc.GetSchemaTypeMapping_Result, error) {
	return cApiCall[*gpyrpc.GetSchemaTypeMapping_Args, *gpyrpc.GetSchemaTypeMapping_Result](c, "KclvmService.GetSchemaTypeMapping", in)
}

func (c *NativeServiceClient) GetFullSchemaType(in *gpyrpc.GetFullSchemaType_Args) (*gpyrpc.GetSchemaType_Result, error) {
	return cApiCall[*gpyrpc.GetFullSchemaType_Args, *gpyrpc.GetSchemaType_Result](c, "KclvmService.GetFullSchemaType", in)
}

func (c *NativeServiceClient) ValidateCode(in *gpyrpc.ValidateCode_Args) (*gpyrpc.ValidateCode_Result, error) {
	return cApiCall[*gpyrpc.ValidateCode_Args, *gpyrpc.ValidateCode_Result](c, "KclvmService.ValidateCode", in)
}

func (c *NativeServiceClient) ListDepFiles(in *gpyrpc.ListDepFiles_Args) (*gpyrpc.ListDepFiles_Result, error) {
	return cApiCall[*gpyrpc.ListDepFiles_Args, *gpyrpc.ListDepFiles_Result](c, "KclvmService.ListDepFiles", in)
}

func (c *NativeServiceClient) LoadSettingsFiles(in *gpyrpc.LoadSettingsFiles_Args) (*gpyrpc.LoadSettingsFiles_Result, error) {
	return cApiCall[*gpyrpc.LoadSettingsFiles_Args, *gpyrpc.LoadSettingsFiles_Result](c, "KclvmService.LoadSettingsFiles", in)
}

func (c *NativeServiceClient) Rename(in *gpyrpc.Rename_Args) (*gpyrpc.Rename_Result, error) {
	return cApiCall[*gpyrpc.Rename_Args, *gpyrpc.Rename_Result](c, "KclvmService.Rename", in)
}

func (c *NativeServiceClient) RenameCode(in *gpyrpc.RenameCode_Args) (*gpyrpc.RenameCode_Result, error) {
	return cApiCall[*gpyrpc.RenameCode_Args, *gpyrpc.RenameCode_Result](c, "KclvmService.RenameCode", in)
}

func (c *NativeServiceClient) Test(in *gpyrpc.Test_Args) (*gpyrpc.Test_Result, error) {
	return cApiCall[*gpyrpc.Test_Args, *gpyrpc.Test_Result](c, "KclvmService.Test", in)
}
