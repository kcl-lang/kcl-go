package loader

import (
	"path/filepath"
	"testing"
)

func TestLoadFileOptions(t *testing.T) {
	options, err := ListFileOptions(filepath.Join(".", "test_data", "options.k"))
	if err != nil {
		t.Errorf("ListFileOptions failed with string source: %v", err)
	}
	if len(options) != 3 {
		t.Errorf("ListFileOptions failed with string source: %v", options)
	}
	if options[2].Type != "str" {
		t.Errorf("ListFileOptions failed with string source: %v", options)
	}
}
