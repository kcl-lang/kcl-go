// Copyright The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package testpluginagent_test

import (
	"testing"

	"kcl-lang.io/kcl-go/pkg/internal/testpluginagent"
	"kcl-lang.io/kcl-go/pkg/kcl"
)

const code = `
import kcl_plugin.testplugin

a = option("a")
b = option("b")
sum = testplugin.add(a, b)
`

// TestRunWithPluginAgent runs KCL code that calls a host function through a
// user-provided plugin agent. The test lives in its own package because the
// native service client is a process-wide singleton initialized with the
// plugin agent of the first client created in the process.
func TestRunWithPluginAgent(t *testing.T) {
	result, err := kcl.Run("main.k",
		kcl.WithCode(code),
		kcl.WithOptions("a=1", "b=2"),
		kcl.WithPluginAgent(testpluginagent.Ptr()),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := result.First().Get("sum"); got != int(3) {
		t.Fatalf("sum = %v (%T), want 3", got, got)
	}
	if calls := testpluginagent.Calls(); calls != 1 {
		t.Fatalf("agent calls = %d, want 1", calls)
	}
	if got := testpluginagent.LastMethod(); got != "kcl_plugin.testplugin.add" {
		t.Fatalf("last method = %q, want %q", got, "kcl_plugin.testplugin.add")
	}
}
