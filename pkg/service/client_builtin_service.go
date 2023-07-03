// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
)

type BuiltinServiceClient = kclvm_runtime.BuiltinServiceClient

func NewBuiltinServiceClient() *BuiltinServiceClient {
	return &BuiltinServiceClient{
		Runtime: kclvm_runtime.GetRuntime(),
	}
}
