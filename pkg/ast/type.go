package ast

// Type is the base interface for all AST types
type Type interface {
	TypeName() string
}

// NamedType represents a named type
type NamedType struct {
	Value struct {
		Identifier *Identifier `json:"identifier"`
	} `json:"value"`
}

func (n *NamedType) TypeName() string { return "Named" }

// AnyType represents a the any type
type AnyType struct{}

func (n *AnyType) TypeName() string { return "Any" }

// BasicType represents a basic type
type BasicType struct {
	Value BasicTypeEnum `json:"value"`
}

func (b *BasicType) TypeName() string { return "Basic" }

type BasicTypeEnum string

const (
	Bool  BasicTypeEnum = "Bool"
	Int   BasicTypeEnum = "Int"
	Float BasicTypeEnum = "Float"
	Str   BasicTypeEnum = "Str"
)

// ListType represents a list type
type ListType struct {
	Value struct {
		InnerType *Node[Type] `json:"inner_type,omitempty"`
	} `json:"value"`
}

func (l *ListType) TypeName() string { return "List" }

// DictType represents a dictionary type
type DictType struct {
	Value struct {
		KeyType   *Node[Type] `json:"key_type,omitempty"`
		ValueType *Node[Type] `json:"value_type,omitempty"`
	} `json:"value"`
}

func (d *DictType) TypeName() string { return "Dict" }

// UnionType represents a union type
type UnionType struct {
	Value struct {
		TypeElements []*Node[Type] `json:"type_elements"`
	} `json:"value"`
}

func (u *UnionType) TypeName() string { return "Union" }

// LiteralType represents a literal type
type LiteralType struct {
	Value LiteralTypeValue `json:"value"`
}

func (l *LiteralType) TypeName() string { return "Literal" }

// LiteralTypeValue is an interface for different literal types
type LiteralTypeValue interface {
	LiteralTypeName() string
}

// BoolLiteralType represents a boolean literal type
type BoolLiteralType bool

func (b *BoolLiteralType) LiteralTypeName() string { return "Bool" }

// IntLiteralType represents an integer literal type
type IntLiteralType struct {
	Value  int                 `json:"value"`
	Suffix *NumberBinarySuffix `json:"binary_suffix,omitempty"`
}

func (i *IntLiteralType) LiteralTypeName() string { return "Int" }

// FloatLiteralType represents a float literal type
type FloatLiteralType float64

func (f *FloatLiteralType) LiteralTypeName() string { return "Float" }

// StrLiteralType represents a string literal type
type StrLiteralType string

func (s *StrLiteralType) LiteralTypeName() string { return "Str" }

// FunctionType represents a function type
type FunctionType struct {
	Value struct {
		ParamsTy []*Node[Type] `json:"params_ty,omitempty"`
		RetTy    *Node[Type]   `json:"ret_ty,omitempty"`
	} `json:"value"`
}

func (f *FunctionType) TypeName() string { return "Function" }
