// Copyright 2021 The KCL Authors. All rights reserved.

package validate

import (
	"strings"
	"testing"
)

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
