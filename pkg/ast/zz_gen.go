// Auto generated. DO NOT EDIT.

package ast

const PyAstMD5 = "1ce576862325454fa47ad18736de06cb" // file: ${KCLVM_SRC_ROOT}/kclvm/kcl/ast/ast.py

const (
	Arguments_TypeName            AstType = "Arguments"
	AssertStmt_TypeName           AstType = "AssertStmt"
	AssignStmt_TypeName           AstType = "AssignStmt"
	AugAssignStmt_TypeName        AstType = "AugAssignStmt"
	BasicType_TypeName            AstType = "BasicType"
	BinaryExpr_TypeName           AstType = "BinaryExpr"
	CallExpr_TypeName             AstType = "CallExpr"
	CheckExpr_TypeName            AstType = "CheckExpr"
	Comment_TypeName              AstType = "Comment"
	CommentGroup_TypeName         AstType = "CommentGroup"
	CompClause_TypeName           AstType = "CompClause"
	Compare_TypeName              AstType = "Compare"
	ConfigEntry_TypeName          AstType = "ConfigEntry"
	ConfigExpr_TypeName           AstType = "ConfigExpr"
	ConfigIfEntryExpr_TypeName    AstType = "ConfigIfEntryExpr"
	Decorator_TypeName            AstType = "Decorator"
	DictComp_TypeName             AstType = "DictComp"
	DictType_TypeName             AstType = "DictType"
	ExprStmt_TypeName             AstType = "ExprStmt"
	FormattedValue_TypeName       AstType = "FormattedValue"
	Identifier_TypeName           AstType = "Identifier"
	IfExpr_TypeName               AstType = "IfExpr"
	IfStmt_TypeName               AstType = "IfStmt"
	ImportStmt_TypeName           AstType = "ImportStmt"
	JoinedString_TypeName         AstType = "JoinedString"
	Keyword_TypeName              AstType = "Keyword"
	LambdaExpr_TypeName           AstType = "LambdaExpr"
	ListComp_TypeName             AstType = "ListComp"
	ListExpr_TypeName             AstType = "ListExpr"
	ListIfItemExpr_TypeName       AstType = "ListIfItemExpr"
	ListType_TypeName             AstType = "ListType"
	Literal_TypeName              AstType = "Literal"
	LiteralType_TypeName          AstType = "LiteralType"
	Module_TypeName               AstType = "Module"
	Name_TypeName                 AstType = "Name"
	NameConstantLit_TypeName      AstType = "NameConstantLit"
	NumberLit_TypeName            AstType = "NumberLit"
	ParenExpr_TypeName            AstType = "ParenExpr"
	QuantExpr_TypeName            AstType = "QuantExpr"
	RuleStmt_TypeName             AstType = "RuleStmt"
	SchemaAttr_TypeName           AstType = "SchemaAttr"
	SchemaExpr_TypeName           AstType = "SchemaExpr"
	SchemaIndexSignature_TypeName AstType = "SchemaIndexSignature"
	SchemaStmt_TypeName           AstType = "SchemaStmt"
	SelectorExpr_TypeName         AstType = "SelectorExpr"
	StarredExpr_TypeName          AstType = "StarredExpr"
	StringLit_TypeName            AstType = "StringLit"
	Subscript_TypeName            AstType = "Subscript"
	Type_TypeName                 AstType = "Type"
	TypeAliasStmt_TypeName        AstType = "TypeAliasStmt"
	UnaryExpr_TypeName            AstType = "UnaryExpr"
	UnificationStmt_TypeName      AstType = "UnificationStmt"
)

func init() {
	_ = _ast_node_factory_map

	_ast_node_factory_map[Arguments_TypeName] = func() Node { return &Arguments{Meta: &Meta{AstType: Arguments_TypeName}} }
	_ast_node_factory_map[AssertStmt_TypeName] = func() Node { return &AssertStmt{Meta: &Meta{AstType: AssertStmt_TypeName}} }
	_ast_node_factory_map[AssignStmt_TypeName] = func() Node { return &AssignStmt{Meta: &Meta{AstType: AssignStmt_TypeName}} }
	_ast_node_factory_map[AugAssignStmt_TypeName] = func() Node { return &AugAssignStmt{Meta: &Meta{AstType: AugAssignStmt_TypeName}} }
	_ast_node_factory_map[BasicType_TypeName] = func() Node { return &BasicType{Meta: &Meta{AstType: BasicType_TypeName}} }
	_ast_node_factory_map[BinaryExpr_TypeName] = func() Node { return &BinaryExpr{Meta: &Meta{AstType: BinaryExpr_TypeName}} }
	_ast_node_factory_map[CallExpr_TypeName] = func() Node { return &CallExpr{Meta: &Meta{AstType: CallExpr_TypeName}} }
	_ast_node_factory_map[CheckExpr_TypeName] = func() Node { return &CheckExpr{Meta: &Meta{AstType: CheckExpr_TypeName}} }
	_ast_node_factory_map[Comment_TypeName] = func() Node { return &Comment{Meta: &Meta{AstType: Comment_TypeName}} }
	_ast_node_factory_map[CommentGroup_TypeName] = func() Node { return &CommentGroup{Meta: &Meta{AstType: CommentGroup_TypeName}} }
	_ast_node_factory_map[CompClause_TypeName] = func() Node { return &CompClause{Meta: &Meta{AstType: CompClause_TypeName}} }
	_ast_node_factory_map[Compare_TypeName] = func() Node { return &Compare{Meta: &Meta{AstType: Compare_TypeName}} }
	_ast_node_factory_map[ConfigEntry_TypeName] = func() Node { return &ConfigEntry{Meta: &Meta{AstType: ConfigEntry_TypeName}} }
	_ast_node_factory_map[ConfigExpr_TypeName] = func() Node { return &ConfigExpr{Meta: &Meta{AstType: ConfigExpr_TypeName}} }
	_ast_node_factory_map[ConfigIfEntryExpr_TypeName] = func() Node { return &ConfigIfEntryExpr{Meta: &Meta{AstType: ConfigIfEntryExpr_TypeName}} }
	_ast_node_factory_map[Decorator_TypeName] = func() Node { return &Decorator{Meta: &Meta{AstType: Decorator_TypeName}} }
	_ast_node_factory_map[DictComp_TypeName] = func() Node { return &DictComp{Meta: &Meta{AstType: DictComp_TypeName}} }
	_ast_node_factory_map[DictType_TypeName] = func() Node { return &DictType{Meta: &Meta{AstType: DictType_TypeName}} }
	_ast_node_factory_map[ExprStmt_TypeName] = func() Node { return &ExprStmt{Meta: &Meta{AstType: ExprStmt_TypeName}} }
	_ast_node_factory_map[FormattedValue_TypeName] = func() Node { return &FormattedValue{Meta: &Meta{AstType: FormattedValue_TypeName}} }
	_ast_node_factory_map[Identifier_TypeName] = func() Node { return &Identifier{Meta: &Meta{AstType: Identifier_TypeName}} }
	_ast_node_factory_map[IfExpr_TypeName] = func() Node { return &IfExpr{Meta: &Meta{AstType: IfExpr_TypeName}} }
	_ast_node_factory_map[IfStmt_TypeName] = func() Node { return &IfStmt{Meta: &Meta{AstType: IfStmt_TypeName}} }
	_ast_node_factory_map[ImportStmt_TypeName] = func() Node { return &ImportStmt{Meta: &Meta{AstType: ImportStmt_TypeName}} }
	_ast_node_factory_map[JoinedString_TypeName] = func() Node { return &JoinedString{Meta: &Meta{AstType: JoinedString_TypeName}} }
	_ast_node_factory_map[Keyword_TypeName] = func() Node { return &Keyword{Meta: &Meta{AstType: Keyword_TypeName}} }
	_ast_node_factory_map[LambdaExpr_TypeName] = func() Node { return &LambdaExpr{Meta: &Meta{AstType: LambdaExpr_TypeName}} }
	_ast_node_factory_map[ListComp_TypeName] = func() Node { return &ListComp{Meta: &Meta{AstType: ListComp_TypeName}} }
	_ast_node_factory_map[ListExpr_TypeName] = func() Node { return &ListExpr{Meta: &Meta{AstType: ListExpr_TypeName}} }
	_ast_node_factory_map[ListIfItemExpr_TypeName] = func() Node { return &ListIfItemExpr{Meta: &Meta{AstType: ListIfItemExpr_TypeName}} }
	_ast_node_factory_map[ListType_TypeName] = func() Node { return &ListType{Meta: &Meta{AstType: ListType_TypeName}} }
	_ast_node_factory_map[Literal_TypeName] = func() Node { return &Literal{Meta: &Meta{AstType: Literal_TypeName}} }
	_ast_node_factory_map[LiteralType_TypeName] = func() Node { return &LiteralType{Meta: &Meta{AstType: LiteralType_TypeName}} }
	_ast_node_factory_map[Module_TypeName] = func() Node { return &Module{Meta: &Meta{AstType: Module_TypeName}} }
	_ast_node_factory_map[Name_TypeName] = func() Node { return &Name{Meta: &Meta{AstType: Name_TypeName}} }
	_ast_node_factory_map[NameConstantLit_TypeName] = func() Node { return &NameConstantLit{Meta: &Meta{AstType: NameConstantLit_TypeName}} }
	_ast_node_factory_map[NumberLit_TypeName] = func() Node { return &NumberLit{Meta: &Meta{AstType: NumberLit_TypeName}} }
	_ast_node_factory_map[ParenExpr_TypeName] = func() Node { return &ParenExpr{Meta: &Meta{AstType: ParenExpr_TypeName}} }
	_ast_node_factory_map[QuantExpr_TypeName] = func() Node { return &QuantExpr{Meta: &Meta{AstType: QuantExpr_TypeName}} }
	_ast_node_factory_map[RuleStmt_TypeName] = func() Node { return &RuleStmt{Meta: &Meta{AstType: RuleStmt_TypeName}} }
	_ast_node_factory_map[SchemaAttr_TypeName] = func() Node { return &SchemaAttr{Meta: &Meta{AstType: SchemaAttr_TypeName}} }
	_ast_node_factory_map[SchemaExpr_TypeName] = func() Node { return &SchemaExpr{Meta: &Meta{AstType: SchemaExpr_TypeName}} }
	_ast_node_factory_map[SchemaIndexSignature_TypeName] = func() Node { return &SchemaIndexSignature{Meta: &Meta{AstType: SchemaIndexSignature_TypeName}} }
	_ast_node_factory_map[SchemaStmt_TypeName] = func() Node { return &SchemaStmt{Meta: &Meta{AstType: SchemaStmt_TypeName}} }
	_ast_node_factory_map[SelectorExpr_TypeName] = func() Node { return &SelectorExpr{Meta: &Meta{AstType: SelectorExpr_TypeName}} }
	_ast_node_factory_map[StarredExpr_TypeName] = func() Node { return &StarredExpr{Meta: &Meta{AstType: StarredExpr_TypeName}} }
	_ast_node_factory_map[StringLit_TypeName] = func() Node { return &StringLit{Meta: &Meta{AstType: StringLit_TypeName}} }
	_ast_node_factory_map[Subscript_TypeName] = func() Node { return &Subscript{Meta: &Meta{AstType: Subscript_TypeName}} }
	_ast_node_factory_map[Type_TypeName] = func() Node { return &Type{Meta: &Meta{AstType: Type_TypeName}} }
	_ast_node_factory_map[TypeAliasStmt_TypeName] = func() Node { return &TypeAliasStmt{Meta: &Meta{AstType: TypeAliasStmt_TypeName}} }
	_ast_node_factory_map[UnaryExpr_TypeName] = func() Node { return &UnaryExpr{Meta: &Meta{AstType: UnaryExpr_TypeName}} }
	_ast_node_factory_map[UnificationStmt_TypeName] = func() Node { return &UnificationStmt{Meta: &Meta{AstType: UnificationStmt_TypeName}} }
}

func (p *Arguments) GetNodeType() AstType            { return Arguments_TypeName }
func (p *Arguments) GetMeta() *Meta                  { return p.Meta }
func (p *Arguments) JSONString() string              { return json_String(p) }
func (p *Arguments) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Arguments) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *AssertStmt) GetNodeType() AstType            { return AssertStmt_TypeName }
func (p *AssertStmt) GetMeta() *Meta                  { return p.Meta }
func (p *AssertStmt) JSONString() string              { return json_String(p) }
func (p *AssertStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *AssertStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *AssignStmt) GetNodeType() AstType            { return AssignStmt_TypeName }
func (p *AssignStmt) GetMeta() *Meta                  { return p.Meta }
func (p *AssignStmt) JSONString() string              { return json_String(p) }
func (p *AssignStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *AssignStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *AugAssignStmt) GetNodeType() AstType            { return AugAssignStmt_TypeName }
func (p *AugAssignStmt) GetMeta() *Meta                  { return p.Meta }
func (p *AugAssignStmt) JSONString() string              { return json_String(p) }
func (p *AugAssignStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *AugAssignStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *BasicType) GetNodeType() AstType            { return BasicType_TypeName }
func (p *BasicType) GetMeta() *Meta                  { return p.Meta }
func (p *BasicType) JSONString() string              { return json_String(p) }
func (p *BasicType) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *BasicType) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *BinaryExpr) GetNodeType() AstType            { return BinaryExpr_TypeName }
func (p *BinaryExpr) GetMeta() *Meta                  { return p.Meta }
func (p *BinaryExpr) JSONString() string              { return json_String(p) }
func (p *BinaryExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *BinaryExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *CallExpr) GetNodeType() AstType            { return CallExpr_TypeName }
func (p *CallExpr) GetMeta() *Meta                  { return p.Meta }
func (p *CallExpr) JSONString() string              { return json_String(p) }
func (p *CallExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *CallExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *CheckExpr) GetNodeType() AstType            { return CheckExpr_TypeName }
func (p *CheckExpr) GetMeta() *Meta                  { return p.Meta }
func (p *CheckExpr) JSONString() string              { return json_String(p) }
func (p *CheckExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *CheckExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Comment) GetNodeType() AstType            { return Comment_TypeName }
func (p *Comment) GetMeta() *Meta                  { return p.Meta }
func (p *Comment) JSONString() string              { return json_String(p) }
func (p *Comment) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Comment) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *CommentGroup) GetNodeType() AstType            { return CommentGroup_TypeName }
func (p *CommentGroup) GetMeta() *Meta                  { return p.Meta }
func (p *CommentGroup) JSONString() string              { return json_String(p) }
func (p *CommentGroup) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *CommentGroup) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *CompClause) GetNodeType() AstType            { return CompClause_TypeName }
func (p *CompClause) GetMeta() *Meta                  { return p.Meta }
func (p *CompClause) JSONString() string              { return json_String(p) }
func (p *CompClause) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *CompClause) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Compare) GetNodeType() AstType            { return Compare_TypeName }
func (p *Compare) GetMeta() *Meta                  { return p.Meta }
func (p *Compare) JSONString() string              { return json_String(p) }
func (p *Compare) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Compare) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ConfigEntry) GetNodeType() AstType            { return ConfigEntry_TypeName }
func (p *ConfigEntry) GetMeta() *Meta                  { return p.Meta }
func (p *ConfigEntry) JSONString() string              { return json_String(p) }
func (p *ConfigEntry) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ConfigEntry) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ConfigExpr) GetNodeType() AstType            { return ConfigExpr_TypeName }
func (p *ConfigExpr) GetMeta() *Meta                  { return p.Meta }
func (p *ConfigExpr) JSONString() string              { return json_String(p) }
func (p *ConfigExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ConfigExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ConfigIfEntryExpr) GetNodeType() AstType            { return ConfigIfEntryExpr_TypeName }
func (p *ConfigIfEntryExpr) GetMeta() *Meta                  { return p.Meta }
func (p *ConfigIfEntryExpr) JSONString() string              { return json_String(p) }
func (p *ConfigIfEntryExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ConfigIfEntryExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Decorator) GetNodeType() AstType            { return Decorator_TypeName }
func (p *Decorator) GetMeta() *Meta                  { return p.Meta }
func (p *Decorator) JSONString() string              { return json_String(p) }
func (p *Decorator) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Decorator) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *DictComp) GetNodeType() AstType            { return DictComp_TypeName }
func (p *DictComp) GetMeta() *Meta                  { return p.Meta }
func (p *DictComp) JSONString() string              { return json_String(p) }
func (p *DictComp) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *DictComp) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *DictType) GetNodeType() AstType            { return DictType_TypeName }
func (p *DictType) GetMeta() *Meta                  { return p.Meta }
func (p *DictType) JSONString() string              { return json_String(p) }
func (p *DictType) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *DictType) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ExprStmt) GetNodeType() AstType            { return ExprStmt_TypeName }
func (p *ExprStmt) GetMeta() *Meta                  { return p.Meta }
func (p *ExprStmt) JSONString() string              { return json_String(p) }
func (p *ExprStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ExprStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *FormattedValue) GetNodeType() AstType            { return FormattedValue_TypeName }
func (p *FormattedValue) GetMeta() *Meta                  { return p.Meta }
func (p *FormattedValue) JSONString() string              { return json_String(p) }
func (p *FormattedValue) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *FormattedValue) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Identifier) GetNodeType() AstType            { return Identifier_TypeName }
func (p *Identifier) GetMeta() *Meta                  { return p.Meta }
func (p *Identifier) JSONString() string              { return json_String(p) }
func (p *Identifier) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Identifier) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *IfExpr) GetNodeType() AstType            { return IfExpr_TypeName }
func (p *IfExpr) GetMeta() *Meta                  { return p.Meta }
func (p *IfExpr) JSONString() string              { return json_String(p) }
func (p *IfExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *IfExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *IfStmt) GetNodeType() AstType            { return IfStmt_TypeName }
func (p *IfStmt) GetMeta() *Meta                  { return p.Meta }
func (p *IfStmt) JSONString() string              { return json_String(p) }
func (p *IfStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *IfStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ImportStmt) GetNodeType() AstType            { return ImportStmt_TypeName }
func (p *ImportStmt) GetMeta() *Meta                  { return p.Meta }
func (p *ImportStmt) JSONString() string              { return json_String(p) }
func (p *ImportStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ImportStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *JoinedString) GetNodeType() AstType            { return JoinedString_TypeName }
func (p *JoinedString) GetMeta() *Meta                  { return p.Meta }
func (p *JoinedString) JSONString() string              { return json_String(p) }
func (p *JoinedString) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *JoinedString) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Keyword) GetNodeType() AstType            { return Keyword_TypeName }
func (p *Keyword) GetMeta() *Meta                  { return p.Meta }
func (p *Keyword) JSONString() string              { return json_String(p) }
func (p *Keyword) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Keyword) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *LambdaExpr) GetNodeType() AstType            { return LambdaExpr_TypeName }
func (p *LambdaExpr) GetMeta() *Meta                  { return p.Meta }
func (p *LambdaExpr) JSONString() string              { return json_String(p) }
func (p *LambdaExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *LambdaExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ListComp) GetNodeType() AstType            { return ListComp_TypeName }
func (p *ListComp) GetMeta() *Meta                  { return p.Meta }
func (p *ListComp) JSONString() string              { return json_String(p) }
func (p *ListComp) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ListComp) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ListExpr) GetNodeType() AstType            { return ListExpr_TypeName }
func (p *ListExpr) GetMeta() *Meta                  { return p.Meta }
func (p *ListExpr) JSONString() string              { return json_String(p) }
func (p *ListExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ListExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ListIfItemExpr) GetNodeType() AstType            { return ListIfItemExpr_TypeName }
func (p *ListIfItemExpr) GetMeta() *Meta                  { return p.Meta }
func (p *ListIfItemExpr) JSONString() string              { return json_String(p) }
func (p *ListIfItemExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ListIfItemExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ListType) GetNodeType() AstType            { return ListType_TypeName }
func (p *ListType) GetMeta() *Meta                  { return p.Meta }
func (p *ListType) JSONString() string              { return json_String(p) }
func (p *ListType) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ListType) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Literal) GetNodeType() AstType            { return Literal_TypeName }
func (p *Literal) GetMeta() *Meta                  { return p.Meta }
func (p *Literal) JSONString() string              { return json_String(p) }
func (p *Literal) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Literal) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *LiteralType) GetNodeType() AstType            { return LiteralType_TypeName }
func (p *LiteralType) GetMeta() *Meta                  { return p.Meta }
func (p *LiteralType) JSONString() string              { return json_String(p) }
func (p *LiteralType) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *LiteralType) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Module) GetNodeType() AstType            { return Module_TypeName }
func (p *Module) GetMeta() *Meta                  { return p.Meta }
func (p *Module) JSONString() string              { return json_String(p) }
func (p *Module) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Module) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Name) GetNodeType() AstType            { return Name_TypeName }
func (p *Name) GetMeta() *Meta                  { return p.Meta }
func (p *Name) JSONString() string              { return json_String(p) }
func (p *Name) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Name) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *NameConstantLit) GetNodeType() AstType            { return NameConstantLit_TypeName }
func (p *NameConstantLit) GetMeta() *Meta                  { return p.Meta }
func (p *NameConstantLit) JSONString() string              { return json_String(p) }
func (p *NameConstantLit) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *NameConstantLit) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *NumberLit) GetNodeType() AstType            { return NumberLit_TypeName }
func (p *NumberLit) GetMeta() *Meta                  { return p.Meta }
func (p *NumberLit) JSONString() string              { return json_String(p) }
func (p *NumberLit) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *NumberLit) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *ParenExpr) GetNodeType() AstType            { return ParenExpr_TypeName }
func (p *ParenExpr) GetMeta() *Meta                  { return p.Meta }
func (p *ParenExpr) JSONString() string              { return json_String(p) }
func (p *ParenExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *ParenExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *QuantExpr) GetNodeType() AstType            { return QuantExpr_TypeName }
func (p *QuantExpr) GetMeta() *Meta                  { return p.Meta }
func (p *QuantExpr) JSONString() string              { return json_String(p) }
func (p *QuantExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *QuantExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *RuleStmt) GetNodeType() AstType            { return RuleStmt_TypeName }
func (p *RuleStmt) GetMeta() *Meta                  { return p.Meta }
func (p *RuleStmt) JSONString() string              { return json_String(p) }
func (p *RuleStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *RuleStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *SchemaAttr) GetNodeType() AstType            { return SchemaAttr_TypeName }
func (p *SchemaAttr) GetMeta() *Meta                  { return p.Meta }
func (p *SchemaAttr) JSONString() string              { return json_String(p) }
func (p *SchemaAttr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *SchemaAttr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *SchemaExpr) GetNodeType() AstType            { return SchemaExpr_TypeName }
func (p *SchemaExpr) GetMeta() *Meta                  { return p.Meta }
func (p *SchemaExpr) JSONString() string              { return json_String(p) }
func (p *SchemaExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *SchemaExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *SchemaIndexSignature) GetNodeType() AstType            { return SchemaIndexSignature_TypeName }
func (p *SchemaIndexSignature) GetMeta() *Meta                  { return p.Meta }
func (p *SchemaIndexSignature) JSONString() string              { return json_String(p) }
func (p *SchemaIndexSignature) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *SchemaIndexSignature) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *SchemaStmt) GetNodeType() AstType            { return SchemaStmt_TypeName }
func (p *SchemaStmt) GetMeta() *Meta                  { return p.Meta }
func (p *SchemaStmt) JSONString() string              { return json_String(p) }
func (p *SchemaStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *SchemaStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *SelectorExpr) GetNodeType() AstType            { return SelectorExpr_TypeName }
func (p *SelectorExpr) GetMeta() *Meta                  { return p.Meta }
func (p *SelectorExpr) JSONString() string              { return json_String(p) }
func (p *SelectorExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *SelectorExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *StarredExpr) GetNodeType() AstType            { return StarredExpr_TypeName }
func (p *StarredExpr) GetMeta() *Meta                  { return p.Meta }
func (p *StarredExpr) JSONString() string              { return json_String(p) }
func (p *StarredExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *StarredExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *StringLit) GetNodeType() AstType            { return StringLit_TypeName }
func (p *StringLit) GetMeta() *Meta                  { return p.Meta }
func (p *StringLit) JSONString() string              { return json_String(p) }
func (p *StringLit) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *StringLit) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Subscript) GetNodeType() AstType            { return Subscript_TypeName }
func (p *Subscript) GetMeta() *Meta                  { return p.Meta }
func (p *Subscript) JSONString() string              { return json_String(p) }
func (p *Subscript) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Subscript) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *Type) GetNodeType() AstType            { return Type_TypeName }
func (p *Type) GetMeta() *Meta                  { return p.Meta }
func (p *Type) JSONString() string              { return json_String(p) }
func (p *Type) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *Type) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *TypeAliasStmt) GetNodeType() AstType            { return TypeAliasStmt_TypeName }
func (p *TypeAliasStmt) GetMeta() *Meta                  { return p.Meta }
func (p *TypeAliasStmt) JSONString() string              { return json_String(p) }
func (p *TypeAliasStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *TypeAliasStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *UnaryExpr) GetNodeType() AstType            { return UnaryExpr_TypeName }
func (p *UnaryExpr) GetMeta() *Meta                  { return p.Meta }
func (p *UnaryExpr) JSONString() string              { return json_String(p) }
func (p *UnaryExpr) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *UnaryExpr) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *UnificationStmt) GetNodeType() AstType            { return UnificationStmt_TypeName }
func (p *UnificationStmt) GetMeta() *Meta                  { return p.Meta }
func (p *UnificationStmt) JSONString() string              { return json_String(p) }
func (p *UnificationStmt) JSONMap() map[string]interface{} { return JSONMap(p) }

func (p *UnificationStmt) GetPosition() (_pos, line, column int) {
	if p.Meta != nil {
		line = p.Meta.Line
		column = p.Meta.Column
		return
	}
	return
}

func (p *AssertStmt) stmt_type()           {}
func (p *AssignStmt) stmt_type()           {}
func (p *AugAssignStmt) stmt_type()        {}
func (p *ExprStmt) stmt_type()             {}
func (p *IfStmt) stmt_type()               {}
func (p *ImportStmt) stmt_type()           {}
func (p *RuleStmt) stmt_type()             {}
func (p *SchemaAttr) stmt_type()           {}
func (p *SchemaIndexSignature) stmt_type() {}
func (p *SchemaStmt) stmt_type()           {}
func (p *TypeAliasStmt) stmt_type()        {}
func (p *UnificationStmt) stmt_type()      {}

func (p *Arguments) expr_type()         {}
func (p *BinaryExpr) expr_type()        {}
func (p *CallExpr) expr_type()          {}
func (p *CheckExpr) expr_type()         {}
func (p *CompClause) expr_type()        {}
func (p *Compare) expr_type()           {}
func (p *ConfigExpr) expr_type()        {}
func (p *ConfigIfEntryExpr) expr_type() {}
func (p *Decorator) expr_type()         {}
func (p *DictComp) expr_type()          {}
func (p *FormattedValue) expr_type()    {}
func (p *Identifier) expr_type()        {}
func (p *IfExpr) expr_type()            {}
func (p *JoinedString) expr_type()      {}
func (p *Keyword) expr_type()           {}
func (p *LambdaExpr) expr_type()        {}
func (p *ListComp) expr_type()          {}
func (p *ListExpr) expr_type()          {}
func (p *ListIfItemExpr) expr_type()    {}
func (p *Literal) expr_type()           {}
func (p *Name) expr_type()              {}
func (p *NameConstantLit) expr_type()   {}
func (p *NumberLit) expr_type()         {}
func (p *ParenExpr) expr_type()         {}
func (p *QuantExpr) expr_type()         {}
func (p *SchemaExpr) expr_type()        {}
func (p *SelectorExpr) expr_type()      {}
func (p *StarredExpr) expr_type()       {}
func (p *StringLit) expr_type()         {}
func (p *Subscript) expr_type()         {}
func (p *UnaryExpr) expr_type()         {}

func (p *BasicType) type_type()   {}
func (p *DictType) type_type()    {}
func (p *ListType) type_type()    {}
func (p *LiteralType) type_type() {}
func (p *Type) type_type()        {}
