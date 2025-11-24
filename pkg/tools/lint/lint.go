// Copyright The KCL Authors. All rights reserved.

package lint

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

func LintPath(paths []string) (results []string, err error) {
	svc := kcl.Service()
	resp, err := svc.LintPath(&gpyrpc.LintPathArgs{
		Paths: paths,
	})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}
