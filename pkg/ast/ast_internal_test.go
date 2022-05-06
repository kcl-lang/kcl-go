// Copyright 2022 The KCL Authors. All rights reserved.

package ast

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	for _, name := range GetTypeNameList() {
		x, ok := NewNode(name)
		if !ok {
			t.Fatalf("NewNode(%q) failed", name)
		}

		if x.GetMeta().AstType != x.GetNodeType() {
			t.Fatal("invalid GetNodeType")
		}

		_ = x.JSONString()
		_ = x.JSONMap()
		_, _, _ = x.GetPosition()

		if x, ok := x.(Stmt); ok {
			x.stmt_type()
		}
		if x, ok := x.(Expr); ok {
			x.expr_type()
		}
		if x, ok := x.(TypeInterface); ok {
			x.type_type()
		}
	}
}
