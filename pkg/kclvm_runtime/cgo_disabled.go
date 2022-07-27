// Copyright 2021 The KCL Authors. All rights reserved.

//go:build !cgo
// +build !cgo

package kclvm_runtime

import (
	"encoding/json"
	"unsafe"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

const cgoEnabled = false

func CGO_kclvm_cli_run(args *gpyrpc.ExecProgram_Args, plugin_agent uintptr) string {
	panic("unsupport cgo")
}

func cgo_kclvm_cli_run(args string, plugin_agent uintptr) string {
	panic("unsupport cgo")
}
