// Copyright 2022 The KCL Authors. All rights reserved.

package list

func ListDepFiles(workDir string, opt *Option) (files []string, err error) {
	if opt == nil {
		opt = &Option{}
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
