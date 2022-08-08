// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_test

import (
	"fmt"
	"log"

	"kusionstack.io/kclvm-go"
)

func assert(v bool, a ...interface{}) {
	if !v {
		a = append([]interface{}{"assert failed"}, a...)
		log.Panic(a...)
	}
}

func ExampleMustRun() {
	yaml := kclvm.MustRun("testdata/main.k", kclvm.WithCode(`name = "kcl"`)).First().YAMLString()
	fmt.Println(yaml)

	// Output:
	// name: kcl
}

func ExampleMustRun_rawYaml() {
	const code = `
b = 1
a = 2
`
	yaml := kclvm.MustRun("testdata/main.k", kclvm.WithCode(code)).GetRawYamlResult()
	fmt.Println(yaml)

	yaml_sorted := kclvm.MustRun("testdata/main.k", kclvm.WithCode(code), kclvm.WithSortKeys(true)).GetRawYamlResult()
	fmt.Println(yaml_sorted)

	// _Output:
	// b: 1
	// a: 2
	//
	// a: 2
	// b: 1
}
func ExampleMustRun_schemaType() {
	const code = `
schema Person:
	name: str = ""

x = Person()
`
	json := kclvm.MustRun("testdata/main.k", kclvm.WithCode(code)).First().JSONString()
	fmt.Println(json)

	json_with_type := kclvm.MustRun("testdata/main.k", kclvm.WithCode(code), kclvm.WithIncludeSchemaTypePath(true)).First().JSONString()
	fmt.Println(json_with_type)

	// _Output:
	// {
	//     "x": {
	//         "name": ""
	//     }
	// }
	// {
	//     "x": {
	//         "@type": "Person",
	//         "name": ""
	//     }
	// }
}

func ExampleMustRun_settings() {
	yaml := kclvm.MustRun("./testdata/app0/kcl.yaml").First().YAMLString()
	fmt.Println(yaml)
}

func ExampleRunFiles() {
	result, _ := kclvm.RunFiles([]string{"./testdata/app0/kcl.yaml"})
	fmt.Println(result.First().YAMLString())
}

func ExampleKCLResult() {
	const k_code = `

name = "kcl"
age = 1
	
schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {name = "kcl-go"}
x1 = Person {age = 101}
`

	result := kclvm.MustRun("testdata/main.k", kclvm.WithCode(k_code)).First()

	fmt.Println("x0.name:", result.Get("x0.name"))
	fmt.Println("x1.age:", result.Get("x1.age"))

	// Output:
	// x0.name: kcl-go
	// x1.age: 101
}

func ExampleKCLResult_Get_struct() {
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

	result := kclvm.MustRun("testdata/main.k", kclvm.WithCode(k_code)).First()

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

func ExampleRun_getField() {
	// run kcl.yaml
	x, err := kclvm.Run("./testdata/app0/kcl.yaml")
	assert(err == nil, err)

	// print deploy_topology[1].zone
	fmt.Println(x.First().Get("deploy_topology.1.zone"))

	// Output:
	// R000A
}

func Example() {
	const k_code = `
import kcl_plugin.hello

name = "kcl"
age = 1

two = hello.add(1, 1)

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {}
x1 = Person {
	age = 101
}
`

	yaml := kclvm.MustRun("testdata/main.k", kclvm.WithCode(k_code)).First().YAMLString()
	fmt.Println(yaml)

	fmt.Println("----")

	result := kclvm.MustRun("./testdata/main.k").First()
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
	// import abc # unable to import
	// import a
	// import a # reimport

	results, err := kclvm.LintPath("testdata/lint/import.k")
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range results {
		fmt.Println(s)
	}

	// _Output:
	// Unable to import abc.
	// a is reimported multiple times.
	// a imported but unused.
}

func ExampleFormatCode() {
	out, err := kclvm.FormatCode(`a  =  1+2`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))

	// _Output:
	// a = 1 + 2
}

func ExampleFormatPath() {
	changedPaths, err := kclvm.FormatPath("testdata/fmt")
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
	x, err := kclvm.Run("hello.k", kclvm.WithCode(code),
		kclvm.WithOptions("name=kcl", "age=1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(x.First().YAMLString())

	// Output:
	// age: 1
	// name: kcl
}
