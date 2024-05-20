package gen

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestGenKcl(t *testing.T) {
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

func TestGenKclFromJsonSchema(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectFilepath string
		expect         string
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
			name:           caseFile.Name(),
			input:          input,
			expectFilepath: expectFilepath,
			expect:         readFileString(t, expectFilepath),
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

func TestGenKclFromJson(t *testing.T) {
	input := filepath.Join("testdata", "json", "input.json")
	expectFilepath := filepath.Join("testdata", "json", "expect.k")
	expect := readFileString(t, expectFilepath)

	var buf bytes.Buffer
	err := GenKcl(&buf, input, nil, &GenKclOptions{})
	if err != nil {
		t.Fatal(err)
	}
	result := buf.Bytes()
	assert2.Equal(t, expect, string(bytes.ReplaceAll(result, []byte("\r\n"), []byte("\n"))))
}

func TestGenKclFromYaml(t *testing.T) {
	type testCase struct {
		name   string
		input  string
		expect string
	}
	var cases []testCase

	casesPath := filepath.Join("testdata", "yaml")
	caseFiles, err := os.ReadDir(casesPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, caseFile := range caseFiles {
		input := filepath.Join(casesPath, caseFile.Name(), "input.yaml")
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

func TestGenKclFromMultipleResourceYaml(t *testing.T) {
	type testCase struct {
		name   string
		input  string
		expect string
	}
	var cases []testCase

	casesPath := filepath.Join("testdata", "yaml2")
	caseFiles, err := os.ReadDir(casesPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, caseFile := range caseFiles {
		input := filepath.Join(casesPath, caseFile.Name(), "input.yaml")
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

func TestGenKclFromProto(t *testing.T) {
	var buf bytes.Buffer
	err := GenKcl(&buf, `./testdata/proto/proto2kcl.proto`, nil, &GenKclOptions{})
	if err != nil {
		t.Fatal(err)
	}

	b, _ := os.ReadFile("./testdata/proto/proto2kcl.k")
	assert2.Equal(t, string(b), buf.String())
}

type TestData = data

func TestGenKclFromJsonAndImports(t *testing.T) {
	file := kclFile{}
	g := &kclGenerator{}
	b := bytes.NewBuffer(nil)
	// Get a KCL config from JSON or YAML string
	data, err := convertKclFromYamlString([]byte(`workload:
  containers:
    nginx:
      image: nginx:v2
  replicas: 2
`))
	if err != nil {
		t.Fatal(err)
	}
	// Add import statements
	importStmt := kImport{
		PkgPath: "models.schema.v1",
		Alias:   "ac",
	}
	file.Imports = append(file.Imports, importStmt)
	configSchemaName := strings.Join([]string{importStmt.PkgName(), "AppConfiguration"}, ".")
	// Add the configuration `app1` at top level.
	file.Config = append(file.Config, config{
		Data:    data,
		IsUnion: true,
		Var:     "app1",
		Name:    configSchemaName,
	})
	// Add the configuration `app2` at top level with more data configs.
	data = append(data, TestData{
		Key: "labels",
		Value: config{
			Name: "Labels",
			Data: []TestData{{
				Key:   "app",
				Value: "nginx",
			}},
		},
	})
	data = append(data, TestData{
		Key: "annotations",
		Value: config{
			Data: []TestData{{
				Key: "app1",
				Value: map[string]any{
					"app": "nginx",
				},
			}, {
				Key: "app2",
				Value: config{
					Name: "Annotation",
					Data: []TestData{{
						Key: "app1",
						Value: map[string]string{
							"app": "nginx",
						},
					},
					},
				},
			}},
		},
	})
	file.Config = append(file.Config, config{
		Data:    data,
		IsUnion: true,
		Var:     "app2",
		Name:    configSchemaName,
	})
	// Generate KCL code.
	g.genKcl(b, file)
	assert2.Equal(t, strings.ReplaceAll(b.String(), "\r\n", "\n"), `"""
This file was generated by the KCL auto-gen tool. DO NOT EDIT.
Editing this file might prove futile when you re-run the KCL auto-gen generate command.
"""
import models.schema.v1 as ac

app1: ac.AppConfiguration {
    workload = {
        containers = {
            nginx = {
                image = "nginx:v2"
            }
        }
        replicas = 2
    }
}
app2: ac.AppConfiguration {
    workload = {
        containers = {
            nginx = {
                image = "nginx:v2"
            }
        }
        replicas = 2
    }
    labels = Labels {
        app = "nginx"
    }
    annotations = {
        app1 = {"app": "nginx"}
        app2 = Annotation {
            app1 = {"app": "nginx"}
        }
    }
}
`)
}

func readFileString(t testing.TB, p string) (content string) {
	data, err := os.ReadFile(p)
	if err != nil {
		t.Errorf("read file failed, %s", err)
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return string(data)
}
