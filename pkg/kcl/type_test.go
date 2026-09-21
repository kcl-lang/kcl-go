// Copyright The KCL Authors. All rights reserved.

//go:build !rpc && cgo
// +build !rpc,cgo

package kcl

import (
	"path/filepath"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// fixtureDir resolves the absolute path to the cross-package fixture
// mirroring crates/api/src/testdata/get_schema_ty_under_path/ on the Rust
// side — `aaa` declares bbb (path) and helloworld (oci) as kcl.mod deps.
func fixtureDir(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("testdata/get_schema_ty_under_path/aaa")
	if err != nil {
		t.Fatalf("resolve fixture dir: %v", err)
	}
	return abs
}

// fixtureSiblingDir resolves a sibling package dir under the same fixture
// root (e.g. "bbb" -> ".../bbb", "helloworld_0.0.1" -> ".../helloworld_0.0.1").
func fixtureSiblingDir(t *testing.T, name string) string {
	t.Helper()
	abs, err := filepath.Abs("testdata/get_schema_ty_under_path/" + name)
	if err != nil {
		t.Fatalf("resolve %s dir: %v", name, err)
	}
	return abs
}

// schemaByName turns a *gpyrpc.SchemaTypes (a slice of *gpyrpc.KclType) into
// a name -> *gpyrpc.KclType map for ergonomic assertions.
func schemaByName(st *gpyrpc.SchemaTypes) map[string]*gpyrpc.KclType {
	out := make(map[string]*gpyrpc.KclType, len(st.GetSchemaType()))
	for _, kt := range st.GetSchemaType() {
		if kt == nil {
			continue
		}
		out[kt.GetSchemaName()] = kt
	}
	return out
}

// pkgKeys returns the package names present in a mapping for nicer
// assertion failure messages.
func pkgKeys(m map[string]*gpyrpc.SchemaTypes) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// TestGetFullSchemaTypeMappingUnderPath is the primary regression test for
// https://github.com/kcl-lang/kcl/issues/1546. The previous single-pkg
// GetFullSchemaTypeMapping could only see `aaa`'s own schemas and dropped
// PkgPath / BaseSchema for anything defined in kcl.mod dependencies.
func TestGetFullSchemaTypeMappingUnderPath(t *testing.T) {
	root := fixtureDir(t)

	// Mirror the Rust doctest in crates/api/src/service/service_impl.rs:
	// pass external packages explicitly so the engine resolves them
	// without needing network access to the OCI registry.
	mapping, err := GetFullSchemaTypeMappingUnderPath(
		[]string{root},
		"",
		WithExternalPkgNameAndPath("bbb", fixtureSiblingDir(t, "bbb")),
		WithExternalPkgNameAndPath("helloworld", fixtureSiblingDir(t, "helloworld_0.0.1")),
	)
	if err != nil {
		t.Fatalf("GetFullSchemaTypeMappingUnderPath: %v", err)
	}

	// Expect at least __main__ (aaa), bbb, helloworld — every package the
	// fixture pulls in. Anything missing means the engine didn't resolve
	// the external dependency.
	assert2.Contains(t, mapping, "__main__", "expected __main__ in mapping; got %v", pkgKeys(mapping))
	assert2.Contains(t, mapping, "bbb", "expected bbb in mapping; got %v", pkgKeys(mapping))
	assert2.Contains(t, mapping, "helloworld", "expected helloworld in mapping; got %v", pkgKeys(mapping))

	// aaa's schema A should live in __main__.
	aaaSchemas := schemaByName(mapping["__main__"])
	assert2.NotNil(t, aaaSchemas["A"], "expected schema A in __main__")
	assert2.Equal(t, "A", aaaSchemas["A"].GetSchemaName())
	assert2.NotEmpty(t, aaaSchemas["A"].GetFilename(), "schema A should have a Filename")

	// bbb defines Base + B. PkgPath must point back into bbb (the
	// dependency package), NOT into aaa's __main__. This is the exact bug
	// the issue describes: previously B's PkgPath came back as "__main__".
	bbbSchemas := schemaByName(mapping["bbb"])
	assert2.NotNil(t, bbbSchemas["Base"], "expected schema Base in bbb")
	assert2.Equal(t, "Base", bbbSchemas["Base"].GetSchemaName())
	assert2.Equal(t, "bbb", bbbSchemas["Base"].GetPkgPath(),
		"Base PkgPath should be 'bbb', not '__main__' (regression for #1546)")

	assert2.NotNil(t, bbbSchemas["B"], "expected schema B in bbb")
	assert2.Equal(t, "B", bbbSchemas["B"].GetSchemaName())
	assert2.Equal(t, "bbb", bbbSchemas["B"].GetPkgPath(),
		"B PkgPath should be 'bbb', not '__main__' (regression for #1546)")
	// B inherits Base — BaseSchema must be non-nil and point back at Base.
	assert2.NotNil(t, bbbSchemas["B"].GetBaseSchema(),
		"B.BaseSchema should be resolved across the kcl.mod boundary (regression for #1546)")
	assert2.Equal(t, "Base", bbbSchemas["B"].GetBaseSchema().GetSchemaName())
	assert2.Equal(t, "bbb", bbbSchemas["B"].GetBaseSchema().GetPkgPath(),
		"B.BaseSchema.PkgPath should be 'bbb', not '__main__' (regression for #1546)")

	// helloworld defines Hello, also a kcl.mod dep of aaa.
	helloSchemas := schemaByName(mapping["helloworld"])
	assert2.NotNil(t, helloSchemas["Hello"], "expected schema Hello in helloworld")
	assert2.Equal(t, "Hello", helloSchemas["Hello"].GetSchemaName())
	assert2.Equal(t, "helloworld", helloSchemas["Hello"].GetPkgPath(),
		"Hello PkgPath should be 'helloworld', not '__main__' (regression for #1546)")
}

// TestGetFullSchemaTypeMappingUnderPath_ExternalPkgs exercises the
// WithExternalPkgAndPath variant — mirrors the Rust doctest in
// crates/api/src/service/service_impl.rs.
func TestGetFullSchemaTypeMappingUnderPath_ExternalPkgs(t *testing.T) {
	bbbDir, err := filepath.Abs("testdata/get_schema_ty_under_path/bbb")
	if err != nil {
		t.Fatalf("resolve bbb dir: %v", err)
	}

	mapping, err := GetFullSchemaTypeMappingUnderPath(
		[]string{bbbDir},
		"",
		WithExternalPkgNameAndPath("bbb", bbbDir),
	)
	if err != nil {
		t.Fatalf("GetFullSchemaTypeMappingUnderPath(WithExternalPkgAndPath): %v", err)
	}

	assert2.Contains(t, mapping, "__main__", "expected __main__ in mapping; got %v", pkgKeys(mapping))
	bbbSchemas := schemaByName(mapping["__main__"])
	assert2.NotNil(t, bbbSchemas["Base"], "expected schema Base in __main__ when using external pkg")
	assert2.NotNil(t, bbbSchemas["B"], "expected schema B in __main__ when using external pkg")
	assert2.Equal(t, "Base", bbbSchemas["B"].GetBaseSchema().GetSchemaName(),
		"B.BaseSchema should resolve across the WithExternalPkgNameAndPath boundary")
}

// TestGetFullSchemaTypeMappingUnderPath_SchemaNameFilter verifies that the
// optional schemaName filter still works when targeting cross-package
// schemas — i.e. asking for "B" only should return B, not Base or Hello.
func TestGetFullSchemaTypeMappingUnderPath_SchemaNameFilter(t *testing.T) {
	root := fixtureDir(t)

	mapping, err := GetFullSchemaTypeMappingUnderPath(
		[]string{root},
		"B",
		WithExternalPkgNameAndPath("bbb", fixtureSiblingDir(t, "bbb")),
		WithExternalPkgNameAndPath("helloworld", fixtureSiblingDir(t, "helloworld_0.0.1")),
	)
	if err != nil {
		t.Fatalf("GetFullSchemaTypeMappingUnderPath(schemaName=B): %v", err)
	}

	totalSchemas := 0
	for _, pkg := range mapping {
		if pkg != nil {
			totalSchemas += len(pkg.GetSchemaType())
		}
	}
	assert2.Equal(t, 1, totalSchemas,
		"schemaName=B should return only B, got %d schemas across packages", totalSchemas)

	bbbSchemas := schemaByName(mapping["bbb"])
	assert2.NotNil(t, bbbSchemas["B"], "expected B in bbb")
	assert2.Nil(t, bbbSchemas["Base"], "Base should be filtered out when schemaName=B")
}