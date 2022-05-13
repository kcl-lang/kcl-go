// Copyright 2021 The KCL Authors. All rights reserved.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"kusionstack.io/kclvm-go/pkg/ast"
	"kusionstack.io/kclvm-go/pkg/compiler/parser"
)

var (
	flagKFile = flag.String("kcl-file", "./testdata/a.k", "set input kcl file")
)

func main() {
	flag.Parse()
	if *flagKFile == "" {
		fmt.Println("no kcl file")
		os.Exit(1)
	}

	f, err := parser.ParseFile(*flagKFile, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(f.JSONString())
	ast.Print(f)
}
