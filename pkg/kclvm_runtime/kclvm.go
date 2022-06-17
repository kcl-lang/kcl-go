// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	_ "embed"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	g_Python3Path = findPython3Path()
	g_KclvmRoot   = findKclvmRoot()
)

var (
	ErrPython3NotFound   = errors.New("python3 not found")
	ErrKclvmRootNotFound = errors.New("kclvm root not found")
)

func InitKclvmRoot(kclvmRoot string) {
	g_KclvmRoot = kclvmRoot
	if runtime.GOOS == "windows" {
		s := filepath.Join(g_KclvmRoot, "kclvm.exe")
		if fi, _ := os.Lstat(s); fi != nil && !fi.IsDir() {
			g_Python3Path = s
		}
	} else {
		s := filepath.Join(g_KclvmRoot, "bin", "kclvm")
		if fi, _ := os.Lstat(s); fi != nil && !fi.IsDir() {
			g_Python3Path = s
		}
	}
}

func GetKclvmRoot() string {
	return g_KclvmRoot
}

func GetKclvmPath() (string, error) {
	if g_Python3Path == "" {
		return "", ErrPython3NotFound
	}
	if g_KclvmRoot == "" {
		return "", ErrKclvmRootNotFound
	}
	return g_Python3Path, nil
}

func MustGetKclvmPath() string {
	s, err := GetKclvmPath()
	if err != nil {
		panic(err)
	}
	return s
}

func findPython3Path() string {
	for _, s := range []string{"kclvm", "python3"} {
		exeName := s
		if runtime.GOOS == "windows" {
			exeName += ".exe"
		}
		if path, err := exec.LookPath(exeName); err == nil {
			return path
		}
	}
	return ""
}

func findKclvmRoot() string {
	kclvm_cli_exe := "kclvm_cli"
	if runtime.GOOS == "windows" {
		kclvm_cli_exe += ".exe"
	}
	if path, err := exec.LookPath(kclvm_cli_exe); err == nil {
		if runtime.GOOS == "windows" {
			return filepath.Dir(path)
		} else {
			return filepath.Dir(filepath.Dir(path))
		}
	}
	return ""
}
