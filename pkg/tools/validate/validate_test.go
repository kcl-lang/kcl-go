// Copyright The KCL Authors. All rights reserved.

package validate

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	ok, err := Validate("./test_data/data.json", "./test_data/schema.k", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("expect: %q, got False", "True")
	}
}

func TestValidateFailed(t *testing.T) {
	ok, err := Validate("./test_data/data-failed.json", "./test_data/schema.k", nil)
	if ok == false && err != nil && strings.Contains(err.Error(), "expected [int], got [int(1) | int(2) | int(3) | str()]") {
		// Test Pass
	} else {
		t.Fatalf("expect: error, got (%v, %v)", ok, err)
	}
}

func TestValidateCode(t *testing.T) {
	data := `{"key": "value"}`
	code := `
schema Person:
    key: str

    check:
        "value" in key  # 'key' is required and 'key' must contain "value"
`

	ok, err := ValidateCode(data, code, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("expect: %q, got False", "True")
	}
}

func TestValidateCodeFail(t *testing.T) {
	data := `{"k": "value"}`
	code := `
schema Person:
    key: str

    check:
        "value" in key  # 'key' is required and 'key' must contain "value"
`

	_, err := ValidateCode(data, code, nil)
	if err == nil {
		t.Fatalf("expect validation error")
	} else if !strings.Contains(err.Error(), "error") {
		t.Fatalf("expect validation error, got %s", err.Error())
	}
}

// TestValidateWithExternalPackages is a regression test for
// https://github.com/kcl-lang/kcl/issues/1877. It ensures that when the
// caller supplies an `ExternalPackages` mapping, the validator can
// resolve imports to schemas declared in those external packages and
// validate data against the resulting type.
func TestValidateWithExternalPackages(t *testing.T) {
	opts := &ValidateOptions{
		ExternalPackages: []ExternalPackage{
			{
				PkgName: "ext",
				PkgPath: "./test_data/external/ext",
			},
		},
	}
	ok, err := Validate(
		"./test_data/external/with_ext/data.json",
		"./test_data/external/with_ext/schema.k",
		opts,
	)
	if err != nil {
		t.Fatalf("expected validation success, got error: %v", err)
	}
	if !ok {
		t.Fatalf("expected validation success, got failure")
	}
}

// TestValidateWithoutExternalPackagesReportsMissingImport ensures the
// historical behaviour is preserved when the caller does not supply an
// `ExternalPackages` mapping: an unresolved external import surfaces as
// a `Cannot find the module` error rather than silently passing.
func TestValidateWithoutExternalPackagesReportsMissingImport(t *testing.T) {
	_, err := Validate(
		"./test_data/external/with_ext/data.json",
		"./test_data/external/with_ext/schema.k",
		nil,
	)
	if err == nil {
		t.Fatalf("expected validation error for unresolved import, got nil")
	}
	if !strings.Contains(err.Error(), "ext") {
		t.Fatalf("expected error to mention missing 'ext' package, got: %v", err)
	}
}

// TestExternalPackagesWithEmptyEntriesIgnored verifies that the
// conversion helper drops empty entries rather than forwarding them to
// the gRPC service, which would otherwise produce confusing errors.
func TestExternalPackagesWithEmptyEntriesIgnored(t *testing.T) {
	opts := &ValidateOptions{
		ExternalPackages: []ExternalPackage{
			{PkgName: "", PkgPath: "./test_data/external/ext"},
			{PkgName: "ext", PkgPath: ""},
		},
	}
	if got := opts.toExternalPkgs(); got != nil {
		t.Fatalf("expected nil gRPC list when only empty entries are provided, got %v", got)
	}
}
