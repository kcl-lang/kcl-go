package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	kcl "kcl-lang.io/kcl-go"
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
/*

## basic types
	┌───────────────────────┬──────────────────────────────────────┐
	│      KCL Types        │      KCL OpenAPI Types (format)      │
	├───────────────────────┼──────────────────────────────────────┤
	│       str             │         string                       │
	├───────────────────────┼──────────────────────────────────────┤
	│       int             │        integer(int64)                │
	├───────────────────────┼──────────────────────────────────────┤
	│      float            │         number(float)                │
	├───────────────────────┼──────────────────────────────────────┤
	│       bool            │             bool                     │
	├───────────────────────┼──────────────────────────────────────┤
	│   number_multiplier   │    string(units.NumberMultiplier)    │
	└───────────────────────┴──────────────────────────────────────┘
## Composite Types
	┌───────────────────────┬───────────────────────────────────────────────────────────────────────────────┐
	│      KCL Types        │      KCL OpenAPI Types (format)                                               │
	├───────────────────────┼───────────────────────────────────────────────────────────────────────────────┤
	│       list            │   type:array, items: itemType                                                 │
	├───────────────────────┼───────────────────────────────────────────────────────────────────────────────┤
	│       dict            │   type: object, additionalProperties: valueType, x-kcl-dict-key-type: keyType │
	├───────────────────────┼───────────────────────────────────────────────────────────────────────────────┤
	│      union            │   type: object, x-kcl-union-types: unionTypes                                 │
	├───────────────────────┼───────────────────────────────────────────────────────────────────────────────┤
	│      schema           │   type: object, properties: propertyTypes, required, x-kcl-type               │
	├───────────────────────┼───────────────────────────────────────────────────────────────────────────────┤
	│    nested schema      │   type: object, ref: jsonRefPath                                              │
	└───────────────────────┴───────────────────────────────────────────────────────────────────────────────┘
*/
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

const oaiV2Ref = "#/definitions/"

// TypeFormat defines possible values of "format" field in Swagger v2.0 spec
type TypeFormat string

const (
	Int64            TypeFormat = "int64"
	Float            TypeFormat = "float"
	NumberMultiplier TypeFormat = "units.NumberMultiplier"
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
		schemaId := strings.TrimPrefix(tpe.Ref, oaiV2Ref)
		return schemaId[strings.LastIndex(schemaId, ".")+1:]
	}
	switch tpe.Type {
	case String:
		if tpe.ReadOnly {
			return tpe.Default
		}
		return typStr
	case Integer:
		if tpe.Format == Int64 {
			if tpe.ReadOnly {
				return tpe.Default
			}
			return typInt
		}
		if tpe.Format == NumberMultiplier {
			if tpe.ReadOnly {
				return tpe.Default
			}
			return string(NumberMultiplier)
		}
		panic(fmt.Errorf("unexpected KCL OpenAPI type and format: %s(%s)", tpe.Type, tpe.Format))
	case Number:
		if tpe.Format == Float {
			if tpe.ReadOnly {
				return tpe.Default
			}
			return typFloat
		}
		panic(fmt.Errorf("unexpected KCL OpenAPI type and format: %s(%s)", tpe.Type, tpe.Format))
	case Bool:
		if tpe.ReadOnly {
			return tpe.Default
		}
		return typBool
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
	return GetPkgDir(base, tpe.KclExtensions.XKclModelType.Import.Package)
}

func GetPkgDir(base string, pkgName string) string {
	return filepath.Join(append([]string{base}, strings.Split(pkgName, ".")...)...)
}

// GetKclOpenAPIType converts the kcl.KclType(the representation of Type in KCL API) to KclOpenAPIType(the representation of Type in KCL Open API)
func GetKclOpenAPIType(pkgPath string, from *kcl.KclType, nested bool) *KclOpenAPIType {
	t := KclOpenAPIType{
		Description: from.Description,
		Default:     from.Default,
	}
	switch from.Type {
	case typInt:
		t.Type = Integer
		t.Format = Int64
		return &t
	case typFloat:
		t.Type = Number
		t.Format = Float
		return &t
	case typBool:
		t.Type = Bool
		return &t
	case typStr:
		t.Type = String
		return &t
	case typAny:
		t.Type = Object
		return &t
	case typNumberMultiplier:
		t.Type = Integer
		t.Format = NumberMultiplier
		return &t
	case typList:
		t.Type = Array
		t.Items = GetKclOpenAPIType(pkgPath, from.Item, true)
		return &t
	case typDict:
		t.Type = Object
		t.AdditionalProperties = GetKclOpenAPIType(pkgPath, from.Item, true)
		t.KclExtensions = &KclExtensions{
			XKclDictKeyType: GetKclOpenAPIType(pkgPath, from.Key, true),
		}
		return &t
	case typSchema:
		id := SchemaId(pkgPath, from)
		if nested {
			// for nested type reference, just return the ref object
			t.Ref = refPath(id)
			return &t
		}
		// resolve schema type
		t.Type = Object
		t.Description = from.SchemaDoc
		t.Properties = make(map[string]*KclOpenAPIType, len(from.Properties))
		for name, fromProp := range from.Properties {
			t.Properties[name] = GetKclOpenAPIType(pkgPath, fromProp, true)
		}
		t.Required = from.Required
		packageName := PackageName(pkgPath, from)
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
		return &t
	case typUnion:
		t.Type = Object
		tps := make([]*KclOpenAPIType, len(from.UnionTypes))
		for i, unionType := range from.UnionTypes {
			tps[i] = GetKclOpenAPIType(pkgPath, unionType, true)
		}
		t.KclExtensions = &KclExtensions{
			XKclUnionTypes: tps,
		}
		return &t
	default:
		if isLit, basicType, litValue := IsLitType(from); isLit {
			t.ReadOnly = true
			t.Enum = []string{litValue}
			t.Default = litValue
			switch basicType {
			case typBool:
				t.Type = Bool
				return &t
			case typInt:
				t.Type = Integer
				t.Format = Int64
				return &t
			case typFloat:
				t.Type = Number
				t.Format = Float
				return &t
			case typStr:
				t.Type = String
				return &t
			case typNumberMultiplier:
				t.Type = Integer
				t.Format = NumberMultiplier
			default:
				panic(fmt.Errorf("unexpected lit type: %s", from.Type))
			}
			return &t
		}
		panic(fmt.Errorf("unexpected KCL type: %s", from.Type))
	}
	return &t
}

// PackageName resolves the package name from the PkgPath and the PkgRoot of the type
func PackageName(pkgPath string, t *kcl.KclType) string {
	// todo after kpm support the correct pkgPath recursively in KclType, refactor the following logic
	// pkgPath is the relative path to the package root path
	// t.PkgPath is the "." joined path from the package root

	// use the resolved pkgPath instead of t.PkgPath if t.PkgPath is __main__
	if t.PkgPath == "__main__" {
		// should be resolved by pkgPath
		if pkgPath == "." {
			return ""
		} else {

			return strings.Join(strings.Split(pkgPath, string(os.PathSeparator)), ".")
		}
	}

	// use resolved t.PkgPath if t.PkgPath is not __main__
	pkgs := strings.Split(t.PkgPath, ".")
	last := pkgs[len(pkgs)-1]
	if filepath.Base(filepath.Dir(t.Filename)) == last {
		return t.PkgPath
	} else {
		return strings.Join(pkgs[:len(pkgs)-1], ".")
	}
}

func SchemaId(pkgPath string, t *kcl.KclType) string {
	pkgName := PackageName(pkgPath, t)
	if pkgName == "" {
		return t.SchemaName
	}
	return fmt.Sprintf("%s.%s", pkgName, t.SchemaName)
}

func refPath(id string) string {
	return fmt.Sprintf("%s%s", oaiV2Ref, id)
}
