// Copyright 2022 The KCL Authors. All rights reserved.

package list

import (
	"context"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

// ListDepFiles return the depend files from the given path. It will scan and parse the kusion applications within the workdir,
// then list depend files of the applications.
func ListDepFiles(workDir string, opt *Option) (files []string, err error) {
	if opt == nil {
		opt = &Option{}
	}

	if opt.RestfulUrl != "" {
		client := service.NewRestClient(opt.RestfulUrl)
		reply, err := client.ListDepFiles(context.Background(), &gpyrpc.ListDepFiles_Args{
			WorkDir: workDir,
		})
		if err != nil {
			return nil, err
		}
		return reply.Files, nil
	}

	pkgroot, pkgpath, err := FindPkgInfo(workDir)
	if err != nil {
		return nil, err
	}

	depParser := NewSingleAppDepParser(pkgroot, *opt)

	for _, s := range depParser.GetAppFiles(pkgpath, opt.FlagAll) {
		if opt.UseAbsPath {
			files = append(files, pkgroot+"/"+s)
		} else {
			files = append(files, s)
		}
	}

	return files, nil
}

// ListUpStreamFiles return a list of upstream depend files from the given path list.
func ListUpStreamFiles(workDir string, opt *DepOption) (deps []string, err error) {
	if opt == nil || opt.Files == nil {
		return nil, nil
	}
	pkgroot, _, err := FindPkgInfo(workDir)
	if err != nil {
		return nil, err
	}
	depParser, err := NewImportDepParser(pkgroot, *opt)
	if err != nil {
		return nil, err
	}
	return depParser.ListUpstreamFiles(), nil
}

// ListDownStreamFiles return a list of downstream depend files from the given changed path list.
func ListDownStreamFiles(workDir string, opt *DepOption) ([]string, error) {
	if opt == nil || opt.Files == nil || opt.ChangedPaths == nil {
		return nil, nil
	}
	pkgroot, _, err := FindPkgInfo(workDir)
	if err != nil {
		return nil, err
	}
	depParser, err := NewImportDepParser(pkgroot, *opt)
	if err != nil {
		return nil, err
	}
	return depParser.ListDownStreamFiles(), nil
}
