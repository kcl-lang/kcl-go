// Copyright The KCL Authors. All rights reserved.

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

func TestMergeWithError(t *testing.T) {
	// Test that errors from options are correctly propagated during merge
	// This is a regression test for issue #515 where path_selector parsing errors
	// were being ignored when merging options
	tmpFile, err := os.CreateTemp("", "kcl-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write an invalid kcl.yaml with path_selector as a scalar instead of array
	_, err = tmpFile.WriteString(`kcl_cli_configs:
  path_selector: data
kcl_options:
  - key: test_option
    value: test_value
`)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Load the invalid settings
	errOpt := WithSettings(tmpFile.Name())
	if errOpt.Err == nil {
		t.Fatal("Expected error when loading invalid settings file, got nil")
	}

	// Merge with another option
	baseOpt := NewOption()
	merged := baseOpt.Merge(errOpt)

	// The error should be propagated
	if merged.Err == nil {
		t.Fatal("Expected error to be propagated after merge, got nil")
	}

	// The error should mention the YAML unmarshaling issue
	if errStr := merged.Err.Error(); !contains(errStr, "cannot unmarshal") {
		t.Fatalf("Expected error about unmarshaling, got: %v", errStr)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
