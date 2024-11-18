package ast

// Stmt is an interface for all statement types
type Stmt interface {
	Type() string
}

// BaseStmt is a struct that all statement types can embed to implement the Stmt interface
type BaseStmt struct {
	StmtType string `json:"type"`
}

func (b BaseStmt) Type() string {
	return b.StmtType
}

// Define all the statement types
type (
	// TypeAliasStmt represents a type alias statement, e.g.
	//
	// type StrOrInt = str | int
	TypeAliasStmt struct {
		BaseStmt
		TypeName  *Node[Identifier] `json:"type_name"`
		TypeValue *Node[string]     `json:"type_value"`
		Ty        *Node[Type]       `json:"ty"`
	}
	// ExprStmt represents a expression statement, e.g.
	//
	// 1
	//
	// """A long string"""
	//
	// 'A string'
	ExprStmt struct {
		BaseStmt
		Exprs []*Node[Expr] `json:"exprs"`
	}
	// UnificationStmt represents a declare statement with the union operator, e.g.
	//
	// data: ASchema {}
	UnificationStmt struct {
		BaseStmt
		Target *Node[Identifier]   `json:"target"`
		Value  *Node[SchemaConfig] `json:"value"`
	}
	// AssignStmt represents an assignment, e.g.
	//
	// a: int = 1
	//
	// a = 1
	//
	// a = b = 1
	AssignStmt struct {
		BaseStmt
		Targets []*Node[Target] `json:"targets"`
		Value   *Node[Expr]     `json:"value"`
		Ty      *Node[Type]     `json:"ty"`
	}
	// AugAssignStmt represents an augmented assignment, e.g.
	//
	// a += 1
	//
	// a -= 1
	AugAssignStmt struct {
		BaseStmt
		Target *Node[Target] `json:"target"`
		Value  *Node[Expr]   `json:"value"`
		Op     AugOp         `json:"op"`
	}
	// AssertStmt represents an assert statement, e.g.
	//
	// assert True if condition, "Assert failed message"
	AssertStmt struct {
		BaseStmt
		Test   *Node[Expr] `json:"test"`
		IfCond *Node[Expr] `json:"if_cond,omitempty"`
		Msg    *Node[Expr] `json:"msg,omitempty"`
	}
	// IfStmt represents an if statement, e.g.
	//
	// if condition1:
	//
	//	if condition2:
	//	    a = 1
	//
	// elif condition3:
	//
	//	b = 2
	//
	// else:
	//
	//	c = 3
	IfStmt struct {
		BaseStmt
		Body   []*Node[Stmt] `json:"body"`
		Cond   *Node[Expr]   `json:"cond"`
		Orelse []*Node[Stmt] `json:"orelse,omitempty"`
	}
	// ImportStmt represents an import statement, e.g.
	//
	// import pkg as pkg_alias
	ImportStmt struct {
		BaseStmt
		Path    *Node[string] `json:"path"`
		Rawpath string        `json:"rawpath"`
		Name    string        `json:"name"`
		Asname  *Node[string] `json:"asname,omitempty"`
		PkgName string        `json:"pkg_name"`
	}
	// SchemaAttr represents schema attribute definitions, e.g.
	//
	// schema SchemaAttrExample:
	//
	//	x: int
	//	y: str
	SchemaAttr struct {
		BaseStmt
		Doc        string             `json:"doc,omitempty"`
		Name       *Node[string]      `json:"name"`
		Op         AugOp              `json:"op,omitempty"`
		Value      *Node[Expr]        `json:"value,omitempty"`
		IsOptional bool               `json:"is_optional"`
		Decorators []*Node[Decorator] `json:"decorators,omitempty"`
		Ty         *Node[Type]        `json:"ty,omitempty"`
	}
	// SchemaStmt represents a schema statement, e.g.
	//
	// schema BaseSchema:
	//
	// schema SchemaExample(BaseSchema)[arg: str]:
	//
	//	"""Schema documents"""
	//	attr?: str = arg
	//	check:
	//	    len(attr) > 3 if attr, "Check failed message"
	//
	// mixin MixinExample for ProtocolExample:
	//
	//	attr: int
	//
	// protocol ProtocolExample:
	//
	//	attr: int
	SchemaStmt struct {
		BaseStmt
		Doc            *Node[string]               `json:"doc,omitempty"`
		Name           *Node[string]               `json:"name"`
		ParentName     *Node[Identifier]           `json:"parent_name,omitempty"`
		ForHostName    *Node[Identifier]           `json:"for_host_name,omitempty"`
		IsMixin        bool                        `json:"is_mixin"`
		IsProtocol     bool                        `json:"is_protocol"`
		Args           *Node[Arguments]            `json:"args,omitempty"`
		Mixins         []*Node[Identifier]         `json:"mixins,omitempty"`
		Body           []*Node[Stmt]               `json:"body,omitempty"`
		Decorators     []*Node[Decorator]          `json:"decorators,omitempty"`
		Checks         []*Node[CheckExpr]          `json:"checks,omitempty"`
		IndexSignature *Node[SchemaIndexSignature] `json:"index_signature,omitempty"`
	}
	// RuleStmt represents a rule statement, e.g.
	//
	// rule RuleExample:
	//
	//	a > 1
	//	b < 0
	RuleStmt struct {
		BaseStmt
		Doc         *Node[string]       `json:"doc,omitempty"`
		Name        *Node[string]       `json:"name"`
		ParentRules []*Node[Identifier] `json:"parent_rules,omitempty"`
		Decorators  []*Node[Decorator]  `json:"decorators,omitempty"`
		Checks      []*Node[CheckExpr]  `json:"checks,omitempty"`
		Args        *Node[Arguments]    `json:"args,omitempty"`
		ForHostName *Node[Identifier]   `json:"for_host_name,omitempty"`
	}
)

// NewTypeAliasStmt creates a new TypeAliasStmt
func NewTypeAliasStmt() *TypeAliasStmt {
	return &TypeAliasStmt{
		BaseStmt: BaseStmt{StmtType: "TypeAlias"},
	}
}

// NewExprStmt creates a new ExprStmt
func NewExprStmt() *ExprStmt {
	return &ExprStmt{
		BaseStmt: BaseStmt{StmtType: "Expr"},
		Exprs:    make([]*Node[Expr], 0),
	}
}

// NewUnificationStmt creates a new UnificationStmt
func NewUnificationStmt() *UnificationStmt {
	return &UnificationStmt{
		BaseStmt: BaseStmt{StmtType: "Unification"},
	}
}

// NewAssignStmt creates a new AssignStmt
func NewAssignStmt() *AssignStmt {
	return &AssignStmt{
		BaseStmt: BaseStmt{StmtType: "Assign"},
	}
}

// NewAugAssignStmt creates a new AugAssignStmt
func NewAugAssignStmt() *AugAssignStmt {
	return &AugAssignStmt{
		BaseStmt: BaseStmt{StmtType: "AugAssign"},
	}
}

// NewAssertStmt creates a new AssertStmt
func NewAssertStmt() *AssertStmt {
	return &AssertStmt{
		BaseStmt: BaseStmt{StmtType: "Assert"},
	}
}

// NewIfStmt creates a new IfStmt
func NewIfStmt() *IfStmt {
	return &IfStmt{
		BaseStmt: BaseStmt{StmtType: "If"},
		Body:     make([]*Node[Stmt], 0),
		Orelse:   make([]*Node[Stmt], 0),
	}
}

// NewImportStmt creates a new ImportStmt
func NewImportStmt() *ImportStmt {
	return &ImportStmt{
		BaseStmt: BaseStmt{StmtType: "Import"},
	}
}

// NewSchemaAttr creates a new SchemaAttr
func NewSchemaAttr() *SchemaAttr {
	return &SchemaAttr{
		BaseStmt:   BaseStmt{StmtType: "SchemaAttr"},
		Decorators: make([]*Node[Decorator], 0),
	}
}

// NewSchemaStmt creates a new SchemaStmt
func NewSchemaStmt() *SchemaStmt {
	return &SchemaStmt{
		BaseStmt:   BaseStmt{StmtType: "Schema"},
		Mixins:     make([]*Node[Identifier], 0),
		Body:       make([]*Node[Stmt], 0),
		Decorators: make([]*Node[Decorator], 0),
		Checks:     make([]*Node[CheckExpr], 0),
	}
}

// NewRuleStmt creates a new RuleStmt
func NewRuleStmt() *RuleStmt {
	return &RuleStmt{
		BaseStmt:    BaseStmt{StmtType: "Rule"},
		ParentRules: make([]*Node[Identifier], 0),
		Decorators:  make([]*Node[Decorator], 0),
		Checks:      make([]*Node[CheckExpr], 0),
	}
}

// SchemaIndexSignature represents a schema index signature, e.g.
//
// schema SchemaIndexSignatureExample:
//
//	[str]: int
type SchemaIndexSignature struct {
	KeyName  *Node[string] `json:"key_name,omitempty"`
	Value    *Node[Expr]   `json:"value,omitempty"`
	AnyOther bool          `json:"any_other"`
	KeyTy    *Node[Type]   `json:"key_ty,omitempty"`
	ValueTy  *Node[Type]   `json:"value_ty,omitempty"`
}

// NewSchemaIndexSignature creates a new SchemaIndexSignature
func NewSchemaIndexSignature() *SchemaIndexSignature {
	return &SchemaIndexSignature{}
}
