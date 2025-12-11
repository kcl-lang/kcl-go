// Copyright 2023 The KCL Authors. All rights reserved.

package gen_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
)

var _ = gen.GenGo

func TestGenGo(t *testing.T) {
	const code = `
import units

type NumberMultiplier = units.NumberMultiplier

schema Person:
	"""Person Example"""
	name: str = "kcl"
	age: int = 2
	friends?: [str] = None
	movies?: {str: Movie} = None

schema Movie:
	desc: str = ""
	size: NumberMultiplier = 2M
	kind?: "Superhero" | "War" | "Unknown"
	unknown1?: int | str = None
	unknown2?: any = None

schema employee(Person):
    bankCard: int
    nationality: str

schema Company:
    name: str
    employees: [employee]
    persons: Person
`
	var buf bytes.Buffer
	err := gen.GenGo(&buf, "hello.k", code, nil)
	if err != nil {
		log.Fatal(err)
	}
	goCode := buf.String()
	fmt.Println(goCode)
	if !strings.Contains(goCode, "movies map[string]*Movie `kcl:\"name=movies,type={str:Movie}\"` // kcl-type: {str:Movie}") {
		panic(fmt.Sprintf("test failed, got %s", goCode))
	}
}

func TestGenGoWithRename(t *testing.T) {
	const code = `
import units

type NumberMultiplier = units.NumberMultiplier

schema Person:
	"""Person Example"""
	name: str = "kcl"
	age: int = 2
	friends?: [str] = None
	movies?: {str: Movie} = None

schema Movie:
	desc: str = ""
	size: NumberMultiplier = 2M
	kind?: "Superhero" | "War" | "Unknown"
	unknown1?: int | str = None
	unknown2?: any = None

schema employee(Person):
    bankCard: int
    nationality: str

schema Company:
    name: str
    employees: [employee]
    persons: Person
`
	var buf bytes.Buffer
	err := gen.GenGo(&buf, "hello.k", code, &gen.GenGoOptions{
		RenameFunc: func(name string) string {
			return strcase.ToCamel(name)
		},
		OmitTag: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	goCode := buf.String()
	fmt.Println(goCode)
	expectedGoCode := `
type Movie struct {
    Desc string  // kcl-type: str
    Size int     // kcl-type: units.NumberMultiplier
    Kind string  // kcl-type: "Superhero"|"War"|"Unknown"
    Unknown1 any // kcl-type: int|str
    Unknown2 any // kcl-type: any
}

type employee struct {
    Name string              // kcl-type: str
    Age int                  // kcl-type: int
    Friends []string         // kcl-type: [str]
    Movies map[string]*Movie // kcl-type: {str:Movie}
    BankCard int             // kcl-type: int
    Nationality string       // kcl-type: str
}

type Company struct {
    Name string           // kcl-type: str
    Employees []*employee // kcl-type: [employee]
    Persons *Person       // kcl-type: Person
}

// Person Example
type Person struct {
    Name string              // kcl-type: str
    Age int                  // kcl-type: int
    Friends []string         // kcl-type: [str]
    Movies map[string]*Movie // kcl-type: {str:Movie}
}
`
	if !strings.Contains(goCode, "Unknown1 any // kcl-type: int|str") {
		panic(fmt.Sprintf("test failed, expected %s got %s", expectedGoCode, goCode))
	}
	if !strings.Contains(goCode, "Unknown2 any // kcl-type: any") {
		panic(fmt.Sprintf("test failed, expected %s got %s", expectedGoCode, goCode))
	}
}

func TestGenGoWithIndexSignatureAndFunction(t *testing.T) {
	const code = `
schema Test:
	name: str
	surname?: str

schema TestMap:
	[name: str]: Test = {name = name}

schema AppConfig:
    maps?: TestMap
    my_func?: (int) -> int
    replicas: int
`
	var buf bytes.Buffer
	err := gen.GenGo(&buf, "hello.k", code, &gen.GenGoOptions{
		RenameFunc: func(name string) string {
			return strcase.ToCamel(name)
		},
		OmitTag: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	goCode := buf.String()
	fmt.Println(goCode)
	expectedGoCode := `
type TestMap map[string]*Test

type AppConfig struct {
    Maps *TestMap        // kcl-type: TestMap
    MyFunc func(int) int // kcl-type: (int) -> int
    Replicas int         // kcl-type: int
}

type Test struct {
    Name string    // kcl-type: str
    Surname string // kcl-type: str
}
`
	if !strings.Contains(goCode, "type TestMap map[string]*Test") {
		panic(fmt.Sprintf("test failed, expected %s got %s", expectedGoCode, goCode))
	}
	if !strings.Contains(goCode, "MyFunc func(int) int // kcl-type: (int) -> int") {
		panic(fmt.Sprintf("test failed, expected %s got %s", expectedGoCode, goCode))
	}
}
