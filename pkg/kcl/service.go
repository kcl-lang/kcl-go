//go:build !rpc && cgo
// +build !rpc,cgo

package kcl

import (
	"kcl-lang.io/lib/go/api"
	"kcl-lang.io/lib/go/native"
)

// Service returns the interaction interface between KCL Go SDK and KCL Rust core.
// When `go build tags=rpc` is opened, use the default RPC interaction logic to avoid CGO usage.
// When closed, use CGO and dynamic libraries to interact.
func Service() api.ServiceClient {
	return native.NewNativeServiceClient()
}
