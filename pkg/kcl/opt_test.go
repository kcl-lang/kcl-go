// Copyright 2021 The KCL Authors. All rights reserved.

package kcl

import (
	"os"
	"testing"
)

func TestWithSettings(t *testing.T) {
	opt := WithSettings("../../testdata/app0/kcl.yaml")
	tAssert(t, opt.Err == nil, opt.JSONString())
}

func TestWithKFilenames(t *testing.T) {
	WithKFilenames("hello.k")
}

func TestWithOptions(t *testing.T) {
	WithOptions("key1=value1", "key2=value2")
}

func TestWithOverrides(t *testing.T) {
	WithOverrides("__main__:name.field1=value1", "__main__:name.field2=value2")
}

func TestWithPrintOverridesAST(t *testing.T) {
	WithPrintOverridesAST(true)
}

func TestWithWorkDir(t *testing.T) {
	wd, _ := os.Getwd()
	WithWorkDir(wd)
}

func TestWithDisableNone(t *testing.T) {
	WithDisableNone(true)
}
