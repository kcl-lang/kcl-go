// Copyright The KCL Authors. All rights reserved.

package lint

import (
	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

func LintPath(paths []string) (results []string, err error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.LintPath(&gpyrpc.LintPath_Args{
		Paths: paths,
	})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}
