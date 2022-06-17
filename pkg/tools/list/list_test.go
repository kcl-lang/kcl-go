// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestListDepFiles(t *testing.T) {
	testListDepFiles(t, nil)
}

func TestListDepFiles_restful(t *testing.T) {
	testListDepFiles(t, &Option{
		RestfulUrl: "http://" + tRestfulAddr,
	})
}

func testListDepFiles(t *testing.T, opt *Option) {
	abspath, err := filepath.Abs("../../../testdata/app0")
	if err != nil {
		t.Fatal(err)
	}

	files, err := ListDepFiles(abspath, opt)
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
