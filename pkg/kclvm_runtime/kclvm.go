// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	_ "embed"
	"errors"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	kclvmArtifact "kusionstack.io/kclvm-artifact-go"
	"kusionstack.io/kclvm-go/pkg/logger"
)

func init() {

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	err := kclvmArtifact.InstallKclvm(gopath)
	if err != nil {
		logger.GetLogger().Warningf("install kclvm failed: %s", err.Error())
	}
	kclvmArtifact.CleanInstall()

	g_KclvmRoot = findKclvmRoot()
}

var (
	g_KclvmRoot string
)

var (
	ErrKclvmRootNotFound = errors.New("kclvm root not found")
)

func InitKclvmRoot(kclvmRoot string) {
	g_KclvmRoot = kclvmRoot
}

// GetKclvmRoot return kclvm root directory, return error if kclvm not found.
func GetKclvmRoot() (string, error) {
	if g_KclvmRoot == "" {
		return "", ErrKclvmRootNotFound
	}
	return g_KclvmRoot, nil
}

// GetKclvmRoot return kclvm root directory, panic if kclvm not found.
func MustGetKclvmRoot() string {
	s, err := GetKclvmRoot()
	if err != nil {
		panic(err)
	}
	return s
}

// GetKclvmPath return kclvm/python3 executable path, return error if not found.
func GetKclvmPath() (string, error) {
	if g_KclvmRoot == "" {
		return "", ErrKclvmRootNotFound
	}
	return g_KclvmRoot, nil
}

// MustGetKclvmPath return kclvm/python3 executable path, panic if not found.
func MustGetKclvmPath() string {
	s, err := GetKclvmPath()
	if err != nil {
		panic(err)
	}
	return s
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
