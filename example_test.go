// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_test

import (
	"fmt"
	"log"

	kcl "kcl-lang.io/kcl-go"
)

func assert(v bool, a ...interface{}) {
	if !v {
		a = append([]interface{}{"assert failed"}, a...)
		log.Panic(a...)
	}
}

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

func ExampleRun_getField() {
	// run kcl.yaml
	x, err := kcl.Run("./testdata/app0/kcl.yaml")
	assert(err == nil, err)

	// print deploy_topology[1].zone
	fmt.Println(x.First().Get("deploy_topology.1.zone"))

	// Output:
	// R000A
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
