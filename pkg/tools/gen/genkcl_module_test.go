package gen

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
	"kcl-lang.io/kcl-go/pkg/kcl"
)

// TestGenKclModule_GoldenTree exercises the layout and contents of the
// multi-file KCL module GenKclModule writes. The fixture's `expect/` tree
// is the golden; the test walks both trees and asserts identical
// relative-path sets and byte-identical contents (normalising `\r\n`).
//
// This is the first `t.TempDir()` usage in pkg/, but the helper is
// available everywhere in the standard library and makes the test
// trivially parallelisable.
func TestGenKclModule_GoldenTree(t *testing.T) {
	outDir := t.TempDir()
	pkgDir, err := filepath.Abs("./testdata/gostruct_module/app")
	if err != nil {
		t.Fatal(err)
	}
	if err := GenKclModule(outDir, pkgDir, nil); err != nil {
		t.Fatal(err)
	}
	goldenDir := "./testdata/gostruct_module/expect"

	wantFiles := make(map[string]string)
	if err := filepath.WalkDir(goldenDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(goldenDir, path)
		if err != nil {
			return err
		}
		wantFiles[rel] = readFileString(t, path)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	gotFiles := make(map[string]string)
	if err := filepath.WalkDir(outDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(outDir, path)
		if err != nil {
			return err
		}
		gotFiles[rel] = readFileString(t, path)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	assert2.Equal(t, len(wantFiles), len(gotFiles),
		"golden tree has %d files, generated has %d", len(wantFiles), len(gotFiles))

	for rel, want := range wantFiles {
		got, ok := gotFiles[rel]
		if !ok {
			t.Errorf("missing generated file %q (relative to %s)", rel, outDir)
			continue
		}
		// Normalise \r\n the same way readFileString does on the
		// golden side so Windows-checkout artefacts don't drift.
		got = strings.ReplaceAll(got, "\r\n", "\n")
		assert2.Equal(t, want, got, "contents mismatch for %s", rel)
	}
	for rel := range gotFiles {
		if _, ok := wantFiles[rel]; !ok {
			t.Errorf("unexpected generated file %q", rel)
		}
	}
}

// TestGenKclModule_Resolves is the semantic assertion #580 asks for:
// after generating a module into a temporary directory, the resulting
// tree must be parseable as a KCL module, the cross-package `import`
// must resolve to a schema in a sibling directory, and the out-of-
// module field must be reported as `any` (no spurious import or file).
func TestGenKclModule_Resolves(t *testing.T) {
	outDir := t.TempDir()
	pkgDir, err := filepath.Abs("./testdata/gostruct_module/app")
	if err != nil {
		t.Fatal(err)
	}
	if err := GenKclModule(outDir, pkgDir, nil); err != nil {
		t.Fatal(err)
	}

	// Point the resolver at the leaf app directory. KCL walks up to
	// the kcl.mod at outDir for the package root, so the dotted
	// import resolves against the generated tree.
	leaf := filepath.Join(outDir, "kcl_lang_io/kcl_go/pkg/tools/gen/testdata/gostruct_module/app")
	mapping, err := kcl.GetFullSchemaTypeMapping([]string{leaf}, "App", *kcl.NewOption())
	if err != nil {
		t.Fatalf("GetFullSchemaTypeMapping failed: %v\n%s", err, fileTree(t, outDir))
	}
	appType, ok := mapping["App"]
	if !assert2.True(t, ok, "expected App schema to be present in mapping") {
		return
	}

	// The cross-package Inner reference must resolve to a schema with
	// SchemaName "Inner"; otherwise the dotted import didn't find the
	// sibling types.k directory.
	inner, ok := appType.Properties["inner"]
	if !assert2.True(t, ok, "expected `inner` property on App") {
		return
	}
	assert2.Equal(t, "schema", inner.Type, "expected `inner` to be a schema reference, got type %q", inner.Type)
	assert2.Equal(t, "Inner", inner.SchemaName,
		"expected Inner reference to resolve through the dotted import; got SchemaName %q", inner.SchemaName)

	// time.Time is short-circuited to str regardless of module root.
	when, ok := appType.Properties["when"]
	if assert2.True(t, ok, "expected `when` property on App") {
		assert2.Equal(t, "str", when.Type)
	}

	// The out-of-module stdlib reference must be `any`.
	meta, ok := appType.Properties["meta"]
	if assert2.True(t, ok, "expected `meta` property on App") {
		assert2.Equal(t, "any", meta.Type,
			"out-of-module field must degrade to any")
	}

	// Resolving from the sibling types/ directory must yield Inner;
	// otherwise the file-form leaf inside the directory didn't match
	// the dotted import path.
	typesLeaf := filepath.Join(outDir, "kcl_lang_io/kcl_go/pkg/tools/gen/testdata/gostruct_module/app/types")
	typesMapping, err := kcl.GetFullSchemaTypeMapping([]string{typesLeaf}, "", *kcl.NewOption())
	if err != nil {
		t.Fatalf("GetFullSchemaTypeMapping from types/ failed: %v", err)
	}
	_, hasInner := typesMapping["Inner"]
	assert2.True(t, hasInner, "expected Inner schema to be discoverable from types/ directory")
}

// TestGenKclModule_ExternalDegradesToAny guards the explicit promise
// from #580: a Go type from outside the caller's module must NOT
// produce a `.k` file or an `import` line; the field type collapses to
// `any` instead. We assert on the text of the generated files rather
// than the resolver because the resolver hides unresolved imports
// behind the module's own error reporting.
func TestGenKclModule_ExternalDegradesToAny(t *testing.T) {
	outDir := t.TempDir()
	pkgDir, err := filepath.Abs("./testdata/gostruct_module/app")
	if err != nil {
		t.Fatal(err)
	}
	if err := GenKclModule(outDir, pkgDir, nil); err != nil {
		t.Fatal(err)
	}
	appK, err := os.ReadFile(filepath.Join(outDir,
		"kcl_lang_io/kcl_go/pkg/tools/gen/testdata/gostruct_module/app/app.k"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(bytes.ReplaceAll(appK, []byte("\r\n"), []byte("\n")))

	// The out-of-module stdlib package path must NOT produce an
	// `import` line. (The Go type name appears in the docstring,
	// which is generated from the source comment and is allowed.)
	assert2.NotContains(t, content, "import encoding",
		"out-of-module Go package path must not produce an import line")

	// The `meta` field must be `any`.
	assert2.Contains(t, content, "meta?: any",
		"out-of-module field must be degraded to any")

	// The generated file tree must not contain any file/import for
	// the out-of-module type.
	for _, file := range []string{"raw_message.k", "encoding_json.k"} {
		_, statErr := os.Stat(filepath.Join(outDir,
			"kcl_lang_io/kcl_go/pkg/tools/gen/testdata/gostruct_module/app", file))
		assert2.True(t, os.IsNotExist(statErr),
			"out-of-module file %q must not be emitted", file)
	}
}

// TestGenKclModule_NotAGoModule covers the error path: when pkgDir is
// outside any Go module (or otherwise yields a nil Module), GenKclModule
// must return a clear error rather than silently produce a broken tree.
func TestGenKclModule_NotAGoModule(t *testing.T) {
	outDir := t.TempDir()
	pkgDir, err := filepath.Abs("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	err = GenKclModule(outDir, pkgDir, nil)
	assert2.Error(t, err, "expected an error when pkgDir is not inside a Go module")
}

// fileTree is a tiny debug helper that dumps the directory layout as
// text. Used by the resolve test on failure so a regression that
// produces a wrong file set is easier to diagnose.
func fileTree(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		b.WriteString(rel)
		b.WriteByte('\n')
		return nil
	})
	return b.String()
}
