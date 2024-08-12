// Copyright The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package kcl_test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"

	kcl "kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const tEnvNumCpu = "KCL_GO_API_TEST_NUM_CPU"

func TestMain(m *testing.M) {
	flag.Parse()

	if s := os.Getenv(tEnvNumCpu); s != "" {
		if x, err := strconv.Atoi(s); err == nil {
			println(tEnvNumCpu, "=", s)
			kcl.InitKclvmRuntime(x)
		}
	}

	os.Exit(m.Run())
}

func TestStreamResult(t *testing.T) {
	file, err := filepath.Abs("./testdata/stream/one_stream.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err := kcl.Run(file, kcl.WithSortKeys(true))
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "a: 1", result.GetRawYamlResult())
	file, err = filepath.Abs("./testdata/stream/two_stream.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err = kcl.Run(file, kcl.WithSortKeys(true))
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "a: 1\n---\nb: 2", result.GetRawYamlResult())
}

func TestWithTypePath(t *testing.T) {
	const code = `
schema App:
	image: str = "default"
	name: str = "app"

a1 = App {
	name = "a1-app"
}
a2 = App {
	image = "a2-image"
	name = "a2-app"
}
`
	const testdata_main_k = "test.k"
	result, err := kcl.Run(testdata_main_k,
		kcl.WithCode(code),
		kcl.WithFullTypePath(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := "App", result.First().Get("a1._type"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "App", result.First().Get("a2._type"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "default", result.First().Get("a1.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "a2-image", result.First().Get("a2.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}

	result, err = kcl.Run(testdata_main_k,
		kcl.WithCode(code),
		kcl.WithIncludeSchemaTypePath(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := "App", result.First().Get("a1._type"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "App", result.First().Get("a2._type"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
}

func TestWithOverrides(t *testing.T) {
	const code = `
schema App:
	image: str = "default"
	name: str = "app"

a1 = App {
	name = "a1-app"
}
a2 = App {
	image = "a2-image"
	name = "a2-app"
}
`
	const testdata_main_k = "testdata/main_.k"
	kfile, err := os.Create(testdata_main_k)
	if err != nil {
		t.Fatal(err)
	}
	kfile.Close()

	result, err := kcl.Run(testdata_main_k,
		kcl.WithCode(code),
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := "default", result.First().Get("a1.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "a2-image", result.First().Get("a2.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}

	os.Remove(testdata_main_k)
	defer os.Remove(testdata_main_k)

	kfile, err = os.Create(testdata_main_k)
	kfile.Close()

	result, err = kcl.Run(testdata_main_k,
		kcl.WithCode(code),
		kcl.WithOverrides("a1.image=\"new-a1-image\""),
		kcl.WithOverrides("a2.image=\"new-a2-image:v123\""),
		kcl.WithOverrides("a2.name-"),
		kcl.WithPrintOverridesAST(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	if expect, got := "new-a1-image", result.First().Get("a1.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "new-a2-image:v123", result.First().Get("a2.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "app", result.First().Get("a2.name"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}

	data, err := os.ReadFile(testdata_main_k)
	if err != nil {
		t.Fatal(err)
	}

	want := strings.TrimSpace(`
schema App:
    image: str = "default"
    name: str = "app"

a1 = App {
    name = "a1-app"
    image = "new-a1-image"
}

a2 = App {image = "new-a2-image:v123"}`)

	got := strings.TrimSpace(string(data))
	got = strings.ReplaceAll(got, "\r\n", "\n")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("golden mismatch (-want +got):\n%s", diff)
	}
}

func TestWithSelectors(t *testing.T) {
	const code = `
schema App:
	image: str = "default"
	name: str = "app"

a1 = App {
	name = "a1-app"
}
a2 = App {
	image = "a2-image"
	name = "a2-app"
}
`
	const testdata_main_k = "testdata/main_selector.k"
	kfile, err := os.Create(testdata_main_k)
	if err != nil {
		t.Fatal(err)
	}
	kfile.Close()

	result, err := kcl.Run(testdata_main_k,
		kcl.WithCode(code),
		kcl.WithSelectors("a1"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := "a1-app", result.First().Get("name"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	os.Remove(testdata_main_k)
	defer os.Remove(testdata_main_k)
}

func TestWithSelectorsReturnString(t *testing.T) {
	const code = `
schema Person:
    labels: {str:str}

alice = Person {
    "labels": {"skin": "yellow"}
}
`
	const testdata_main_k = "testdata/main_selector.k"
	kfile, err := os.Create(testdata_main_k)
	if err != nil {
		t.Fatal(err)
	}
	kfile.Close()

	result, err := kcl.Run(testdata_main_k,
		kcl.WithCode(code),
		kcl.WithSelectors("alice.labels.skin"),
	)
	assert2.Equal(t, err, nil)
	resInStr, err := result.ToString()
	assert2.Equal(t, err, nil)
	assert2.Equal(t, resInStr, "yellow")

	resInBool, err := result.ToBool()
	assert2.Equal(t, err.Error(), "failed to convert result to *bool: type mismatch")
	assert2.Equal(t, resInBool, (*bool)(nil))
	resInF64, err := result.ToFloat64()
	assert2.Equal(t, err.Error(), "failed to convert result to *float64: type mismatch")
	assert2.Equal(t, resInF64, (*float64)(nil))
	resInList, err := result.ToList()
	assert2.Equal(t, err.Error(), "failed to convert result to *[]interface {}: type mismatch")
	assert2.Equal(t, resInList, []interface{}(nil))
	resInMap, err := result.ToMap()
	assert2.Equal(t, err.Error(), "failed to convert result to *map[string]interface {}: type mismatch")
	assert2.Equal(t, resInMap, map[string]interface{}(map[string]interface{}(nil)))

	os.Remove(testdata_main_k)
	defer os.Remove(testdata_main_k)
}

func TestWithKFilenames(t *testing.T) {
	kcl.WithKFilenames("/testdata/main.k")
}

func TestWithOptions(t *testing.T) {
	kcl.WithOptions("key1=value1", "key2=value2")
}

func TestWithSettings(t *testing.T) {
	kcl.WithSettings("a_settings.yml")
}

func TestWithWorkDir(t *testing.T) {
	wd, _ := os.Getwd()
	kcl.WithWorkDir(wd)
}

func TestWithDisableNone(t *testing.T) {
	kcl.WithDisableNone(true)
}

func TestFormatCode(t *testing.T) {
	result, err := kcl.FormatCode("a=1")
	if err != nil {
		t.Error(err)
	}
	assert2.Equalf(t, string(result), "a = 1\n", "format result unexpected: expect: a = 1\n, actual: %s", result)
}

func TestGetSchemaType(t *testing.T) {
	result, err := kcl.GetSchemaType("test.k", "@info(\"v\", k=\"v\")\nschema Person:\n    name: str", "")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, []*gpyrpc.KclType{
		{
			Filename:   "test.k",
			PkgPath:    "__main__",
			Type:       "schema",
			SchemaName: "Person",
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type: "str",
					Line: 1,
				},
			},
			Required: []string{"name"},
			Decorators: []*gpyrpc.Decorator{
				{
					Name:      "info",
					Arguments: []string{"\"v\""},
					Keywords: map[string]string{
						"k": "\"v\"",
					},
				},
			},
		},
	}, result)
	p, err := filepath.Abs("./testdata/main.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err = kcl.GetSchemaType(p, nil, "Person")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, []*gpyrpc.KclType{
		{
			Filename:   p,
			PkgPath:    "__main__",
			Type:       "schema",
			SchemaName: "Person",
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type:    "str",
					Line:    1,
					Default: "\"kcl\"",
				},
				"age": {
					Type:    "int",
					Line:    2,
					Default: "1",
				},
			},
			Required: []string{"name", "age"},
		},
	}, result)
}

func TestGetSchemaTypeMapping(t *testing.T) {
	result, err := kcl.GetSchemaTypeMapping("test.k", "schema Person:\n    name: str\n\nschema Sub(Person):\n    count: int\n", "")
	if err != nil {
		t.Fatal(err)
	}
	personSchema := gpyrpc.KclType{
		Filename:   "test.k",
		PkgPath:    "__main__",
		Type:       "schema",
		SchemaName: "Person",
		Properties: map[string]*gpyrpc.KclType{
			"name": {
				Type: "str",
				Line: 1,
			},
		},
		Required: []string{"name"},
	}
	assert2.Equal(t, map[string]*gpyrpc.KclType{
		"Person": &personSchema,
		"Sub": {
			Filename:   "test.k",
			PkgPath:    "__main__",
			Type:       "schema",
			SchemaName: "Sub",
			BaseSchema: &personSchema,
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type: "str",
					Line: 1,
				},
				"count": {
					Type: "int",
					Line: 2,
				},
			},
			Required: []string{"count", "name"},
		},
	}, result)

	p, err := filepath.Abs("./testdata/main.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err = kcl.GetSchemaTypeMapping(p, nil, "Person")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, map[string]*gpyrpc.KclType{
		"Person": {
			Filename:   p,
			PkgPath:    "__main__",
			Type:       "schema",
			SchemaName: "Person",
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type:    "str",
					Line:    1,
					Default: "\"kcl\"",
				},
				"age": {
					Type:    "int",
					Line:    2,
					Default: "1",
				},
			},
			Required: []string{"name", "age"},
		},
	}, result)
}
func TestListUpStreamFiles(t *testing.T) {
	files, err := kcl.ListUpStreamFiles("./testdata/", &kcl.ListDepsOptions{Files: []string{"main.k", "app0/before/base.k", "app0/main.k"}})
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"app0/sub",
		"app0/sub/sub.k",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}

func TestListDepFiles(t *testing.T) {
	files, err := kcl.ListDepFiles("./testdata/app0", nil)
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

func TestTestAPI(t *testing.T) {
	result, err := kcl.Test(&kcl.TestOptions{
		PkgList: []string{"./testdata/test_module/..."},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, len(result.Info), 2)
	assert2.Equal(t, result.Info[0].ErrMessage, "")
	assert2.Equal(t, strings.Contains(result.Info[1].ErrMessage, "Error"), true, result.Info[1].ErrMessage)
}

func TestWithExternalpkg(t *testing.T) {
	absPath1, err := filepath.Abs("./testdata/external/external_1/")
	assert2.Equal(t, nil, err)
	absPath2, err := filepath.Abs("./testdata/external/external_2/")
	assert2.Equal(t, nil, err)
	opt := kcl.WithExternalPkgs("external_1="+absPath1, "external_2="+absPath2)
	result, err := kcl.Run("./testdata/import-external/main.k", opt)
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "{\"a\": \"Hello External_1 World!\", \"b\": \"Hello External_2 World!\"}", result.GetRawJsonResult())
	assert2.Equal(t, "a: Hello External_1 World!\nb: Hello External_2 World!", result.GetRawYamlResult())
}

func TestWithSortKeys(t *testing.T) {
	file, err := filepath.Abs("./testdata/test_plan/main.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err := kcl.Run(file, kcl.WithSortKeys(true))
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "a: 2\nb: 1", result.GetRawYamlResult())
}

func TestWithShowHidden(t *testing.T) {
	file, err := filepath.Abs("./testdata/test_plan/main.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err := kcl.Run(file, kcl.WithShowHidden(true))
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "b: 1\na: 2\n_c: 3", result.GetRawYamlResult())
	result, err = kcl.Run(file, kcl.WithShowHidden(true), kcl.WithSortKeys(true))
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "_c: 3\na: 2\nb: 1", result.GetRawYamlResult())
}

func TestWithLogger(t *testing.T) {
	file, err := filepath.Abs("./testdata/test_print/main.k")
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	result, err := kcl.Run(file, kcl.WithLogger(&buf))
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "hello: world", result.GetRawYamlResult())
	assert2.Equal(t, "Hello world\n", buf.String())
}
