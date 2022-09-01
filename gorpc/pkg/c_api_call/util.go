package capicall

import (
	"fmt"
	"runtime"
	"strings"

	"kusionstack.io/kclvm-go/gorpc/pkg/3rdparty/dlopen"
	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
)

func loadKclvmServiceCapiLib() *dlopen.LibHandle {

	kclvmRoot, err := kclvm_runtime.GetKclvmRoot()
	if err != nil {
		panic(fmt.Errorf("kclvm_capi : can't find kclvm_capi lib path"))
	}
	libPaths := []string{}
	sysType := runtime.GOOS
	libSuffix := ".so"

	if sysType == "darwin" {
		libSuffix = ".dylib"
	} else if sysType == "windows" {
		libSuffix = ".dll"
	}

	libName := "libkclvm_capi" + libSuffix
	kclvmLibPath := strings.Join([]string{kclvmRoot, "lib", libName}, "/")
	libPaths = append(libPaths, kclvmLibPath)

	h, err := dlopen.GetHandle(libPaths)

	if err != nil {
		panic(fmt.Errorf(`couldn't get a handle to kclvm_capi library: %v`, err))
	}

	runtime.SetFinalizer(h, func(x *dlopen.LibHandle) {
		x.Close()
	})

	return h
}
