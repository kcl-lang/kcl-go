//go:build !native
// +build !native

package kcl

import "kcl-lang.io/kcl-go/pkg/service"

// Service returns the interaction interface between KCL Go SDK and KCL Rust core.
// When `go build tags=native` is opened, use CGO and dynamic libraries to interact.
// When closed, use the default RPC interaction logic to avoid CGO usage.
func Service() service.KclvmService {
	return service.NewKclvmServiceClient()
}
