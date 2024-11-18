package ast

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements custom JSON unmarshaling for Stmt
func UnmarshalStmt(data []byte) (Stmt, error) {
	var base BaseStmt

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	var stmt Stmt
	switch base.StmtType {
	case "TypeAlias":
		stmt = &TypeAliasStmt{}
	case "Expr":
		stmt = &ExprStmt{}
	case "Unification":
		stmt = &UnificationStmt{}
	case "Assign":
		stmt = &AssignStmt{}
	case "AugAssign":
		stmt = &AugAssignStmt{}
	case "Assert":
		stmt = &AssertStmt{}
	case "If":
		stmt = &IfStmt{}
	case "Import":
		stmt = &ImportStmt{}
	case "SchemaAttr":
		stmt = &SchemaAttr{}
	case "Schema":
		stmt = &SchemaStmt{}
	case "Rule":
		stmt = &RuleStmt{}
	default:
		return nil, fmt.Errorf("unknown statement type: %s", base.StmtType)
	}

	if err := json.Unmarshal(data, stmt); err != nil {
		return nil, err
	}

	return stmt, nil
}

// UnmarshalExprJSON implements custom JSON unmarshaling for Expr
func UnmarshalExpr(data []byte) (Expr, error) {
	var base BaseExpr
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	var expr Expr
	switch base.ExprType {
	case "Target":
		expr = &TargetExpr{}
	case "Identifier":
		expr = &IdentifierExpr{}
	case "Unary":
		expr = &UnaryExpr{}
	case "Binary":
		expr = &BinaryExpr{}
	case "If":
		expr = &IfExpr{}
	case "Selector":
		expr = &SelectorExpr{}
	case "Call":
		expr = &CallExpr{}
	case "Paren":
		expr = &ParenExpr{}
	case "Quant":
		expr = &QuantExpr{}
	case "List":
		expr = &ListExpr{}
	case "ListIfItem":
		expr = &ListIfItemExpr{}
	case "ListComp":
		expr = &ListComp{}
	case "Starred":
		expr = &StarredExpr{}
	case "DictComp":
		expr = &DictComp{}
	case "ConfigIfEntry":
		expr = &ConfigIfEntryExpr{}
	case "CompClause":
		expr = &CompClause{}
	case "Schema":
		expr = &SchemaExpr{}
	case "Config":
		expr = &ConfigExpr{}
	case "Lambda":
		expr = &LambdaExpr{}
	case "Subscript":
		expr = &Subscript{}
	case "Compare":
		expr = &Compare{}
	case "NumberLit":
		expr = &NumberLit{}
	case "StringLit":
		expr = &StringLit{}
	case "NameConstantLit":
		expr = &NameConstantLit{}
	case "JoinedString":
		expr = &JoinedString{}
	case "FormattedValue":
		expr = &FormattedValue{}
	case "Missing":
		expr = &MissingExpr{}
	default:
		return nil, fmt.Errorf("unknown expression type: %s", base.ExprType)
	}

	if err := json.Unmarshal(data, expr); err != nil {
		return nil, err
	}

	return expr, nil
}

// MarshalJSON implements the json.Marshaler interface
func (n *Node[Stmt]) MarshalJSON() ([]byte, error) {
	type Alias Node[Stmt]
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (n *Node[T]) UnmarshalJSON(data []byte) error {
	var objMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &objMap); err != nil {
		return err
	}

	if data, ok := objMap["id"]; ok {
		if err := json.Unmarshal(data, &n.ID); err != nil {
			return err
		}
	}
	if data, ok := objMap["filename"]; ok {
		if err := json.Unmarshal(data, &n.Pos.Filename); err != nil {
			return err
		}
	}
	if data, ok := objMap["line"]; ok {
		if err := json.Unmarshal(data, &n.Pos.Line); err != nil {
			return err
		}
	}
	if data, ok := objMap["column"]; ok {
		if err := json.Unmarshal(data, &n.Pos.Column); err != nil {
			return err
		}
	}
	if data, ok := objMap["end_line"]; ok {
		if err := json.Unmarshal(data, &n.Pos.EndLine); err != nil {
			return err
		}
	}
	if data, ok := objMap["end_column"]; ok {
		if err := json.Unmarshal(data, &n.Pos.EndColumn); err != nil {
			return err
		}
	}

	nodeData, ok := objMap["node"]
	if !ok {
		return fmt.Errorf("missing 'node' field")
	}
	if expr, err := UnmarshalExpr(nodeData); err == nil {
		if n.Node, ok = expr.(T); ok {
			return nil
		}
	}
	if stmt, err := UnmarshalStmt(nodeData); err == nil {
		if n.Node, ok = stmt.(T); ok {
			return nil
		}
	}
	if memberOrIndex, err := UnmarshalMemberOrIndex(nodeData); err == nil {
		if n.Node, ok = memberOrIndex.(T); ok {
			return nil
		}
	}
	if numberLit, err := UnmarshalNumberLitValue(nodeData); err == nil {
		if n.Node, ok = numberLit.(T); ok {
			return nil
		}
	}
	if ty, err := UnmarshalType(nodeData); err == nil {
		if n.Node, ok = ty.(T); ok {
			return nil
		}
	} else {
		otherNode := new(T)
		if err := json.Unmarshal(nodeData, otherNode); err != nil {
			return err
		}
		n.Node = *otherNode
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for AssignStmt
func (a *AssignStmt) MarshalJSON() ([]byte, error) {
	type Alias AssignStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for AssignStmt
func (a *AssignStmt) UnmarshalJSON(data []byte) error {
	type Alias AssignStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for ExprContext
func (e ExprContext) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for ExprContext
func (e *ExprContext) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Load":
		*e = Load
	case "Store":
		*e = Store
	default:
		return fmt.Errorf("invalid ExprContext: %s", s)
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for SchemaConfig
func (s *SchemaConfig) MarshalJSON() ([]byte, error) {
	type Alias SchemaConfig
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SchemaConfig
func (s *SchemaConfig) UnmarshalJSON(data []byte) error {
	type Alias SchemaConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for Keyword
func (k *Keyword) MarshalJSON() ([]byte, error) {
	type Alias Keyword
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(k),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Keyword
func (k *Keyword) UnmarshalJSON(data []byte) error {
	type Alias Keyword
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(k),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// UnmarshalMemberOrIndex is a helper function to unmarshal JSON into MemberOrIndex
func UnmarshalMemberOrIndex(data []byte) (MemberOrIndex, error) {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	var result MemberOrIndex
	switch base.Type {
	case "Member":
		result = &Member{}
	case "Index":
		result = &Index{}
	default:
		return nil, fmt.Errorf("unknown MemberOrIndex type: %s", base.Type)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}

// MarshalJSON implements custom JSON marshaling for Target
func (t *Target) MarshalJSON() ([]byte, error) {
	type Alias Target
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Target
func (t *Target) UnmarshalJSON(data []byte) error {
	type Alias Target
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Unmarshal Paths
	var rawPaths []json.RawMessage
	if err := json.Unmarshal(data, &struct {
		Paths *[]json.RawMessage `json:"paths"`
	}{Paths: &rawPaths}); err != nil {
		return err
	}

	t.Paths = make([]*MemberOrIndex, len(rawPaths))
	for i, rawPath := range rawPaths {
		path, err := UnmarshalMemberOrIndex(rawPath)
		if err != nil {
			return fmt.Errorf("error unmarshaling path %d in Target: %v", i, err)
		}
		t.Paths[i] = &path
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for AugAssignStmt
func (a *AugAssignStmt) MarshalJSON() ([]byte, error) {
	type Alias AugAssignStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for AugAssignStmt
func (a *AugAssignStmt) UnmarshalJSON(data []byte) error {
	type Alias AugAssignStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for AssertStmt
func (a *AssertStmt) MarshalJSON() ([]byte, error) {
	type Alias AssertStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for AssertStmt
func (a *AssertStmt) UnmarshalJSON(data []byte) error {
	type Alias AssertStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for IfStmt
func (i *IfStmt) MarshalJSON() ([]byte, error) {
	type Alias IfStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for IfStmt
func (i *IfStmt) UnmarshalJSON(data []byte) error {
	type Alias IfStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for TypeAliasStmt
func (t *TypeAliasStmt) MarshalJSON() ([]byte, error) {
	type Alias TypeAliasStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for TypeAliasStmt
func (t *TypeAliasStmt) UnmarshalJSON(data []byte) error {
	type Alias TypeAliasStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for ImportStmt
func (i *ImportStmt) MarshalJSON() ([]byte, error) {
	type Alias ImportStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ImportStmt
func (i *ImportStmt) UnmarshalJSON(data []byte) error {
	type Alias ImportStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for SchemaAttr
func (s *SchemaAttr) MarshalJSON() ([]byte, error) {
	type Alias SchemaAttr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SchemaAttr
func (s *SchemaAttr) UnmarshalJSON(data []byte) error {
	type Alias SchemaAttr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Decorators is initialized
	if s.Decorators == nil {
		s.Decorators = make([]*Node[Decorator], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for Decorator
func (d *Decorator) MarshalJSON() ([]byte, error) {
	type Alias Decorator
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Decorator
func (d *Decorator) UnmarshalJSON(data []byte) error {
	type Alias Decorator
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Args and Keywords are initialized
	if d.Args == nil {
		d.Args = make([]*Node[Expr], 0)
	}
	if d.Keywords == nil {
		d.Keywords = make([]*Node[Keyword], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for Arguments
func (a *Arguments) MarshalJSON() ([]byte, error) {
	type Alias Arguments
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Arguments
func (a *Arguments) UnmarshalJSON(data []byte) error {
	type Alias Arguments
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if a.Args == nil {
		a.Args = make([]*Node[Identifier], 0)
	}
	if a.Defaults == nil {
		a.Defaults = make([]*Node[Expr], 0)
	}
	if a.TyList == nil {
		a.TyList = make([]*Node[Type], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for CheckExpr
func (c *CheckExpr) MarshalJSON() ([]byte, error) {
	type Alias CheckExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for CheckExpr
func (c *CheckExpr) UnmarshalJSON(data []byte) error {
	type Alias CheckExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for SchemaStmt
func (s *SchemaStmt) MarshalJSON() ([]byte, error) {
	type Alias SchemaStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "SchemaStmt",
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SchemaStmt
func (s *SchemaStmt) UnmarshalJSON(data []byte) error {
	type Alias SchemaStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if s.Mixins == nil {
		s.Mixins = make([]*Node[Identifier], 0)
	}
	if s.Body == nil {
		s.Body = make([]*Node[Stmt], 0)
	}
	if s.Decorators == nil {
		s.Decorators = make([]*Node[Decorator], 0)
	}
	if s.Checks == nil {
		s.Checks = make([]*Node[CheckExpr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for SchemaIndexSignature
func (s *SchemaIndexSignature) MarshalJSON() ([]byte, error) {
	type Alias SchemaIndexSignature
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SchemaIndexSignature
func (s *SchemaIndexSignature) UnmarshalJSON(data []byte) error {
	type Alias SchemaIndexSignature
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for RuleStmt
func (r *RuleStmt) MarshalJSON() ([]byte, error) {
	type Alias RuleStmt
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for RuleStmt
func (r *RuleStmt) UnmarshalJSON(data []byte) error {
	type Alias RuleStmt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if r.ParentRules == nil {
		r.ParentRules = make([]*Node[Identifier], 0)
	}
	if r.Decorators == nil {
		r.Decorators = make([]*Node[Decorator], 0)
	}
	if r.Checks == nil {
		r.Checks = make([]*Node[CheckExpr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for TargetExpr
func (t *TargetExpr) MarshalJSON() ([]byte, error) {
	type Alias TargetExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for TargetExpr
func (t *TargetExpr) UnmarshalJSON(data []byte) error {
	type Alias TargetExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Paths is initialized
	if t.Paths == nil {
		t.Paths = make([]MemberOrIndex, 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for IdentifierExpr
func (i *IdentifierExpr) MarshalJSON() ([]byte, error) {
	type Alias IdentifierExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for IdentifierExpr
func (i *IdentifierExpr) UnmarshalJSON(data []byte) error {
	type Alias IdentifierExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Names is initialized
	if i.Names == nil {
		i.Names = make([]*Node[string], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for UnaryExpr
func (u *UnaryExpr) MarshalJSON() ([]byte, error) {
	type Alias UnaryExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for UnaryExpr
func (u *UnaryExpr) UnmarshalJSON(data []byte) error {
	type Alias UnaryExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for IfExpr
func (i *IfExpr) MarshalJSON() ([]byte, error) {
	type Alias IfExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for IfExpr
func (i *IfExpr) UnmarshalJSON(data []byte) error {
	type Alias IfExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for SelectorExpr
func (s *SelectorExpr) MarshalJSON() ([]byte, error) {
	type Alias SelectorExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SelectorExpr
func (s *SelectorExpr) UnmarshalJSON(data []byte) error {
	type Alias SelectorExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for CallExpr
func (c *CallExpr) MarshalJSON() ([]byte, error) {
	type Alias CallExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for CallExpr
func (c *CallExpr) UnmarshalJSON(data []byte) error {
	type Alias CallExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if c.Args == nil {
		c.Args = make([]*Node[Expr], 0)
	}
	if c.Keywords == nil {
		c.Keywords = make([]*Node[Keyword], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for ParenExpr
func (p *ParenExpr) MarshalJSON() ([]byte, error) {
	type Alias ParenExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ParenExpr
func (p *ParenExpr) UnmarshalJSON(data []byte) error {
	type Alias ParenExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for QuantExpr
func (q *QuantExpr) MarshalJSON() ([]byte, error) {
	type Alias QuantExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(q),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for QuantExpr
func (q *QuantExpr) UnmarshalJSON(data []byte) error {
	type Alias QuantExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(q),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if q.Variables == nil {
		q.Variables = make([]*Node[Identifier], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for ListExpr
func (l *ListExpr) MarshalJSON() ([]byte, error) {
	type Alias ListExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ListExpr
func (l *ListExpr) UnmarshalJSON(data []byte) error {
	type Alias ListExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if l.Elts == nil {
		l.Elts = make([]*Node[Expr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for ListIfItemExpr
func (l *ListIfItemExpr) MarshalJSON() ([]byte, error) {
	type Alias ListIfItemExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ListIfItemExpr
func (l *ListIfItemExpr) UnmarshalJSON(data []byte) error {
	type Alias ListIfItemExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if l.Exprs == nil {
		l.Exprs = make([]*Node[Expr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for ListComp
func (l *ListComp) MarshalJSON() ([]byte, error) {
	type Alias ListComp
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ListComp
func (l *ListComp) UnmarshalJSON(data []byte) error {
	type Alias ListComp
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if l.Generators == nil {
		l.Generators = make([]*Node[CompClause], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for StarredExpr
func (s *StarredExpr) MarshalJSON() ([]byte, error) {
	type Alias StarredExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for StarredExpr
func (s *StarredExpr) UnmarshalJSON(data []byte) error {
	type Alias StarredExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for DictComp
func (d *DictComp) MarshalJSON() ([]byte, error) {
	type Alias DictComp
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for DictComp
func (d *DictComp) UnmarshalJSON(data []byte) error {
	type Alias DictComp
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if d.Generators == nil {
		d.Generators = make([]*Node[CompClause], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for ConfigEntry
func (c *ConfigEntry) MarshalJSON() ([]byte, error) {
	type Alias ConfigEntry
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ConfigEntry
func (c *ConfigEntry) UnmarshalJSON(data []byte) error {
	type Alias ConfigEntry
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for ConfigIfEntryExpr
func (c *ConfigIfEntryExpr) MarshalJSON() ([]byte, error) {
	type Alias ConfigIfEntryExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ConfigIfEntryExpr
func (c *ConfigIfEntryExpr) UnmarshalJSON(data []byte) error {
	type Alias ConfigIfEntryExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if c.Items == nil {
		c.Items = make([]*Node[ConfigEntry], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for CompClause
func (c *CompClause) MarshalJSON() ([]byte, error) {
	type Alias CompClause
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for CompClause
func (c *CompClause) UnmarshalJSON(data []byte) error {
	type Alias CompClause
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if c.Targets == nil {
		c.Targets = make([]*Node[Identifier], 0)
	}
	if c.Ifs == nil {
		c.Ifs = make([]*Node[Expr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for SchemaExpr
func (s *SchemaExpr) MarshalJSON() ([]byte, error) {
	type Alias SchemaExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for SchemaExpr
func (s *SchemaExpr) UnmarshalJSON(data []byte) error {
	type Alias SchemaExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if s.Args == nil {
		s.Args = make([]*Node[Expr], 0)
	}
	if s.Kwargs == nil {
		s.Kwargs = make([]*Node[Keyword], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for ConfigExpr
func (c *ConfigExpr) MarshalJSON() ([]byte, error) {
	type Alias ConfigExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for ConfigExpr
func (c *ConfigExpr) UnmarshalJSON(data []byte) error {
	type Alias ConfigExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if c.Items == nil {
		c.Items = make([]*Node[ConfigEntry], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for LambdaExpr
func (l *LambdaExpr) MarshalJSON() ([]byte, error) {
	type Alias LambdaExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for LambdaExpr
func (l *LambdaExpr) UnmarshalJSON(data []byte) error {
	type Alias LambdaExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if l.Body == nil {
		l.Body = make([]*Node[Stmt], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for Subscript
func (s *Subscript) MarshalJSON() ([]byte, error) {
	type Alias Subscript
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Subscript
func (s *Subscript) UnmarshalJSON(data []byte) error {
	type Alias Subscript
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for Compare
func (c *Compare) MarshalJSON() ([]byte, error) {
	type Alias Compare
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Compare
func (c *Compare) UnmarshalJSON(data []byte) error {
	type Alias Compare
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slices are initialized
	if c.Ops == nil {
		c.Ops = make([]CmpOp, 0)
	}
	if c.Comparators == nil {
		c.Comparators = make([]*Node[Expr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for NumberLit
func (n *NumberLit) MarshalJSON() ([]byte, error) {
	type Alias NumberLit
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for NumberLit
func (n *NumberLit) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type         string              `json:"type"`
		BinarySuffix *NumberBinarySuffix `json:"binary_suffix,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Type != "NumberLit" {
		return fmt.Errorf("unexpected type for NumberLit: %s", aux.Type)
	}
	// Unmarshal the Value field based on its actual type
	var rawValue json.RawMessage
	if err := json.Unmarshal(data, &struct {
		Value *json.RawMessage `json:"value"`
	}{Value: &rawValue}); err != nil {
		return err
	}

	value, err := UnmarshalNumberLitValue(rawValue)
	if err != nil {
		return err
	}
	n.Value = value
	n.BinarySuffix = aux.BinarySuffix

	return nil
}

// UnmarshalNumberLitValue unmarshals JSON data into a NumberLitValue
func UnmarshalNumberLitValue(data []byte) (NumberLitValue, error) {
	var temp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	var result NumberLitValue
	switch temp.Type {
	case "Int":
		var intValue IntNumberLitValue
		if err := json.Unmarshal(data, &intValue); err != nil {
			return nil, err
		}
		result = &intValue
	case "Float":
		var floatValue FloatNumberLitValue
		if err := json.Unmarshal(data, &floatValue); err != nil {
			return nil, err
		}
		result = &floatValue
	default:
		return nil, fmt.Errorf("unknown NumberLitValue type: %s", temp.Type)
	}

	return result, nil
}

// MarshalJSON implements custom JSON marshaling for NumberLitValue
func MarshalNumberLitValue(v NumberLitValue) ([]byte, error) {
	switch value := v.(type) {
	case *IntNumberLitValue:
		return json.Marshal(struct {
			Type  string `json:"type"`
			Value int64  `json:"value"`
		}{
			Type:  "Int",
			Value: value.Value,
		})
	case *FloatNumberLitValue:
		return json.Marshal(struct {
			Type  string  `json:"type"`
			Value float64 `json:"value"`
		}{
			Type:  "Float",
			Value: value.Value,
		})
	default:
		return nil, fmt.Errorf("unknown NumberLitValue type: %T", v)
	}
}

// MarshalJSON implements custom JSON marshaling for StringLit
func (s *StringLit) MarshalJSON() ([]byte, error) {
	type Alias StringLit
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for StringLit
func (s *StringLit) UnmarshalJSON(data []byte) error {
	type Alias StringLit
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for NameConstantLit
func (n *NameConstantLit) MarshalJSON() ([]byte, error) {
	type Alias NameConstantLit
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for NameConstantLit
func (n *NameConstantLit) UnmarshalJSON(data []byte) error {
	type Alias NameConstantLit
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling for JoinedString
func (j *JoinedString) MarshalJSON() ([]byte, error) {
	type Alias JoinedString
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(j),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for JoinedString
func (j *JoinedString) UnmarshalJSON(data []byte) error {
	type Alias JoinedString
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(j),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure slice is initialized
	if j.Values == nil {
		j.Values = make([]*Node[Expr], 0)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for FormattedValue
func (f *FormattedValue) MarshalJSON() ([]byte, error) {
	type Alias FormattedValue
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for FormattedValue
func (f *FormattedValue) UnmarshalJSON(data []byte) error {
	type Alias FormattedValue
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON implements custom JSON marshaling for MissingExpr
func (m *MissingExpr) MarshalJSON() ([]byte, error) {
	type Alias MissingExpr
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for MissingExpr
func (m *MissingExpr) UnmarshalJSON(data []byte) error {
	type Alias MissingExpr
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	return json.Unmarshal(data, &aux)
}

// UnmarshalType is a helper function to unmarshal JSON into Type
func UnmarshalType(data []byte) (Type, error) {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	var result Type
	switch base.Type {
	case "Named":
		result = &NamedType{}
	case "Any":
		result = &AnyType{}
	case "Basic":
		result = &BasicType{}
	case "List":
		result = &ListType{}
	case "Dict":
		result = &DictType{}
	case "Union":
		result = &UnionType{}
	case "Literal":
		literalType := &LiteralType{}
		var literalBase struct {
			Value struct {
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"value"`
		}
		if err := json.Unmarshal(data, &literalBase); err != nil {
			return nil, err
		}

		var literalValue LiteralTypeValue
		switch literalBase.Value.Type {
		case "Bool":
			literalValue = new(BoolLiteralType)
		case "Int":
			literalValue = &IntLiteralType{}
		case "Float":
			literalValue = new(FloatLiteralType)
		case "Str":
			literalValue = new(StrLiteralType)
		default:
			return nil, fmt.Errorf("unknown LiteralType: %s", literalBase.Value.Type)
		}

		if err := json.Unmarshal(literalBase.Value.Value, literalValue); err != nil {
			return nil, err
		}

		literalType.Value = literalValue
		return literalType, nil
	case "Function":
		result = &FunctionType{}
	default:
		return nil, fmt.Errorf("unknown Type: %s", base.Type)
	}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}
