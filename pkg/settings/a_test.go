// Copyright 2021 The KCL Authors. All rights reserved.

package settings

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadFile(t *testing.T) {
	f, err := LoadFile("../../testdata/app0/kcl.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	f.To_ExecProgram_Args()
	_ = f
}

func _TestLoadFile_testFormat(t *testing.T) {
	// kcl_options: -D key1=val -D key2=2 -D key3=4.4 -D key4=[1,2,3] -D key5={'key':'value'} -S app.value
	f, err := LoadFile("../../../test/grammar/option/complex_type_option/settings.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	f.To_ExecProgram_Args()
	_ = f
}

func TestLoadFile_xtype(t *testing.T) {
	_, err := LoadFile("settings.yaml", "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadFile("settings.yaml", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	_, err = LoadFile("settings.yaml", []byte(""))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadFile_to_ExecProgram_Args(t *testing.T) {
	const s = `
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

func tAssert(t *testing.T, ok bool, a ...interface{}) {
	if !ok {
		t.Helper()
		t.Fatal(a...)
	}
}

func tAssertEQ(t *testing.T, x, y interface{}) {
	if !reflect.DeepEqual(x, y) {
		t.Helper()
		t.Fatalf("not equal:\n  x = %v\n  y = %v\n", x, y)
	}
}
