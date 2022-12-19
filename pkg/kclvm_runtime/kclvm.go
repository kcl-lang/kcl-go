// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	_ "embed"
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	kclPlugin "kusionstack.io/kcl-plugin"
	kclvmArtifact "kusionstack.io/kclvm-artifact"
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

	g_Python3Path = findPython3Path()
	g_KclvmRoot = findKclvmRoot()
	kclvmPluginPath := filepath.Join(g_KclvmRoot, "plugins")
	if runtime.GOOS == "windows" {
		kclvmPluginPath = filepath.Join(g_KclvmRoot, "bin", "plugins")
	}

	_, err = os.Stat(kclvmPluginPath)

	if os.IsNotExist(err) {
		err = kclPlugin.InstallPlugins(kclvmPluginPath)
		if err != nil {
			panic(fmt.Errorf("install kclvm plugins failed: %s", err.Error()))
		}
	}
}

var (
	g_Python3Path string
	g_KclvmRoot   string
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
	if g_Python3Path == "" {
		return "", ErrPython3NotFound
	}
	if g_KclvmRoot == "" {
		return "", ErrKclvmRootNotFound
	}
	return g_Python3Path, nil
}

// MustGetKclvmPath return kclvm/python3 executable path, panic if not found.
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
	kclvm_cli_exe := "kclvm"
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
