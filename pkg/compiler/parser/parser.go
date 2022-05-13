// Copyright 2021 The KCL Authors. All rights reserved.

// Package parser implements a parser for KCL source files.
package parser

import (
	"fmt"
	"io"
	"os"

	"kusionstack.io/kclvm-go/pkg/ast"
	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func ParseExpr(x string, opts ...Option) (f ast.Expr, err error) {
	panic("TODO")
}

func ParseFile(filename string, src interface{}, opts ...Option) (f *ast.File, err error) {
	var all_opts options
	for _, opt := range opts {
		opt.apply(&all_opts)
	}
	_ = all_opts

	if src == nil {
		src, err = os.ReadFile(filename)
		if err != nil {
			return
		}
	}

	var code string
	switch src := src.(type) {
	case []byte:
		code = string(src)
	case string:
		code = string(src)
	case io.Reader:
		d, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		code = string(d)
	default:
		return nil, fmt.Errorf("unsupported src type: %T", src)
	}
	if code == "" {
		return new(ast.File), nil
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.ParseFile_AST(&gpyrpc.ParseFile_AST_Args{
		Filename:   filename,
		SourceCode: code,
	})
	if err != nil {
		return nil, err
	}

	f = &ast.File{JSON: resp.AstJson}

	m, err := ast.DecodeModule(filename, resp.AstJson)
	if err != nil {
		return nil, err
	}

	f.Module = m
	return f, nil
}
