package kcl

import (
	"kcl-lang.io/lib/go/api"
	"kcl-lang.io/lib/go/native"
)

// Service returns the interaction interface between KCL Go SDK and KCL Rust core.
func Service() api.ServiceClient {
	return native.NewNativeServiceClient()
}

// serviceWithPluginAgent returns the KCL service client, built with the
// given plugin agent pointer when one is set on the option.
func serviceWithPluginAgent(pluginAgent *uint64) api.ServiceClient {
	if pluginAgent == nil {
		return Service()
	}
	return native.NewNativeServiceClientWithPluginAgent(*pluginAgent)
}
