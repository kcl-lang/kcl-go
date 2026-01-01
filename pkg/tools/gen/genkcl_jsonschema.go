package gen

import (
	"encoding/json"
	"io"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"kcl-lang.io/kcl-go/pkg/3rdparty/jsonschema"
	"kcl-lang.io/kcl-go/pkg/logger"
	"kcl-lang.io/kcl-go/pkg/source"
)

type CastingOption int

const (
	OriginalName CastingOption = iota
	SnakeCase
	CamelCase
)

type context struct {
	imports       map[string]struct{}
	resultMap     map[string]convertResult
	paths         []string
	castingOption CastingOption
}

type convertContext struct {
	context
	rootSchema *jsonschema.Schema
	// pathObjects is used to avoid infinite loop when converting recursive schema
	// TODO: support recursive schema
	pathObjects []*jsonschema.Schema
}

type convertResult struct {
	IsSchema    bool
	Name        string
	Description string
	schema
	property
}

func convertPropertyName(name string, option CastingOption) string {
	switch option {
	case SnakeCase:
		return strcase.ToSnake(name)
	case CamelCase:
		return strcase.ToCamel(name)
	default:
		return name
	}
}

func (k *kclGenerator) genSchemaFromJsonSchema(w io.Writer, filename string, src any) error {
	code, err := source.ReadSource(filename, src)
	if err != nil {
		return err
	}
	js := &jsonschema.Schema{}
	if err = js.UnmarshalJSON(code); err != nil {
		return err
	}
	// convert json schema to kcl schema
	ctx := convertContext{
		rootSchema: js,
		context: context{
			resultMap: make(map[string]convertResult),
			imports:   make(map[string]struct{}),
			paths:     []string{},
		},
		pathObjects: []*jsonschema.Schema{},
	}
	kclSch := kclFile{}
	result := convertSchemaFromJsonSchema(&ctx, js,
		strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)))
	if result.IsSchema {
		kclSch.Schemas = append(kclSch.Schemas, result.schema)
	}
	for _, imp := range getSortedKeys(ctx.imports) {
		kclSch.Imports = append(kclSch.Imports, kImport{PkgPath: imp})
	}
	for _, key := range getSortedKeys(ctx.resultMap) {
		if ctx.resultMap[key].IsSchema {
			kclSch.Schemas = append(kclSch.Schemas, ctx.resultMap[key].schema)
		}
	}
	// Generate kcl schema code
	return k.genKcl(w, kclSch)
}

func convertSchemaFromJsonSchema(ctx *convertContext, s *jsonschema.Schema, name string) convertResult {
	// in jsonschema, type is one of True, False and Object
	// we only convert Object type
	if s.SchemaType != jsonschema.SchemaTypeObject {
		return convertResult{IsSchema: false}
	}

	// For the name of the result, we prefer $id, then name in the function parameter.
	// if none of them exists, "AnonymousType" as default
	if id, ok := s.Keywords["$id"].(*jsonschema.ID); ok {
		lastSlashIndex := strings.LastIndex(string(*id), "/")
		name = strings.Replace(string(*id)[lastSlashIndex+1:], ".json", "", -1)
	}
	if name == "" {
		name = "AnonymousType"
	}
	result := convertResult{IsSchema: false, Name: name}
	if objectExists(ctx.pathObjects, s) {
		result.Type = typePrimitive(typAny)
		return result
	}
	ctx.paths = append(ctx.paths, name)
	ctx.pathObjects = append(ctx.pathObjects, s)
	defer func() {
		ctx.paths = ctx.paths[:len(ctx.paths)-1]
		ctx.pathObjects = ctx.pathObjects[:len(ctx.pathObjects)-1]
	}()

	isArray := false
	isJsonNullType := false
	reference := ""
	typeList := typeUnion{}
	required := make(map[string]struct{})
	hasTypeKeyword := false // Track if we've seen a type keyword
	for i := 0; i < len(s.OrderedKeywords); i++ {
		k := s.OrderedKeywords[i]
		switch v := s.Keywords[k].(type) {
		case *jsonschema.Title:
		case *jsonschema.Comment:
		case *jsonschema.SchemaURI:
		case *jsonschema.ID:
		case *jsonschema.Description:
			result.Description = string(*v)
		case *jsonschema.Type:
			hasTypeKeyword = true
			if len(v.Vals) == 1 {
				switch v.Vals[0] {
				case "object":
					result.IsSchema = true
					continue
				case "array":
					isArray = true
					continue
				case "null":
					isJsonNullType = true
				}
			}
			typeList.Items = append(typeList.Items, jsonTypesToKclTypes(v.Vals))
		case *jsonschema.Items:
			if !v.Single {
				logger.GetLogger().Warningf("unsupported multiple items: %#v", v)
				break
			}
			for i, val := range v.Schemas {
				item := convertSchemaFromJsonSchema(ctx, val, "items"+strconv.Itoa(i))
				if item.IsSchema {
					ctx.resultMap[item.schema.Name] = item
					typeList.Items = append(typeList.Items, typeCustom{Name: item.schema.Name})
				} else {
					typeList.Items = append(typeList.Items, item.Type)
				}
			}
		case *jsonschema.Required:
			for _, key := range []string(*v) {
				required[key] = struct{}{}
			}
		case *jsonschema.Properties:
			result.IsSchema = true
			for _, prop := range *v {
				key := prop.Key
				val := prop.Value
				propSch := convertSchemaFromJsonSchema(ctx, val, key)
				_, propSch.Required = required[key]
				if propSch.IsSchema {
					ctx.resultMap[propSch.schema.Name] = propSch
				}
				result.Properties = append(result.Properties, propSch.property)
				if !propSch.IsSchema {
					for _, validate := range propSch.Validations {
						validate.Name = propSch.property.Name
						validate.Required = propSch.property.Required
						result.Validations = append(result.Validations, validate)
					}
				}
			}
		case *jsonschema.PatternProperties:
			result.IsSchema = true
			canConvert := true
			if result.HasIndexSignature {
				canConvert = false
				logger.GetLogger().Warningf("failed to convert patternProperties: already has index signature.")
			}
			if len(*v) != 1 {
				canConvert = false
				logger.GetLogger().Warningf("unsupported multiple patternProperties.")
			}
			result.HasIndexSignature = true
			result.IndexSignature = indexSignature{
				Type: typePrimitive(typAny),
			}
			for i, prop := range *v {
				val := prop.Schema
				propSch := convertSchemaFromJsonSchema(ctx, val, "patternProperties"+strconv.Itoa(i))
				if propSch.IsSchema {
					ctx.resultMap[propSch.schema.Name] = propSch
				}
				if canConvert {
					result.IndexSignature = indexSignature{
						Alias: "key",
						Type:  propSch.property.Type,
						Validations: []validation{
							{
								Required: true,
								Name:     "key",
								Regex:    prop.Re,
							},
						},
					}
					ctx.imports["regex"] = struct{}{}
				}
			}
		case *jsonschema.Default:
			result.HasDefault = true
			result.DefaultValue = v.Data
		case *jsonschema.Enum:
			typeList.Items = make([]typeInterface, 0, len(*v))
			for _, val := range *v {
				unmarshalledVal := any(nil)
				err := json.Unmarshal(val, &unmarshalledVal)
				if err != nil {
					logger.GetLogger().Warningf("failed to unmarshal enum value: %s", err)
					continue
				}
				typeList.Items = append(typeList.Items, typeValue{
					Value: unmarshalledVal,
				})
			}
		case *jsonschema.Const:
			unmarshalledVal := any(nil)
			err := json.Unmarshal(*v, &unmarshalledVal)
			if err != nil {
				logger.GetLogger().Warningf("failed to unmarshal const value: %s", err)
				continue
			}
			typeList.Items = []typeInterface{typeValue{Value: unmarshalledVal}}
			result.HasDefault = true
			result.DefaultValue = unmarshalledVal
			// Add const as validation only if there's also a type keyword
			// (e.g., type: string + const: "value" should generate "field == value" check)
			if hasTypeKeyword {
				_, req := required[name]
				result.Validations = append(result.Validations, validation{
					Name:       name,
					Required:   req,
					ConstValue: unmarshalledVal,
				})
			}
		case *jsonschema.Defs:
		case *jsonschema.Ref:
			refSch := v.ResolveRef(ctx.rootSchema)
			if refSch == nil || refSch.OrderedKeywords == nil {
				logger.GetLogger().Warningf("failed to resolve ref: %s", v.Reference)
				continue
			}
			schs := []*jsonschema.Schema{refSch}
			for i := 0; i < len(schs); i++ {
				sch := schs[i]
				for _, key := range sch.OrderedKeywords {
					// If not existed in the current schema, inherit from the ref schema.
					if _, ok := s.Keywords[key]; !ok {
						s.OrderedKeywords = append(s.OrderedKeywords, key)
						s.Keywords[key] = sch.Keywords[key]
					} else {
						switch v := sch.Keywords[key].(type) {
						case *jsonschema.Ref:
							refSch := v.ResolveRef(ctx.rootSchema)
							if refSch == nil || refSch.OrderedKeywords == nil {
								logger.GetLogger().Warningf("failed to resolve ref: %s, path: %s", v.Reference, strings.Join(ctx.paths, "/"))
								continue
							}
							schs = append(schs, refSch)
						case *jsonschema.Properties:
							props := *s.Keywords[key].(*jsonschema.Properties)
							for _, p := range *v {
								if r, _ := props.Get(p.Key); r == nil {
									props = append(props, p)
								}
							}
							s.Keywords[key] = &props
						case *jsonschema.AdditionalProperties:
							prop := *s.Keywords[key].(*jsonschema.AdditionalProperties)
							s.Keywords[key] = &prop
						case *jsonschema.PropertyNames:
							prop := *s.Keywords[key].(*jsonschema.PropertyNames)
							s.Keywords[key] = &prop
						case *jsonschema.Required:
							reqs := *s.Keywords[key].(*jsonschema.Required)
							reqs = append(*v, reqs...)
							s.Keywords[key] = &reqs
						case *jsonschema.Items:
							items := *s.Keywords[key].(*jsonschema.Items)
							items.Schemas = append(v.Schemas, items.Schemas...)
							s.Keywords[key] = &items
						default:
							logger.GetLogger().Warningf("failed to merge ref: unsupported keyword %s in ref, path: %s", key, strings.Join(ctx.paths, "/"))
						}
					}
				}
			}
			reference = v.Reference
			sort.SliceStable(s.OrderedKeywords[i+1:], func(i, j int) bool {
				return jsonschema.GetKeywordOrder(s.OrderedKeywords[i]) < jsonschema.GetKeywordOrder(s.OrderedKeywords[j])
			})
		case *jsonschema.AdditionalProperties:
			switch v.SchemaType {
			case jsonschema.SchemaTypeObject:
				sch := convertSchemaFromJsonSchema(ctx, (*jsonschema.Schema)(v), "additionalProperties")
				if sch.IsSchema {
					ctx.resultMap[sch.schema.Name] = sch
				}
				result.HasIndexSignature = true
				result.IndexSignature = indexSignature{
					Type: sch.Type,
				}
			case jsonschema.SchemaTypeTrue:
				result.HasIndexSignature = true
				result.IndexSignature = indexSignature{
					Type: typePrimitive(typAny),
				}
			case jsonschema.SchemaTypeFalse:
			}
		case *jsonschema.PropertyNames:
			if result.HasIndexSignature && result.IndexSignature.Alias != "" {
				var validations []validation
				for _, key := range v.OrderedKeywords {
					switch v := v.Keywords[key].(type) {
					case *jsonschema.Minimum:
						validations = append(validations, validation{
							Name:             result.IndexSignature.Alias,
							Required:         true,
							Minimum:          (*float64)(v),
							ExclusiveMinimum: false,
						})
					case *jsonschema.Maximum:
						validations = append(validations, validation{
							Name:             result.IndexSignature.Alias,
							Required:         true,
							Maximum:          (*float64)(v),
							ExclusiveMaximum: false,
						})
					case *jsonschema.ExclusiveMinimum:
						validations = append(validations, validation{
							Name:             result.IndexSignature.Alias,
							Required:         true,
							Minimum:          (*float64)(v),
							ExclusiveMinimum: true,
						})
					case *jsonschema.ExclusiveMaximum:
						validations = append(validations, validation{
							Name:             result.IndexSignature.Alias,
							Required:         true,
							Maximum:          (*float64)(v),
							ExclusiveMaximum: true,
						})
					case *jsonschema.MinLength:
						validations = append(validations, validation{
							Name:      result.IndexSignature.Alias,
							Required:  true,
							MinLength: (*int)(v),
						})
					case *jsonschema.MaxLength:
						validations = append(validations, validation{
							Name:      result.IndexSignature.Alias,
							Required:  true,
							MaxLength: (*int)(v),
						})
					case *jsonschema.Pattern:
						validations = append(validations, validation{
							Name:     result.IndexSignature.Alias,
							Required: true,
							Regex:    (*regexp.Regexp)(v),
						})
						ctx.imports["regex"] = struct{}{}
					case *jsonschema.MultipleOf:
						vInt := int(*v)
						if float64(vInt) != float64(*v) {
							logger.GetLogger().Warningf("unsupported multipleOf value: %f", *v)
							continue
						}
						result.Validations = append(result.Validations, validation{
							Name:       result.IndexSignature.Alias,
							Required:   true,
							MultiplyOf: &vInt,
						})
					case *jsonschema.UniqueItems:
						if *v {
							result.Validations = append(result.Validations, validation{
								Name:     result.IndexSignature.Alias,
								Required: true,
								Unique:   true,
							})
						}
					case *jsonschema.MinItems:
						result.Validations = append(result.Validations, validation{
							Name:      result.IndexSignature.Alias,
							Required:  true,
							MinLength: (*int)(v),
						})
					case *jsonschema.MaxItems:
						result.Validations = append(result.Validations, validation{
							Name:      result.IndexSignature.Alias,
							Required:  true,
							MaxLength: (*int)(v),
						})
					default:

					}
				}
				result.IndexSignature.Validations = append(result.IndexSignature.Validations, validations...)
			}
		case *jsonschema.Minimum:
			result.Validations = append(result.Validations, validation{
				Minimum:          (*float64)(v),
				ExclusiveMinimum: false,
			})
		case *jsonschema.Maximum:
			result.Validations = append(result.Validations, validation{
				Maximum:          (*float64)(v),
				ExclusiveMaximum: false,
			})
		case *jsonschema.ExclusiveMinimum:
			result.Validations = append(result.Validations, validation{
				Minimum:          (*float64)(v),
				ExclusiveMinimum: true,
			})
		case *jsonschema.ExclusiveMaximum:
			result.Validations = append(result.Validations, validation{
				Maximum:          (*float64)(v),
				ExclusiveMaximum: true,
			})
		case *jsonschema.MinLength:
			result.Validations = append(result.Validations, validation{
				MinLength: (*int)(v),
			})
		case *jsonschema.MaxLength:
			result.Validations = append(result.Validations, validation{
				MaxLength: (*int)(v),
			})
		case *jsonschema.Pattern:
			result.Validations = append(result.Validations, validation{
				Regex: (*regexp.Regexp)(v),
			})
			ctx.imports["regex"] = struct{}{}
		case *jsonschema.MultipleOf:
			vInt := int(*v)
			if float64(vInt) != float64(*v) {
				logger.GetLogger().Warningf("unsupported multipleOf value: %f", *v)
				continue
			}
			result.Validations = append(result.Validations, validation{
				MultiplyOf: &vInt,
			})
		case *jsonschema.UniqueItems:
			if *v {
				result.Validations = append(result.Validations, validation{
					Unique: true,
				})
			}
		case *jsonschema.MinItems:
			result.Validations = append(result.Validations, validation{
				MinLength: (*int)(v),
			})
		case *jsonschema.MaxItems:
			result.Validations = append(result.Validations, validation{
				MaxLength: (*int)(v),
			})
		case *jsonschema.OneOf:
			for i, val := range *v {
				item := convertSchemaFromJsonSchema(ctx, val, "oneOf"+strconv.Itoa(i))
				if item.IsSchema {
					ctx.resultMap[item.schema.Name] = item
					typeList.Items = append(typeList.Items, typeCustom{Name: item.schema.Name})
				} else if !item.isJsonNullType {
					typeList.Items = append(typeList.Items, item.Type)
				}
			}
		case *jsonschema.AllOf:
			schs := *v
			var validations []*validation
			_, req := required[name]
			for i := 0; i < len(schs); i++ {
				sch := schs[i]
				for _, key := range sch.OrderedKeywords {
					switch v := sch.Keywords[key].(type) {
					case *jsonschema.Minimum:
						validations = append(validations, &validation{
							Name:             name,
							Required:         req,
							Minimum:          (*float64)(v),
							ExclusiveMinimum: false,
						})
					case *jsonschema.Maximum:
						validations = append(validations, &validation{
							Name:             name,
							Required:         req,
							Maximum:          (*float64)(v),
							ExclusiveMaximum: false,
						})
					case *jsonschema.ExclusiveMinimum:
						validations = append(validations, &validation{
							Name:             name,
							Required:         req,
							Minimum:          (*float64)(v),
							ExclusiveMinimum: true,
						})
					case *jsonschema.ExclusiveMaximum:
						validations = append(validations, &validation{
							Name:             name,
							Required:         req,
							Maximum:          (*float64)(v),
							ExclusiveMaximum: true,
						})
					case *jsonschema.MinLength:
						validations = append(validations, &validation{
							Name:      name,
							Required:  req,
							MinLength: (*int)(v),
						})
					case *jsonschema.MaxLength:
						validations = append(validations, &validation{
							Name:      name,
							Required:  req,
							MaxLength: (*int)(v),
						})
					case *jsonschema.Pattern:
						validations = append(validations, &validation{
							Name:     name,
							Required: req,
							Regex:    (*regexp.Regexp)(v),
						})
						ctx.imports["regex"] = struct{}{}
					case *jsonschema.MultipleOf:
						vInt := int(*v)
						if float64(vInt) != float64(*v) {
							logger.GetLogger().Warningf("unsupported multipleOf value: %f", *v)
							continue
						}
						result.Validations = append(result.Validations, validation{
							Name:       name,
							Required:   req,
							MultiplyOf: &vInt,
						})
					case *jsonschema.UniqueItems:
						if *v {
							result.Validations = append(result.Validations, validation{
								Name:     name,
								Required: req,
								Unique:   true,
							})
						}
					case *jsonschema.MinItems:
						result.Validations = append(result.Validations, validation{
							Name:      name,
							Required:  req,
							MinLength: (*int)(v),
						})
					case *jsonschema.MaxItems:
						result.Validations = append(result.Validations, validation{
							Name:      name,
							Required:  req,
							MaxLength: (*int)(v),
						})
					default:
						if _, ok := s.Keywords[key]; !ok {
							s.OrderedKeywords = append(s.OrderedKeywords, key)
							s.Keywords[key] = sch.Keywords[key]
						} else {
							switch v := sch.Keywords[key].(type) {
							case *jsonschema.Ref:
								refSch := v.ResolveRef(ctx.rootSchema)
								if refSch == nil || refSch.OrderedKeywords == nil {
									logger.GetLogger().Warningf("failed to resolve ref: %s", v.Reference)
									continue
								}
								schs = append(schs, refSch)
							case *jsonschema.Properties:
								props := *s.Keywords[key].(*jsonschema.Properties)
								for _, p := range *v {
									if r, _ := props.Get(p.Key); r == nil {
										props = append(props, p)
									}
								}
								s.Keywords[key] = &props
							case *jsonschema.AdditionalProperties:
								prop := *s.Keywords[key].(*jsonschema.AdditionalProperties)
								s.Keywords[key] = &prop
							case *jsonschema.PropertyNames:
								prop := *s.Keywords[key].(*jsonschema.PropertyNames)
								s.Keywords[key] = &prop
							case *jsonschema.Items:
								items := *s.Keywords[key].(*jsonschema.Items)
								items.Schemas = append(v.Schemas, items.Schemas...)
								s.Keywords[key] = &items
							case *jsonschema.Required:
								reqs := *s.Keywords[key].(*jsonschema.Required)
								reqs = append(reqs, *v...)
								s.Keywords[key] = &reqs
							default:
								logger.GetLogger().Warningf("failed to merge allOf: unsupported keyword %s in allOf, path: %s", key, strings.Join(ctx.paths, "/"))
							}
						}
					}
				}
			}
			if len(validations) > 0 {
				result.Validations = append(result.Validations, validation{
					AllOf: validations,
				})
			}
			sort.SliceStable(s.OrderedKeywords[i+1:], func(i, j int) bool {
				return jsonschema.GetKeywordOrder(s.OrderedKeywords[i]) < jsonschema.GetKeywordOrder(s.OrderedKeywords[j])
			})
		case *jsonschema.AnyOf:
			// anyOf is similar to oneOf but allows more than one schema to match
			// We treat it as a union type for type-level anyOf
			// If all schemas only contain validations (no explicit types), convert to AnyOf validation
			schs := *v
			var allValidationsOnly = true
			var anyOfValidations []*validation
			_, req := required[name]

			// Check if this is a required-constraints anyOf (e.g., "field1 or field2 is required")
			var requiredFields []string
			var hasRequiredConstraints = true
			for _, val := range schs {
				hasRequired := false
				hasOtherKeywords := false
				for _, key := range val.OrderedKeywords {
					if r, ok := val.Keywords[key].(*jsonschema.Required); ok && len(*r) == 1 {
						requiredFields = append(requiredFields, string((*r)[0]))
						hasRequired = true
					} else {
						// Check if it's a metadata keyword
						if _, ok := val.Keywords[key].(*jsonschema.Description); ok {
							continue
						}
						if _, ok := val.Keywords[key].(*jsonschema.Title); ok {
							continue
						}
						if _, ok := val.Keywords[key].(*jsonschema.Comment); ok {
							continue
						}
						hasOtherKeywords = true
					}
				}
				if !hasRequired || hasOtherKeywords {
					hasRequiredConstraints = false
					break
				}
			}

			if hasRequiredConstraints && len(requiredFields) > 0 {
				// Generate a check using 'or' operator for required fields
				// e.g., field1 or field2 or field3
				result.Validations = append(result.Validations, validation{
					Name:        "",
					Required:    false,
					AnyOfFields: requiredFields,
				})
				break
			}

			for i, val := range schs {
				// Check if this schema is validation-only (format, pattern, const, etc.)
				// const is NEVER validation-only - it goes into type union
				// format/pattern are validation-only ONLY when there's no type keyword
				hasOnlyValidation := false
				var validationType string // "format", "pattern"
				var validationValue any
				var hasTypeKeyword bool

				for _, key := range val.OrderedKeywords {
					switch v := val.Keywords[key].(type) {
					case *jsonschema.Format:
						// Mark as potential validation-only, will be confirmed below
						if !hasTypeKeyword {
							hasOnlyValidation = true
							validationType = "format"
							validationValue = v
						}
					case *jsonschema.Pattern:
						// Mark as potential validation-only, will be confirmed below
						if !hasTypeKeyword {
							hasOnlyValidation = true
							validationType = "pattern"
							validationValue = (*regexp.Regexp)(v)
						}
					case *jsonschema.Const:
						// const is NOT validation-only - goes to type union
						hasOnlyValidation = false
					case *jsonschema.Type:
						// Has type keyword, so format/pattern are not validation-only
						hasTypeKeyword = true
						hasOnlyValidation = false
					case *jsonschema.Description, *jsonschema.Title, *jsonschema.Comment:
						// Metadata keywords are allowed
						continue
					default:
						// Other keywords mean this is not validation-only
						hasOnlyValidation = false
					}
					if !hasOnlyValidation && validationType == "" && !hasTypeKeyword {
						// Break only if we haven't found a validation type
						break
					}
				}

				if hasOnlyValidation {
					switch validationType {
					case "format":
						// Convert format to regex validation
						format := string(*(validationValue.(*jsonschema.Format)))
						var regexPattern *regexp.Regexp
						switch format {
						case "hostname":
							regexPattern = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]))*$`)
						case "f5ip":
							regexPattern = regexp.MustCompile(`^[a-fA-F0-9]{1,3}\.[a-fA-F0-9]{1,3}\.[a-fA-F0-9]{1,3}\.[a-fA-F0-9]{1,3}$`)
						default:
							// For other formats, don't treat as validation-only
							hasOnlyValidation = false
						}
						if regexPattern != nil {
							anyOfValidations = append(anyOfValidations, &validation{
								Name:  name,
								Regex: regexPattern,
							})
							ctx.imports["regex"] = struct{}{}
						}
					case "pattern":
						// Use pattern directly (already a *regexp.Regexp)
						anyOfValidations = append(anyOfValidations, &validation{
							Name:  name,
							Regex: validationValue.(*regexp.Regexp),
						})
						ctx.imports["regex"] = struct{}{}
					}
				}

				if !hasOnlyValidation {
					// Process normally
					item := convertSchemaFromJsonSchema(ctx, val, "anyOf"+strconv.Itoa(i))
					if item.IsSchema {
						// Has schema definition, not a validation-only schema
						allValidationsOnly = false
						ctx.resultMap[item.schema.Name] = item
						typeList.Items = append(typeList.Items, typeCustom{Name: item.schema.Name})
					} else if !item.isJsonNullType {
						// Check if this is a validation-only schema (no type specified, only validations)
						if len(item.Validations) > 0 && (item.Type == typePrimitive(typAny) || item.Type == typePrimitive("")) {
							// This is a validation-only schema, add to AnyOf validations
							for _, v := range item.Validations {
								v.Name = name
								v.Required = req
								anyOfValidations = append(anyOfValidations, &v)
							}
						} else {
							// Has a type definition, not validation-only
							// Exception: const values (typeValue or typeUnion with only typeValue) are part of the type union, don't mark as non-validation-only
							isOnlyConstValue := false
							if _, ok := item.Type.(typeValue); ok {
								isOnlyConstValue = true
							} else if tu, ok := item.Type.(typeUnion); ok && len(tu.Items) == 1 {
								if _, ok := tu.Items[0].(typeValue); ok {
									isOnlyConstValue = true
								}
							}
							if !isOnlyConstValue {
								allValidationsOnly = false
							}
							typeList.Items = append(typeList.Items, item.Type)
						}
					}
				}
			}

			// If all schemas are validation-only, create an AnyOf validation
			if allValidationsOnly && len(anyOfValidations) > 0 {
				result.Validations = append(result.Validations, validation{
					Name:     name,
					Required: req,
					AnyOf:    anyOfValidations,
				})
			}
		case *jsonschema.Not:
			// not negates a schema validation
			var notValidation *validation
			_, req := required[name]
			for _, key := range (*v).OrderedKeywords {
				switch val := (*v).Keywords[key].(type) {
				case *jsonschema.Pattern:
					notValidation = &validation{
						Name:     name,
						Required: req,
						Regex:    (*regexp.Regexp)(val),
					}
					ctx.imports["regex"] = struct{}{}
				case *jsonschema.Minimum:
					notValidation = &validation{
						Name:             name,
						Required:         req,
						Minimum:          (*float64)(val),
						ExclusiveMinimum: false,
					}
				case *jsonschema.Maximum:
					notValidation = &validation{
						Name:             name,
						Required:         req,
						Maximum:          (*float64)(val),
						ExclusiveMaximum: false,
					}
				case *jsonschema.ExclusiveMinimum:
					notValidation = &validation{
						Name:             name,
						Required:         req,
						Minimum:          (*float64)(val),
						ExclusiveMinimum: true,
					}
				case *jsonschema.ExclusiveMaximum:
					notValidation = &validation{
						Name:             name,
						Required:         req,
						Maximum:          (*float64)(val),
						ExclusiveMaximum: true,
					}
				case *jsonschema.MinLength:
					notValidation = &validation{
						Name:      name,
						Required:  req,
						MinLength: (*int)(val),
					}
				case *jsonschema.MaxLength:
					notValidation = &validation{
						Name:      name,
						Required:  req,
						MaxLength: (*int)(val),
					}
				}
			}
			if notValidation != nil {
				result.Validations = append(result.Validations, validation{
					Name:     notValidation.Name,
					Required: notValidation.Required,
					Not:      notValidation,
				})
			}
		case *jsonschema.ReadOnly:
			// Do nothing for the readOnly keyword.
			logger.GetLogger().Infof("unsupported keyword: %s, path: %s, omit it", k, strings.Join(ctx.paths, "/"))
		case *jsonschema.Format:
			format := string(*v)
			// Determine validation name and required status
			var validationName string
			var required bool
			if len(ctx.paths) >= 2 {
				validationName = ctx.paths[len(ctx.paths)-1]
				required = result.property.Required
			} else {
				validationName = result.Name
				required = true
			}
			var regexPattern *regexp.Regexp
			switch format {
			case "date-time":
				regexPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`)
			case "email":
				regexPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			case "hostname":
				regexPattern = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]))*$`)
			case "ipv4":
				regexPattern = regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
			case "ipv6":
				regexPattern = regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])))$`)
			case "uri":
				regexPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+-.]*://[^/?#]+(?:/[^?#]*)?(?:\?[^#]*)?(?:#.*)?$`)
			case "uuid":
				regexPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
			default:
				logger.GetLogger().Warningf("unsupported format: %s, path: %s", format, strings.Join(ctx.paths, "/"))
				regexPattern = nil
			}
			if regexPattern != nil {
				result.Validations = append(result.Validations, validation{
					Name:     validationName,
					Required: required,
					Regex:    regexPattern,
				})
				result.Type = typePrimitive(typStr)
				ctx.imports["regex"] = struct{}{} // Ensure regex import is included in KCL
			}
		default:
			logger.GetLogger().Warningf("unsupported keyword: %s, path: %s", k, strings.Join(ctx.paths, "/"))
		}
	}

	if result.IsSchema {
		// We use the reference schema id as the generated schema name
		if reference != "" {
			lastSlashIndex := strings.LastIndex(reference, "/")
			result.schema.Name = convertPropertyName(strings.Replace(string(reference)[lastSlashIndex+1:], ".json", "", -1), CamelCase)
		} else {
			var s strings.Builder
			for _, p := range ctx.paths {
				s.WriteString(strcase.ToCamel(p))
			}
			result.schema.Name = s.String()
		}
		result.schema.Description = result.Description
		typeList.Items = append(typeList.Items, typeCustom{Name: result.schema.Name})
		if len(result.Properties) == 0 && !result.HasIndexSignature {
			result.HasIndexSignature = true
			result.IndexSignature = indexSignature{Type: typePrimitive(typAny)}
		}
	}
	if len(typeList.Items) != 0 {
		if isArray {
			result.Type = typeArray{Items: typeList}
		} else {
			result.Type = typeList
		}
	} else {
		result.Type = typePrimitive(typAny)
	}
	result.isJsonNullType = isJsonNullType
	if result.HasIndexSignature && len(result.IndexSignature.Validations) > 0 {
		result.Validations = append(result.Validations, result.IndexSignature.Validations...)
	}
	// Update AllOf validation required fields
	for i := range result.Validations {
		for j := range result.Validations[i].AllOf {
			result.Validations[i].AllOf[j].Name = result.Validations[i].Name
			result.Validations[i].AllOf[j].Required = result.Validations[i].Required
		}
	}

	result.property.Name = convertPropertyName(result.Name, ctx.castingOption)
	result.property.Description = result.Description
	return result
}

func jsonTypesToKclTypes(t []string) typeInterface {
	var kclTypes typeUnion
	for _, v := range t {
		// Skip the `type | null` format.
		if v != "null" {
			kclTypes.Items = append(kclTypes.Items, jsonTypeToKclType(v))
		}
	}
	// If no any items in the union types, return the `any` type.
	if len(kclTypes.Items) == 0 {
		return typePrimitive(typAny)
	}
	return kclTypes
}

func jsonTypeToKclType(t string) typeInterface {
	switch t {
	case "string":
		return typePrimitive(typStr)
	case "boolean", "bool":
		return typePrimitive(typBool)
	case "integer":
		return typePrimitive(typInt)
	case "number":
		return typePrimitive(typFloat)
	case "array":
		return typeArray{Items: typePrimitive(typAny)}
	case "object":
		return typePrimitive(typAny)
	case "null":
		return typePrimitive(typAny)
	default:
		logger.GetLogger().Warningf("unknown type: %s, use the any type", t)
		return typePrimitive(typAny)
	}
}

func objectExists(objs []*jsonschema.Schema, obj *jsonschema.Schema) bool {
	for _, o := range objs {
		if reflect.DeepEqual(o, obj) {
			return true
		}
	}
	return false
}
