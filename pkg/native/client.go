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
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

var libInit sync.Once

var lib *dlopen.LibHandle

type NativeServiceClient struct {
	client *C.kclvm_service
}

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

func (c *NativeServiceClient) cApiCall(callName string, in proto.Message, out protoreflect.ProtoMessage) error {
	type Validator interface {
		Validate() error
	}

	if callName == "" {
		return errors.New("kclvm service c api call: empty method name")
	}

	if in == nil {
		return errors.New("kclvm service c api call: unknown proto input type")
	}

	if out == nil {
		return errors.New("kclvm service c api call: unknown proto output type")
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

	msg := C.GoString(cOut)

	if strings.HasPrefix(msg, "ERROR:") {
		return errors.New(strings.TrimPrefix(msg, "ERROR:"))
	}

	err = proto.Unmarshal([]byte(msg), out)
	if err != nil {
		return errors.New(msg)
	}

	return nil
}

func (c *NativeServiceClient) ExecProgram(in *gpyrpc.ExecProgram_Args) (*gpyrpc.ExecProgram_Result, error) {
	if in == nil {
		in = new(gpyrpc.ExecProgram_Args)
	}

	out := new(gpyrpc.ExecProgram_Result)
	err := c.cApiCall("KclvmService.ExecProgram", in, out)

	return out, err
}

func (c *NativeServiceClient) BuildProgram(in *gpyrpc.BuildProgram_Args) (*gpyrpc.BuildProgram_Result, error) {
	if in == nil {
		in = new(gpyrpc.BuildProgram_Args)
	}
	out := new(gpyrpc.BuildProgram_Result)
	err := c.cApiCall("KclvmService.BuildProgram", in, out)

	return out, err
}

func (c *NativeServiceClient) ExecArtifact(in *gpyrpc.ExecArtifact_Args) (*gpyrpc.ExecProgram_Result, error) {
	if in == nil {
		in = new(gpyrpc.ExecArtifact_Args)
	}
	out := new(gpyrpc.ExecProgram_Result)
	err := c.cApiCall("KclvmService.ExecArtifact", in, out)

	return out, err
}
