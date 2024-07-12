package parser

import (
	"fmt"
	"io"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// ParseFileASTJson parses the source code from the specified file or Reader
// and returns the JSON representation of the Abstract Syntax Tree (AST).
// The source code can be provided directly as a string or []byte,
// or indirectly via a filename or an io.Reader.
// If src is nil, the function reads the content from the provided filename.
func ParseFileASTJson(filename string, src interface{}) (result string, err error) {
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
	resp, err := svc.ParseFile(&gpyrpc.ParseFile_Args{
		Path:   filename,
		Source: code,
	})
	if err != nil {
		return "", err
	}
	return resp.AstJson, nil
}
