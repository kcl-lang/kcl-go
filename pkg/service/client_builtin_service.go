//go:build rpc || !cgo
// +build rpc !cgo

// Copyright The KCL Authors. All rights reserved.

package service

import (
	"kcl-lang.io/kcl-go/pkg/runtime"
)

type BuiltinServiceClient = runtime.BuiltinServiceClient

func NewBuiltinServiceClient() *BuiltinServiceClient {
	return &BuiltinServiceClient{
		Runtime: runtime.GetRuntime(),
	}
}
