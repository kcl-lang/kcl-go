// Copyright The KCL Authors. All rights reserved.

package settings

import (
	"strings"
	"testing"
)

func TestLoadFile(t *testing.T) {
	f, err := LoadFile("../../testdata/app0/kcl.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	f.To_ExecProgramArgs()
	_ = f
}

func _TestLoadFile_testFormat(t *testing.T) {
	// kcl_options: -D key1=val -D key2=2 -D key3=4.4 -D key4=[1,2,3] -D key5={'key':'value'} -S app.value
	f, err := LoadFile("../../../test/grammar/option/complex_type_option/settings.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	f.To_ExecProgramArgs()
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

func tAssert(t *testing.T, ok bool, a ...interface{}) {
	if !ok {
		t.Helper()
		t.Fatal(a...)
	}
}
