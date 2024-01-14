package parser

import (
	"strings"
	"testing"
	"time"
)

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
