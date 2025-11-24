// Copyright The KCL Authors. All rights reserved.

package override

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const (
	DeleteAction         = "Delete"
	CreateOrUpdateAction = "CreateOrUpdate"
)

func OverrideFile(file string, specs, importPaths []string) (result bool, err error) {
	svc := kcl.Service()
	resp, err := svc.OverrideFile(&gpyrpc.OverrideFileArgs{
		File:        file,
		Specs:       specs,
		ImportPaths: importPaths,
	})
	if err != nil {
		return false, err
	}
	return resp.Result, nil
}
