// Copyright 2023 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package plugin

import (
	"kcl-lang.io/lib/go/plugin"
)

var (
	GetInvokeJsonProxyPtr = plugin.GetInvokeJsonProxyPtr
	Invoke                = plugin.Invoke
	InvokeJson            = plugin.InvokeJson
)
