// Copyright 2021 The KCL Authors. All rights reserved.

package validate

import (
	"testing"
)

func TestValidateCode(t *testing.T) {
	t.Skip("unsupport cgo")
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
