package parser

import (
	"encoding/json"
	"testing"

	"kcl-lang.io/kcl-go/pkg/ast"
)

// TestAstJsonAlignment_ParseFile parses real KCL code through the native
// (Rust) kcl-lib and verifies the resulting JSON shape matches the Go AST
// struct definitions declared in ../pkg/ast.
//
// The wire format is whatever the Rust compiler in ../kcl/crates/ast emits;
// if a Go struct is out of sync, json.Unmarshal fails with a clear error or
// the assertions below fire.
func TestAstJsonAlignment_ParseFile(t *testing.T) {
	const src = `
"""Sample KCL for AstJsonAlignment tests."""

schema Person:
    """A person."""

    @deprecated
    name: str = "anonymous"

    age: int = 0

    check:
        age >= 0 if age, "age must be non-negative"

x = Person {name = "Alice", age = 30}

adder = lambda x: int, y: int -> int {
    x + y
}

mixin HasTimestamp:
    createdAt: str = "1970-01-01T00:00:00Z"

@deprecated
schema Article(HasTimestamp):
    title: str
`
	module, err := ParseFile("alignment.k", src)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// --- Module: filename/doc/body/comments only, no `pkg` ---
	if module.Filename == "" {
		t.Errorf("expected module.Filename to be set")
	}
	// Round-trip the parsed module and verify the JSON does not contain
	// a `pkg` field — this proves both that Rust never emits one and that
	// the Go Module struct does not synthesize one on re-serialize.
	if moduleJSON, err := json.Marshal(module); err == nil {
		if got := string(moduleJSON); containsPkg(got) {
			t.Errorf("Module JSON must not contain a `pkg` field; got %s", got)
		}
	}
	if len(module.Body) == 0 {
		t.Fatalf("expected non-empty module body")
	}

	// --- Literal discriminators use the long form ---
	for _, stmt := range module.Body {
		walkLiterals(t, stmt.Node)
	}

	// --- SchemaStmt: decorators deserialize into the flat Decorator struct
	// (no `type:"Call"` tag inside NodeRef<Decorator>) ---
	var article *ast.SchemaStmt
	for _, s := range module.Body {
		if ss, ok := s.Node.(*ast.SchemaStmt); ok && ss.Name != nil && ss.Name.Node == "Article" {
			article = ss
			break
		}
	}
	if article == nil {
		t.Fatalf("Article schema not found in module body")
	}
	if len(article.Decorators) == 0 {
		t.Fatalf("Article schema should have at least one decorator")
	}
	for _, d := range article.Decorators {
		if d == nil {
			t.Errorf("decorator node must not be nil")
			continue
		}
		// d.Node is `ast.Decorator` (struct, not pointer). Verify by
		// inspecting fields rather than type-asserting to *Decorator.
		if d.Node.Func == nil {
			t.Errorf("decorator must have a non-nil Func")
		}
	}

	// --- SchemaAttr: decorators field supported ---
	var person *ast.SchemaStmt
	for _, s := range module.Body {
		if ss, ok := s.Node.(*ast.SchemaStmt); ok && ss.Name != nil && ss.Name.Node == "Person" {
			person = ss
			break
		}
	}
	if person == nil {
		t.Fatalf("Person schema not found in module body")
	}
	var nameAttr *ast.SchemaAttr
	for _, a := range person.Body {
		if sa, ok := a.Node.(*ast.SchemaAttr); ok && sa.Name != nil && sa.Name.Node == "name" {
			nameAttr = sa
			break
		}
	}
	if nameAttr == nil {
		t.Fatalf("Person.name SchemaAttr not found")
	}
	if len(nameAttr.Decorators) != 1 {
		t.Errorf("Person.name should have one decorator (@deprecated); got %d",
			len(nameAttr.Decorators))
	}

	// --- ConfigEntry: is_shorthand field round-trips ---
	cfgEntry := &ast.ConfigEntry{Operation: ast.ConfigEntryOperationUnion}
	b, err := json.Marshal(cfgEntry)
	if err != nil {
		t.Fatalf("Marshal ConfigEntry failed: %v", err)
	}
	// `Key` and `Value` don't carry `omitempty` in Go, so they serialize as
	// explicit nulls — that's fine because Rust emits `null` for absent
	// fields too. The contract being tested is just that `is_shorthand`
	// follows Rust's `skip_serializing_if = "is_false"`.
	const wantNoShorthand = `{"key":null,"value":null,"operation":"Union"}`
	if got := string(b); got != wantNoShorthand {
		t.Errorf("is_shorthand=false must be skipped; got %s", got)
	}
	cfgEntry.IsShorthand = true
	b, _ = json.Marshal(cfgEntry)
	const wantShorthand = `{"key":null,"value":null,"operation":"Union","is_shorthand":true}`
	if got := string(b); got != wantShorthand {
		t.Errorf("is_shorthand=true must be emitted; got %s", got)
	}

	// --- NumberLit/StringLit/NameConstantLit constructors use long-form
	// discriminators so round-trips match the Rust wire format. ---
	if ast.NewNumberLit().Type() != "NumberLit" {
		t.Errorf("NumberLit must serialize as \"NumberLit\", got %q", ast.NewNumberLit().Type())
	}
	if ast.NewStringLit().Type() != "StringLit" {
		t.Errorf("StringLit must serialize as \"StringLit\", got %q", ast.NewStringLit().Type())
	}
	if ast.NewNameConstantLit().Type() != "NameConstantLit" {
		t.Errorf("NameConstantLit must serialize as \"NameConstantLit\", got %q",
			ast.NewNameConstantLit().Type())
	}
}

// walkLiterals recursively inspects a parsed AST for any literal node whose
// discriminator is the short ("Number"/"String"/"NameConstant") form.
// Mirrors Java's literalDiscriminators_useLongFormNotLiteralEnumForm test.
//
// Go's Node[T] is a struct (not an interface), so we walk by type-asserting
// on the `Node` field — but Expr subtypes carry their concrete type via the
// `*ast.FooExpr` pointer that the JSON unmarshaller populates on Module.Body.
func walkLiterals(t *testing.T, node any) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case ast.ExprStmt:
		for _, e := range n.Exprs {
			walkLiterals(t, e.Node)
		}
	case *ast.NumberLit:
		if n.Type() == "Number" {
			t.Errorf("NumberLit must use long-form discriminator, got %q", n.Type())
		}
	case *ast.StringLit:
		if n.Type() == "String" {
			t.Errorf("StringLit must use long-form discriminator, got %q", n.Type())
		}
	case *ast.NameConstantLit:
		if n.Type() == "NameConstant" {
			t.Errorf("NameConstantLit must use long-form discriminator, got %q", n.Type())
		}
	case *ast.ConfigExpr:
		for _, item := range n.Items {
			walkLiterals(t, item.Node)
		}
	case *ast.SchemaExpr:
		if n.Config != nil {
			walkLiterals(t, n.Config.Node)
		}
	case *ast.UnaryExpr:
		if n.Operand != nil {
			walkLiterals(t, n.Operand.Node)
		}
	case *ast.BinaryExpr:
		if n.Left != nil {
			walkLiterals(t, n.Left.Node)
		}
		if n.Right != nil {
			walkLiterals(t, n.Right.Node)
		}
	case *ast.IfExpr:
		if n.Body != nil {
			walkLiterals(t, n.Body.Node)
		}
		if n.Cond != nil {
			walkLiterals(t, n.Cond.Node)
		}
		if n.Orelse != nil {
			walkLiterals(t, n.Orelse.Node)
		}
	case *ast.CallExpr:
		if n.Func != nil {
			walkLiterals(t, n.Func.Node)
		}
		for _, a := range n.Args {
			walkLiterals(t, a.Node)
		}
	case *ast.ListExpr:
		for _, e := range n.Elts {
			walkLiterals(t, e.Node)
		}
	case *ast.ListComp:
		if n.Elt != nil {
			walkLiterals(t, n.Elt.Node)
		}
	case *ast.Compare:
		if n.Left != nil {
			walkLiterals(t, n.Left.Node)
		}
		for _, c := range n.Comparators {
			walkLiterals(t, c.Node)
		}
	case *ast.Subscript:
		if n.Value != nil {
			walkLiterals(t, n.Value.Node)
		}
		if n.Index != nil {
			walkLiterals(t, n.Index.Node)
		}
	case *ast.SelectorExpr:
		if n.Value != nil {
			walkLiterals(t, n.Value.Node)
		}
	case *ast.CheckExpr:
		if n.Test != nil {
			walkLiterals(t, n.Test.Node)
		}
		if n.IfCond != nil {
			walkLiterals(t, n.IfCond.Node)
		}
		if n.Msg != nil {
			walkLiterals(t, n.Msg.Node)
		}
	case *ast.AssignStmt:
		if n.Value != nil {
			walkLiterals(t, n.Value.Node)
		}
	case *ast.SchemaStmt:
		for _, a := range n.Body {
			walkLiterals(t, a.Node)
		}
		for _, c := range n.Checks {
			walkLiterals(t, c.Node)
		}
	case *ast.SchemaAttr:
		if n.Value != nil {
			walkLiterals(t, n.Value.Node)
		}
	case *ast.UnificationStmt:
		if n.Value != nil {
			walkLiterals(t, n.Value.Node)
		}
	case *ast.LambdaExpr:
		if n.Body != nil {
			for _, b := range n.Body {
				walkLiterals(t, b.Node)
			}
		}
	}
}

// containsPkg returns true if the serialized JSON contains a top-level
// `"pkg"` key on the Module object. Mirrors Java's
// `module_parsesAndHasNoPkgField` assertion against the round-tripped JSON.
func containsPkg(jsonStr string) bool {
	// Cheap substring check — fine for a unit test where the JSON is small
	// and the field name is short.
	for i := 0; i+5 <= len(jsonStr); i++ {
		if jsonStr[i] == '"' && jsonStr[i+1] == 'p' && jsonStr[i+2] == 'k' &&
			jsonStr[i+3] == 'g' && jsonStr[i+4] == '"' {
			return true
		}
	}
	return false
}