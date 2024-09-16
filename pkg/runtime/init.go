// Copyright The KCL Authors. All rights reserved.

package runtime

import (
	"fmt"
	"runtime"
	"sync"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

var (
	Debug      bool
	rpcRuntime *Runtime
	once       sync.Once
)

const tip = "Tip: Have you used a binary version of KCL in your PATH that is not consistent with the KCL Go SDK? You can upgrade or reduce the KCL version or delete the KCL in your PATH"

func InitRuntime(maxProc int) {
	once.Do(func() { initRuntime(maxProc) })
}

func GetRuntime() *Runtime {
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

	rpcRuntime = NewRuntime(int(maxProc), "kclvm_cli", "server")
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
			fmt.Printf("If not, you can run `rm -r %s/bin` to fix this issue\n", MustGetKclvmPath())
			panic(err)
		}
	}
}
