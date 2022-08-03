// Copyright 2021 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package kclvm_runtime

/*
#include <stdint.h>
#include <stdlib.h>

const char* cgo_kclvm_cli_run(uint64_t f, char *args, char *plugin_agent) {
	typedef const char* (*kclvm_cli_run_fn_t)(const char *, const char *);
	kclvm_cli_run_fn_t kclvm_cli_run = (kclvm_cli_run_fn_t)((void*)(f));
	return (*kclvm_cli_run)(args, plugin_agent);
}
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"unsafe"

	"kusionstack.io/kclvm-go/pkg/3rdparty/dlopen"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

const cgoEnabled = true

func CGO_kclvm_cli_run(args *gpyrpc.ExecProgram_Args, plugin_agent uintptr) string {
	b, err := json.MarshalIndent(args, "", "    ")
	if err != nil {
		panic(err)
	}

	return cgo_kclvm_cli_run(string(b), plugin_agent)
}

func cgo_kclvm_cli_run(args string, plugin_agent uintptr) string {
	var libpath string
	switch runtime.GOOS {
	case "darwin":
		libpath = filepath.Join(MustGetKclvmRoot(), "bin", "libkclvm_cli_cdylib.dylib")
	case "linux":
		libpath = filepath.Join(MustGetKclvmRoot(), "bin", "libkclvm_cli_cdylib.so")
	default:
		panic("unsupport GOOS: " + runtime.GOOS)
	}

	h, err := dlopen.GetHandle([]string{libpath})
	if err != nil {
		panic(fmt.Errorf(`couldn't get a handle to the library: %v`, err))
	}
	defer h.Close()

	f, err := h.GetSymbolPointer("kclvm_cli_run")
	if err != nil {
		panic(fmt.Errorf(`couldn't get symbol %q: %v`, f, err))
	}

	cargs := C.CString(args)
	defer C.free(unsafe.Pointer(cargs))

	cs := C.cgo_kclvm_cli_run(C.uint64_t(uintptr(unsafe.Pointer(h))), cargs, (*C.char)(unsafe.Pointer(plugin_agent)))
	if cs == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cs))

	return C.GoString(cs)
}
