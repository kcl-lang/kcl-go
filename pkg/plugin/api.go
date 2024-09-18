// Copyright 2023 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package plugin

import (
	"kcl-lang.io/lib/go/plugin"
)

var (
	// Register register a new kcl plugin.
	RegisterPlugin = plugin.RegisterPlugin
	// GetPlugin get plugin object by name.
	GetPlugin = plugin.GetPlugin
	// GetMethodSpec get plugin method by name.
	GetMethodSpec = plugin.GetMethodSpec
	// ResetPlugin reset all kcl plugin state.
	ResetPlugin = plugin.ResetPlugin
)
