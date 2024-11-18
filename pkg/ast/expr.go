package ast

// Expr is an interface for all expression types
type Expr interface {
	Type() string
}

// BaseExpr is a struct that all expression types can embed to implement the Expr interface
type BaseExpr struct {
	ExprType string `json:"type"`
}

func (b BaseExpr) Type() string {
	return b.ExprType
}

// Define all the expression types
type (
	// TargetExpr represents a target expression, e.g.
	//
	// a
	// b
	// _c
	// a["b"][0].c
	TargetExpr struct {
		BaseExpr
		Name    *Node[string]   `json:"name"`
		Pkgpath string          `json:"pkgpath,omitempty"`
		Paths   []MemberOrIndex `json:"paths,omitempty"`
	}
	// IdentifierExpr represents an identifier expression, e.g.
	//
	// a
	// b
	// _c
	// pkg.a
	IdentifierExpr struct {
		BaseExpr
		Identifier
		// Add fields specific to IdentifierExpr
	}
	// UnaryExpr represents a unary expression, e.g.
	//
	// +1
	// -2
	// ~3
	// not True
	UnaryExpr struct {
		BaseExpr
		Op      UnaryOp     `json:"op"`
		Operand *Node[Expr] `json:"operand"`
	}
	// BinaryExpr represents a binary expression, e.g.
	//
	// 1 + 1
	// 3 - 2
	// 5 / 2
	// a is None
	BinaryExpr struct {
		BaseExpr
		Left  *Node[Expr] `json:"left"`
		Op    BinOp       `json:"op"`
		Right *Node[Expr] `json:"right"`
	}
	// IfExpr represents an if expression, e.g.
	//
	// 1 if condition else 2
	IfExpr struct {
		BaseExpr
		Body   *Node[Expr] `json:"body"`
		Cond   *Node[Expr] `json:"cond"`
		Orelse *Node[Expr] `json:"orelse"`
	}
	// SelectorExpr represents a selector expression, e.g.
	//
	// x.y
	// x?.y
	SelectorExpr struct {
		BaseExpr
		Value       *Node[Expr]       `json:"value"`
		Attr        *Node[Identifier] `json:"attr"`
		Ctx         ExprContext       `json:"ctx"`
		HasQuestion bool              `json:"has_question"`
	}
	// CallExpr represents a function call expression, e.g.
	//
	// func1()
	// func2(1)
	// func3(x=2)
	CallExpr struct {
		BaseExpr
		Func     *Node[Expr]      `json:"func"`
		Args     []*Node[Expr]    `json:"args"`
		Keywords []*Node[Keyword] `json:"keywords"`
	}
	// ParenExpr represents a parenthesized expression, e.g.
	//
	// 1 + (2 - 3)
	ParenExpr struct {
		BaseExpr
		Expr *Node[Expr] `json:"expr"`
	}
	// QuantExpr represents a quantifier expression, e.g.
	//
	// all x in collection {x > 0}
	// any y in collection {y < 0}
	// map x in collection {x + 1}
	// filter x in collection {x > 1}
	QuantExpr struct {
		BaseExpr
		Target    *Node[Expr]         `json:"target"`
		Variables []*Node[Identifier] `json:"variables"`
		Op        QuantOperation      `json:"op"`
		Test      *Node[Expr]         `json:"test"`
		IfCond    *Node[Expr]         `json:"if_cond"`
		Ctx       ExprContext         `json:"ctx"`
	}
	// ListExpr represents a list expression, e.g.
	//
	// [1, 2, 3]
	// [1, if True: 2, 3]
	ListExpr struct {
		BaseExpr
		Elts []*Node[Expr] `json:"elts"`
		Ctx  ExprContext   `json:"ctx"`
	}
	// ListIfItemExpr represents a list if-item expression, e.g.
	//
	// [1, if True: 2, 3]
	ListIfItemExpr struct {
		BaseExpr
		IfCond *Node[Expr]   `json:"if_cond"`
		Exprs  []*Node[Expr] `json:"exprs"`
		Orelse *Node[Expr]   `json:"orelse"`
	}
	// ListComp represents a list comprehension expression, e.g.
	//
	// [x ** 2 for x in [1, 2, 3]]
	ListComp struct {
		BaseExpr
		Elt        *Node[Expr]         `json:"elt"`
		Generators []*Node[CompClause] `json:"generators"`
	}
	// StarredExpr represents a starred expression, e.g.
	//
	// [1, 2, *[3, 4]]
	StarredExpr struct {
		BaseExpr
		Value *Node[Expr] `json:"value"`
		Ctx   ExprContext `json:"ctx"`
	}
	// DictComp represents a dictionary comprehension expression, e.g.
	//
	// {k: v + 1 for k, v in {k1 = 1, k2 = 2}}
	DictComp struct {
		BaseExpr
		Entry      ConfigEntry         `json:"entry"`
		Generators []*Node[CompClause] `json:"generators"`
	}
	// ConfigIfEntryExpr represents a conditional configuration entry, e.g.
	//
	//	{
	//	  k1 = 1
	//	  if condition:
	//	    k2 = 2
	//	}
	ConfigIfEntryExpr struct {
		BaseExpr
		IfCond *Node[Expr]          `json:"if_cond"`
		Items  []*Node[ConfigEntry] `json:"items"`
		Orelse *Node[Expr]          `json:"orelse"`
	}
	// CompClause represents a comprehension clause, e.g.
	//
	// i, a in [1, 2, 3] if i > 1 and a > 1
	CompClause struct {
		BaseExpr
		Targets []*Node[Identifier] `json:"targets"`
		Iter    *Node[Expr]         `json:"iter"`
		Ifs     []*Node[Expr]       `json:"ifs"`
	}
	// SchemaExpr represents a schema expression, e.g.
	//
	//	ASchema(arguments) {
	//	  attr1 = 1
	//	  attr2 = BSchema {attr3 = 2}
	//	}
	SchemaExpr struct {
		BaseExpr
		Name   *Node[Identifier] `json:"name"`
		Args   []*Node[Expr]     `json:"args"`
		Kwargs []*Node[Keyword]  `json:"kwargs"`
		Config *Node[Expr]       `json:"config"`
	}
	// ConfigExpr represents a configuration expression, e.g.
	//
	//	{
	//	  attr1 = 1
	//	  attr2 += [0, 1]
	//	  attr3: {key = value}
	//	}
	ConfigExpr struct {
		BaseExpr
		Items []*Node[ConfigEntry] `json:"items"`
	}
	// LambdaExpr represents a lambda expression, e.g.
	//
	//	lambda x, y {
	//	  z = 2 * x
	//	  z + y
	//	}
	LambdaExpr struct {
		BaseExpr
		Args     *Node[Arguments] `json:"args"`
		Body     []*Node[Stmt]    `json:"body"`
		ReturnTy *Node[Type]      `json:"return_ty"`
	}
	// Subscript represents a subscript expression, e.g.
	//
	// a[0]
	// b["k"]
	// c?[1]
	// d[1:2:n]
	Subscript struct {
		BaseExpr
		Value       *Node[Expr] `json:"value"`
		Index       *Node[Expr] `json:"index"`
		Lower       *Node[Expr] `json:"lower"`
		Upper       *Node[Expr] `json:"upper"`
		Step        *Node[Expr] `json:"step"`
		Ctx         ExprContext `json:"ctx"`
		HasQuestion bool        `json:"has_question"`
	}
	// Compare represents a comparison expression, e.g.
	//
	// 0 < a < 10
	// b is not None
	// c != d
	Compare struct {
		BaseExpr
		Left        *Node[Expr]   `json:"left"`
		Ops         []CmpOp       `json:"ops"`
		Comparators []*Node[Expr] `json:"comparators"`
	}
	// NumberLit represents a number literal, e.g.
	//
	// 1
	// 2.0
	// 1m
	// 1K
	// 1Mi
	NumberLit struct {
		BaseExpr
		BinarySuffix *NumberBinarySuffix `json:"binary_suffix,omitempty"`
		Value        NumberLitValue      `json:"value"`
	}
	// StringLit represents a string literal, e.g.
	//
	// "string literal"
	// """long string literal"""
	StringLit struct {
		BaseExpr
		IsLongString bool   `json:"is_long_string"`
		RawValue     string `json:"raw_value"`
		Value        string `json:"value"`
	}
	// NameConstantLit represents a name constant literal, e.g.
	//
	// True
	// False
	// None
	// Undefined
	NameConstantLit struct {
		BaseExpr
		Value NameConstant `json:"value"`
	}
	// JoinedString represents a joined string, e.g. abc in the string interpolation "${var1} abc ${var2}"
	JoinedString struct {
		BaseExpr
		IsLongString bool          `json:"is_long_string"`
		Values       []*Node[Expr] `json:"values"`
		RawValue     string        `json:"raw_value"`
	}
	// FormattedValue represents a formatted value, e.g. var1 and var2 in the string interpolation "${var1} abc ${var2}"
	FormattedValue struct {
		BaseExpr
		IsLongString bool        `json:"is_long_string"`
		Value        *Node[Expr] `json:"value"`
		FormatSpec   string      `json:"format_spec"`
	}
	// MissingExpr is a placeholder for error recovery
	MissingExpr struct {
		BaseExpr
	}
)

// NewTargetExpr creates a new TargetExpr
func NewTargetExpr() *TargetExpr {
	return &TargetExpr{
		BaseExpr: BaseExpr{ExprType: "Target"},
		Paths:    make([]MemberOrIndex, 0),
	}
}

// NewIdentifierExpr creates a new IdentifierExpr
func NewIdentifierExpr() *IdentifierExpr {
	return &IdentifierExpr{
		BaseExpr: BaseExpr{ExprType: "Identifier"},
		Identifier: Identifier{
			Names: make([]*Node[string], 0),
		},
	}
}

// NewUnaryExpr creates a new UnaryExpr
func NewUnaryExpr() *UnaryExpr {
	return &UnaryExpr{
		BaseExpr: BaseExpr{ExprType: "Unary"},
	}
}

// NewBinaryExpr creates a new BinaryExpr
func NewBinaryExpr() *BinaryExpr {
	return &BinaryExpr{
		BaseExpr: BaseExpr{ExprType: "Binary"},
	}
}

// NewIfExpr creates a new IfExpr
func NewIfExpr() *IfExpr {
	return &IfExpr{
		BaseExpr: BaseExpr{ExprType: "If"},
	}
}

// NewCallExpr creates a new CallExpr
func NewCallExpr() *CallExpr {
	return &CallExpr{
		BaseExpr: BaseExpr{ExprType: "Call"},
		Args:     make([]*Node[Expr], 0),
		Keywords: make([]*Node[Keyword], 0),
	}
}

// NewParenExpr creates a new ParenExpr
func NewParenExpr() *ParenExpr {
	return &ParenExpr{
		BaseExpr: BaseExpr{ExprType: "Paren"},
	}
}

// NewQuantExpr creates a new QuantExpr
func NewQuantExpr() *QuantExpr {
	return &QuantExpr{
		BaseExpr:  BaseExpr{ExprType: "Quant"},
		Variables: make([]*Node[Identifier], 0),
	}
}

// NewListExpr creates a new ListExpr
func NewListExpr() *ListExpr {
	return &ListExpr{
		BaseExpr: BaseExpr{ExprType: "List"},
		Elts:     make([]*Node[Expr], 0),
	}
}

// NewListIfItemExpr creates a new ListIfItemExpr
func NewListIfItemExpr() *ListIfItemExpr {
	return &ListIfItemExpr{
		BaseExpr: BaseExpr{ExprType: "ListIfItem"},
		Exprs:    make([]*Node[Expr], 0),
	}
}

// NewListComp creates a new ListComp
func NewListComp() *ListComp {
	return &ListComp{
		BaseExpr:   BaseExpr{ExprType: "ListComp"},
		Generators: make([]*Node[CompClause], 0),
	}
}

// NewStarredExpr creates a new StarredExpr
func NewStarredExpr() *StarredExpr {
	return &StarredExpr{
		BaseExpr: BaseExpr{ExprType: "Starred"},
	}
}

// NewDictComp creates a new DictComp
func NewDictComp() *DictComp {
	return &DictComp{
		BaseExpr:   BaseExpr{ExprType: "DictComp"},
		Generators: make([]*Node[CompClause], 0),
	}
}

// NewCompClause creates a new CompClause
func NewCompClause() *CompClause {
	return &CompClause{
		Targets: make([]*Node[Identifier], 0),
		Ifs:     make([]*Node[Expr], 0),
	}
}

// NewSchemaExpr creates a new SchemaExpr
func NewSchemaExpr() *SchemaExpr {
	return &SchemaExpr{
		BaseExpr: BaseExpr{ExprType: "Schema"},
		Args:     make([]*Node[Expr], 0),
		Kwargs:   make([]*Node[Keyword], 0),
	}
}

// NewConfigExpr creates a new ConfigExpr
func NewConfigExpr() *ConfigExpr {
	return &ConfigExpr{
		BaseExpr: BaseExpr{ExprType: "Config"},
		Items:    make([]*Node[ConfigEntry], 0),
	}
}

// NewLambdaExpr creates a new LambdaExpr
func NewLambdaExpr() *LambdaExpr {
	return &LambdaExpr{
		BaseExpr: BaseExpr{ExprType: "Lambda"},
		Body:     make([]*Node[Stmt], 0),
	}
}

// NewSubscript creates a new Subscript
func NewSubscript() *Subscript {
	return &Subscript{
		BaseExpr: BaseExpr{ExprType: "Subscript"},
	}
}

// NewCompare creates a new Compare
func NewCompare() *Compare {
	return &Compare{
		BaseExpr:    BaseExpr{ExprType: "Compare"},
		Ops:         make([]CmpOp, 0),
		Comparators: make([]*Node[Expr], 0),
	}
}

// NewNumberLit creates a new NumberLit
func NewNumberLit() *NumberLit {
	return &NumberLit{
		BaseExpr: BaseExpr{ExprType: "Number"},
	}
}

// NewStringLit creates a new StringLit with default values
func NewStringLit() *StringLit {
	return &StringLit{
		BaseExpr:     BaseExpr{ExprType: "String"},
		Value:        "",
		RawValue:     "\"\"",
		IsLongString: false,
	}
}

// NewNameConstantLit creates a new NameConstantLit
func NewNameConstantLit() *NameConstantLit {
	return &NameConstantLit{
		BaseExpr: BaseExpr{ExprType: "NameConstant"},
	}
}

// NewJoinedString creates a new JoinedString
func NewJoinedString() *JoinedString {
	return &JoinedString{
		BaseExpr: BaseExpr{ExprType: "JoinedString"},
		Values:   make([]*Node[Expr], 0),
	}
}

// NewFormattedValue creates a new FormattedValue
func NewFormattedValue() *FormattedValue {
	return &FormattedValue{
		BaseExpr: BaseExpr{ExprType: "FormattedValue"},
	}
}

// NewMissingExpr creates a new MissingExpr
func NewMissingExpr() *MissingExpr {
	return &MissingExpr{
		BaseExpr: BaseExpr{ExprType: "MissingExpr"},
	}
}

// Identifier represents an identifier, e.g.
//
// a
// b
// _c
// pkg.a
type Identifier struct {
	Names   []*Node[string] `json:"names"`
	Pkgpath string          `json:"pkgpath"`
	Ctx     ExprContext     `json:"ctx"`
}

// ExprContext denotes the value context in the expression. e.g.,
//
// The context of 'a' in 'a = b' is Store
//
// The context of 'b' in 'a = b' is Load
type ExprContext int

const (
	Load ExprContext = iota
	Store
)

// String returns the string representation of ExprContext
func (e ExprContext) String() string {
	return [...]string{"Load", "Store"}[e]
}

// SchemaConfig represents a schema configuration, e.g.
//
//	ASchema(arguments) {
//	  attr1 = 1
//	  attr2 = BSchema {attr3 = 2}
//	}
type SchemaConfig struct {
	Name   *Node[Identifier] `json:"name"`
	Args   []*Node[Expr]     `json:"args"`
	Kwargs []*Node[Keyword]  `json:"kwargs"`
	Config *Node[Expr]       `json:"config"`
}

// NewSchemaConfig creates a new SchemaConfig
func NewSchemaConfig() *SchemaConfig {
	return &SchemaConfig{
		Args:   make([]*Node[Expr], 0),
		Kwargs: make([]*Node[Keyword], 0),
	}
}

// Keyword represents a keyword argument, e.g.
//
// arg = value
type Keyword struct {
	Arg   *Node[Identifier] `json:"arg"`
	Value *Node[Expr]       `json:"value"`
}

// Target represents a target in an assignment, e.g.
//
// a
// b
// _c
// a["b"][0].c
type Target struct {
	Name    *Node[string]    `json:"name"`
	Pkgpath string           `json:"pkgpath"`
	Paths   []*MemberOrIndex `json:"paths"`
}

// MemberOrIndex is the base interface for member or index expression
//
// a.<member>
// b[<index>]
type MemberOrIndex interface {
	Type() string
}

// Member represents a member access
type Member struct {
	Value *Node[string] `json:"value"`
}

// Type returns the type of Member
func (m *Member) Type() string {
	return "Member"
}

// Index represents an index access
type Index struct {
	Value *Node[Expr] `json:"value"`
}

// Type returns the type of Index
func (i *Index) Type() string {
	return "Index"
}

// Decorator represents a decorator, e.g.
//
// deprecated(strict=True)
type Decorator struct {
	Func     *Node[Expr]      `json:"func"`
	Args     []*Node[Expr]    `json:"args,omitempty"`
	Keywords []*Node[Keyword] `json:"keywords,omitempty"`
}

// NewDecorator creates a new Decorator
func NewDecorator() *Decorator {
	return &Decorator{
		Args:     make([]*Node[Expr], 0),
		Keywords: make([]*Node[Keyword], 0),
	}
}

// Arguments represents function arguments, e.g.
//
//	lambda x: int = 1, y: int = 1 {
//	    x + y
//	}
type Arguments struct {
	Args     []*Node[Identifier] `json:"args"`
	Defaults []*Node[Expr]       `json:"defaults,omitempty"` // Slice can contain nil to represent Rust's Vec<Option<Node<Expr>>>
	TyList   []*Node[Type]       `json:"ty_list,omitempty"`  // Slice can contain nil to represent Rust's Vec<Option<Node<Type>>>
}

// NewArguments creates a new Arguments
func NewArguments() *Arguments {
	return &Arguments{
		Args:     make([]*Node[Identifier], 0),
		Defaults: make([]*Node[Expr], 0),
		TyList:   make([]*Node[Type], 0),
	}
}

// CheckExpr represents a check expression, e.g.
//
// len(attr) > 3 if attr, "Check failed message"
type CheckExpr struct {
	Test   *Node[Expr] `json:"test"`
	IfCond *Node[Expr] `json:"if_cond,omitempty"`
	Msg    *Node[Expr] `json:"msg,omitempty"`
}

// NewCheckExpr creates a new CheckExpr
func NewCheckExpr() *CheckExpr {
	return &CheckExpr{}
}

// ConfigEntry represents a configuration entry, e.g.
//
//	{
//	  attr1 = 1
//	  attr2 += [0, 1]
//	  attr3: {key = value}
//	}
type ConfigEntry struct {
	Key       *Node[Expr]          `json:"key"`
	Value     *Node[Expr]          `json:"value"`
	Operation ConfigEntryOperation `json:"operation"`
}

// NewConfigEntry creates a new ConfigEntry
func NewConfigEntry() *ConfigEntry {
	return &ConfigEntry{}
}

// NewConfigIfEntryExpr creates a new ConfigIfEntryExpr
func NewConfigIfEntryExpr() *ConfigIfEntryExpr {
	return &ConfigIfEntryExpr{
		BaseExpr: BaseExpr{ExprType: "ConfigIfEntry"},
		Items:    make([]*Node[ConfigEntry], 0),
	}
}

// NumberLitValue represents the value of a number literal
type NumberLitValue interface {
	Type() string
}

// IntNumberLitValue represents an integer number literal value
type IntNumberLitValue struct {
	Value int64 `json:"value"`
}

// Type returns the type of the number literal value
func (i *IntNumberLitValue) Type() string {
	return "Int"
}

// FloatNumberLitValue represents a float number literal value
type FloatNumberLitValue struct {
	Value float64 `json:"value"`
}

// Type returns the type of the number literal value
func (f *FloatNumberLitValue) Type() string {
	return "Float"
}
