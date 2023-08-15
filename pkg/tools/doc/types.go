package doc

import (
	"fmt"
	"path/filepath"
	"strings"

	kcl "kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
)

// SwaggerV2Spec defines KCL OpenAPI Spec based on Swagger v2.0
type SwaggerV2Spec struct {
	Definitions map[string]*KclOpenAPIType `json:"definitions"`
	Paths       map[string]interface{}     `json:"paths"`
	Swagger     string                     `json:"swagger"`
	Info        SpecInfo                   `json:"info"`
}

// SpecInfo defines KCL package info
type SpecInfo struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// KclOpenAPIType defines the KCL representation of SchemaObject field in Swagger v2.0.
// And the mapping between kcl type and the Kcl OpenAPI type is:
//
// ## Basic Types
//
// | KCL Type   |  KCL OpenAPI Type (format) |
//
// | ---------- | -------------------------- |
//
// |    str     |        string              |
//
// |    int     |       integer(int64)       |
//
// |   float    |       number(float)        |
//
// |   bool     |       bool                 |
//
// ## Composite Types
//
// |  KCL Types    |     KCL OpenAPI Types         |
//
// |  ----------   |   --------------------------  |
//
// |     list      |  type:array, items: itemType  |
//
// |     dict      | type: object, additionalProperties: valueType, x-kcl-dict-key-type: keyType |
//
// |     union     | type: object, x-kcl-union-types: unionTypes                   |
//
// |    schema     | type: object, properties: propertyTypes, required, x-kcl-type |
//
// | nested schema | type: object, ref: jsonRefPath |
type KclOpenAPIType struct {
	Type                 SwaggerTypeName            // object, string, array, integer, number, bool
	Format               TypeFormat                 // type format
	Default              string                     // default value
	Enum                 []string                   // enum values
	ReadOnly             bool                       // readonly
	Description          string                     // description
	Properties           map[string]*KclOpenAPIType // schema properties
	Required             []string                   // list of required schema property names
	Items                *KclOpenAPIType            // list item type
	AdditionalProperties *KclOpenAPIType            // dict value type
	Example              string                     // example
	ExternalDocs         string                     // externalDocs
	KclExtensions        *KclExtensions             // x-kcl- extensions
	Ref                  string                     // reference to schema path
}

// SwaggerTypeName defines possible values of "type" field in Swagger v2.0 spec
type SwaggerTypeName string

const (
	Object  SwaggerTypeName = "object"
	String  SwaggerTypeName = "string"
	Array   SwaggerTypeName = "array"
	Integer SwaggerTypeName = "integer"
	Number  SwaggerTypeName = "number"
	Bool    SwaggerTypeName = "bool"
)

// KclTypeName defines the KCL representation of basic types
type KclTypeName string

const (
	Str      KclTypeName = "str"
	BoolKcl  KclTypeName = "bool"
	Int      KclTypeName = "int"
	FloatKcl KclTypeName = "float"
)

// TypeFormat defines possible values of "format" field in Swagger v2.0 spec
type TypeFormat string

const (
	Int64 TypeFormat = "int64"
	Float TypeFormat = "float"
)

// KclExtensions defines all the KCL specific extensions patched to OpenAPI
type KclExtensions struct {
	XKclModelType   *XKclModelType    `json:"x-kcl-type,omitempty"`
	XKclDecorators  *XKclDecorators   `json:"x-kcl-decorators,omitempty"`
	XKclUnionTypes  []*KclOpenAPIType `json:"x-kcl-union-types,omitempty"`
	XKclDictKeyType *KclOpenAPIType   `json:"x-kcl-dict-key-type,omitempty"` // dict key type
}

// XKclModelType defines the `x-kcl-type` extension
type XKclModelType struct {
	Type   string              `json:"type,omitempty"`   // schema short name
	Import *KclModelImportInfo `json:"import,omitempty"` // import information
}

// KclModelImportInfo defines how to import the current type
type KclModelImportInfo struct {
	Package string `json:"package,omitempty"` // import package path
	Alias   string `json:"alias,omitempty"`   // import alias
}

// XKclDecorators defines the `x-kcl-decorators` extension
type XKclDecorators struct {
	Name      string
	Arguments []string
	Keywords  map[string]string
}

// GetKclTypeName get the string representation of a KclOpenAPIType
func (tpe *KclOpenAPIType) GetKclTypeName(omitAny bool) string {
	if tpe.Ref != "" {
		return tpe.Ref[strings.LastIndex(tpe.Ref, ".")+1:]
	}
	switch tpe.Type {
	case String:
		if tpe.ReadOnly {
			return tpe.Default
		}
		return string(Str)
	case Integer:
		if tpe.Format == Int64 {
			if tpe.ReadOnly {
				return tpe.Default
			}
			return string(Int)
		}
		panic(fmt.Errorf("unexpected KCL OpenAPI type and format: %s(%s)", tpe.Type, tpe.Format))
	case Number:
		if tpe.Format == Float {
			if tpe.ReadOnly {
				return tpe.Default
			}
			return string(FloatKcl)
		}
		panic(fmt.Errorf("unexpected KCL OpenAPI type and format: %s(%s)", tpe.Type, tpe.Format))
	case Bool:
		if tpe.ReadOnly {
			return tpe.Default
		}
		return string(BoolKcl)
	case Array:
		return fmt.Sprintf("[%s]", tpe.Items.GetKclTypeName(true))
	case Object:
		if tpe.AdditionalProperties != nil {
			// dict type
			if tpe.KclExtensions.XKclDictKeyType.isAnyType() && tpe.AdditionalProperties.isAnyType() {
				return "{}"
			}
			return fmt.Sprintf("{%s:%s}", tpe.KclExtensions.XKclDictKeyType.GetKclTypeName(true), tpe.AdditionalProperties.GetKclTypeName(true))
		}
		if tpe.KclExtensions != nil && len(tpe.KclExtensions.XKclUnionTypes) > 0 {
			// union type
			tpes := make([]string, len(tpe.KclExtensions.XKclUnionTypes))
			for i, unionType := range tpe.KclExtensions.XKclUnionTypes {
				tpes[i] = unionType.GetKclTypeName(true)
			}
			return strings.Join(tpes, " | ")
		}
		if tpe.isAnyType() {
			if omitAny {
				return ""
			} else {
				return "any"
			}
		}
	default:
		panic(fmt.Errorf("unexpected KCL OpenAPI type: %s", tpe.Type))
	}
	return string(tpe.Type)
}

// isAnyType checks if a KclOpenAPIType is any type
func (tpe *KclOpenAPIType) isAnyType() bool {
	return tpe.Type == Object && tpe.Properties == nil && tpe.AdditionalProperties == nil && tpe.Ref == "" && (tpe.KclExtensions == nil || tpe.KclExtensions.XKclUnionTypes == nil)
}

func (tpe *KclOpenAPIType) GetSchemaPkgDir(base string) string {
	return filepath.Join(append([]string{base}, strings.Split(tpe.KclExtensions.XKclModelType.Import.Package, ".")...)...)
}

// GetKclOpenAPIType converts the kcl.KclType(the representation of Type in KCL API) to KclOpenAPIType(the representation of Type in KCL Open API)
func GetKclOpenAPIType(from *kcl.KclType, defs map[string]*KclOpenAPIType, nested bool) *KclOpenAPIType {
	t := KclOpenAPIType{
		Description: from.Description,
		Default:     from.Default,
	}
	switch from.Type {
	case string(Int):
		t.Type = Integer
		t.Format = Int64
		return &t
	case string(FloatKcl):
		t.Type = Number
		t.Format = Float
		return &t
	case string(BoolKcl):
		t.Type = Bool
		return &t
	case string(Str):
		t.Type = String
		return &t
	case "list":
		t.Type = Array
		t.Items = GetKclOpenAPIType(from.Item, defs, true)
		return &t
	case "dict":
		t.Type = Object
		t.AdditionalProperties = GetKclOpenAPIType(from.Item, defs, true)
		t.KclExtensions = &KclExtensions{
			XKclDictKeyType: GetKclOpenAPIType(from.Key, defs, true),
		}
		return &t
	case "schema":
		id := SchemaId(from)
		if _, ok := defs[id]; ok {
			// skip converting if schema existed
			t.Ref = refPath(id)
			return &t
		}

		// resolve type and add to definitions
		defs[id] = &t
		t.Type = Object
		t.Properties = make(map[string]*KclOpenAPIType, len(from.Properties))
		for name, fromProp := range from.Properties {
			t.Properties[name] = GetKclOpenAPIType(fromProp, defs, true)
		}
		t.Required = from.Required
		packageName := from.PkgPath
		if from.PkgPath == "__main__" {
			packageName = ""
		}
		t.KclExtensions = &KclExtensions{
			XKclModelType: &XKclModelType{
				Import: &KclModelImportInfo{
					Package: packageName,
					Alias:   filepath.Base(from.Filename),
				},
				Type: from.SchemaName,
			},
		}
		// todo newT.Example = from.Examples
		// todo newT.KclExtensions.XKclDecorators = from.Decorators
		// todo externalDocs(see also)

		if nested {
			return &KclOpenAPIType{
				Description: from.Description,
				Ref:         refPath(id),
			}
		} else {
			return &t
		}
	case "union":
		t.Type = Object
		tps := make([]*KclOpenAPIType, len(from.UnionTypes))
		for i, unionType := range from.UnionTypes {
			tps[i] = GetKclOpenAPIType(unionType, defs, true)
		}
		t.KclExtensions = &KclExtensions{
			XKclUnionTypes: tps,
		}
		return &t
	case "any":
		t.Type = Object
		return &t
	default:
		if isLit, basicType, litValue := gen.IsLitType(from); isLit {
			t.ReadOnly = true
			t.Enum = []string{litValue}
			t.Default = litValue
			switch basicType {
			case string(BoolKcl):
				t.Type = Bool
				return &t
			case string(Int):
				t.Type = Integer
				t.Format = Int64
				return &t
			case string(FloatKcl):
				t.Type = Number
				t.Format = Float
				return &t
			case string(Str):
				t.Type = String
				return &t
			default:
				panic(fmt.Errorf("unexpected lit type: %s", from.Type))
			}
			return &t
		}
		panic(fmt.Errorf("unexpected KCL type: %s", from.Type))
	}
	return &t
}

func SchemaId(t *kcl.KclType) string {
	return fmt.Sprintf("%s.%s", t.PkgPath, t.SchemaName)
}
func refPath(id string) string {
	return fmt.Sprintf("#/definitions/%s", id)
}
