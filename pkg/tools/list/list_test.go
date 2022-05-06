// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"reflect"
	"sort"
	"testing"
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
