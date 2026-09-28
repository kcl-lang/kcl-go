// Copyright The KCL Authors. All rights reserved.

//go:build cgo

package kcl_test

import (
	"os"
	"path/filepath"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	kcl "kcl-lang.io/kcl-go"
)

// TestParseFileReexport exercises the top-level ParseFile wrapper, which
// delegates to pkg/parser and returns the AST as a Go structure.
func TestParseFileReexport(t *testing.T) {
	m, err := kcl.ParseFile("test.k", "schema Person:\n    name: str\n")
	assert2.NoError(t, err)
	if assert2.NotNil(t, m) {
		assert2.Equal(t, "test.k", m.Filename)
		assert2.NotEmpty(t, m.Body)
	}

	// Reading from a real file path must work too.
	dir := t.TempDir()
	file := filepath.Join(dir, "person.k")
	assert2.NoError(t, os.WriteFile(file, []byte("schema Person:\n    name: str\n"), 0o644))
	m, err = kcl.ParseFile(file, nil)
	assert2.NoError(t, err)
	if assert2.NotNil(t, m) {
		assert2.NotEmpty(t, m.Filename)
		assert2.NotEmpty(t, m.Body)
	}
}

// TestGetFullSchemaTypeReexports exercises the top-level wrappers around the
// full-schema-type APIs: GetFullSchemaType, GetFullSchemaTypeMapping and
// GetFullSchemaTypeMappingUnderPath.
func TestGetFullSchemaTypeReexports(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "person.k")
	assert2.NoError(t, os.WriteFile(file, []byte("schema Person:\n    name: str\n    age: int = 18\n"), 0o644))

	types, err := kcl.GetFullSchemaType([]string{file}, "")
	assert2.NoError(t, err)
	if assert2.Len(t, types, 1) {
		assert2.Equal(t, "Person", types[0].GetSchemaName())
	}

	mapping, err := kcl.GetFullSchemaTypeMapping([]string{file}, "")
	assert2.NoError(t, err)
	if assert2.Contains(t, mapping, "Person") {
		assert2.Equal(t, "schema", mapping["Person"].GetType())
	}

	underPath, err := kcl.GetFullSchemaTypeMappingUnderPath([]string{file}, "")
	assert2.NoError(t, err)
	if assert2.Contains(t, underPath, "__main__") {
		var names []string
		for _, st := range underPath["__main__"].GetSchemaType() {
			names = append(names, st.GetSchemaName())
		}
		assert2.Contains(t, names, "Person")
	}
}
