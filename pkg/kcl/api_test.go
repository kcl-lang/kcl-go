// Copyright 2021 The KCL Authors. All rights reserved.

package kcl

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"kcl-lang.io/kcl-go/pkg/tools/list"
)

var _ = fmt.Sprint

const case_path = "../../testdata/main.k"

func TestRun_kcl_yaml(t *testing.T) {
	const s = "../../testdata/app0/kcl.yaml"
	_, err := RunFiles([]string{s})
	tAssert(t, err == nil, err)
}

func TestRun(t *testing.T) {
	const k_code = `
name = "kcl"
i = 123
f = 1.5
`

	result, err := RunFiles([]string{case_path}, WithCode(k_code))
	tAssert(t, err == nil, err)
	tAssert(t, result != nil)

	opts := WithCode(k_code)
	opts.Merge(WithKFilenames(case_path))
	result, err = RunWithOpts(opts)
	tAssert(t, result != nil)
	tAssert(t, err == nil, err)

	result, err = Run(case_path, WithCode(k_code))
	tAssert(t, err == nil, err)
	tAssert(t, result.Len() > 0)
	tAssert(t, result.First().Get("name") == "kcl")

	var s string
	var i int
	var f float64

	_, err = result.Get(0).GetValue("name", &s)
	tAssert(t, err == nil, err)
	tAssert(t, s == "kcl", s)

	_, err = result.Get(0).GetValue("i", &i)
	tAssert(t, err == nil, err)
	tAssert(t, i == 123, i)

	_, err = result.Get(0).GetValue("f", &f)
	tAssert(t, err == nil, err)
	tAssert(t, f == 1.5, f)

	_, err = result.Tail().GetValue("name", &s)
	tAssert(t, err == nil, err)
	tAssert(t, s == "kcl", s)

	result.First().YAMLString()
	result.Tail().JSONString()
}

// go test -run=TestRun_failed
func TestRun_failed(t *testing.T) {
	_, err := Run(case_path, WithCode(`x = {`))
	tAssert(t, err != nil, err)
}

func TestGetSchemaType(t *testing.T) {
	const k_code = `a=1`

	result, err := GetSchemaType("main.k", k_code, "")
	tAssert(t, err == nil)
	_ = result
}

func TestListUpstreamFiles(t *testing.T) {
	deps, err := list.ListUpStreamFiles("./testdata/complicate", &list.DepOptions{Files: []string{"appops/projectA/base/base.k", "appops/projectA/dev/main.k", "base/render/server/server_render.k"}})
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"base/frontend/server/server.k",
		"base/frontend/container/container.k",
		"base/frontend/container/container_port.k",
		"base/frontend/server",
		"base/frontend/container",
	}

	sort.Strings(deps)
	sort.Strings(expect)

	if !reflect.DeepEqual(deps, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, deps)
	}
}
