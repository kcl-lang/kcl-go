package ast

// Module is an abstract syntax tree for a single KCL file.
type Module struct {
	Filename string           `json:"filename"`
	Pkg      string           `json:"pkg"`
	Doc      *Node[string]    `json:"doc"`
	Body     []*Node[Stmt]    `json:"body"`
	Comments []*Node[Comment] `json:"comments"`
}

// NewModule creates a new Module instance
func NewModule() *Module {
	return &Module{
		Body:     make([]*Node[Stmt], 0),
		Comments: make([]*Node[Comment], 0),
	}
}

// Node is the file, line and column number information that all AST nodes need to contain.
// In fact, column and end_column are the counts of character. For example, `\t` is counted as 1 character,
// so it is recorded as 1 here, but generally col is 4.
type Node[T any] struct {
	ID   AstIndex `json:"id,omitempty"`
	Node T        `json:"node,omitempty"`
	Pos
}

// AstIndex represents a unique identifier for AST nodes.
type AstIndex string

// Pos denotes the struct tuple (filename, line, column, end_line, end_column).
type Pos struct {
	Filename  string `json:"filename,omitempty"`
	Line      int64  `json:"line,omitempty"`
	Column    int64  `json:"column,omitempty"`
	EndLine   int64  `json:"end_line,omitempty"`
	EndColumn int64  `json:"end_column,omitempty"`
}

// Comment node.
type Comment struct {
	Text string
}
