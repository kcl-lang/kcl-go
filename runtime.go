//go:build rpc || !cgo
// +build rpc !cgo

package kcl

import "kcl-lang.io/kcl-go/pkg/runtime"

// InitKclvmPath init kclvm path.
func InitKclvmPath(kclvmRoot string) {
	runtime.InitKclvmRoot(kclvmRoot)
}

// InitKclvmRuntime init kclvm process.
func InitKclvmRuntime(n int) {
	runtime.InitRuntime(n)
}
