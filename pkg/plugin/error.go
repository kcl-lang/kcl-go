// Copyright 2023 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package plugin

import (
	"kcl-lang.io/lib/go/plugin"
)

type PanicInfo = plugin.PanicInfo

type BacktraceFrame = plugin.BacktraceFrame

var (
	JSONError = plugin.JSONError
)
