// Copyright 2021 The KCL Authors. All rights reserved.

package lint

import (
	"strings"
	"testing"
)

func TestLintPath(t *testing.T) {
	t.Skip("unsupport cgo")
	expect := `hello imported but unused.`

	ss, err := LintPath("./testdata/a.k")
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) != 1 {
		t.Fatalf("expect: %q, got empty", expect)
	}
	if !strings.HasSuffix(ss[0], expect) {
		t.Fatalf("expect: %q, got %q", expect, ss[0])
	}
}
