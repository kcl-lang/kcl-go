// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	"fmt"
	"runtime"
	"sync"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var (
	Debug        bool
	pyrpcRuntime *Runtime
	once         sync.Once
)

func InitRuntime(maxProc int) {
	once.Do(func() { initRuntime(maxProc) })
}

func GetRuntime() *Runtime {
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

	pyrpcRuntime = NewRuntime(int(maxProc), MustGetKclvmPath(), "-m", "kclvm.program.rpc-server")
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
