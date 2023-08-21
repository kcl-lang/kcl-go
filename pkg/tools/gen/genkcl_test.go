package gen

import (
	"bytes"
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
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

func TestGenKclFromJson(t *testing.T) {
	type testCase struct {
		name   string
		input  string
		expect string
	}
	var cases []testCase

	casesPath := filepath.Join("testdata", "jsonschema")
	caseFiles, err := os.ReadDir(casesPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, caseFile := range caseFiles {
		input := filepath.Join(casesPath, caseFile.Name(), "input.json")
		expectFilepath := filepath.Join(casesPath, caseFile.Name(), "expect.k")
		cases = append(cases, testCase{
			name:   caseFile.Name(),
			input:  input,
			expect: readFileString(t, expectFilepath),
		})
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := GenKcl(&buf, testcase.input, nil, &GenKclOptions{})
			if err != nil {
				t.Fatal(err)
			}
			result := buf.Bytes()
			assert2.Equal(t, testcase.expect, string(bytes.ReplaceAll(result, []byte("\r\n"), []byte("\n"))))
		})
	}
}

func TestGenKclFromTerraform(t *testing.T) {
	input := filepath.Join("testdata", "terraform", "schema.json")
	expectFilepath := filepath.Join("testdata", "terraform", "expect.k")
	expect := readFileString(t, expectFilepath)

	var buf bytes.Buffer
	err := GenKcl(&buf, input, nil, &GenKclOptions{})
	if err != nil {
		t.Fatal(err)
	}
	result := buf.Bytes()
	assert2.Equal(t, expect, string(bytes.ReplaceAll(result, []byte("\r\n"), []byte("\n"))))
}

func readFileString(t testing.TB, p string) (content string) {
	data, err := os.ReadFile(p)
	if err != nil {
		t.Errorf("read file failed, %s", err)
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return string(data)
}
