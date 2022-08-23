// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var (
	Debug        bool
	pyrpcRuntime *Runtime
	gorpcRuntime *Runtime
	once         sync.Once
)

func InitRuntime(maxProc int) {
	once.Do(func() { initRuntime(maxProc) })
}

func GetRuntime() *Runtime {
	once.Do(func() { initRuntime(0) })
	if strings.EqualFold(os.Getenv("KCLVM_SERVICE_CLIENT_HANDLER"), "native") {
		return gorpcRuntime
	} else {
		return pyrpcRuntime
	}
}

func GetGoRuntime() *Runtime {
	once.Do(func() { initRuntime(0) })
	return gorpcRuntime
}

func GetPyRuntime() *Runtime {
	once.Do(func() { initRuntime(0) })
	return pyrpcRuntime
}

func initRuntime(maxProc int) {
	if maxProc <= 0 {
		maxProc = 2
	}
	if maxProc > runtime.NumCPU()*2 {
		maxProc = runtime.NumCPU() * 2
	}

	if g_Python3Path == "" {
		panic(ErrPython3NotFound)
	}
	if g_KclvmRoot == "" {
		panic(ErrKclvmRootNotFound)
	}

	if strings.HasSuffix(g_Python3Path, "kclvm") || strings.HasSuffix(g_Python3Path, "kclvm.exe") {
		os.Setenv("PYTHONHOME", "")
		os.Setenv("PYTHONPATH", "")
	} else {
		os.Setenv("PYTHONHOME", "")
		os.Setenv("PYTHONPATH", filepath.Join(g_KclvmRoot, "lib", "site-packages"))
	}
	//todo add kclvm_capi gorpc server
	gorpcRuntime = NewRuntime(int(maxProc), findKclvmRoot()+"/bin/gorpc")
	pyrpcRuntime = NewRuntime(int(maxProc), MustGetKclvmPath(), "-m", "kclvm.program.rpc-server")
	gorpcRuntime.Start()
	pyrpcRuntime.Start()

	client := &BuiltinServiceClient{
		Runtime: pyrpcRuntime,
	}

	// ping
	{
		args := &gpyrpc.Ping_Args{Value: "ping: kcl-go rest-server"}
		resp, err := client.Ping(args)
		if err != nil {
			fmt.Println("KclvmRuntime: ping failed")
			fmt.Println("kclvm path:", MustGetKclvmPath())
			panic(err)
		}
		if resp.Value != args.Value {
			fmt.Println("KclvmRuntime: ping failed, resp =", resp)
			fmt.Println("kclvm path:", MustGetKclvmPath())
			panic("ping failed")
		}
	}
}
