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

func TestGenKclFromToml(t *testing.T) {
	input := filepath.Join("testdata", "toml", "input.toml")
	expectFilepath := filepath.Join("testdata", "toml", "expect.k")
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
	t.Run("not UseIntegersForNumbers", func(t *testing.T) {
		var buf bytes.Buffer
		err := GenKcl(&buf, `./testdata/proto/proto2kcl.proto`, nil, &GenKclOptions{UseIntegersForNumbers: false})
		if err != nil {
			t.Fatal(err)
		}

		b, _ := os.ReadFile("./testdata/proto/proto2kcl.k")
		assert2.Equal(t, string(b), buf.String())
	})

	t.Run("UseIntegersForNumbers", func(t *testing.T) {
		var buf bytes.Buffer
		err := GenKcl(&buf, `./testdata/proto/proto2kcl.proto`, nil, &GenKclOptions{UseIntegersForNumbers: true})
		if err != nil {
			t.Fatal(err)
		}

		b, _ := os.ReadFile("./testdata/proto/proto2kcl_for_num.k")
		assert2.Equal(t, string(b), buf.String())
	})
}

func TestGenKclFromTextProto(t *testing.T) {
	generator := TextProtoGenerator{}
	got, err := generator.GenFromSchemaFile("./testdata/textproto/data.textproto", nil, "./testdata/textproto/schema.k", nil, "Config")
	assert2.Nil(t, err)
	expected := &config{
		Name: "Config",
		Data: []data{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2.0},
			{Key: "c", Value: true},
			{Key: "d", Value: "value"},
			{Key: "empty1", Value: []interface{}(nil)},
			{Key: "empty2", Value: []interface{}(nil)},
			{Key: "int1", Value: []interface{}{1, 2}},
			{Key: "int2", Value: []interface{}{1, 2}},
			{Key: "int3", Value: []interface{}{1, 2}},
			{Key: "string1", Value: []interface{}{"a", "b"}},
			{Key: "float1", Value: []interface{}{100.0, 1.0, 0.0}},
			{Key: "map", Value: []data{{Key: "foo", Value: 2}, {Key: "bar", Value: 3}}},
		},
	}
	assert2.Equal(t, expected, got)
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
        app1 = {app = "nginx"}
        app2 = Annotation {
            app1 = {app = "nginx"}
        }
    }
}
`)
}

func TestMarshalKcl(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "bool true",
			input:    true,
			expected: "True",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "False",
		},
		{
			name:     "int",
			input:    123,
			expected: "123",
		},
		{
			name:     "float",
			input:    123.456,
			expected: "123.456",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "\"hello\"",
		},
		{
			name:  "array",
			input: [3]int{1, 2, 3},
			expected: `[
    1
    2
    3
]`,
		},
		{
			name:  "slice",
			input: []string{"A", "B", "C"},
			expected: `[
    "A"
    "B"
    "C"
]`,
		},
		{
			name:  "map",
			input: map[string]interface{}{"name": "example", "age": 30},
			expected: `{
    age = 30
    name = "example"
}`,
		},
		{
			name:  "struct",
			input: struct{ Name string }{"John"},
			expected: `{
    Name = "John"
}`,
		},
		{
			name: "nested",
			input: map[string]interface{}{
				"details": map[string]interface{}{
					"age": 30,
					"address": map[string]interface{}{
						"city":    "New York",
						"zipcode": 10001,
					},
				},
				"active": true,
			},
			expected: "{\n    active = True\n    details = {\n        address = {\n            city = \"New York\"\n            zipcode = 10001\n        }\n        age = 30\n    }\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, string(got))
			}
		})
	}
}

func readFileString(t testing.TB, p string) (content string) {
	data, err := os.ReadFile(p)
	if err != nil {
		t.Errorf("read file failed, %s", err)
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return string(data)
}
