package kcl

import (
	"kcl-lang.io/lib/go/api"
	"kcl-lang.io/lib/go/native"
)

// Service returns the interaction interface between KCL Go SDK and KCL Rust core.
func Service() api.ServiceClient {
	return native.NewNativeServiceClient()
}
