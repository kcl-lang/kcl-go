package native

import (
	"fmt"
	"path/filepath"
	"runtime"

	"kcl-lang.io/kcl-go/pkg/3rdparty/dlopen"
	kcl_runtime "kcl-lang.io/kcl-go/pkg/runtime"
	"kcl-lang.io/kcl-go/pkg/utils"
)

const libName = "kclvm_cli_cdylib"

func loadServiceNativeLib() *dlopen.LibHandle {
	root := kcl_runtime.MustGetKclvmRoot()
	libPaths := []string{}
	sysType := runtime.GOOS
	fullLibName := "lib" + libName + ".so"

	if sysType == "darwin" {
		fullLibName = "lib" + libName + ".dylib"
	} else if sysType == "windows" {
		fullLibName = libName + ".dll"
	}

	libPath := filepath.Join(root, "bin", fullLibName)
	if !utils.FileExists(libPath) {
		libPath = filepath.Join(root, "lib", fullLibName)
	}

	libPaths = append(libPaths, libPath)

	h, err := dlopen.GetHandle(libPaths)

	if err != nil {
		panic(fmt.Errorf(`couldn't get a handle to kcl native library: %v`, err))
	}

	runtime.SetFinalizer(h, func(x *dlopen.LibHandle) {
		x.Close()
	})

	return h
}
