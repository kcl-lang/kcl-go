// Copyright The KCL Authors. All rights reserved.

package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindPkgInfo find the pkg information(1. the pkg root 2. the pkg path of current workdir)
func FindPkgInfo(workDir string) (pkgroot, pkgpath string, err error) {
	// fix ${env}
	if idxEnvStart := strings.Index(workDir, "${"); idxEnvStart >= 0 {
		if idxEnvEnd := strings.Index(workDir, "}"); idxEnvEnd > idxEnvStart {
			envKey := workDir[idxEnvStart+2 : idxEnvEnd-1]
			workDir = strings.Replace(workDir, fmt.Sprintf("${%s}", envKey), os.Getenv(envKey), 1)
		}
	}

	var wd = workDir
	if wd == "" {
		if x, _ := os.Getwd(); x != "" {
			wd = x
		}
	}
	if abs, _ := filepath.Abs(wd); abs != "" {
		wd = abs
	}
	if wd == "" {
		return "", "", fmt.Errorf("not found pkg root")
	}

	// try load ${pwd}/.../kcl.mod
	pkgroot = wd
	for pkgroot != "" {
		kModPath := filepath.Join(pkgroot, "kcl.mod")
		if fi, _ := os.Stat(kModPath); fi != nil {
			pkgpath, err = filepath.Rel(pkgroot, wd)
			pkgroot = filepath.ToSlash(pkgroot)
			pkgpath = filepath.ToSlash(pkgpath)
			return
		}
		pkgroot = filepath.Dir(pkgroot)
		if pkgroot == "" || pkgroot == "/" || pkgroot == filepath.Dir(pkgroot) {
			break
		}
	}

	// failed
	return "", "", fmt.Errorf("pkgroot: not found")
}
