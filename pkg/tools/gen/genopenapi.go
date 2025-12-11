package gen

import (
	"errors"
	"fmt"
	htmlTmpl "html/template"

	"github.com/getkin/kin-openapi/openapi3"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"

	"os"
	"path/filepath"
	"strings"

	kcl "kcl-lang.io/kcl-go"

	api "kcl-lang.io/kcl-go/pkg/kcl"
)

// No kcl files
var ErrNoKclFiles = errors.New("No input KCL files")

const (
	ExtensionKclType        = "x-kcl-type"
	ExtensionKclDecorators  = "x-kcl-decorators"
	ExtensionKclUnionTypes  = "x-kcl-union-types"
	ExtensionKclDictKeyType = "x-kcl-dict-key-type"
)

// An additional field 'Name' is added to the original 'KclType'.
//
// 'Name' is the name of the kcl type.
//
// 'RelPath' is the relative path to the package home path.
type NamedKclType struct {
	*gpyrpc.KclType
	Name    string
	RelPath string
}

// IsSchema returns true if the type is schema.
func IsSchema(kt *NamedKclType) bool {
	return kt.Type == "schema"
}

// IsSchemaType returns true if the type is schema type.
func IsSchemaType(kt *NamedKclType) bool {
	return IsSchema(kt) && kt.SchemaName == kt.Name
}

// IsSchemaInstance returns true if the type is schema instance.
func IsSchemaInstance(kt *NamedKclType) bool {
	return IsSchema(kt) && kt.SchemaName != kt.Name
}

// IsSchemaNamed returns true if the type is schema and the name is equal to the given name.
func IsSchemaNamed(kt *NamedKclType, name string) bool {
	return IsSchema(kt) && kt.Name == name
}

// Get all schema type mapping within the root path.
func GetSchemaTypeMappingByPath(root string) (map[string]map[string]*NamedKclType, error) {
	schemaTypes := make(map[string]map[string]*NamedKclType)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			filteredTypeMap := make(map[string]*NamedKclType)
			opts := kcl.NewOption()
			schemaTypeMap, err := api.GetFullSchemaTypeMapping([]string{path}, "", *opts)
			if err != nil && err.Error() != ErrNoKclFiles.Error() {
				return err
			}
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if len(schemaTypeMap) != 0 && schemaTypeMap != nil {
				for kName, kType := range schemaTypeMap {
					kTy := &NamedKclType{
						KclType: kType,
						Name:    kName,
						RelPath: relPath,
					}
					if IsSchemaType(kTy) {
						filteredTypeMap[kName] = kTy
					}
				}
				if len(filteredTypeMap) > 0 {
					schemaTypes[relPath] = filteredTypeMap
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return schemaTypes, nil
}

// ExportOpenAPIV3Spec exports open api v3 spec of a kcl package
func ExportOpenAPIV3Spec(pkgPath string) (*openapi3.T, error) {
	s, err := ExportSwaggerV2Spec(pkgPath)
	if err != nil {
		return nil, err
	}
	return SwaggerV2ToOpenAPIV3Spec(s), nil
}

// ExportOpenAPITypeToSchema exports open api v3 schema ref from the kcl type.
func ExportOpenAPITypeToSchema(ty *KclOpenAPIType) *openapi3.SchemaRef {
	types := openapi3.Types([]string{string(ty.Type)})
	s := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        &types,
			Format:      string(ty.Format),
			Default:     ty.Default,
			Enum:        ty.GetAnyEnum(),
			ReadOnly:    ty.ReadOnly,
			Description: ty.Description,
			Properties:  make(openapi3.Schemas),
			Required:    ty.Required,
			Extensions:  ty.GetExtensionsMapping(),
		},
		Ref: ty.Ref,
	}
	for i, t := range ty.Properties {
		s.Value.Properties[i] = ExportOpenAPITypeToSchema(t)
	}
	if ty.Items != nil {
		s.Value.Items = ExportOpenAPITypeToSchema(ty.Items)
	}
	if ty.AdditionalProperties != nil {
		s.Value.AdditionalProperties = openapi3.AdditionalProperties{
			Schema: ExportOpenAPITypeToSchema(ty.AdditionalProperties),
		}
	}
	if ty.Examples != nil && len(ty.Examples) > 0 {
		s.Value.Example = ty.Examples
	}
	if len(ty.ExternalDocs) > 0 {
		s.Value.ExternalDocs = &openapi3.ExternalDocs{
			Description: ty.ExternalDocs,
		}
	}
	return s
}

// ExportSwaggerV2Spec extracts the swagger v2 representation of a
// kcl package without the external dependencies in kcl.mod
func ExportSwaggerV2Spec(path string) (*SwaggerV2Spec, error) {
	spec := &SwaggerV2Spec{
		Swagger:     "2.0",
		Definitions: make(map[string]*KclOpenAPIType),
		Paths:       map[string]interface{}{},
	}
	pkgMapping, err := GetSchemaTypeMappingByPath(path)
	if err != nil {
		return spec, err
	}
	// package path -> package
	for packagePath, p := range pkgMapping {
		// schema name -> schema type
		for _, t := range p {
			id := SchemaId(packagePath, t.KclType)
			spec.Definitions[id] = GetKclOpenAPIType(packagePath, t.KclType, false)
		}
	}
	return spec, nil
}

// SwaggerV2ToOpenAPIV3Spec converts swagger v2 spec to open api v3 spec.
func SwaggerV2ToOpenAPIV3Spec(s *SwaggerV2Spec) *openapi3.T {
	t := &openapi3.T{
		OpenAPI: "3.0",
		Info: &openapi3.Info{
			Title:       s.Info.Title,
			Description: s.Info.Description,
			Version:     s.Info.Version,
		},
		Paths: &openapi3.Paths{},
		Components: &openapi3.Components{
			Schemas: make(openapi3.Schemas),
		},
	}
	for i, d := range s.Definitions {
		t.Components.Schemas[i] = ExportOpenAPITypeToSchema(d)
	}
	return t
}

// ExportSwaggerV2SpecString exports swagger v2 spec of a kcl package
func ExportSwaggerV2SpecString(pkgPath string) (string, error) {
	spec, err := ExportSwaggerV2Spec(pkgPath)
	if err != nil {
		return "", err
	}
	return jsonString(spec), nil
}

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
	Description string `json:"description,omitempty"`
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
	Type                 SwaggerTypeName            `json:"type,omitempty"`                 // object, string, array, integer, number, bool
	Format               TypeFormat                 `json:"format,omitempty"`               // type format
	Default              string                     `json:"default,omitempty"`              // default value
	Enum                 []string                   `json:"enum,omitempty"`                 // enum values
	ReadOnly             bool                       `json:"readOnly,omitempty"`             // readonly
	Description          string                     `json:"description,omitempty"`          // description
	Properties           map[string]*KclOpenAPIType `json:"properties,omitempty"`           // schema properties
	Required             []string                   `json:"required,omitempty"`             // list of required schema property names
	Items                *KclOpenAPIType            `json:"items,omitempty"`                // list item type
	AdditionalProperties *KclOpenAPIType            `json:"additionalProperties,omitempty"` // dict value type
	Examples             map[string]KclExample      `json:"examples,omitempty"`             // examples
	ExternalDocs         string                     `json:"externalDocs,omitempty"`         // externalDocs
	Ref                  string                     `json:"ref,omitempty"`                  // reference to schema path
	BaseSchema           *KclOpenAPIType            `json:"baseSchema,omitempty"`
	ReferencedBy         []string                   `json:"referencedBy,omitempty"` // schemas referncing this schema
	*KclExtensions                                  // x-kcl- extensions
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

// KclExample defines the example code snippet of the schema
type KclExample struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
}

// KclExtensions defines all the KCL specific extensions patched to OpenAPI
type KclExtensions struct {
	XKclModelType    *XKclModelType    `json:"x-kcl-type,omitempty"`
	XKclDecorators   XKclDecorators    `json:"x-kcl-decorators,omitempty"`
	XKclUnionTypes   []*KclOpenAPIType `json:"x-kcl-union-types,omitempty"`
	XKclFunctionType *XKclFunctionType `json:"x-kcl-func-type,omitempty"`
	XKclDictKeyType  *KclOpenAPIType   `json:"x-kcl-dict-key-type,omitempty"` // dict key type
}

// XKclModelType defines the `x-kcl-type` extension
type XKclModelType struct {
	Type   string              `json:"type,omitempty"`   // schema short name
	Import *KclModelImportInfo `json:"import,omitempty"` // import information
}

type XKclFunctionType struct {
	Params   []*KclOpenAPIType `json:"params,omitempty"`
	ReturnTy *KclOpenAPIType   `json:"return_ty,omitempty"`
}

// KclModelImportInfo defines how to import the current type
type KclModelImportInfo struct {
	Package string `json:"package,omitempty"` // import package path
	Alias   string `json:"alias,omitempty"`   // import alias
}

// XKclDecorators defines the `x-kcl-decorators` extension
type XKclDecorators []*XKclDecorator

// XKclDecorator definition
type XKclDecorator struct {
	Name      string            `json:"name,omitempty"`
	Arguments []string          `json:"arguments,omitempty"`
	Keywords  map[string]string `json:"keywords,omitempty"`
}

// GetKclTypeName get the string representation of a KclOpenAPIType
func (tpe *KclOpenAPIType) GetKclTypeName(omitAny bool, addLink bool, escapeHtml bool) string {
	if tpe.Ref != "" {
		schemaId := Ref2SchemaId(tpe.Ref)
		schemaName := schemaId[strings.LastIndex(schemaId, ".")+1:]
		if addLink {
			return fmt.Sprintf("[%s](#%s)", schemaName, strings.ToLower(schemaName))
		} else {
			return schemaName
		}
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
		return fmt.Sprintf("[%s]", tpe.Items.GetKclTypeName(true, addLink, escapeHtml))
	case Object:
		if tpe.AdditionalProperties != nil {
			// dict type
			if tpe.KclExtensions.XKclDictKeyType.isAnyType() && tpe.AdditionalProperties.isAnyType() {
				return "{}"
			}
			return fmt.Sprintf("{%s:%s}", tpe.KclExtensions.XKclDictKeyType.GetKclTypeName(true, addLink, escapeHtml), tpe.AdditionalProperties.GetKclTypeName(true, addLink, escapeHtml))
		}
		if tpe.KclExtensions != nil && len(tpe.KclExtensions.XKclUnionTypes) > 0 {
			// union type
			tpes := make([]string, len(tpe.KclExtensions.XKclUnionTypes))
			for i, unionType := range tpe.KclExtensions.XKclUnionTypes {
				tpes[i] = unionType.GetKclTypeName(true, addLink, escapeHtml)
			}
			if escapeHtml {
				return strings.Join(tpes, htmlTmpl.HTMLEscapeString(" \\| "))
			} else {
				return strings.Join(tpes, " | ")
			}
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

func (tpe *KclOpenAPIType) GetExtensionsMapping() map[string]interface{} {
	m := make(map[string]interface{})
	if tpe.KclExtensions != nil {
		if tpe.XKclModelType != nil {
			m[ExtensionKclType] = tpe.XKclModelType
		}
		if tpe.XKclDecorators != nil {
			m[ExtensionKclDecorators] = tpe.XKclDecorators
		}
		if tpe.XKclUnionTypes != nil {
			m[ExtensionKclUnionTypes] = tpe.XKclUnionTypes
		}
		if tpe.XKclDictKeyType != nil {
			m[ExtensionKclDictKeyType] = tpe.XKclDictKeyType
		}
	}
	return m
}

func (tpe *KclOpenAPIType) GetAnyEnum() []interface{} {
	e := make([]interface{}, 0)
	for _, v := range tpe.Enum {
		e = append(e, v)
	}
	return e
}

func GetPkgDir(base string, pkgName string) string {
	return filepath.Join(append([]string{base}, strings.Split(pkgName, ".")...)...)
}

// GetKclOpenAPIType converts the kcl.KclType(the representation of Type in KCL API) to KclOpenAPIType(the representation of Type in KCL Open API)
func GetKclOpenAPIType(pkgPath string, from *kcl.KclType, nested bool) *KclOpenAPIType {
	var baseSchema *KclOpenAPIType

	if nested && from.BaseSchema != nil {
		baseSchema = GetKclOpenAPIType(pkgPath, from.BaseSchema, true)
		baseSchema.ReferencedBy = append(baseSchema.ReferencedBy, SchemaId(pkgPath, from))
	}

	t := KclOpenAPIType{
		Description: from.Description,
		Default:     from.Default,
		BaseSchema:  baseSchema,
	}

	// Get decorators
	decorators := from.GetDecorators()
	if len(decorators) > 0 {
		t.KclExtensions = &KclExtensions{
			XKclDecorators: make(XKclDecorators, 0),
		}
	}
	for _, d := range decorators {
		t.KclExtensions.XKclDecorators = append(t.KclExtensions.XKclDecorators, &XKclDecorator{
			Name:      d.Name,
			Arguments: d.Arguments,
			Keywords:  d.Keywords,
		})
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
		ty := GetKclOpenAPIType(pkgPath, from.Key, true)
		if t.KclExtensions == nil {
			t.KclExtensions = &KclExtensions{
				XKclDictKeyType: ty,
			}
		} else {
			t.KclExtensions.XKclDictKeyType = ty
		}
		return &t
	case typSchema:
		id := SchemaId(pkgPath, from)
		if nested {
			// for nested type reference, just return the ref object
			t.Ref = SchemaId2Ref(id)
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
		ty := &XKclModelType{
			Import: &KclModelImportInfo{
				Package: packageName,
				Alias:   filepath.Base(from.Filename),
			},
			Type: from.SchemaName,
		}
		if t.KclExtensions == nil {
			t.KclExtensions = &KclExtensions{
				XKclModelType: ty,
			}
		} else {
			t.KclExtensions.XKclModelType = ty
		}
		t.Examples = make(map[string]KclExample, len(from.GetExamples()))
		for name, example := range from.GetExamples() {
			t.Examples[name] = KclExample{
				Summary:     example.Summary,
				Description: example.Description,
				Value:       example.Value,
			}
		}
		if from.IndexSignature != nil {
			t.AdditionalProperties = GetKclOpenAPIType(pkgPath, from.IndexSignature.Val, nested)
		}
		// todo externalDocs(see also)
		return &t
	case typUnion:
		t.Type = Object
		tps := make([]*KclOpenAPIType, len(from.UnionTypes))
		for i, unionType := range from.UnionTypes {
			tps[i] = GetKclOpenAPIType(pkgPath, unionType, true)
		}
		if t.KclExtensions == nil {
			t.KclExtensions = &KclExtensions{
				XKclUnionTypes: tps,
			}
		} else {
			t.KclExtensions.XKclUnionTypes = tps
		}
		return &t
	case typFunction:
		t.Type = Object
		paramsTypes := make([]*KclOpenAPIType, len(from.Function.Params))
		for i, param := range from.Function.Params {
			paramsTypes[i] = GetKclOpenAPIType(pkgPath, param.Ty, true)
		}
		returnTy := GetKclOpenAPIType(pkgPath, from.Function.ReturnTy, true)
		if t.KclExtensions == nil {
			t.KclExtensions = &KclExtensions{
				XKclFunctionType: &XKclFunctionType{
					Params:   paramsTypes,
					ReturnTy: returnTy,
				},
			}
		} else {
			t.KclExtensions.XKclFunctionType = &XKclFunctionType{
				Params:   paramsTypes,
				ReturnTy: returnTy,
			}
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
}

// PackageName resolves the package name from the PkgPath and the PkgRoot of the type
func PackageName(pkgPath string, t *kcl.KclType) string {
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

func SchemaId2Ref(id string) string {
	return fmt.Sprintf("%s%s", oaiV2Ref, id)
}

func Ref2SchemaId(ref string) string {
	return strings.TrimPrefix(ref, oaiV2Ref)
}
