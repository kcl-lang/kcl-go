package ast

// TODO: add more nodes from https://github.com/kcl-lang/kcl/blob/main/kclvm/ast/src/ast.rs

// Pos denotes the struct tuple (filename, line, column, end_line, end_column).
type Pos struct {
	Filename  string `json:"filename"`
	Line      uint64 `json:"line"`
	Column    uint64 `json:"column"`
	EndLine   uint64 `json:"end_line"`
	EndColumn uint64 `json:"end_column"`
}

// Node is the file, line, and column number information that all AST nodes need to contain.
type Node interface {
	Pos() Pos
	Index() string
}

// AstIndex represents a unique identifier for AST nodes.
type AstIndex string

// Comment node.
type Comment struct {
	Text string
}
