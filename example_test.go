// Copyright The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package kcl_test

import (
	"fmt"
	"log"

	kcl "kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/native"
	"kcl-lang.io/kcl-go/pkg/parser"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

func ExampleMustRun() {
	yaml := kcl.MustRun("testdata/main.k", kcl.WithCode(`name = "kcl"`)).First().YAMLString()
	fmt.Println(yaml)

	// Output:
	// name: kcl
}

func ExampleMustRun_rawYaml() {
	const code = `
b = 1
a = 2
`
	yaml := kcl.MustRun("testdata/main.k", kcl.WithCode(code)).GetRawYamlResult()
	fmt.Println(yaml)

	yaml_sorted := kcl.MustRun("testdata/main.k", kcl.WithCode(code), kcl.WithSortKeys(true)).GetRawYamlResult()
	fmt.Println(yaml_sorted)

	// Output:
	// b: 1
	// a: 2
	// a: 2
	// b: 1
}
func ExampleMustRun_schemaType() {
	const code = `
schema Person:
	name: str = ""

x = Person()
`
	json := kcl.MustRun("testdata/main.k", kcl.WithCode(code)).First().JSONString()
	fmt.Println(json)

	// Output:
	// {
	//     "x": {
	//         "name": ""
	//     }
	// }
}

func ExampleMustRun_settings() {
	yaml := kcl.MustRun("./testdata/app0/kcl.yaml").First().YAMLString()
	fmt.Println(yaml)
}

func ExampleRunFiles() {
	result, _ := kcl.RunFiles([]string{"./testdata/app0/kcl.yaml"})
	fmt.Println(result.First().YAMLString())
}

func ExampleKCLResult() {
	const k_code = `
name = "kcl"
age = 1
	
two = 2
	
schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {name = "kcl-go"}
x1 = Person {age = 101}
`

	result := kcl.MustRun("testdata/main.k", kcl.WithCode(k_code)).First()

	fmt.Println("x0.name:", result.Get("x0.name"))
	fmt.Println("x1.age:", result.Get("x1.age"))

	// Output:
	// x0.name: kcl-go
	// x1.age: 101
}

func ExampleKCLResult_get() {
	const k_code = `
schema Person:
    name: str = "kcl"
    age: int = 1
    X: int = 2

x = {
    "a": Person {age = 101}
    "b": 123
}
`

	result := kcl.MustRun("testdata/main.k", kcl.WithCode(k_code)).First()

	var person struct {
		Name string
		Age  int
	}
	fmt.Printf("person: %+v\n", result.Get("x.a", &person))
	fmt.Printf("person: %+v\n", person)

	// Output:
	// person: &{Name:kcl Age:101}
	// person: {Name:kcl Age:101}
}

func Example() {
	const k_code = `
name = "kcl"
age = 1

two = 2

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {}
x1 = Person {
	age = 101
}
`

	yaml := kcl.MustRun("testdata/main.k", kcl.WithCode(k_code)).First().YAMLString()
	fmt.Println(yaml)

	fmt.Println("----")

	result := kcl.MustRun("./testdata/main.k").First()
	fmt.Println(result.JSONString())

	fmt.Println("----")
	fmt.Println("x0.name:", result.Get("x0.name"))
	fmt.Println("x1.age:", result.Get("x1.age"))

	fmt.Println("----")

	var person struct {
		Name string
		Age  int
	}
	fmt.Printf("person: %+v\n", result.Get("x1", &person))
}

func ExampleLintPath() {
	// import a
	// import a # reimport

	results, err := kcl.LintPath([]string{"testdata/lint/import.k"})
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range results {
		fmt.Println(s)
	}

	// Output:
	// Module 'a' is reimported multiple times
	// Module 'a' imported but unused
	// Module 'a' imported but unused
}

func ExampleFormatCode() {
	out, err := kcl.FormatCode(`a  =  1+2`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

	// Output:
	// a = 1 + 2
}

func ExampleFormatPath() {
	changedPaths, err := kcl.FormatPath("testdata/fmt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(changedPaths)
}

func ExampleWithOptions() {
	const code = `
name = option("name")
age = option("age")
`
	x, err := kcl.Run("hello.k", kcl.WithCode(code),
		kcl.WithOptions("name=kcl", "age=1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(x.First().YAMLString())

	// Output:
	// age: 1
	// name: kcl
}

func ExampleParseFile() {
	result, err := parser.ParseFile("testdata/main.k", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func ExampleParseProgram() {
	result, err := kcl.ParseProgram(&kcl.ParseProgramArgs{
		Paths: []string{"testdata/main.k"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func ExampleLoadPackage() {
	result, err := kcl.LoadPackage(&kcl.LoadPackageArgs{
		ParseArgs: &kcl.ParseProgramArgs{
			Paths: []string{"testdata/main.k"},
		},
		ResolveAst: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func ExampleListVariables() {
	result, err := kcl.ListVariables(&kcl.ListVariablesArgs{
		Files: []string{"testdata/main.k"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func ExampleListOptions() {
	result, err := kcl.ListOptions(&kcl.ListOptionsArgs{
		Paths: []string{"testdata/option/main.k"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func ExampleUpdateDependencies() {
	// [package]
	// name = "mod_update"
	// edition = "0.0.1"
	// version = "0.0.1"
	//
	// [dependencies]
	// helloworld = { oci = "oci://ghcr.io/kcl-lang/helloworld", tag = "0.1.0" }
	// flask = { git = "https://github.com/kcl-lang/flask-demo-kcl-manifests", commit = "ade147b" }

	result, err := kcl.UpdateDependencies(&gpyrpc.UpdateDependenciesArgs{
		ManifestPath: "testdata/update_dependencies",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func ExampleUpdateDependencies_execProgram() {
	// [package]
	// name = "mod_update"
	// edition = "0.0.1"
	// version = "0.0.1"
	//
	// [dependencies]
	// helloworld = { oci = "oci://ghcr.io/kcl-lang/helloworld", tag = "0.1.0" }
	// flask = { git = "https://github.com/kcl-lang/flask-demo-kcl-manifests", commit = "ade147b" }

	svc := native.NewNativeServiceClient()

	result, err := svc.UpdateDependencies(&gpyrpc.UpdateDependenciesArgs{
		ManifestPath: "testdata/update_dependencies",
	})
	if err != nil {
		log.Fatal(err)
	}

	// import helloworld
	// import flask
	// a = helloworld.The_first_kcl_program
	// fmt.Println(result.ExternalPkgs)

	execResult, err := svc.ExecProgram(&gpyrpc.ExecProgramArgs{
		KFilenameList: []string{"testdata/update_dependencies/main.k"},
		ExternalPkgs:  result.ExternalPkgs,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(execResult.YamlResult)

	// Output:
	// a: Hello World!
}
