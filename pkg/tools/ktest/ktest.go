// Copyright 2021 The KCL Authors. All rights reserved.

// Package ktest defines helper functions for kcl-test.
package ktest

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"kusionstack.io/kclvm-go/pkg/logger"
	"kusionstack.io/kclvm-go/pkg/utils"
)

var klog = logger.GetLogger()

// TODO: add test

func RunTest(packages string, opt Options) error {
	if opt.Debug {
		klog.SetLevel("DEBUG")
	}

	klog.Debugf("packages=%v, opt=%+v\n", packages, opt)

	pkglist := getPkgList(packages, opt)
	if len(pkglist) == 0 {
		if !opt.QuietMode {
			fmt.Printf("kcl-go: warning: \"%s\" matched no packages\n", packages)
			fmt.Println("no packages to test")
		}
		return nil
	}

	defer func() {
		for _, pkg := range pkglist {
			os.RemoveAll(filepath.Join(pkg, "__pycache__"))
		}
	}()

	var lastErr error
	for _, pkg := range pkglist {
		test, err := loadKTestSuit(pkg, opt)
		if err != nil {
			lastErr = err
		}
		if errx := test.RunTest(); errx != nil {
			lastErr = errx
		}
	}
	if lastErr != nil {
		return lastErr
	}

	return nil
}

func getPkgList(pkgpath string, opt Options) []string {
	if pkgpath == "" {
		pkgpath, _ = os.Getwd()
		return []string{pkgpath}
	}

	var includeSubPkg bool
	if strings.HasSuffix(pkgpath, "/...") {
		includeSubPkg = true
		pkgpath = pkgpath[:len(pkgpath)-len("/...")]
	}
	if pkgpath != "." && strings.HasSuffix(pkgpath, ".") {
		return nil
	}
	if pkgpath == "" {
		return nil
	}

	// ktest ./pkg/...
	switch {
	case strings.HasPrefix(pkgpath, "."):
		wd, _ := os.Getwd()
		pkgpath = filepath.Join(wd, pkgpath)
	case filepath.IsAbs(pkgpath):
		// skip
	default:
		if !strings.ContainsAny(pkgpath, `\/`) {
			pkgpath = strings.ReplaceAll(pkgpath, ".", "/")
		}

		wd, _ := os.Getwd()
		pkgroot, _ := utils.FindPkgRoot(wd)
		if pkgroot != "" {
			pkgpath = filepath.Join(pkgroot, pkgpath)
		} else {
			pkgpath = filepath.Join(wd, pkgpath)
		}
	}

	if !includeSubPkg {
		return []string{pkgpath}
	}

	var (
		dirList []string
		dirMap  = make(map[string]bool)
	)
	filepath.Walk(pkgpath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".k") && !strings.HasPrefix(filepath.Base(path), "_") {
			isIgnoredDir := false
			relPath, _ := filepath.Rel(pkgpath, path)
			for _, s := range filepath.SplitList(relPath) {
				if strings.HasPrefix(s, "_") {
					isIgnoredDir = true
					break
				}
			}

			if !isIgnoredDir {
				if dir := filepath.Dir(path); !dirMap[dir] {
					dirList = append(dirList, dir)
					dirMap[dir] = true
				}
			}
		}
		return nil
	})

	return dirList
}
