// Copyright The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package plugin

import (
	"kcl-lang.io/lib/go/plugin"
)

// Plugin represents a KCL Plugin with metadata and methods.
// It contains the plugin name, version, a reset function, and a map of methods.
type Plugin = plugin.Plugin

// MethodSpec defines the specification for a KCL Plugin method.
// It includes the method type and the body function which executes the method logic.
type MethodSpec = plugin.MethodSpec

// MethodType describes the type of a KCL Plugin method's arguments, keyword arguments, and result.
// It specifies the types of positional arguments, keyword arguments, and the result type.
type MethodType = plugin.MethodType

// MethodArgs represents the arguments passed to a KCL Plugin method.
// It includes a list of positional arguments and a map of keyword arguments.
type MethodArgs = plugin.MethodArgs

// MethodResult represents the result returned from a KCL Plugin method.
// It holds the value of the result.
type MethodResult = plugin.MethodResult

var (
	// ParseMethodArgs parses JSON strings for positional and keyword arguments
	// and returns a MethodArgs object.
	// args_json: JSON string of positional arguments
	// kwargs_json: JSON string of keyword arguments
	ParseMethodArgs = plugin.ParseMethodArgs
)
