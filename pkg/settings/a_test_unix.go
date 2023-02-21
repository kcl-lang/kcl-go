//go:build linux || darwin
// +build linux darwin

package settings

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadFile_to_ExecProgram_Args(t *testing.T) {
	s := `
kcl_cli_configs:
  file:
    - /abs_file.k
    - sub_main.k
    - ${KCL_MOD}/file2.k
    - ../../base/base.k
  disable_none: false
`
	f, err := LoadFile("./sub/settings.yaml", []byte(s))
	if err != nil {
		t.Fatal(err)
	}

	pwd, _ := os.Getwd()
	x := f.To_ExecProgram_Args()

	tAssertEQ(t, len(x.KFilenameList), 4)
	tAssertEQ(t, x.KFilenameList[0], "/abs_file.k")
	tAssertEQ(t, x.KFilenameList[1], filepath.Join(pwd, "sub", "sub_main.k"))
	tAssertEQ(t, x.KFilenameList[2], filepath.Join(pwd, "file2.k"))
	tAssertEQ(t, x.KFilenameList[3], filepath.Join(pwd, "..", "base", "base.k"))
}

func tAssertEQ(t *testing.T, x, y interface{}) {
	if !reflect.DeepEqual(x, y) {
		t.Helper()
		t.Fatalf("not equal:\n  x = %v\n  y = %v\n", x, y)
	}
}
