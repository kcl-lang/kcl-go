// Copyright 2022 The KCL Authors. All rights reserved.

package ast_test

import (
	"fmt"

	"kusionstack.io/kclvm-go/pkg/ast"
)

func Example() {
	m, err := ast.DecodeModule("../compiler/parser/testdata/a.k.ast.json", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(ast.JSONString(m))
}
