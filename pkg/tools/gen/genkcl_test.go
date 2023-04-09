package gen

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestGenKcl(t *testing.T) {
	const code = `
package main
type employee struct {
	name        string            ` + "`kcl:\"name=name,type=str\"`" + `           // kcl-type: str
	age         int               ` + "`kcl:\"name=age,type=int\"`" + `            // kcl-type: int
	friends     []string          ` + "`kcl:\"name=friends,type=[str]\"`" + `      // kcl-type: [str]
	movies      map[string]*Movie ` + "`kcl:\"name=movies,type={str:Movie}\"`" + ` // kcl-type: {str:Movie}
	bankCard    int               ` + "`kcl:\"nam=bankCard,type=int\",abc:\"name=bankCard,type=int\"`" + `       // kcl-type: int
	nationality string            ` + "`kcl:\"name=nationality,type=str\"`" + `    // kcl-type: str
}`
	var buf bytes.Buffer
	opts := &GenKclOptions{
		ParseFromTag: false,
	}
	err := GenKcl(&buf, "./testdata/genkcldata.go", nil, opts)
	// err := GenKcl(&buf, "demo", code, opts)
	if err != nil {
		log.Fatal(err)
	}
	kclCode := buf.String()
	fmt.Println("###############")
	fmt.Print(kclCode)
	expectedKclCodeFromTag := `
schema Person:
    """Person Example"""
    name: str
    age: int
    friends: [str]
    movies: {str:Movie}
    MapInterface: {str:{str:any}}
    Ep: employee
    Com: Company
    StarInt: int
    StarMap: {str:str}
    Inter: any

schema Movie:
    desc: str
    size: units.NumberMultiplier
    kind?: "Superhero"|"War"|"Unknown"
    unknown1?: int|str
    unknown2?: any

schema employee:
    name: str
    age: int
    friends: [str]
    movies: {str:Movie}
    bankCard: int
    nationality: str

schema Company:
    name: str
    employees: [employee]
    persons: Person

`
	expectedKclCodeFromField := `
schema Person:
    """Person Example"""
    name: str
    age: int
    friends: [str]
    movies: {str:Movie}
    MapInterface: {str:{str:any}}
    Ep: employee
    Com: Company
    StarInt: int
    StarMap: {str:str}
    Inter: any

schema Movie:
    desc: str
    size: int
    kind: str
    unknown1: any
    unknown2: any

schema employee:
    name: str
    age: int
    friends: [str]
    movies: {str:Movie}
    bankCard: int
    nationality: str

schema Company:
    name: str
    employees: [employee]
    persons: Person

`
	if opts.ParseFromTag {
		if kclCode != expectedKclCodeFromTag {
			panic(fmt.Sprintf("test failed, expected %s got %s", expectedKclCodeFromTag, kclCode))
		}
	} else {
		if kclCode != expectedKclCodeFromField {
			panic(fmt.Sprintf("test failed, expected %s got %s", expectedKclCodeFromField, kclCode))
		}
	}

}
