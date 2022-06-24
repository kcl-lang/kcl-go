// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
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

func TestListDepFiles_restful(t *testing.T) {
	var args = gpyrpc.ListDepFiles_Args{
		WorkDir:       "../../../testdata/app0",
		UseAbsPath:    false,
		IncludeAll:    false,
		UseFastParser: true,
	}
	var result gpyrpc.ListDepFiles_Result

	var err = service.CallRestMethod(
		"http://"+tRestServerAddr, "KclvmService.ListDepFiles",
		&args, &result,
	)
	if err != nil {
		t.Fatal(err)
	}

	files := result.Files

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

func TestListDepFiles_restfulFailed(t *testing.T) {
	var args = gpyrpc.ListDepFiles_Args{
		WorkDir:       "../../../testdata/app0-failed",
		UseAbsPath:    false,
		IncludeAll:    false,
		UseFastParser: true,
	}
	var result gpyrpc.ListDepFiles_Result

	var err = service.CallRestMethod(
		"http://"+tRestServerAddr, "KclvmService.ListDepFiles",
		&args, &result,
	)
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	expectErrMsg := "package app0-failed/sub_not_found: no kcl file"
	if !strings.Contains(err.Error(), expectErrMsg) {
		t.Fatalf("expect %q, got %q", expectErrMsg, err)
	}
}
