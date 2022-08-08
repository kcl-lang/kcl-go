// Copyright 2021 The KCL Authors. All rights reserved.

package ktest

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("PYTHONDONTWRITEBYTECODE", "1")
}

func TestPlugin(t *testing.T) {
	t.Skip("unsupport cgo")
	err := RunTest("./testdata/kcl_plugin/...", Options{
		QuietMode: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlugin_hello_failed(t *testing.T) {
	err := RunTest("./testdata/kcl_plugin_failed/hello", Options{
		QuietMode: true,
	})
	if err == nil {
		t.Fatal("expect error, got nil")
	}
}

func TestMainK(t *testing.T) {
	err := RunTest("./testdata/app/...", Options{
		QuietMode: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}
