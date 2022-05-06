// Copyright 2021 The KCL Authors. All rights reserved.

package lint

import (
	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func LintPath(path string) (results []string, err error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.LintPath(&gpyrpc.LintPath_Args{
		Path: path,
	})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}
