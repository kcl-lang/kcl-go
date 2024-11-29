// Copyright 2023 The KCL Authors. All rights reserved.

package gen_test

import (
	"bytes"
	"fmt"
	"log"
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
	/*
			expectedGoCode := `
		// Person Example
		type Person struct {
		    name string              // kcl-type: str
		    age int                  // kcl-type: int
		    friends []string         // kcl-type: [str]
		    movies map[string]*Movie // kcl-type: {str:Movie}
		}

		type Movie struct {
		    desc string          // kcl-type: str
		    size int             // kcl-type: units.NumberMultiplier
		    kind string          // kcl-type: "Superhero"|"War"|"Unknown"
		    unknown1 interface{} // kcl-type: int|str
		    unknown2 interface{} // kcl-type: any
		}
		`
			if goCode != expectedGoCode {
				panic(fmt.Sprintf("test failed, expected %s got %s", expectedGoCode, goCode))
			}
	*/
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
	})
	if err != nil {
		log.Fatal(err)
	}
	goCode := buf.String()
	fmt.Println(goCode)
	/*
			expectedGoCode := `
		// Person Example
		type Person struct {
		    name string              // kcl-type: str
		    age int                  // kcl-type: int
		    friends []string         // kcl-type: [str]
		    movies map[string]*Movie // kcl-type: {str:Movie}
		}

		type Movie struct {
		    desc string          // kcl-type: str
		    size int             // kcl-type: units.NumberMultiplier
		    kind string          // kcl-type: "Superhero"|"War"|"Unknown"
		    unknown1 interface{} // kcl-type: int|str
		    unknown2 interface{} // kcl-type: any
		}
		`
			if goCode != expectedGoCode {
				panic(fmt.Sprintf("test failed, expected %s got %s", expectedGoCode, goCode))
			}
	*/
}
