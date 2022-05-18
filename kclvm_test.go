// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_test

import (
	"flag"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"

	"kusionstack.io/kclvm-go"
	"kusionstack.io/kclvm-go/pkg/kcl"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

const tEnvNumCpu = "KCLVM_GO_API_TEST_NUM_CPU"

func TestMain(m *testing.M) {
	flag.Parse()

	if s := os.Getenv(tEnvNumCpu); s != "" {
		if x, err := strconv.Atoi(s); err == nil {
			println(tEnvNumCpu, "=", s)
			kclvm.InitKclvmRuntime(x)
		}
	}

	os.Exit(m.Run())
}

func TestRunFiles(t *testing.T) {
	_, err := kclvm.RunFiles([]string{"./testdata/app0/kcl.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = kclvm.RunFiles([]string{"./testdata/app0/kcl.yaml"})
	if err != nil {
		t.Fatal(err)
	}

	chErr := make(chan error, 3)

	var wg sync.WaitGroup
	for i := 0; i < cap(chErr); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := kclvm.RunFiles([]string{"./testdata/app0/kcl.yaml"})
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

	result, err := kclvm.Run(testdata_main_k,
		kclvm.WithCode(code),
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

	result, err = kclvm.Run(testdata_main_k,
		kclvm.WithCode(code),
		kclvm.WithOverrides(":a1.image=\"new-a1-image\""),
		kclvm.WithOverrides("__main__:a2.image=\"new-a2-image:v123\""),
		kclvm.WithPrintOverridesAST(true),
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

a2 = App {
    image = "new-a2-image:v123"
    name = "a2-app"
}`)
	got := strings.TrimSpace(string(data))

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("golden mismatch (-want +got):\n%s", diff)
	}
}

func _BenchmarkRunFilesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := kclvm.RunFiles([]string{"./testdata/app0/kcl.yaml"})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestWithKFilenames(t *testing.T) {
	kclvm.WithKFilenames("/testdata/main.k")
}

func TestWithOptions(t *testing.T) {
	kclvm.WithOptions("key1=value1", "key2=value2")
}

func TestWithSettings(t *testing.T) {
	kclvm.WithSettings("a_settings.yml")
}

func TestWithWorkDir(t *testing.T) {
	wd, _ := os.Getwd()
	kclvm.WithWorkDir(wd)
}

func TestWithDisableNone(t *testing.T) {
	kclvm.WithDisableNone(true)
}

func TestFormatCode(t *testing.T) {
	result, err := kclvm.FormatCode("a=1")
	if err != nil {
		t.Error(err)
	}
	assert2.Equalf(t, string(result), "a = 1\n", "format result unexpected: expect: a = 1\n, actual: %s", result)
}

func TestPlugin(t *testing.T) {
	const code = `
import kcl_plugin.hello as hello

a = hello.add(1, 2)
`
	_, err := kclvm.Run("testdata/main.k",
		kclvm.WithCode(code),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEvalCode(t *testing.T) {
	_, err := kclvm.EvalCode(`name = "kcl"`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSchemaType(t *testing.T) {
	result, err := kclvm.GetSchemaType("", "schema Person:\n    name: str", "")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, []*gpyrpc.KclType{
		{
			Type:       "schema",
			SchemaName: "Person",
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type: "str",
					Line: 1,
				},
			},
			Required: []string{"name"},
		},
	}, result)
	result, err = kcl.GetSchemaType("./testdata/main.k", "", "Person")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, []*gpyrpc.KclType{
		{
			Type:       "schema",
			SchemaName: "Person",
			Properties: map[string]*gpyrpc.KclType{
				"name": {
					Type:    "str",
					Line:    1,
					Default: "kcl",
				},
				"age": {
					Type:    "int",
					Line:    2,
					Default: "1",
				},
			},
			Required: []string{"age", "name"},
		},
	}, result)
}

func TestListUpStreamFiles(t *testing.T) {
	files, err := kclvm.ListUpStreamFiles("./testdata/", &kclvm.ListDepsOption{Files: []string{"main.k", "app0/before/base.k", "app0/main.k"}})
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"main.k",
		"app0/before/base.k",
		"app0/main.k",
		"app0/sub",
		"app0/sub/sub.k",
		"kcl_plugin/hello",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}

func TestListDepFiles(t *testing.T) {
	files, err := kclvm.ListDepFiles("./testdata/app0", nil)
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
