package capicall

/*
#include <stdlib.h>
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
	"kusionstack.io/kclvm-go/gorpc/pkg/3rdparty/dlopen"
	"kusionstack.io/kclvm-go/gorpc/pkg/kcl_plugin"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var libInit sync.Once

var lib *dlopen.LibHandle

type PROTOCAPI_KclvmServiceClient struct {
	client *C.kclvm_service
}

func PROTOCAPI_NewKclvmServiceClient() *PROTOCAPI_KclvmServiceClient {
	libInit.Do(func() {
		lib = loadKclvmServiceCapiLib()
	})
	c := new(PROTOCAPI_KclvmServiceClient)
	c.client = NewKclvmService(C.longlong(kcl_plugin.GetInvokeJsonProxyPtr()))
	runtime.SetFinalizer(c, func(x *PROTOCAPI_KclvmServiceClient) {
		DeleteKclvmService(x.client)
		x.client = nil
	})
	return c
}

func (c *PROTOCAPI_KclvmServiceClient) cApiCall(callName string, in proto.Message, out protoreflect.ProtoMessage) error {
	type Validator interface {
		Validate() error
	}

	if callName == "" {
		return errors.New("kclvm service c api call : empty method name")
	}

	if in == nil {
		return errors.New("kclvm service c api call : unknown proto input type")
	}

	if out == nil {
		return errors.New("kclvm service c api call : unknown proto output type")
	}

	if x, ok := in.(Validator); ok {
		if err := x.Validate(); err != nil {
			return err
		}
	}
	inBytes, err := proto.Marshal(in)
	if err != nil {
		return err
	}

	cCallName := C.CString(callName)

	defer C.free(unsafe.Pointer(cCallName))

	cIn := C.CString(string(inBytes))

	defer C.free(unsafe.Pointer(cIn))

	cOut := KclvmServiceCall(c.client, cCallName, cIn)

	defer KclvmServiceFreeString(cOut)

	goOut := C.GoString(cOut)

	if strings.HasPrefix(goOut, "KCLVM_CAPI_CALL_ERROR") {
		return errors.New(goOut)
	}

	return proto.Unmarshal([]byte(goOut), out)
}

func (c *PROTOCAPI_KclvmServiceClient) Ping(in *gpyrpc.Ping_Args, out *gpyrpc.Ping_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.Ping_Args)
	}

	err = c.cApiCall("KclvmService.Ping", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ParseFile_LarkTree(in *gpyrpc.ParseFile_LarkTree_Args, out *gpyrpc.ParseFile_LarkTree_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ParseFile_LarkTree_Args)
	}

	err = c.cApiCall("KclvmService.ParseFile_LarkTree", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ParseFile_AST(in *gpyrpc.ParseFile_AST_Args, out *gpyrpc.ParseFile_LarkTree_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ParseFile_AST_Args)
	}

	err = c.cApiCall("KclvmService.ParseFile_AST", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ParseProgram_AST(in *gpyrpc.ParseProgram_AST_Args, out *gpyrpc.ParseProgram_AST_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ParseProgram_AST_Args)
	}

	err = c.cApiCall("KclvmService.ParseProgram_AST", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ExecProgram(in *gpyrpc.ExecProgram_Args, out *gpyrpc.ExecProgram_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ExecProgram_Args)
	}

	err = c.cApiCall("KclvmService.ExecProgram", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ResetPlugin(in *gpyrpc.ResetPlugin_Args, out *gpyrpc.ResetPlugin_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ResetPlugin_Args)
	}

	err = c.cApiCall("KclvmService.ResetPlugin", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) FormatCode(in *gpyrpc.FormatCode_Args, out *gpyrpc.FormatCode_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.FormatCode_Args)
	}

	err = c.cApiCall("KclvmService.FormatCode", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) FormatPath(in *gpyrpc.FormatPath_Args, out *gpyrpc.FormatPath_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.FormatPath_Args)
	}

	err = c.cApiCall("KclvmService.FormatPath", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) LintPath(in *gpyrpc.LintPath_Args, out *gpyrpc.LintPath_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.LintPath_Args)
	}

	err = c.cApiCall("KclvmService.LintPath", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) OverrideFile(in *gpyrpc.OverrideFile_Args, out *gpyrpc.OverrideFile_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.OverrideFile_Args)
	}

	err = c.cApiCall("KclvmService.OverrideFile", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) EvalCode(in *gpyrpc.EvalCode_Args, out *gpyrpc.EvalCode_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.EvalCode_Args)
	}

	err = c.cApiCall("KclvmService.EvalCode", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ResolveCode(in *gpyrpc.ResolveCode_Args, out *gpyrpc.ResolveCode_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ResolveCode_Args)
	}

	err = c.cApiCall("KclvmService.ResolveCode", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) GetSchemaType(in *gpyrpc.GetSchemaType_Args, out *gpyrpc.GetSchemaType_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.GetSchemaType_Args)
	}

	err = c.cApiCall("KclvmService.GetSchemaType", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ValidateCode(in *gpyrpc.ValidateCode_Args, out *gpyrpc.ValidateCode_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ValidateCode_Args)
	}

	err = c.cApiCall("KclvmService.ValidateCode", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) SpliceCode(in *gpyrpc.SpliceCode_Args, out *gpyrpc.SpliceCode_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.SpliceCode_Args)
	}

	err = c.cApiCall("KclvmService.SpliceCode", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) Complete(in *gpyrpc.Complete_Args, out *gpyrpc.Complete_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.Complete_Args)
	}

	err = c.cApiCall("KclvmService.Complete", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) GoToDef(in *gpyrpc.GoToDef_Args, out *gpyrpc.GoToDef_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.GoToDef_Args)
	}

	err = c.cApiCall("KclvmService.GoToDef", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) DocumentSymbol(in *gpyrpc.DocumentSymbol_Args, out *gpyrpc.DocumentSymbol_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.DocumentSymbol_Args)
	}

	err = c.cApiCall("KclvmService.DocumentSymbol", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) Hover(in *gpyrpc.Hover_Args, out *gpyrpc.Hover_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.Hover_Args)
	}

	err = c.cApiCall("KclvmService.Hover", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ListDepFiles(in *gpyrpc.ListDepFiles_Args, out *gpyrpc.ListDepFiles_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ListDepFiles_Args)
	}

	err = c.cApiCall("KclvmService.ListDepFiles", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ListUpStreamFiles(in *gpyrpc.ListUpStreamFiles_Args, out *gpyrpc.ListUpStreamFiles_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ListUpStreamFiles_Args)
	}

	err = c.cApiCall("KclvmService.ListUpStreamFiles", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) ListDownStreamFiles(in *gpyrpc.ListDownStreamFiles_Args, out *gpyrpc.ListUpStreamFiles_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.ListDownStreamFiles_Args)
	}

	err = c.cApiCall("KclvmService.ListDownStreamFiles", in, out)

	return
}

func (c *PROTOCAPI_KclvmServiceClient) LoadSettingsFiles(in *gpyrpc.LoadSettingsFiles_Args, out *gpyrpc.LoadSettingsFiles_Result) (err error) {
	if in == nil {
		in = new(gpyrpc.LoadSettingsFiles_Args)
	}

	err = c.cApiCall("KclvmService.LoadSettingsFiles", in, out)

	return
}
