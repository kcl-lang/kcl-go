// Copyright 2021 The KCL Authors. All rights reserved.

package override

import (
	"testing"
)

func TestOverrideFile(t *testing.T) {
	t.Skip("unsupport cgo")
	_, err := OverrideFile("./testdata/test.k", []string{":config.image=\"image/image:v1\""}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}
