package doc

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDocRender(t *testing.T) {
	tcases := [...]struct {
		source *KclOpenAPIType
		expect string
	}{
		{
			source: &KclOpenAPIType{
				Description: "Description of Schema Person",
				Properties: map[string]*KclOpenAPIType{
					"name": {
						Type:        "string",
						Description: "name of the person",
					},
				},
				Required: []string{"name"},
				KclExtensions: &KclExtensions{
					XKclModelType: &XKclModelType{
						Import: &KclModelImportInfo{
							Package: "models",
							Alias:   "person.k",
						},
						Type: "Person",
					},
				},
			},

			expect: `## Schema Person

Description of Schema Person

### Attributes

**name** *required*

` + "`" + `str` + "`" + `

name of the person


## Source Files

- [models.Person](models.person.k)
`,
		},
	}

	context := GenContext{
		Format:           Markdown,
		IgnoreDeprecated: true,
	}

	for _, tcase := range tcases {
		content, err := context.renderContent(tcase.source)
		if err != nil {
			t.Errorf("render failed, err: %s", err)
		}
		expect := tcase.expect
		if runtime.GOOS == "windows" {
			expect = strings.ReplaceAll(tcase.expect, "\n", "\r\n")
		}
		assert.Equal(t, expect, string(content), "render content mismatch")
	}
}

func TestDocGenerate(t *testing.T) {
	tCases := initTestCases(t)
	for _, tCase := range tCases {
		genContext := GenContext{
			PackagePath:      tCase.PackagePath,
			Format:           Markdown,
			IgnoreDeprecated: false,
			Target:           tCase.GotMd,
		}
		err := genContext.GenDoc()
		if err != nil {
			t.Fatalf("generate failed: %s", err)
		}
		// check the content of expected and actual
		err = CompareDir(tCase.ExpectMd, tCase.GotMd)
		if err != nil {
			t.Fatal(err)
		}
		// if test failed, keep generate files for checking
		os.RemoveAll(genContext.Target)
	}
}

const testdataDir = "testdata"

func initTestCases(t *testing.T) []*TestCase {
	var tcases []*TestCase
	paths, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatal("missing testdata directory")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("get work directory failed")
	}
	for _, path := range paths {
		tcases = append(tcases, &TestCase{
			PackagePath: filepath.Join(testdataDir, path.Name()),
			ExpectMd:    filepath.Join(cwd, testdataDir, path.Name(), "md"),
			ExpectHtml:  filepath.Join(cwd, testdataDir, path.Name(), "html"),
			GotMd:       filepath.Join(cwd, testdataDir, path.Name(), "md_got"),
			GotHtml:     filepath.Join(cwd, testdataDir, path.Name(), "html_got"),
		})
	}
	return tcases
}

type TestCase struct {
	PackagePath string
	ExpectMd    string
	ExpectHtml  string
	GotMd       string
	GotHtml     string
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
			return CompareDir(aPath, bPath)
		}
		linesA, err := readLines(aPath)
		if err != nil {
			return fmt.Errorf("failed to readlins from %s when compare files", aPath)
		}
		linesB, err := readLines(bPath)
		if err != nil {
			return fmt.Errorf("failed to readlins from %s when compare files", bPath)
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
