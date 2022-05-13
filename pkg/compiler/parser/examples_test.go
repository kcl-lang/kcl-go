// Copyright 2021 The KCL Authors. All rights reserved.

package parser_test

import (
	"fmt"

	"kusionstack.io/kclvm-go/pkg/compiler/parser"
)

func ExampleParseFile() {
	const hello_k = `
# Copyright 2020 The KCL Authors. All rights reserved.

import some.pkg as pkgName

schema Person:
	name: str = 'kcl'
	age: int = 1

go = Person {
	name: 'golang'
}

if go.name == 'golang':
	print("hello KCL")
`

	f, err := parser.ParseFile("hello.k", hello_k)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(f.Module)
}
