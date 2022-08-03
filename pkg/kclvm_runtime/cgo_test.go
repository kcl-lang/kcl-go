package kclvm_runtime

import (
	"testing"

	"kusionstack.io/kclvm-go/pkg/settings"
)

func TestLoadFile(t *testing.T) {
	f, err := settings.LoadFile("../../testdata/app0/kcl.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	args := f.To_ExecProgram_Args()
	CGO_kclvm_cli_run(args, 0)
}
