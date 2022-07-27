//go:build ignore
// +build ignore

package main

import (
	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func main() {
	s := kclvm_runtime.CGO_kclvm_cli_run(&gpyrpc.ExecProgram_Args{}, 0)
	fmt.Println(s)
}
