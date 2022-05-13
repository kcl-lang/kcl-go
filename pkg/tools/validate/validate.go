// Copyright 2021 The KCL Authors. All rights reserved.

package validate

import (
	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

type ValidateOptions struct {
	Schema        string
	AttributeName string
	Format        string
}

func ValidateCode(data, code string, opt *ValidateOptions) (ok bool, err error) {
	if opt == nil {
		opt = &ValidateOptions{}
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.ValidateCode(&gpyrpc.ValidateCode_Args{
		Data:          data,
		Code:          code,
		Schema:        opt.Schema,
		AttributeName: opt.AttributeName,
		Format:        opt.Format,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
