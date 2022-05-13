// Copyright 2021 The KCL Authors. All rights reserved.

package override

import (
	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func OverrideFile(file string, specs, importPaths []string) (result bool, err error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.OverrideFile(&gpyrpc.OverrideFile_Args{
		File:        file,
		Specs:       specs,
		ImportPaths: importPaths,
	})
	if err != nil {
		return false, err
	}
	return resp.Result, nil
}
