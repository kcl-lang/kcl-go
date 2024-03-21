// Copyright The KCL Authors. All rights reserved.

package kcl

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kcl-go/pkg/tools/list"
)

var _ = fmt.Sprint

const case_path = "../../testdata/main.k"

func TestKCLResultMap(t *testing.T) {
	var data = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	result := NewResult(data)
	m, _ := result.ToMap()
	tAssert(t, m["key1"] == "value1", m)
	tAssert(t, m["key2"] == "value2", m)
}

func TestKCLResultInt(t *testing.T) {
	result := NewResult(1)
	m, _ := result.ToInt()
	tAssert(t, *m == 1)
}

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

func TestGetFullSchemaType(t *testing.T) {
	testPath := filepath.Join(".", "testdata", "get_schema_ty")
	tys, err := GetFullSchemaType(
		[]string{filepath.Join(testPath, "aaa")},
		"",
		WithExternalPkgs(fmt.Sprintf("bbb=%s", filepath.Join(testPath, "bbb"))),
	)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(tys), 1)
	assert.Equal(t, tys[0].Filename, filepath.Join("testdata", "get_schema_ty", "bbb", "main.k"))
	assert.Equal(t, tys[0].SchemaName, "B")
}
