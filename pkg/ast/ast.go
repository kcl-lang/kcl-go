// Copyright 2022 The KCL Authors. All rights reserved.

package ast

type Node interface {
	GetNodeType() AstType
	GetMeta() *Meta

	JSONString() string
	JSONMap() map[string]interface{}

	GetPosition() (_pos, line, column int)
}

type Stmt interface {
	Node
	stmt_type()
}
type Expr interface {
	Node
	expr_type()
}
type TypeInterface interface {
	Node
	type_type()
}

type Meta struct {
	AstType          AstType `json:"_ast_type,omitempty"`
	Filename         string  `json:"filename,omitempty"`
	RelativeFilename string  `json:"relative_filename"`
	Line             int     `json:"line,omitempty"`
	Column           int     `json:"column,omitempty"`
	EndLine          int     `json:"end_line,omitempty"`
	EndColumn        int     `json:"end_column,omitempty"`
}

// kcl main.k -D name=value
type CmdArgSpec struct {
	Name string `json:"name"`
}

// kcl main.k -O pkgpath:path.to.field=field_value
type CmdOverrideSpec struct {
	Pkgpath    string `json:"value"`
	FieldPath  string `json:"field_path"`
	FieldValue string `json:"field_value"`
}

type Name struct {
	*Meta `json:",omitempty"`

	Value string `json:"value"`
}

type TypeAliasStmt struct {
	*Meta `json:",omitempty"`

	TypeName  string `json:"type_name"`
	TypeValue *Type  `json:"type_value"`
}

type ExprStmt struct {
	*Meta `json:",omitempty"`

	Exprs []Expr `json:"exprs"`
}

type UnificationStmt struct {
	*Meta `json:",omitempty"`

	Target *Identifier `json:"target"`
	Value  *SchemaExpr `json:"value"`
}

type AssignStmt struct {
	*Meta `json:",omitempty"`

	Targets            []*Identifier `json:"targets"`
	Value              Expr          `json:"value"`
	TypeAnnotation     string        `json:"type_annotation"`
	TypeAnnotationNode *Type         `json:"type_annotation_node"`
}

type AugAssignStmt struct {
	*Meta `json:",omitempty"`

	Op     string      `json:"op"`
	Target *Identifier `json:"target"`
	Value  Expr        `json:"value"`
}

type AssertStmt struct {
	*Meta `json:",omitempty"`

	Test   Expr `json:"test"`
	IfCond Expr `json:"if_cond"`
	Msg    Expr `json:"msg"`
}

type IfStmt struct {
	*Meta `json:",omitempty"`

	Cond     Expr     `json:"cond"`
	Body     []Stmt   `json:"body"`
	ElifCond []Expr   `json:"elif_cond"`
	ElifBody [][]Stmt `json:"elif_body"`
	ElseBody []Stmt   `json:"else_body"`
}

type ImportStmt struct {
	*Meta `json:",omitempty"`

	Path       string  `json:"path"`
	Name       string  `json:"name"`
	AsName     string  `json:"asname"`
	PathNodes  []*Name `json:"path_nodes,omitempty"`
	AsNameNode *Name   `json:"as_name_node,omitempty"`
	RawPath    string  `json:"rawpath"`
}

type SchemaIndexSignature struct {
	*Meta `json:",omitempty"`

	KeyName       string `json:"key_name"`
	KeyType       string `json:"key_type"`
	ValueType     string `json:"value_type"`
	Value         Expr   `json:"value"`
	AnyOther      bool   `json:"any_other"`
	NameNode      *Name  `json:"name_node,omitempty"`
	ValueTypeNode *Type  `json:"value_type_node"`
}

type SchemaAttr struct {
	*Meta `json:",omitempty"`

	Doc        string       `json:"doc"`
	Name       string       `json:"name"`
	TypeStr    string       `json:"type_str"`
	Op         string       `json:"op"`
	Value      Expr         `json:"value"`
	IsFinal    bool         `json:"is_final"`
	IsOptional bool         `json:"is_optional"`
	Decorators []*Decorator `json:"decorators"`
	NameNode   *Name        `json:"name_node,omitempty"`
	TypeNode   *Type        `json:"type_node,omitempty"`
}

type SchemaStmt struct {
	*Meta `json:",omitempty"`

	Doc            string                `json:"doc"`
	Name           string                `json:"name"`
	ParentName     *Identifier           `json:"parent_name"`
	ForHostName    *Identifier           `json:"for_host_name"`
	IsRelaxed      bool                  `json:"is_relaxed"`
	IsMixin        bool                  `json:"is_mixin"`
	IsProtocol     bool                  `json:"is_protocol"`
	Args           *Arguments            `json:"args"`
	Mixins         []*Identifier         `json:"mixins"`
	Body           []Stmt                `json:"body"`
	Decorators     []*Decorator          `json:"decorators"`
	Checks         []*CheckExpr          `json:"checks"`
	IndexSignature *SchemaIndexSignature `json:"index_signature"`
	NameNode       *Name                 `json:"name_node,omitempty"`
}

type RuleStmt struct {
	*Meta `json:",omitempty"`

	Doc  string `json:"doc"`
	Name string `json:"name"`

	ParentRules []*Identifier `json:"parent_rules"`
	Decorators  []*Decorator  `json:"decorators"`
	Checks      []*CheckExpr  `json:"checks"`
	NameNode    *Name         `json:"name_node,omitempty"`
	Args        *Arguments    `json:"args"`
	ForHostName *Identifier   `json:"for_host_name"`
}

type IfExpr struct {
	*Meta `json:",omitempty"`

	Body   Expr `json:"body"`
	Cond   Expr `json:"cond"`
	OrElse Expr `json:"orelse"`
}

type UnaryExpr struct {
	*Meta `json:",omitempty"`

	Op      string `json:"op"`
	Operand Expr   `json:"operand"`
}

type BinaryExpr struct {
	*Meta `json:",omitempty"`

	Left  Expr   `json:"left"`
	Op    string `json:"op"`
	Right Expr   `json:"right"`
}

type SelectorExpr struct {
	*Meta `json:",omitempty"`

	Value       Expr        `json:"value"`
	Attr        *Identifier `json:"attr"`
	Ctx         string      `json:"ctx"`
	HasQuestion bool        `json:"has_question"`
}

type CallExpr struct {
	*Meta `json:",omitempty"`

	Func     Expr       `json:"func"`
	Args     []Expr     `json:"args"`
	Keywords []*Keyword `json:"keywords"`
}

type ParenExpr struct {
	*Meta `json:",omitempty"`

	Expr Expr `json:"expr"`
}

type QuantExpr struct {
	*Meta `json:",omitempty"`

	Target    Expr          `json:"target"`
	Variables []*Identifier `json:"variables"`
	Op        int           `json:"op"` // string?
	CheckTest Expr          `json:"check_test"`
	IfCond    Expr          `json:"if_cond"`
	Ctx       string        `json:"ctx"`
}

type ListExpr struct {
	*Meta `json:",omitempty"`

	Elts []Expr `json:"elts"`
	Ctx  string `json:"ctx"`
}

type ListIfItemExpr struct {
	*Meta `json:",omitempty"`

	IfCond Expr   `json:"if_cond"`
	Exprs  []Expr `json:"exprs"`
	Orelse Expr   `json:"orelse"`
}

type ListComp struct {
	*Meta `json:",omitempty"`

	Elt        Expr          `json:"elt"`
	Generators []*CompClause `json:"generators"`
}

type StarredExpr struct {
	*Meta `json:",omitempty"`

	Value Expr   `json:"value"`
	Ctx   string `json:"ctx"`
}

type DictComp struct {
	*Meta `json:",omitempty"`

	Key        Expr          `json:"key"`
	Value      Expr          `json:"value"`
	Generators []*CompClause `json:"generators"`
}

type ConfigIfEntryExpr struct {
	*Meta `json:",omitempty"`

	IfCond     Expr   `json:"if_cond"`
	Keys       []Expr `json:"keys"`
	Values     []Expr `json:"values"`
	Operations []int  `json:"operations"` // string?
	Orelse     Expr   `json:"orelse"`
}

type CompClause struct {
	*Meta `json:",omitempty"`

	Targets []*Identifier `json:"targets"`
	Iter    Expr          `json:"iter"`
	Ifs     []Expr        `json:"ifs"`
}

type SchemaExpr struct {
	*Meta `json:",omitempty"`

	Name   *Identifier `json:"name"`
	Args   []Expr      `json:"args"`
	Kwargs []*Keyword  `json:"kwargs"`
	Config *ConfigExpr `json:"config"`
}

type ConfigExpr struct {
	*Meta `json:",omitempty"`

	Items []*ConfigEntry `json:"items"`
}

type ConfigEntry struct {
	*Meta `json:",omitempty"`

	Key         Expr `json:"key"`
	Value       Expr `json:"value"`
	Operation   int  `json:"operation"`
	InsertIndex int  `json:"insert_index"`
}

type CheckExpr struct {
	*Meta `json:",omitempty"`

	CheckTest Expr `json:"check_test"`
	IfCond    Expr `json:"if_cond"`
	Msg       Expr `json:"msg"`
}

type LambdaExpr struct {
	*Meta `json:",omitempty"`

	Args           *Arguments `json:"args"`
	ReturnTypeStr  string     `json:"return_type_str"`
	ReturnTypeNode *Type      `json:"return_type_node"`
	Body           []Stmt     `json:"body"`
}

type Decorator struct {
	*Meta `json:",omitempty"`

	Name *Identifier `json:"name"`
	Args *CallExpr   `json:"args"`
}

type Subscript struct {
	*Meta `json:",omitempty"`

	Value       Expr   `json:"value"`
	Index       Expr   `json:"index"`
	Lower       Expr   `json:"lower"`
	Upper       Expr   `json:"upper"`
	Step        Expr   `json:"step"`
	Ctx         string `json:"ctx"`
	HasQuestion bool   `json:"has_question"`
}

type Keyword struct {
	*Meta `json:",omitempty"`

	Arg   *Identifier `json:"arg"`
	Value Expr        `json:"value"`
}

type Arguments struct {
	*Meta `json:",omitempty"`

	Args                   []*Identifier `json:"args"`
	Defaults               []Expr        `json:"defaults"`
	TypeAnnotationList     []string      `json:"type_annotation_list,omitempty"`
	TypeAnnotationNodeList []*Type       `json:"type_annotation_node_list,omitempty"`
}

type Compare struct {
	*Meta `json:",omitempty"`

	Left        Expr     `json:"left"`
	Ops         []string `json:"ops"`
	Comparators []Expr   `json:"comparators"`
}

type Identifier struct {
	*Meta `json:",omitempty"`

	Names     []string `json:"names"`
	Pkgpath   string   `json:"pkgpath"`
	Ctx       string   `json:"ctx"`
	NameNodes []*Name  `json:"name_nodes,omitempty"`
}

type Literal struct {
	*Meta `json:",omitempty"`

	Value interface{} `json:"value"`
}

type NumberLit struct {
	*Meta `json:",omitempty"`

	Value        float64 `json:"value"`
	BinarySuffix string  `json:"binary_suffix"`
}

type StringLit struct {
	*Meta `json:",omitempty"`

	Value        string `json:"value"`
	IsLongString bool   `json:"is_long_string"`
	RawValue     string `json:"raw_value"`
}

type NameConstantLit struct {
	*Meta `json:",omitempty"`

	Value string `json:"value"` // True, False, None
}

type JoinedString struct {
	*Meta `json:",omitempty"`

	IsLongString bool   `json:"is_long_string"`
	Values       []Expr `json:"values"` // StringLit, FormattedValue
	RawValue     string `json:"raw_value"`
}

type FormattedValue struct {
	*Meta `json:",omitempty"`

	IsLongString bool   `json:"is_long_string"`
	Value        Expr   `json:"value"`
	FormatSpec   string `json:"format_spec"`
}

type Comment struct {
	*Meta `json:",omitempty"`

	Text string `json:"text"`
}

type CommentGroup struct {
	*Meta `json:",omitempty"`

	Comments []*Comment `json:"comments"`
}

type Type struct {
	*Meta `json:",omitempty"`

	TypeElements []TypeInterface `json:"type_elements,omitempty"`
	PlainTypeStr string          `json:"plain_type_str"`
}

type BasicType struct {
	*Meta `json:",omitempty"`

	TypeElements []TypeInterface `json:"type_elements,omitempty"`
	TypeName     string          `json:"type_name"`
}

type ListType struct {
	*Meta `json:",omitempty"`

	InnerType    TypeInterface `json:"inner_type"`
	PlainTypeStr string        `json:"plain_type_str"`
}

type DictType struct {
	*Meta `json:",omitempty"`

	KeyType      TypeInterface `json:"key_type"`
	ValueType    TypeInterface `json:"value_type"`
	PlainTypeStr string        `json:"plain_type_str"`
}

type LiteralType struct {
	*Meta `json:",omitempty"`

	PlainValue  string     `json:"plain_value"`
	ValueType   string     `json:"value_type"`
	StringValue *StringLit `json:"string_value"`
	NumberValue *NumberLit `json:"number_value"`
}

type Module struct {
	*Meta `json:",omitempty"`

	Pkg  string `json:"pkg"`
	Body []Stmt `json:"body"`
	Doc  string `json:"doc"`
	Name string `json:"name"`

	GlobalNames []string            `json:"global_names"`
	LocalNames  map[string][]string `json:"local_names"`
	Comments    []*Comment          `json:"comments"`
}

type Program struct {
	Root string               `json:"root"`
	Main string               `json:"main"`
	Pkgs map[string][]*Module `json:"pkgs"`

	CmdArgs      []*CmdArgSpec      `json:"cmd_args"`
	CmdOverrides []*CmdOverrideSpec `json:"cmd_overrides"`
}

type File struct {
	JSON   string
	Module *Module
}
