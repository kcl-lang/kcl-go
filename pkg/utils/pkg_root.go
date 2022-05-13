// Copyright 2020 The KCL Authors. All rights reserved.

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GoodPkgPath(path string) (string, error) {
	if x, _ := filepath.Abs(path); x != "" {
		path = x
	}
	if strings.HasSuffix(path, ".k") {
		path = filepath.Dir(path)
	}
	pkgRoot, err := FindPkgRoot(path)
	if err != nil {
		return "", err
	}

	pkgPath, err := filepath.Rel(pkgRoot, path)
	if err != nil {
		return "", err
	}
	return pkgPath, nil
}

func FindPkgRoot(workDir string) (string, error) {
	wd := workDir
	if wd == "" {
		if x, _ := os.Getwd(); x != "" {
			wd = x
		}
	}
	if abs, _ := filepath.Abs(wd); abs != "" {
		wd = abs
	}

	if wd == "" {
		return "", fmt.Errorf("not found pkgroot")
	}

	// try load ${pwd}/.../kcl.mod
	pkgroot := wd
	for pkgroot != "" {
		kModPath := filepath.Join(pkgroot, "kcl.mod")
		if fi, _ := os.Stat(kModPath); fi != nil {
			return pkgroot, nil
		}
		parentDir := filepath.Dir(pkgroot)
		if parentDir == pkgroot {
			break
		}
		pkgroot = parentDir
	}

	// failed
	return "", fmt.Errorf("not found pkgroot")
}
