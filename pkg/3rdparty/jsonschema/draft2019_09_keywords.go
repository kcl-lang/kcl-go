package jsonschema

// LoadDraft2019_09 loads the Keywords for schema validation
// based on draft2019_09
// this is also the default keyword set loaded automatically
// if no other is loaded
func LoadDraft2019_09() {
	// core Keywords
	RegisterKeyword("$schema", NewSchemaURI)
	RegisterKeyword("$id", NewID)
	RegisterKeyword("description", NewDescription)
	RegisterKeyword("title", NewTitle)
	RegisterKeyword("$comment", NewComment)
	RegisterKeyword("examples", NewExamples)
	RegisterKeyword("readOnly", NewReadOnly)
	RegisterKeyword("writeOnly", NewWriteOnly)
	RegisterKeyword("$ref", NewRef)
	RegisterKeyword("$recursiveRef", NewRecursiveRef)
	RegisterKeyword("$anchor", NewAnchor)
	RegisterKeyword("$recursiveAnchor", NewRecursiveAnchor)
	RegisterKeyword("$defs", NewDefs)
	RegisterKeyword("definitions", NewDefs)
	RegisterKeyword("default", NewDefault)

	SetKeywordOrder("$ref", 0)
	SetKeywordOrder("$recursiveRef", 0)

	// standard Keywords
	RegisterKeyword("type", NewType)
	RegisterKeyword("enum", NewEnum)
	RegisterKeyword("const", NewConst)

	// numeric Keywords
	RegisterKeyword("multipleOf", NewMultipleOf)
	RegisterKeyword("maximum", NewMaximum)
	RegisterKeyword("exclusiveMaximum", NewExclusiveMaximum)
	RegisterKeyword("minimum", NewMinimum)
	RegisterKeyword("exclusiveMinimum", NewExclusiveMinimum)

	// string Keywords
	RegisterKeyword("maxLength", NewMaxLength)
	RegisterKeyword("minLength", NewMinLength)
	RegisterKeyword("pattern", NewPattern)

	// boolean Keywords
	RegisterKeyword("allOf", NewAllOf)
	RegisterKeyword("anyOf", NewAnyOf)
	RegisterKeyword("oneOf", NewOneOf)
	RegisterKeyword("not", NewNot)

	SetKeywordOrder("allOf", -1)

	// object Keywords
	RegisterKeyword("properties", NewProperties)
	RegisterKeyword("patternProperties", NewPatternProperties)
	RegisterKeyword("additionalProperties", NewAdditionalProperties)
	RegisterKeyword("required", NewRequired)
	RegisterKeyword("propertyNames", NewPropertyNames)
	RegisterKeyword("maxProperties", NewMaxProperties)
	RegisterKeyword("minProperties", NewMinProperties)
	RegisterKeyword("dependentSchemas", NewDependentSchemas)
	RegisterKeyword("dependentRequired", NewDependentRequired)
	RegisterKeyword("unevaluatedProperties", NewUnevaluatedProperties)

	SetKeywordOrder("properties", 2)
	SetKeywordOrder("additionalProperties", 3)
	SetKeywordOrder("unevaluatedProperties", 4)

	// array Keywords
	RegisterKeyword("items", NewItems)
	RegisterKeyword("additionalItems", NewAdditionalItems)
	RegisterKeyword("maxItems", NewMaxItems)
	RegisterKeyword("minItems", NewMinItems)
	RegisterKeyword("uniqueItems", NewUniqueItems)
	RegisterKeyword("contains", NewContains)
	RegisterKeyword("maxContains", NewMaxContains)
	RegisterKeyword("minContains", NewMinContains)
	RegisterKeyword("unevaluatedItems", NewUnevaluatedItems)

	SetKeywordOrder("maxContains", 2)
	SetKeywordOrder("minContains", 2)
	SetKeywordOrder("additionalItems", 3)
	SetKeywordOrder("unevaluatedItems", 4)

	// conditional Keywords
	RegisterKeyword("if", NewIf)
	RegisterKeyword("then", NewThen)
	RegisterKeyword("else", NewElse)

	SetKeywordOrder("then", 2)
	SetKeywordOrder("else", 2)

	//optional formats
	RegisterKeyword("format", NewFormat)
}
