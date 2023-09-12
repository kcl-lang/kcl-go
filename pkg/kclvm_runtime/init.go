// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

var (
	Debug              bool
	rpcRuntime         *Runtime
	once               sync.Once
	UseKCLPluginEnvVar = "KCL_GO_USE_PLUGIN"
)

const tip = "Tip: Have you used a binary version of KCL in your PATH that is not consistent with the KCL Go SDK? You can upgrade or reduce the KCL version or delete the KCL in your PATH"

func InitRuntime(maxProc int) {
	once.Do(func() { initRuntime(maxProc) })
}

func GetRuntime() *Runtime {
	once.Do(func() { initRuntime(0) })
	return rpcRuntime
}

func GetPyRuntime() *Runtime {
	once.Do(func() { initRuntime(0) })
	return rpcRuntime
}

func initRuntime(maxProc int) {
	if maxProc <= 0 {
		maxProc = 2
	}
	if maxProc > runtime.NumCPU()*2 {
		maxProc = runtime.NumCPU() * 2
	}

	if g_KclvmRoot == "" {
		panic(ErrKclvmRootNotFound)
	}

	if os.Getenv(UseKCLPluginEnvVar) != "" {
		os.Setenv("PYTHONHOME", "")
		os.Setenv("PYTHONPATH", filepath.Join(g_KclvmRoot, "lib", "site-packages"))
		rpcRuntime = NewRuntime(int(maxProc), MustGetKclvmPath(), "-m", "kclvm.program.rpc-server")
	} else {
		rpcRuntime = NewRuntime(int(maxProc), "kclvm_cli", "server")
	}

	rpcRuntime.Start()

	client := &BuiltinServiceClient{
		Runtime: rpcRuntime,
	}

	// ping
	{
		args := &gpyrpc.Ping_Args{Value: "ping: kcl-go rest-server"}
		resp, err := client.Ping(args)
		if err != nil || resp.Value != args.Value {
			fmt.Println("Init kcl runtime failed, path: ", MustGetKclvmPath())
			fmt.Println(tip)
			panic(err)
		}
	}
}
