// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_test

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"

	kcl "kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const tEnvNumCpu = "KCLVM_GO_API_TEST_NUM_CPU"

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

func TestRunFiles(t *testing.T) {
	_, err := kcl.RunFiles([]string{"./testdata/app0/kcl.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = kcl.RunFiles([]string{"./testdata/app0/kcl.yaml"})
	if err != nil {
		t.Fatal(err)
	}

	chErr := make(chan error, 3)

	var wg sync.WaitGroup
	for i := 0; i < cap(chErr); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := kcl.RunFiles([]string{"./testdata/app0/kcl.yaml"})
			chErr <- e
		}()
	}
	wg.Wait()

	for i := 0; i < cap(chErr); i++ {
		if e := <-chErr; e != nil {
			t.Fatal(e)
		}
	}
}

func TestIndlu(t *testing.T) {
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
	const testdata_main_k = "testdata/main_include_schema_type_path.k"
	kfile, err := os.Create(testdata_main_k)
	if err != nil {
		t.Fatal(err)
	}
	kfile.Close()

	result, err := kcl.Run(testdata_main_k,
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
	if expect, got := "default", result.First().Get("a1.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}
	if expect, got := "a2-image", result.First().Get("a2.image"); expect != got {
		t.Fatalf("expect = %v, got = %v", expect, got)
	}

	os.Remove(testdata_main_k)
	defer os.Remove(testdata_main_k)
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
		kcl.WithOverrides(":a1.image=\"new-a1-image\""),
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

func _BenchmarkRunFilesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := kcl.RunFiles([]string{"./testdata/app0/kcl.yaml"})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
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
					Type:       "str",
					Line:       1,
					Properties: map[string]*gpyrpc.KclType{},
					Required:   []string{},
					UnionTypes: []*gpyrpc.KclType{},
					Decorators: []*gpyrpc.Decorator{},
					Examples:   map[string]*gpyrpc.Example{},
				},
			},
			Required:   []string{"name"},
			UnionTypes: []*gpyrpc.KclType{},
			Decorators: []*gpyrpc.Decorator{
				{
					Name:      "info",
					Arguments: []string{"\"v\""},
					Keywords: map[string]string{
						"k": "\"v\"",
					},
				},
			},
			Examples: map[string]*gpyrpc.Example{},
		},
	}, result)
	p, err := filepath.Abs("./testdata/main.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err = kcl.GetSchemaType(p, "", "Person")
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
					Type:       "str",
					Line:       1,
					Default:    "\"kcl\"",
					Properties: map[string]*gpyrpc.KclType{},
					Required:   []string{},
					UnionTypes: []*gpyrpc.KclType{},
					Decorators: []*gpyrpc.Decorator{},
					Examples:   map[string]*gpyrpc.Example{},
				},
				"age": {
					Type:       "int",
					Line:       2,
					Default:    "1",
					Properties: map[string]*gpyrpc.KclType{},
					Required:   []string{},
					UnionTypes: []*gpyrpc.KclType{},
					Decorators: []*gpyrpc.Decorator{},
					Examples:   map[string]*gpyrpc.Example{},
				},
			},
			Required:   []string{"name", "age"},
			UnionTypes: []*gpyrpc.KclType{},
			Decorators: []*gpyrpc.Decorator{},
			Examples:   map[string]*gpyrpc.Example{},
		},
	}, result)
}

func TestGetSchemaTypeMapping(t *testing.T) {
	result, err := kcl.GetSchemaTypeMapping("test.k", "schema Person:\n    name: str", "")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, map[string]*gpyrpc.KclType{
		"Person": {
			Filename:   "test.k",
			PkgPath:    "__main__",
			Type:       "schema",
			SchemaName: "Person",
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type:       "str",
					Line:       1,
					Properties: map[string]*gpyrpc.KclType{},
					Required:   []string{},
					UnionTypes: []*gpyrpc.KclType{},
					Decorators: []*gpyrpc.Decorator{},
					Examples:   map[string]*gpyrpc.Example{},
				},
			},
			Required:   []string{"name"},
			UnionTypes: []*gpyrpc.KclType{},
			Decorators: []*gpyrpc.Decorator{},
			Examples:   map[string]*gpyrpc.Example{},
		},
	}, result)

	p, err := filepath.Abs("./testdata/main.k")
	if err != nil {
		t.Fatal(err)
	}
	result, err = kcl.GetSchemaTypeMapping(p, "", "Person")
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
					Type:       "str",
					Line:       1,
					Default:    "\"kcl\"",
					Properties: map[string]*gpyrpc.KclType{},
					Required:   []string{},
					UnionTypes: []*gpyrpc.KclType{},
					Decorators: []*gpyrpc.Decorator{},
					Examples:   map[string]*gpyrpc.Example{},
				},
				"age": {
					Type:       "int",
					Line:       2,
					Default:    "1",
					Properties: map[string]*gpyrpc.KclType{},
					Required:   []string{},
					UnionTypes: []*gpyrpc.KclType{},
					Decorators: []*gpyrpc.Decorator{},
					Examples:   map[string]*gpyrpc.Example{},
				},
			},
			Required:   []string{"name", "age"},
			UnionTypes: []*gpyrpc.KclType{},
			Decorators: []*gpyrpc.Decorator{},
			Examples:   map[string]*gpyrpc.Example{},
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
	absPath1, err := filepath.Abs("./testdata_external/external_1/")
	assert2.Equal(t, nil, err)
	absPath2, err := filepath.Abs("./testdata_external/external_2/")
	assert2.Equal(t, nil, err)
	opt := kcl.WithExternalPkgs("external_1="+absPath1, "external_2="+absPath2)
	result, err := kcl.Run("./testdata/import-external/main.k", opt)
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "[{\"a\": \"Hello External_1 World!\", \"b\": \"Hello External_2 World!\"}]", result.GetRawJsonResult())
	assert2.Equal(t, "a: Hello External_1 World!\nb: Hello External_2 World!", result.GetRawYamlResult())
}
