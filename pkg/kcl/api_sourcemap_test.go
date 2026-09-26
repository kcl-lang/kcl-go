// Copyright The KCL Authors. All rights reserved.

package kcl

import (
	"testing"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// TestExecResultToKCLResultSourcemap verifies that the source map returned
// by the runtime is surfaced through KCLResultList.GetSourcemap without
// needing a live libkcl runtime.
func TestExecResultToKCLResultSourcemap(t *testing.T) {
	sourcemap := "{\"version\":3,\"file\":\"out.map\"}"
	resp := &gpyrpc.ExecProgramResult{
		JsonResult: `{"a": 1}`,
		YamlResult: "a: 1",
		Sourcemap:  &sourcemap,
	}
	result, err := ExecResultToKCLResult(nil, resp, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := result.GetSourcemap(); got != sourcemap {
		t.Fatalf("expected sourcemap %q, got %q", sourcemap, got)
	}

	// No source map requested: the result must carry an empty source map.
	resp.Sourcemap = nil
	result, err = ExecResultToKCLResult(nil, resp, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := result.GetSourcemap(); got != "" {
		t.Fatalf("expected empty sourcemap, got %q", got)
	}
}
