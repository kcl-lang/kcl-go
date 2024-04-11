// Copyright The KCL Authors. All rights reserved.

package list

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDepFiles(t *testing.T) {
	files, err := ListDepFiles("../../../testdata/app0", nil)
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"main.k",
		"app0/before/base.k",
		"app0/main.k",
		"app0/sub/sub.k",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}

func TestListDepPackages(t *testing.T) {
	files, err := ListDepPackages("./testdata/module_with_external/", &Option{
		ExcludeExternalPackage: true,
		ExcludeBuiltin:         true,
		IgnoreImportError:      true,
	})
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"pkg1",
		"pkg2",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}

func TestListDepFiles_failed(t *testing.T) {
	_, err := ListDepFiles("../../../testdata/app0-failed", nil)
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	expectErrMsg := "package app0-failed/sub_not_found: no kcl file"
	if !strings.Contains(err.Error(), expectErrMsg) {
		t.Fatalf("expect %q, got %q", expectErrMsg, err)
	}
}

func TestNewDepParser(t *testing.T) {

	dp := NewDepParser("../../../testdata/import_builtin", Option{
		KclYaml:     "kcl.yaml",
		ProjectYaml: "project.yaml",
	})

	assert.Equal(t, dp.err, nil)
	assert.Equal(t, len(dp.pkgFilesMap), 2)
	assert.Equal(t, dp.pkgFilesMap["entry"], []string{"entry/main.k"})
	assert.Equal(t, dp.pkgFilesMap["sub"], []string{"sub/main.k"})
}
