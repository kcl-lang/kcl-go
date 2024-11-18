package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"kcl-lang.io/kcl-go/pkg/ast"
)

func TestParseFile(t *testing.T) {
	// Example: Test with string source
	src := `schema Name:
    name: str

a: int = 1
b = 2Ki
c = 1 + 1
d = "123"
e: 123 = 123
f: "Red" | "Blue" = "Red"
n1 = Name()
n2 = Name {name = "name"}
n3: Name {name = "name"}
schema Container:
    name: str = "main"
    command?: [str]
    ports?: [ContainerPort]

schema Person:
    name?: any

version = "dev"

appConfiguration = xxx.xxxAppConfiguration {
    mainContainer = container.Main {
        readinessProbe = probe_tpl.defaultReadinessProbe
        env : [
            e.Env {
                name: "a"
                value: "b"
            },
        ] + values._envs
    }
}
` // Sample KCL source code
	module, err := ParseFile("", src)
	if err != nil {
		t.Errorf("ParseFile failed with string source: %v and error: %v", src, err)
	}
	if module == nil {
		t.Errorf("Expected non-empty AST JSON with string source")
	} else {
		schemaStmt := module.Body[0].Node.(*ast.SchemaStmt)
		if len(schemaStmt.Body) != 1 {
			t.Errorf("wrong schema stmt body count")
		}
		simpleAssignStmt := module.Body[1].Node.(*ast.AssignStmt)
		if simpleAssignStmt.Value.Node.(*ast.NumberLit).Value.(*ast.IntNumberLitValue).Value != 1 {
			t.Errorf("wrong assign stmt literal value")
		}
		schemaAssignStmt := module.Body[8].Node.(*ast.AssignStmt)
		if len(schemaAssignStmt.Value.Node.(*ast.SchemaExpr).Config.Node.(*ast.ConfigExpr).Items) != 1 {
			t.Errorf("wrong assign stmt schema entry count")
		}
		schemaUnificationStmt := module.Body[9].Node.(*ast.UnificationStmt)
		if len(schemaUnificationStmt.Value.Node.Config.Node.(*ast.ConfigExpr).Items) != 1 {
			t.Errorf("wrong assign stmt schema entry count")
		}
	}
}

func TestParseFileInTheWholeRepo(t *testing.T) {
	root := filepath.Join(".", "..", "..")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".k") {
			testParseFile(t, path)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error walking the path %v: %v", root, err)
	}
}

func testParseFile(t *testing.T, path string) {
	var err error
	var content string

	t.Logf("Start parse file: %s", path)
	module, err := ParseFile(path, content)
	if err != nil {
		t.Errorf("ParseFile failed for %s: %v", path, err)
		return
	}
	if module == nil {
		t.Errorf("Expected non-empty AST JSON for %s", path)
		return
	}
	t.Logf("Successfully parsed file: %s", path)
}

// TestParseFileASTJson tests the ParseFileASTJson function with various input sources.
func TestParseFileASTJson(t *testing.T) {
	// Example: Test with string source
	src := `schema Name:
    name: str
	
n = Name {name = "name"}` // Sample KCL source code
	astJson, err := ParseFileASTJson("", src)
	if err != nil {
		t.Errorf("ParseFileASTJson failed with string source: %v", err)
	}
	if astJson == "" {
		t.Errorf("Expected non-empty AST JSON with string source")
	}

	// Example: Test with byte slice source
	srcBytes := []byte(src)
	astJson, err = ParseFileASTJson("", srcBytes)
	if err != nil {
		t.Errorf("ParseFileASTJson failed with byte slice source: %v", err)
	}
	if astJson == "" {
		t.Errorf("Expected non-empty AST JSON with byte slice source")
	}

	startTime := time.Now()
	// Example: Test with io.Reader source
	srcReader := strings.NewReader(src)
	astJson, err = ParseFileASTJson("", srcReader)
	if err != nil {
		t.Errorf("ParseFileASTJson failed with io.Reader source: %v", err)
	}
	if astJson == "" {
		t.Errorf("Expected non-empty AST JSON with io.Reader source")
	}
	elapsed := time.Since(startTime)
	t.Logf("ParseFileASTJson took %s", elapsed)
	if !strings.Contains(astJson, "Schema") {
		t.Errorf("Expected Schema Node AST JSON with io.Reader source")
	}
	if !strings.Contains(astJson, "Assign") {
		t.Errorf("Expected Assign Node AST JSON with io.Reader source")
	}
}
