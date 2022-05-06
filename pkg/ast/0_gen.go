// Copyright 2022 The KCL Authors. All rights reserved.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
	"text/template"

	"kusionstack.io/kclvm-go/pkg/tools/kclvm_tool"
)

func main() {
	spec := buildAstSpec()

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmpl))
	err := t.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(buf.String())
}

type AstSpec struct {
	NodeNameList []string
	StmtNameList []string
	ExprNameList []string
	TypeNameList []string

	PyAstPath string
	PyAstMD5  string
}

func buildAstSpec() *AstSpec {
	p := NewPyAstTypeParser()

	spec := &AstSpec{
		PyAstPath: p.GetPyAstPath(),
		PyAstMD5:  p.GetPyAstMD5(),
	}

	for _, s := range p.GetNodeTypeList() {
		spec.NodeNameList = append(spec.NodeNameList, string(s))
	}

	for _, s := range spec.NodeNameList {
		switch {
		case p.IsStmtType(s):
			spec.StmtNameList = append(spec.StmtNameList, s)
		case p.IsExprType(s):
			spec.ExprNameList = append(spec.ExprNameList, s)
		case p.IsTypeType(s):
			spec.TypeNameList = append(spec.TypeNameList, s)
		}
	}

	sort.Strings(spec.NodeNameList)
	sort.Strings(spec.StmtNameList)
	sort.Strings(spec.ExprNameList)
	sort.Strings(spec.TypeNameList)

	return spec
}

const tmpl = `
{{- $root := . -}}

// Auto generated. DO NOT EDIT.

package ast

const PyAstMD5 = "{{$root.PyAstMD5}}" // file: ${KCLVM_SRC_ROOT}/kclvm/{{$root.PyAstPath}}

const ({{range $_, $name := .NodeNameList}}
	{{$name}}_TypeName AstType = "{{$name}}"
{{- end}}
)

func init() {
	_ = _ast_node_factory_map

	{{range $_, $name := .NodeNameList}}
	_ast_node_factory_map[{{$name}}_TypeName] = func() Node { return &{{$name}}{Meta: &Meta{AstType: {{$name}}_TypeName}} }
	{{- end}}
}

{{range $_, $name := .NodeNameList}}
func (p *{{$name}}) GetNodeType() AstType { return {{$name}}_TypeName }
func (p *{{$name}}) GetMeta() *Meta { return p.Meta }
func (p *{{$name}}) JSONString() string {return json_String(p)}
func (p *{{$name}}) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *{{$name}}) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}
{{end}}

{{range $_, $name := .StmtNameList}}
func (p *{{$name}}) stmt_type() {}
{{- end}}

{{range $_, $name := .ExprNameList}}
func (p *{{$name}}) expr_type() {}
{{- end}}

{{range $_, $name := .TypeNameList}}
func (p *{{$name}}) type_type() {}
{{- end}}
`

type PyAstTypeParser struct {
	ast_py      string
	ast_py_path string
	ast_py_md5  string
}

func NewPyAstTypeParser() *PyAstTypeParser {
	const ast_py_filename = "kcl/ast/ast.py"

	ast_py := GetAstPy(ast_py_filename)
	return &PyAstTypeParser{
		ast_py:      ast_py,
		ast_py_path: ast_py_filename,
		ast_py_md5:  fmt.Sprintf("%x", md5.Sum([]byte(ast_py))),
	}
}

func (p *PyAstTypeParser) GetPyAstPath() string {
	return p.ast_py_path
}

func (p *PyAstTypeParser) GetPyAstMD5() string {
	return p.ast_py_md5
}

func (p *PyAstTypeParser) GetNodeTypeList() []string {
	var lines []string
	for _, line := range strings.Split(p.ast_py, "\n") {
		if s := strings.TrimSpace(line); s != "" {
			lines = append(lines, s)
		}
	}

	var m = make(map[string]string)
	for i, line := range lines {
		// self._ast_type = "Expr"
		if strings.Contains(line, "self._ast_type") {
			typeName := p.GetNodeType(lines, i)
			if _, ok := m[typeName]; ok {
				log.Printf("dup: %s\n", typeName)
			}

			switch typeName {
			case "Stmt", "Expr":
				continue
			}

			m[typeName] = typeName
		}
	}

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (p *PyAstTypeParser) GetNodeType(lines []string, lineno int) string {
	line := lines[lineno]
	idx0 := strings.Index(line, `"`)
	idx1 := strings.LastIndex(line, `"`)
	if idx0 > 0 && idx1 > idx0 {
		return line[idx0+1 : idx1]
	}
	log.Fatal("invalie line:", line)
	return ""
}

func (p *PyAstTypeParser) IsStmtType(name string) bool {
	if strings.HasSuffix(name, "Stmt") {
		return true
	}
	if strings.Contains(p.ast_py, fmt.Sprintf("%s(Stmt)", name)) {
		return true
	}
	return false
}
func (p *PyAstTypeParser) IsExprType(name string) bool {
	if strings.HasSuffix(name, "Expr") {
		return true
	}
	if strings.Contains(p.ast_py, fmt.Sprintf("%s(Expr)", name)) {
		return true
	}
	if strings.Contains(p.ast_py, fmt.Sprintf("%s(Expr,", name)) {
		return true
	}
	if strings.Contains(p.ast_py, fmt.Sprintf("%s(Literal)", name)) {
		return true
	}
	return false
}
func (p *PyAstTypeParser) IsTypeType(name string) bool {
	if strings.HasSuffix(name, "Type") {
		return true
	}
	if strings.Contains(p.ast_py, fmt.Sprintf("%s(Type)", name)) {
		return true
	}
	return false
}

func GetAstPy(ast_py_filename string) string {
	d, err := fs.ReadFile(kclvm_tool.GetFS(), ast_py_filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(d)
}
