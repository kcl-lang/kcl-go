package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"kcl-lang.io/kcl-go/pkg/ast"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type ParseProgramArgs = gpyrpc.ParseProgramArgs
type ParseProgramResult = gpyrpc.ParseProgramResult

// ParseFileASTJson parses the source code from the specified file or Reader
// and returns the JSON representation of the Abstract Syntax Tree (AST).
// The source code can be provided directly as a string or []byte,
// or indirectly via a filename or an io.Reader.
// If src is nil, the function reads the content from the provided filename.
func ParseFileASTJson(filename string, src any) (result string, err error) {
	var code string
	if src != nil {
		switch src := src.(type) {
		case []byte:
			code = string(src)
		case string:
			code = src
		case io.Reader:
			d, err := io.ReadAll(src)
			if err != nil {
				return "", err
			}
			code = string(d)
		default:
			return "", fmt.Errorf("unsupported src type: %T", src)
		}
	}
	svc := kcl.Service()
	resp, err := svc.ParseFile(&gpyrpc.ParseFileArgs{
		Path:   filename,
		Source: code,
	})
	if err != nil {
		return "", err
	}
	return resp.AstJson, nil
}

// ParseFile parses the source code from the specified file or Reader
// and returns the Go structure representation of the Abstract Syntax
// Tree (AST). The source code can be provided directly as a string or
// []byte, or indirectly via a filename or an io.Reader. If src is nil,
// the function reads the content from the provided filename.
func ParseFile(filename string, src any) (m *ast.Module, err error) {
	astJson, err := ParseFileASTJson(filename, src)
	if err != nil {
		return nil, err
	}
	m = ast.NewModule()
	err = json.Unmarshal([]byte(astJson), m)
	if err != nil {
		return nil, err
	}
	return
}

// Parse KCL program with entry files and return the AST JSON string.
func ParseProgram(args *ParseProgramArgs) (*ParseProgramResult, error) {
	svc := kcl.Service()
	return svc.ParseProgram(args)
}
