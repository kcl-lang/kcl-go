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
