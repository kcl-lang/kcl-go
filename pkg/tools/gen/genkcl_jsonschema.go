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
					// If a property is required and is of object type (becomes a
					// sub-schema), default-initialize it with an empty instance of
					// the sub-schema. This makes `Parent{}` valid when every
					// required inner field has its own default, instead of failing
					// with `attribute 'X' of Parent is required`. When some inner
					// required field has no default, the failure is preserved with
					// a more precise error pointing at the inner schema.
					if propSch.Required && !propSch.property.HasDefault {
						propSch.property.HasDefault = true
						propSch.property.DefaultValue = schemaInstantiation{SchemaName: propSch.schema.Name}
					}
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
								existing, found := props.Get(p.Key)
								if !found {
									props = append(props, p)
								} else if len(existing.Keywords) == 0 {
									// Replace empty placeholder schema with the
									// more specific schema from the allOf member.
									for i := range props {
										if props[i].Key == p.Key {
											props[i].Value = p.Value
											break
										}
									}
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
						case *jsonschema.Type:
							// type already set; keep existing.
						default:
							logger.GetLogger().Warningf("failed to merge ref: unsupported keyword %s in ref, path: %s", key, strings.Join(ctx.paths, "/"))
						}
					}
				}
			}
			reference = v.Reference
			remaining := s.OrderedKeywords[i+1:]
			sort.SliceStable(remaining, func(a, b int) bool {
				return jsonschema.GetKeywordOrder(remaining[a]) < jsonschema.GetKeywordOrder(remaining[b])
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
					case *jsonschema.Ref:
						refSch := v.ResolveRef(ctx.rootSchema)
						if refSch == nil || refSch.OrderedKeywords == nil {
							logger.GetLogger().Warningf("failed to resolve ref: %s", v.Reference)
							continue
						}
						schs = append(schs, refSch)
					default:
						if _, ok := s.Keywords[key]; !ok {
							s.OrderedKeywords = append(s.OrderedKeywords, key)
							s.Keywords[key] = sch.Keywords[key]
						} else {
							switch v := sch.Keywords[key].(type) {
							case *jsonschema.Properties:
								props := *s.Keywords[key].(*jsonschema.Properties)
								for _, p := range *v {
									existing, found := props.Get(p.Key)
									if !found {
										props = append(props, p)
									} else if len(existing.Keywords) == 0 {
										// Replace empty placeholder schema with the
										// more specific schema from the allOf member.
										for i := range props {
											if props[i].Key == p.Key {
												props[i].Value = p.Value
												break
											}
										}
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
							case *jsonschema.Type:
								// Multiple allOf members commonly all declare "type": "object".
								// The type was already stored from the first member; keep it and
								// do not warn — having duplicate type declarations is not an error.
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
			remaining := s.OrderedKeywords[i+1:]
			sort.SliceStable(remaining, func(a, b int) bool {
				return jsonschema.GetKeywordOrder(remaining[a]) < jsonschema.GetKeywordOrder(remaining[b])
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
		case *jsonschema.If:
			// `if` drives the whole if/then/else triple, so read its two
			// companion keywords here rather than in their own cases. This
			// mirrors how the jsonschema package itself evaluates the triple
			// (see keywords_conditional.go).
			_, req := required[name]
			// A `type`-only `if` subschema (e.g. `if: {type: "string"}` on a
			// property declared as `type: ["string", "null"]`) narrows which
			// values the branches apply to. That typeof guard is only sound
			// when the enclosing schema itself declares a union type, so
			// pass that fact down.
			allowTypeGuard := false
			if t, ok := s.Keywords["type"].(*jsonschema.Type); ok && len(t.Vals) > 1 {
				allowTypeGuard = true
			}
			cond := jsonSchemaSubConstraints(ctx, (*jsonschema.Schema)(v), name, req, required, allowTypeGuard)
			if condExpr(cond).Expr == "" {
				// Nothing in the `if` subschema maps onto a KCL expression, so
				// the guard cannot be expressed and the branches would silently
				// apply unconditionally. Skip the triple instead.
				logger.GetLogger().Warningf("unsupported if condition, path: %s", strings.Join(ctx.paths, "/"))
				continue
			}
			var thenVal, elseVal *validation
			if thenKW, ok := s.Keywords["then"]; ok {
				if t, ok := thenKW.(*jsonschema.Then); ok {
					thenVal = jsonSchemaSubConstraints(ctx, (*jsonschema.Schema)(t), name, req, required, false)
				}
			}
			if elseKW, ok := s.Keywords["else"]; ok {
				if e, ok := elseKW.(*jsonschema.Else); ok {
					elseVal = jsonSchemaSubConstraints(ctx, (*jsonschema.Schema)(e), name, req, required, false)
				}
			}
			if condExpr(thenVal).Expr == "" && condExpr(elseVal).Expr == "" {
				// A bare `if` without usable branches constrains nothing.
				continue
			}
			result.Validations = append(result.Validations, validation{
				Name:     name,
				Required: req,
				IfCond:   cond,
				Then:     thenVal,
				Else:     elseVal,
			})
		case *jsonschema.Then, *jsonschema.Else:
			// Consumed by the *jsonschema.If case above.
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

// jsonSchemaSubConstraints collapses the scalar constraint keywords of an
// `if`/`then`/`else` subschema into a single validation targeting `name`.
// Every constraint it collects is later joined with `and` by condExpr.
//
// Only keywords with a direct KCL expression equivalent are collected;
// annotation-only keywords such as `default` carry no validation meaning and
// are skipped. A nil return means the subschema constrains nothing usable.
//
// Unlike the top-level conversion path, this helper also descends into a
// `properties` keyword and folds the first single-value constraint of each
// property (e.g. `country: { const: "US" }`) into the produced expression.
// That lets `if`/`then`/`else` reference the values of individual fields of
// the surrounding schema, which is the typical shape in real-world JSON
// Schemas.
//
// parentRequired lists every field name the surrounding schema declared as
// required; it is used to mark each sub-constraint so that condExpr can
// guard optional-field references and avoid running comparisons against an
// undefined value at runtime.
//
// allowTypeGuard permits a `type` keyword in the subschema to become a
// `typeof(field) == "<type>"` guard. It is only set for the `if` part of a
// triple whose target field declares a union type, where the guard selects
// the union member the branches apply to.
func jsonSchemaSubConstraints(ctx *convertContext, sch *jsonschema.Schema, name string, required bool, parentRequired map[string]struct{}, allowTypeGuard bool) *validation {
	if sch == nil {
		return nil
	}
	v := &validation{Name: name, Required: required}
	for _, key := range sch.OrderedKeywords {
		switch kw := sch.Keywords[key].(type) {
		case *jsonschema.Type:
			// A single-typed `if` subschema on a property whose JSON Schema
			// type is a union (e.g. `type: ["string", "null"]` together with
			// `if: {type: "string"}`) narrows which values the branches
			// apply to; it becomes `typeof(field) == "str"`. Without the
			// guard the `then` constraints would also run against None,
			// where e.g. `len(None)` fails at runtime.
			if allowTypeGuard && len(kw.Vals) == 1 {
				if tn := jsonTypeToKclTypeName(kw.Vals[0]); tn != "" {
					v.TypeName = tn
				}
			}
			// Otherwise we deliberately do not emit a
			// `typeof(...) == "<type>"` guard. The field already has a
			// concrete KCL type from the property declaration, so a value
			// assigned to it can only be that type (or assignment fails
			// earlier). Adding a typeof check would spuriously fail when
			// the user writes `price = 500` (typeof returns "int" for
			// integer literals even though the field is float). Constraint
			// keywords that do not have a corresponding field declaration
			// in KCL — e.g. a brand new property introduced inside `then`
			// — still get checked by the surrounding schema declaration in
			// the parent schema.
		case *jsonschema.Const:
			var constVal any
			if err := json.Unmarshal(*kw, &constVal); err == nil {
				v.ConstValue = constVal
			}
		case *jsonschema.Minimum:
			v.Minimum = (*float64)(kw)
		case *jsonschema.Maximum:
			v.Maximum = (*float64)(kw)
		case *jsonschema.ExclusiveMinimum:
			v.Minimum = (*float64)(kw)
			v.ExclusiveMinimum = true
		case *jsonschema.ExclusiveMaximum:
			v.Maximum = (*float64)(kw)
			v.ExclusiveMaximum = true
		case *jsonschema.MinLength:
			v.MinLength = (*int)(kw)
		case *jsonschema.MaxLength:
			v.MaxLength = (*int)(kw)
		case *jsonschema.Pattern:
			v.Regex = (*regexp.Regexp)(kw)
			ctx.imports["regex"] = struct{}{}
		case *jsonschema.MultipleOf:
			i := int(*kw)
			if float64(i) == float64(*kw) {
				v.MultiplyOf = &i
			}
		case *jsonschema.Properties:
			// JSON Schema `if`/`then`/`else` subschemas commonly constrain
			// fields of the surrounding schema via `properties.<field>`. We
			// need to descend into them so the rendered KCL expression
			// refers to the field values, e.g. `currency == "JPY"`. The
			// required flag of each property (looked up against the parent
			// schema's required list when known) is propagated so that
			// condExpr can guard optional-field references.
			for _, prop := range *kw {
				_, isReq := parentRequired[prop.Key]
				if propSub := jsonSchemaSubConstraints(ctx, prop.Value, prop.Key, isReq, parentRequired, false); propSub != nil {
					if condExpr(propSub).Expr != "" {
						v.SubConstraints = append(v.SubConstraints, propSub)
					}
				}
			}
		}
	}
	if condExpr(v).Expr == "" {
		return nil
	}
	return v
}

func jsonTypesToKclTypes(t []string) typeInterface {
	var kclTypes typeUnion
	for _, v := range t {
		// Skip the `type | null` format. KCL has no `None` type annotation
		// (None is not a valid union member); nullability of a value is
		// expressed by the optional `?` attribute modifier instead, and
		// `None`/`Undefined` are assignable to every type.
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

// jsonTypeToKclTypeName maps a JSON Schema type name onto the type name
// KCL's `typeof` returns for values of that type (`typeof` reports `list`
// and `dict` rather than `[any]` or `any`, so it cannot reuse
// jsonTypeToKclType). An empty string means the type has no typeof
// equivalent and no guard can be built.
func jsonTypeToKclTypeName(t string) string {
	switch t {
	case "string":
		return typStr
	case "boolean", "bool":
		return typBool
	case "integer":
		return typInt
	case "number":
		return typFloat
	case "array":
		return typList
	case "object":
		return typDict
	case "null":
		return typNone
	}
	return ""
}

func objectExists(objs []*jsonschema.Schema, obj *jsonschema.Schema) bool {
	for _, o := range objs {
		if reflect.DeepEqual(o, obj) {
			return true
		}
	}
	return false
}

// schemaInstantiation is a default value rendered as a zero-value
// instantiation of a generated sub-schema, e.g. `Inner{}`. It is used so that
// a required object-typed property (which becomes a separate KCL schema) can
// be initialized from the parent schema's default when every required inner
// field has its own default.
type schemaInstantiation struct {
	SchemaName string
}

// MarshalKcl implements the gen.Marshaler interface so that
// formatValue/walkValue emit the bytes verbatim into the generated KCL source.
func (s schemaInstantiation) MarshalKcl() ([]byte, error) {
	return []byte(s.SchemaName + "{}"), nil
}
