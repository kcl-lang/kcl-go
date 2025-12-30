package gen

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestIndexContent(t *testing.T) {
	rootPkg := KclPackage{
		Name: "test",
		SubPackageList: []*KclPackage{
			{
				Name: "sub1",
				SubPackageList: []*KclPackage{
					{
						Name: "sub2",
						SchemaList: []*KclOpenAPIType{
							{
								KclExtensions: &KclExtensions{
									XKclModelType: &XKclModelType{
										Type: "C",
									},
								},
							},
						},
					},
				},
				SchemaList: []*KclOpenAPIType{
					{
						KclExtensions: &KclExtensions{
							XKclModelType: &XKclModelType{
								Type: "B",
							},
						},
					},
				},
			},
		},
		SchemaList: []*KclOpenAPIType{
			{
				KclExtensions: &KclExtensions{
					XKclModelType: &XKclModelType{
						Type: "A",
					},
				},
			},
		},
	}
	tCases := []struct {
		root   KclPackage
		expect string
	}{
		{
			root: rootPkg,
			expect: `- [A](#a)
- sub1
  - [B](#b)
  - sub2
    - [C](#c)
`,
		},
	}
	for _, tCase := range tCases {
		got := tCase.root.getIndexContent(0, "  ")
		assert2.Equal(t, tCase.expect, got)
	}
}

func TestDocGenerate(t *testing.T) {
	tCases := initTestCases(t)
	for _, tCase := range tCases {
		// create target directory
		err := os.MkdirAll(tCase.GotMd, 0755)
		if err != nil {
			t.Fatal(err)
		}
		genOpts := GenOpts{
			Path:             tCase.PackagePath,
			Format:           string(Markdown),
			Target:           tCase.GotMd,
			IgnoreDeprecated: false,
			EscapeHtml:       true,
			TemplateDir:      tCase.TmplPath,
		}
		genContext, err := genOpts.ValidateComplete()
		if err != nil {
			t.Fatal(err)
		}
		err = genContext.GenDoc()
		if err != nil {
			t.Fatalf("generate failed: %s", err)
		}
		// check the content of expected and actual
		err = CompareDir(tCase.ExpectMd, filepath.Join(tCase.GotMd, "docs"))
		if err != nil {
			t.Fatal(err)
		}
		// if test failed, keep generate files for checking
		os.RemoveAll(tCase.GotMd)
	}
}

func TestPopulateReferencedBy(t *testing.T) {
	schemaA := &KclOpenAPIType{
		Type:        "object",
		Description: "Schema A",
		Properties: map[string]*KclOpenAPIType{
			"propA": {Type: "string"},
		},
	}
	schemaB := &KclOpenAPIType{
		Type:        "object",
		Description: "Schema B",
		Properties: map[string]*KclOpenAPIType{
			"propB": {Type: "object"},
		},
	}
	schemaC := &KclOpenAPIType{
		Type:        "object",
		Description: "Schema C",
		Properties: map[string]*KclOpenAPIType{
			"propC": {Type: "object"},
		},
	}

	schemaB.Properties["refA"] = schemaA
	schemaC.Properties["refB"] = schemaB

	pkg := KclPackage{
		SchemaList: []*KclOpenAPIType{schemaA, schemaB, schemaC},
	}

	pkg.populateReferencedBy()

	assert2.ElementsMatch(t, schemaA.ReferencedBy, []string{"Schema B"})
	assert2.ElementsMatch(t, schemaB.ReferencedBy, []string{"Schema C"})
	assert2.Empty(t, schemaC.ReferencedBy)
}

func initTestCases(t *testing.T) []*TestCase {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("get work directory failed")
	}

	testdataDir := filepath.Join("testdata", "doc")
	sourcePkgs := []string{
		"k8s",
		"pkg",
		"reimport",
		"index_func",
	}
	tcases := make([]*TestCase, len(sourcePkgs))

	for i, p := range sourcePkgs {
		packageDir := filepath.Join(cwd, testdataDir, p)
		tcases[i] = &TestCase{
			PackagePath: packageDir,
			ExpectMd:    filepath.Join(packageDir, "md"),
			ExpectHtml:  filepath.Join(packageDir, "html"),
			GotMd:       filepath.Join(packageDir, "md_got"),
			GotHtml:     filepath.Join(packageDir, "html_got"),
		}
	}
	return tcases
}

type TestCase struct {
	PackagePath string
	ExpectMd    string
	ExpectHtml  string
	GotMd       string
	GotHtml     string
	TmplPath    string
}

func CompareDir(a string, b string) error {
	dirA, err := os.ReadDir(a)
	if err != nil {
		return fmt.Errorf("read dir %s failed when comparing with %s", a, b)
	}
	dirB, err := os.ReadDir(b)
	if err != nil {
		return fmt.Errorf("read dir %s failed when comparing with %s", b, a)
	}
	if len(dirA) != len(dirB) {
		return fmt.Errorf("dirs contains different number of files:\n%s: %v\n%s: %v", a, len(dirA), b, len(dirB))
	}
	for _, fA := range dirA {
		// check if the same file exist in dirB
		aPath := filepath.Join(a, fA.Name())
		bPath := filepath.Join(b, fA.Name())
		_, err := os.Open(bPath)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file %s exist in %s, but missing in %s", fA.Name(), a, b)
		}
		if err != nil {
			return fmt.Errorf("open file failed when compare, file path: %s", bPath)
		}
		if fA.IsDir() {
			err := CompareDir(aPath, bPath)
			if err != nil {
				return err
			}
			continue
		}
		linesA, err := readLines(aPath)
		if err != nil {
			return fmt.Errorf("failed to readlines from %s when compare files", aPath)
		}
		linesB, err := readLines(bPath)
		if err != nil {
			return fmt.Errorf("failed to readlines from %s when compare files", bPath)
		}
		for i, line := range linesA {
			if line != linesB[i] {
				lineNo := i + 1
				return fmt.Errorf(
					"file content different: \n%s:%v:%s\n%s:%v:%s",
					aPath, lineNo, line, bPath, lineNo, linesB[i],
				)
			}
		}
		if len(linesA) < len(linesB) {
			return fmt.Errorf("file content different, contains more lines in file %s:%v - %v:\n%s", aPath, len(linesA), len(linesB), strings.Join(linesB[len(linesA):], "\n"))
		}
	}
	return nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func TestDeepEscapePipesInSchema_EnumPipeEscaping(t *testing.T) {
	schema := &KclOpenAPIType{
		Enum: []string{"A|B", "C", "D|E|F"},
	}

	escaped := deepEscapePipesInSchema(schema)

	assert2.Equal(t, []string{"A\\|B", "C", "D\\|E\\|F"}, escaped.Enum)
	assert2.Equal(t, schema.Default, escaped.Default)
	assert2.Equal(t, schema.Description, escaped.Description)
	assert2.Equal(t, schema.ExternalDocs, escaped.ExternalDocs)
	assert2.Equal(t, schema.Ref, escaped.Ref)
}
