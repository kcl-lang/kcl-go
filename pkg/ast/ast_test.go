// Copyright 2022 The KCL Authors. All rights reserved.

package ast_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"kusionstack.io/kclvm-go/pkg/ast"
)

func init() {
	ast.DebugMode = true
}

func TestBuildAST(t *testing.T) {
	filename := "../compiler/parser/testdata/a.k.ast.json"

	m, err := ast.DecodeModule(filename, nil)
	if err != nil {
		t.Fatal(err)
	}

	tAssertEQ(t, m.AstType, ast.Module_TypeName)
	tAssertEQ(t, m.Line, 1)
	tAssertEQ(t, m.Column, 1)
	tAssertEQ(t, m.EndLine, 24)
	tAssertEQ(t, m.EndColumn, 23)
	tAssertEQ(t, m.Filename, "testdata/a.k")
	tAssertEQ(t, m.Pkg, "")

	tAssert(t, len(m.Body) > 0)
	tAssert(t, m.Body[0].(*ast.ImportStmt).AstType == ast.ImportStmt_TypeName)

	want, err := ast.LoadJson(filename, nil)
	tAssert(t, err == nil)

	got, err := ast.LoadJson(filename, ast.JSONString(m))
	tAssert(t, err == nil)

	if diff := cmp.Diff(want, got, tCmpEquateEmpty()); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func tCmpEquateEmpty() cmp.Option {
	isEmpty := func(x interface{}) bool {
		if x == nil {
			return true
		}
		vx := reflect.ValueOf(x)
		switch {
		case x == nil:
			return true
		default:
			switch vx.Kind() {
			case reflect.Slice, reflect.Map:
				return vx.Len() == 0
			case reflect.String:
				return vx.Len() == 0
			}
		}
		return false
	}

	return cmp.FilterValues(
		func(x, y interface{}) bool {
			return isEmpty(x) && isEmpty(y)
		},
		cmp.Comparer(func(_, _ interface{}) bool {
			return true
		}),
	)
}

func tAssert(tb testing.TB, condition bool, args ...interface{}) {
	if !condition {
		tb.Helper()
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("Assert failed, %s", msg)
		} else {
			tb.Fatalf("Assert failed")
		}
	}
}

func tAssertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	if !condition {
		tb.Helper()
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatalf("tAssert failed, %s", msg)
		} else {
			tb.Fatalf("tAssert failed")
		}
	}
}

func tAssertEQ(t *testing.T, x, y interface{}) {
	if !reflect.DeepEqual(x, y) {
		t.Helper()
		t.Fatalf("not equal:\n  x = %v\n  y = %v\n", x, y)
	}
}
