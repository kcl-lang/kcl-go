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

func init() {
	os.Setenv("PYTHONHOME", "")
	os.Setenv("PYTHONPATH", "")
}

var kclvmPath = findKclvm_exePath()

var ErrKclvmNotFound = errors.New("kclvm not found")

func InitKclvmPath(kclvmRoot string) {
	if runtime.GOOS == "windows" {
		kclvmPath = filepath.Join(kclvmRoot, "kclvm.exe")
	} else {
		kclvmPath = filepath.Join(kclvmRoot, "bin", "kclvm")
	}
}

func GetKclvmPath() (string, error) {
	if kclvmPath == "" {
		return "", ErrKclvmNotFound
	}
	return kclvmPath, nil
}

func MustGetKclvmPath() string {
	s, err := GetKclvmPath()
	if err != nil {
		panic(err)
	}
	return s
}

func findKclvm_exePath() string {
	kclvmName := "kclvm"
	if runtime.GOOS == "windows" {
		kclvmName += ".exe"
	}

	exePath := getExeDir()
	if exePath == "" {
		return kclvmName
	}

	if fi, _ := os.Stat(filepath.Join(exePath, kclvmName)); fi != nil && !fi.IsDir() {
		return filepath.Join(exePath, kclvmName)
	}

	if path, err := exec.LookPath(kclvmName); err == nil {
		return path
	}

	for _, dir := range []string{
		"/usr/local/python3.7/bin",
		"/usr/local/python3.8/bin",
		"/usr/local/python3.9/bin",
		"C:/python3.7",
	} {
		if fi, _ := os.Stat(filepath.Join(dir, kclvmName)); fi != nil && !fi.IsDir() {
			return filepath.Join(dir, kclvmName)
		}
	}

	return ""
}

func getExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	if s, _ := filepath.Abs(exePath); s != "" {
		exePath = s
	}

	return filepath.Dir(exePath)
}
