// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	"context"
	_ "embed"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gofrs/flock"
	artifact "kcl-lang.io/kcl-artifact-go"
	"kusionstack.io/kclvm-go/pkg/logger"
	"kusionstack.io/kclvm-go/pkg/path"
)

func init() {
	// Get the install lib path.
	path := path.LibPath()
	// Acquire a file lock for process synchronization
	lockPath := filepath.Join(path, "init.lock")
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		logger.GetLogger().Warningf("install kclvm failed: %s", err.Error())
	}
	// Install lib
	err = artifact.InstallKclvm(path)
	if err != nil {
		logger.GetLogger().Warningf("install kclvm failed: %s", err.Error())
	}
	artifact.CleanInstall()

	g_KclvmRoot = findKclvmRoot()
}

var (
	g_KclvmRoot          string
	ErrKclvmRootNotFound = errors.New("kclvm root not found, please ensure kcl is in your PATH")
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
