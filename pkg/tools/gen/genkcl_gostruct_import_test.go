package gen

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
	"kcl-lang.io/kcl-go/pkg/ast"
	"kcl-lang.io/kcl-go/pkg/parser"
)

// TestGoPkgToKclImport covers the path-shape contract of goPkgToKclImport:
// every segment of the dotted output must be a valid KCL identifier
// (matching crates/parser/src/parser/stmt.rs:436 → parse_identifier).
//
// This is the regression assertion for the bug that #575 shipped: the
// previous implementation only split on '/' and emitted identifiers like
// `kcl-lang` containing hyphens, which are invalid in KCL and caused the
// cross_pkg golden test to silently pass through the parser's
// error-recovery path.
func TestGoPkgToKclImport(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain",
			in:   "github.com/foo/bar",
			want: "github_com.foo.bar",
		},
		{
			name: "hyphen",
			in:   "kcl-lang.io/kcl-go/pkg/tools/gen",
			want: "kcl_lang_io.kcl_go.pkg.tools.gen",
		},
		{
			name: "dot_in_segment",
			in:   "example.com/foo.bar/baz",
			want: "example_com.foo_bar.baz",
		},
		{
			name: "digit_leading_segment",
			in:   "github.com/99designs/gqlgen",
			want: "github_com._99designs.gqlgen",
		},
		{
			name: "keyword_segment",
			in:   "github.com/foo/schema/bar",
			want: "github_com.foo._schema.bar",
		},
		{
			name: "empty_segment",
			in:   "github.com//foo",
			want: "github_com._.foo",
		},
		{
			name: "trailing_slash",
			in:   "github.com/foo/",
			want: "github_com.foo._",
		},
		{
			name: "single_segment",
			in:   "kcl",
			want: "kcl",
		},
		{
			name: "underscore_kept",
			in:   "github.com/_99designs/_init",
			want: "github_com._99designs._init",
		},
	}
	segRegexp := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := goPkgToKclImport(tc.in)
			assert2.Equal(t, tc.want, got)
			for _, seg := range strings.Split(got, ".") {
				if seg == "" {
					continue
				}
				assert2.Truef(t, segRegexp.MatchString(seg),
					"segment %q in %q is not a valid KCL identifier", seg, got)
				assert2.Falsef(t, isKclKeyword(seg),
					"segment %q in %q collides with a KCL keyword", seg, got)
			}
		})
	}
}

// TestAliasForGoPkg covers the post-sanitisation alias-derivation logic
// added alongside the PR 1 sanitiser fix. We feed both slashed (raw Go)
// and dotted (post-goPkgToKclImport) inputs because both are accepted by
// the public contract.
func TestAliasForGoPkg(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain_last_segment",
			in:   "github.com/foo/bar",
			want: "Bar",
		},
		{
			name: "hyphenated_last_segment",
			in:   "kcl-lang.io/kcl-go/pkg/tools/gen",
			want: "Gen",
		},
		{
			name: "dotted_input",
			in:   "kcl_lang_io.kcl_go.pkg.tools.gen",
			want: "Gen",
		},
		{
			name: "digit_leading_segment",
			in:   "github.com/99designs/gqlgen",
			want: "Gqlgen",
		},
		{
			name: "keyword_segment",
			in:   "github.com/foo/schema",
			want: "Schema",
		},
		{
			name: "hyphen_only_segment",
			in:   "github.com/foo/-",
			want: "Imported__",
		},
		{
			name: "empty_path",
			in:   "",
			want: "Imported",
		},
		{
			name: "single_segment_keyword",
			in:   "schema",
			want: "Schema",
		},
		{
			name: "keyword_with_capitalisation_post_camel",
			in:   "true",
			want: "Imported_True",
		},
		{
			name: "digit_leading_post_camel",
			in:   "_99designs",
			want: "Imported_99Designs",
		},
	}
	segRegexp := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := aliasForGoPkg(tc.in)
			assert2.Equal(t, tc.want, got)
			assert2.Truef(t, segRegexp.MatchString(got),
				"alias %q is not a valid KCL identifier", got)
			assert2.Falsef(t, isKclKeyword(got),
				"alias %q collides with a KCL keyword", got)
		})
	}
}

// TestGenKclFromGoStructCrossPkgCollision guards PR 2: when two distinct
// Go package paths sanitise to the same preferred KCL import alias (both
// `internal/a/types` and `internal/b/types` want `Types`), the generator
// must rename the second to `Types2` so neither reference is silently
// shadowed. We feed the generated output to pkg/parser.ParseFile and
// assert two distinct `ImportStmt`s with non-colliding `asname`s, then
// verify the alias-qualified field types use the same names.
func TestGenKclFromGoStructCrossPkgCollision(t *testing.T) {
	oneFile := false
	var buf bytes.Buffer
	err := GenKcl(&buf, "./testdata/gostruct/cross_pkg_collision/main.go", nil, &GenKclOptions{
		Mode:    ModeGoStruct,
		OneFile: &oneFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	got := string(bytes.ReplaceAll(buf.Bytes(), []byte("\r\n"), []byte("\n")))
	expect := readFileString(t, "./testdata/gostruct/cross_pkg_collision/expect.k")
	assert2.Equal(t, expect, got)

	mod, err := parser.ParseFile("cross_pkg_collision.k", buf.Bytes())
	if err != nil {
		t.Fatalf("generated KCL failed to parse: %v\n%s", err, got)
	}

	var imports []*ast.ImportStmt
	for _, node := range mod.Body {
		if node == nil {
			continue
		}
		if stmt, ok := node.Node.(*ast.ImportStmt); ok {
			imports = append(imports, stmt)
		}
	}
	assert2.Equal(t, 2, len(imports), "expected exactly two import statements; got %d:\n%s", len(imports), got)

	seen := make(map[string]string)
	for _, imp := range imports {
		if imp.Asname == nil {
			t.Fatalf("expected every import to carry an `as` alias; rawpath=%q", imp.Rawpath)
		}
		name := imp.Asname.Node
		if prev, dup := seen[name]; dup {
			t.Fatalf("import alias %q is used by both rawpath %q and %q", name, prev, imp.Rawpath)
		}
		seen[name] = imp.Rawpath
	}
	assert2.Equal(t, "Types", imports[0].Asname.Node,
		"the first import (a/types) should keep the preferred alias")
	assert2.Equal(t, "Types2", imports[1].Asname.Node,
		"the second import (b/types) should be suffixed with 2 to disambiguate")

	// The alias-qualified type names in the schema must match the alias
	// each `import` line was assigned; otherwise a regression in either
	// direction (imports loop vs. typeName external branch) would slip
	// past the parser-only check above.
	assert2.Contains(t, got, "Types.Tagged")
	assert2.Contains(t, got, "Types2.Labeled")
}

// TestGenKclFromGoStructCrossPkgIsParseable is the semantic assertion the
// PR 1 fix was missing in the upstream golden test: the generated KCL
// must parse as exactly one top-level ImportStmt whose rawpath matches
// the sanitised Go path, and zero top-level Expr statements. If a Go
// module path contained characters that are illegal in KCL identifiers
// (hyphens, dots, leading digits, …) the parser would silently split
// the malformed `import` into an ImportStmt with rawpath `"kcl"` plus a
// junk ExprStmt for the remainder, and TestParseFileInTheWholeRepo would
// still report success because the parser error-recovers.
//
// This test feeds the generated output to pkg/parser.ParseFile and walks
// the AST, so it is the assertion that would have caught the original
// PR #575 bug.
func TestGenKclFromGoStructCrossPkgIsParseable(t *testing.T) {
	oneFile := false
	var buf bytes.Buffer
	err := GenKcl(&buf, "./testdata/gostruct/cross_pkg/main.go", nil, &GenKclOptions{
		Mode:    ModeGoStruct,
		OneFile: &oneFile,
	})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := parser.ParseFile("cross_pkg.k", buf.Bytes())
	if err != nil {
		t.Fatalf("generated KCL failed to parse: %v\n%s", err, buf.String())
	}

	var importCount, exprCount int
	var foundImport *ast.ImportStmt
	for _, node := range mod.Body {
		if node == nil {
			continue
		}
		switch stmt := node.Node.(type) {
		case *ast.ImportStmt:
			importCount++
			foundImport = stmt
		case *ast.ExprStmt:
			exprCount++
		}
	}
	assert2.Equal(t, 1, importCount,
		"expected exactly one top-level ImportStmt, got %d:\n%s",
		importCount, buf.String())
	assert2.Equal(t, 0, exprCount,
		"expected zero top-level ExprStmt nodes; a parser-recovered junk expression indicates the import path was malformed:\n%s",
		buf.String())
	if assert2.NotNil(t, foundImport, "expected one ImportStmt to inspect") {
		assert2.Equal(t,
			"kcl_lang_io.kcl_go.pkg.tools.gen.testdata.gostruct.cross_pkg.types",
			foundImport.Rawpath,
			"the import rawpath must match the sanitised Go path; got %q",
			foundImport.Rawpath)
		assert2.NotNil(t, foundImport.Asname,
			"expected the import to carry an `as` alias")
		if foundImport.Asname != nil {
			assert2.Equal(t, "Types", foundImport.Asname.Node,
				"the alias should derive from the last path component")
		}
	}
}
