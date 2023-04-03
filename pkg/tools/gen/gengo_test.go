// Copyright 2022 The KCL Authors. All rights reserved.

package gen_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"kusionstack.io/kclvm-go/pkg/tools/gen"
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
`
	var buf bytes.Buffer
	err := gen.GenGo(&buf, "hello.k", code, nil)
	if err != nil {
		log.Fatal(err)
	}
	goCode := buf.String()
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
}
