package ast

import (
	"encoding/json"
	"fmt"
)

// AugOp represents augmented assignment operations
type AugOp string

const (
	AugOpAssign   AugOp = "="
	AugOpAdd      AugOp = "+="
	AugOpSub      AugOp = "-="
	AugOpMul      AugOp = "*="
	AugOpDiv      AugOp = "/="
	AugOpMod      AugOp = "%="
	AugOpPow      AugOp = "**="
	AugOpFloorDiv AugOp = "//="
	AugOpLShift   AugOp = "<<="
	AugOpRShift   AugOp = ">>="
	AugOpBitXor   AugOp = "^="
	AugOpBitAnd   AugOp = "&="
	AugOpBitOr    AugOp = "|="
)

// Symbol returns the string representation of the AugOp
func (a AugOp) Symbol() string {
	return string(a)
}

// ToBinOp converts AugOp to BinOp, if possible
func (a AugOp) ToBinOp() (BinOp, error) {
	switch a {
	case AugOpAdd:
		return BinOpAdd, nil
	case AugOpSub:
		return BinOpSub, nil
	case AugOpMul:
		return BinOpMul, nil
	case AugOpDiv:
		return BinOpDiv, nil
	case AugOpMod:
		return BinOpMod, nil
	case AugOpPow:
		return BinOpPow, nil
	case AugOpFloorDiv:
		return BinOpFloorDiv, nil
	case AugOpLShift:
		return BinOpLShift, nil
	case AugOpRShift:
		return BinOpRShift, nil
	case AugOpBitXor:
		return BinOpBitXor, nil
	case AugOpBitAnd:
		return BinOpBitAnd, nil
	case AugOpBitOr:
		return BinOpBitOr, nil
	default:
		return "", fmt.Errorf("AugOp cannot be converted to BinOp")
	}
}

// UnaryOp represents a unary operator
type UnaryOp string

const (
	UnaryOpUAdd   UnaryOp = "+"
	UnaryOpUSub   UnaryOp = "-"
	UnaryOpInvert UnaryOp = "~"
	UnaryOpNot    UnaryOp = "not"
)

// Symbol returns the string representation of the unary operator
func (op UnaryOp) Symbol() string {
	return string(op)
}

// AllUnaryOps returns all possible UnaryOp values
func AllUnaryOps() []UnaryOp {
	return []UnaryOp{
		UnaryOpUAdd,
		UnaryOpUSub,
		UnaryOpInvert,
		UnaryOpNot,
	}
}

// UnaryOpFromSymbol returns the UnaryOp corresponding to the given symbol
func UnaryOpFromSymbol(symbol string) (UnaryOp, bool) {
	switch symbol {
	case "+":
		return UnaryOpUAdd, true
	case "-":
		return UnaryOpUSub, true
	case "~":
		return UnaryOpInvert, true
	case "not":
		return UnaryOpNot, true
	default:
		return "", false
	}
}

// BinOp represents a binary operator
type BinOp string

const (
	BinOpAdd      BinOp = "+"
	BinOpSub      BinOp = "-"
	BinOpMul      BinOp = "*"
	BinOpDiv      BinOp = "/"
	BinOpMod      BinOp = "%"
	BinOpPow      BinOp = "**"
	BinOpFloorDiv BinOp = "//"
	BinOpLShift   BinOp = "<<"
	BinOpRShift   BinOp = ">>"
	BinOpBitXor   BinOp = "^"
	BinOpBitAnd   BinOp = "&"
	BinOpBitOr    BinOp = "|"
	BinOpAnd      BinOp = "and"
	BinOpOr       BinOp = "or"
	BinOpAs       BinOp = "as"
)

// Symbol returns the string representation of the binary operator
func (op BinOp) Symbol() string {
	return string(op)
}

// AllBinOps returns all possible BinOp values
func AllBinOps() []BinOp {
	return []BinOp{
		BinOpAdd,
		BinOpSub,
		BinOpMul,
		BinOpDiv,
		BinOpMod,
		BinOpPow,
		BinOpFloorDiv,
		BinOpLShift,
		BinOpRShift,
		BinOpBitXor,
		BinOpBitAnd,
		BinOpBitOr,
		BinOpAnd,
		BinOpOr,
		BinOpAs,
	}
}

// BinOpFromSymbol returns the BinOp corresponding to the given symbol
func BinOpFromSymbol(symbol string) (BinOp, bool) {
	switch symbol {
	case "+":
		return BinOpAdd, true
	case "-":
		return BinOpSub, true
	case "*":
		return BinOpMul, true
	case "/":
		return BinOpDiv, true
	case "%":
		return BinOpMod, true
	case "**":
		return BinOpPow, true
	case "//":
		return BinOpFloorDiv, true
	case "<<":
		return BinOpLShift, true
	case ">>":
		return BinOpRShift, true
	case "^":
		return BinOpBitXor, true
	case "&":
		return BinOpBitAnd, true
	case "|":
		return BinOpBitOr, true
	case "and":
		return BinOpAnd, true
	case "or":
		return BinOpOr, true
	case "as":
		return BinOpAs, true
	default:
		return "", false
	}
}

// CmpOp represents a comparison operator
type CmpOp string

const (
	CmpOpEq    CmpOp = "=="
	CmpOpNotEq CmpOp = "!="
	CmpOpLt    CmpOp = "<"
	CmpOpLtE   CmpOp = "<="
	CmpOpGt    CmpOp = ">"
	CmpOpGtE   CmpOp = ">="
	CmpOpIs    CmpOp = "is"
	CmpOpIn    CmpOp = "in"
	CmpOpNotIn CmpOp = "not in"
	CmpOpNot   CmpOp = "not"
	CmpOpIsNot CmpOp = "is not"
)

// Symbol returns the string representation of the comparison operator
func (c CmpOp) Symbol() string {
	return string(c)
}

// AllCmpOps returns all possible CmpOp values
func AllCmpOps() []CmpOp {
	return []CmpOp{
		CmpOpEq,
		CmpOpNotEq,
		CmpOpLt,
		CmpOpLtE,
		CmpOpGt,
		CmpOpGtE,
		CmpOpIs,
		CmpOpIn,
		CmpOpNotIn,
		CmpOpNot,
		CmpOpIsNot,
	}
}

// CmpOpFromString returns the CmpOp corresponding to the given string
func CmpOpFromString(s string) (CmpOp, bool) {
	switch s {
	case "==":
		return CmpOpEq, true
	case "!=":
		return CmpOpNotEq, true
	case "<":
		return CmpOpLt, true
	case "<=":
		return CmpOpLtE, true
	case ">":
		return CmpOpGt, true
	case ">=":
		return CmpOpGtE, true
	case "is":
		return CmpOpIs, true
	case "in":
		return CmpOpIn, true
	case "not in":
		return CmpOpNotIn, true
	case "not":
		return CmpOpNot, true
	case "is not":
		return CmpOpIsNot, true
	default:
		return "", false
	}
}

// QuantOperation represents the operation of a quantifier expression
type QuantOperation string

const (
	QuantOperationAll    QuantOperation = "All"
	QuantOperationAny    QuantOperation = "Any"
	QuantOperationFilter QuantOperation = "Filter"
	QuantOperationMap    QuantOperation = "Map"
)

// String returns the string representation of the QuantOperation
func (qo QuantOperation) String() string {
	return string(qo)
}

// AllQuantOperations returns all possible QuantOperation values
func AllQuantOperations() []QuantOperation {
	return []QuantOperation{
		QuantOperationAll,
		QuantOperationAny,
		QuantOperationFilter,
		QuantOperationMap,
	}
}

// QuantOperationFromString returns the QuantOperation corresponding to the given string
func QuantOperationFromString(s string) (QuantOperation, bool) {
	switch s {
	case "All":
		return QuantOperationAll, true
	case "Any":
		return QuantOperationAny, true
	case "Filter":
		return QuantOperationFilter, true
	case "Map":
		return QuantOperationMap, true
	default:
		return "", false
	}
}

// ConfigEntryOperation represents the operation of a configuration entry
type ConfigEntryOperation string

const (
	ConfigEntryOperationUnion    ConfigEntryOperation = "Union"
	ConfigEntryOperationOverride ConfigEntryOperation = "Override"
	ConfigEntryOperationInsert   ConfigEntryOperation = "Insert"
)

// String returns the string representation of the ConfigEntryOperation
func (c ConfigEntryOperation) String() string {
	return string(c)
}

// Value returns the integer value of the ConfigEntryOperation
func (c ConfigEntryOperation) Value() int {
	switch c {
	case ConfigEntryOperationUnion:
		return 0
	case ConfigEntryOperationOverride:
		return 1
	case ConfigEntryOperationInsert:
		return 2
	default:
		panic(fmt.Sprintf("unknown operation: %v", c))
	}
}

// Symbol returns the symbol representation of the ConfigEntryOperation
func (c ConfigEntryOperation) Symbol() string {
	switch c {
	case ConfigEntryOperationUnion:
		return ":"
	case ConfigEntryOperationOverride:
		return "="
	case ConfigEntryOperationInsert:
		return "+="
	default:
		panic(fmt.Sprintf("unknown operation: %v", c))
	}
}

// ConfigEntryOperationFromString returns the ConfigEntryOperation corresponding to the given string
func ConfigEntryOperationFromString(s string) (ConfigEntryOperation, error) {
	switch s {
	case "Union":
		return ConfigEntryOperationUnion, nil
	case "Override":
		return ConfigEntryOperationOverride, nil
	case "Insert":
		return ConfigEntryOperationInsert, nil
	default:
		return ConfigEntryOperation("Unknown"), fmt.Errorf("unknown ConfigEntryOperation: %s", s)
	}
}

// AllConfigEntryOperations returns all possible ConfigEntryOperation values
func AllConfigEntryOperations() []ConfigEntryOperation {
	return []ConfigEntryOperation{
		ConfigEntryOperationUnion,
		ConfigEntryOperationOverride,
		ConfigEntryOperationInsert,
	}
}

// NumberBinarySuffix represents the binary suffix of a number
type NumberBinarySuffix string

const (
	NumberBinarySuffixN  NumberBinarySuffix = "n"
	NumberBinarySuffixU  NumberBinarySuffix = "u"
	NumberBinarySuffixM  NumberBinarySuffix = "m"
	NumberBinarySuffixK  NumberBinarySuffix = "k"
	NumberBinarySuffixKU NumberBinarySuffix = "K"
	NumberBinarySuffixMU NumberBinarySuffix = "M"
	NumberBinarySuffixG  NumberBinarySuffix = "G"
	NumberBinarySuffixT  NumberBinarySuffix = "T"
	NumberBinarySuffixP  NumberBinarySuffix = "P"
	NumberBinarySuffixKi NumberBinarySuffix = "Ki"
	NumberBinarySuffixMi NumberBinarySuffix = "Mi"
	NumberBinarySuffixGi NumberBinarySuffix = "Gi"
	NumberBinarySuffixTi NumberBinarySuffix = "Ti"
	NumberBinarySuffixPi NumberBinarySuffix = "Pi"
)

// Value returns the string representation of the NumberBinarySuffix
func (n NumberBinarySuffix) Value() string {
	return string(n)
}

// AllNumberBinarySuffixes returns all possible NumberBinarySuffix values
func AllNumberBinarySuffixes() []NumberBinarySuffix {
	return []NumberBinarySuffix{
		NumberBinarySuffixN,
		NumberBinarySuffixU,
		NumberBinarySuffixM,
		NumberBinarySuffixK,
		NumberBinarySuffixKU,
		NumberBinarySuffixMU,
		NumberBinarySuffixG,
		NumberBinarySuffixT,
		NumberBinarySuffixP,
		NumberBinarySuffixKi,
		NumberBinarySuffixMi,
		NumberBinarySuffixGi,
		NumberBinarySuffixTi,
		NumberBinarySuffixPi,
	}
}

// AllNumberBinarySuffixNames returns all names of NumberBinarySuffix
func AllNumberBinarySuffixNames() []string {
	return []string{"n", "u", "m", "k", "K", "M", "G", "T", "P", "Ki", "Mi", "Gi", "Ti", "Pi", "i"}
}

// NumberBinarySuffixFromString returns the NumberBinarySuffix corresponding to the given string
func NumberBinarySuffixFromString(s string) (NumberBinarySuffix, bool) {
	switch s {
	case "n":
		return NumberBinarySuffixN, true
	case "u":
		return NumberBinarySuffixU, true
	case "m":
		return NumberBinarySuffixM, true
	case "k":
		return NumberBinarySuffixK, true
	case "K":
		return NumberBinarySuffixKU, true
	case "M":
		return NumberBinarySuffixMU, true
	case "G":
		return NumberBinarySuffixG, true
	case "T":
		return NumberBinarySuffixT, true
	case "P":
		return NumberBinarySuffixP, true
	case "Ki":
		return NumberBinarySuffixKi, true
	case "Mi":
		return NumberBinarySuffixMi, true
	case "Gi":
		return NumberBinarySuffixGi, true
	case "Ti":
		return NumberBinarySuffixTi, true
	case "Pi":
		return NumberBinarySuffixPi, true
	default:
		return "", false
	}
}

// NameConstant represents a name constant, e.g.
//
// True
// False
// None
// Undefined
type NameConstant string

const (
	NameConstantTrue      NameConstant = "True"
	NameConstantFalse     NameConstant = "False"
	NameConstantNone      NameConstant = "None"
	NameConstantUndefined NameConstant = "Undefined"
)

// Symbol returns the symbol for each constant
func (n NameConstant) Symbol() string {
	return string(n)
}

// JSONValue returns the JSON value for each constant
func (n NameConstant) JSONValue() string {
	switch n {
	case NameConstantTrue:
		return "true"
	case NameConstantFalse:
		return "false"
	case NameConstantNone, NameConstantUndefined:
		return "null"
	default:
		panic(fmt.Sprintf("unknown NameConstant: %s", n))
	}
}

// MarshalJSON implements custom JSON marshaling for NameConstant
func (n NameConstant) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Symbol())
}

// UnmarshalJSON implements custom JSON unmarshaling for NameConstant
func (n *NameConstant) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "True":
		*n = NameConstantTrue
	case "False":
		*n = NameConstantFalse
	case "None":
		*n = NameConstantNone
	case "Undefined":
		*n = NameConstantUndefined
	default:
		return fmt.Errorf("unknown NameConstant: %s", s)
	}
	return nil
}

// AllNameConstants returns all possible NameConstant values
func AllNameConstants() []NameConstant {
	return []NameConstant{
		NameConstantTrue,
		NameConstantFalse,
		NameConstantNone,
		NameConstantUndefined,
	}
}
