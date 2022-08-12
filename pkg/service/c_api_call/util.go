package capicall

import (
	"fmt"
	"runtime"
	"strings"

	"kusionstack.io/kclvm-go/pkg/3rdparty/dlopen"
	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
)

func loadKclvmServiceCapiLib() *dlopen.LibHandle {

	_, srcFilePath, _, ok := runtime.Caller(1)
	kclvmRoot, err := kclvm_runtime.GetKclvmRoot()
	if !ok && err != nil {
		panic(fmt.Errorf("kclvm_capi : can't find kclvm_capi lib path"))
	}
	libPaths := []string{}
	sysType := runtime.GOOS
	archType := runtime.GOARCH
	osWithArch := sysType + "-" + archType
	libSuffix := ".so"

	if sysType == "darwin" {
		libSuffix = ".dylib"
	} else if sysType == "windows" {
		libSuffix = ".dll"
	}

	libName := "libkclvm_capi" + libSuffix
	if ok {
		curDir := srcFilePath[:strings.LastIndex(srcFilePath, "/")]

		srcLibPath := strings.Join([]string{curDir, "packaged/lib", osWithArch, libName}, "/")
		libPaths = append(libPaths, srcLibPath)
	}

	if err == nil {
		kclvmLibPath := strings.Join([]string{kclvmRoot, "lib", libName}, "/")
		libPaths = append(libPaths, kclvmLibPath)
	}

	h, err := dlopen.GetHandle(libPaths)

	if err != nil {
		panic(err)
	}

	runtime.SetFinalizer(h, func(x *dlopen.LibHandle) {
		x.Close()
	})

	return h
}
