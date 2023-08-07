package gen

import (
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"kcl-lang.io/kcl-go/pkg/3rdparty/jsonschema"
	"kcl-lang.io/kcl-go/pkg/logger"
)

type convertContext struct {
	imports   map[string]struct{}
	resultMap map[string]convertResult
}

type convertResult struct {
	IsSchema    bool
	Name        string
	Description string
	schema
	property
}

func (k *kclGenerator) genSchemaFromJsonSchema(w io.Writer, filename string, src interface{}) error {
	code, err := readSource(filename, src)
	if err != nil {
		return err
	}
	js := &jsonschema.Schema{}
	if err = js.UnmarshalJSON(code); err != nil {
		return err
	}
	// use Validate to trigger the evaluation of json schema
	js.Validate(context.Background(), nil)

	// convert json schema to kcl schema
	ctx := convertContext{resultMap: make(map[string]convertResult)}
	result := convertSchemaFromJsonSchema(ctx, js,
		strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)))
	if !result.IsSchema {
		panic("result is not schema")
	}
	kclSch := kclSchema{
		Imports: []string{},
		Schemas: []schema{result.schema},
	}
	for imp := range ctx.imports {
		kclSch.Imports = append(kclSch.Imports, imp)
	}
	for _, res := range ctx.resultMap {
		if res.IsSchema {
			kclSch.Schemas = append(kclSch.Schemas, res.schema)
		}
	}

	// generate kcl schema code
	return k.genKclSchema(w, kclSch)
}

func convertSchemaFromJsonSchema(ctx convertContext, s *jsonschema.Schema, name string) convertResult {
	// in jsonschema, type is one of True, False and Object
	// we only convert Object type
	if s.SchemaType != jsonschema.SchemaTypeObject {
		return convertResult{IsSchema: false}
	}

	result := convertResult{IsSchema: false, Name: name}
	if result.Name == "" {
		result.Name = "MyType"
	}

	isArray := false
	typeList := typeUnion{}
	required := make(map[string]struct{})
	for _, k := range s.OrderedKeywords {
		switch v := s.Keywords[k].(type) {
		case *jsonschema.Title:
		case *jsonschema.Comment:
		case *jsonschema.SchemaURI:
		case *jsonschema.ID:
			// if the schema has ID, use it as the name
			lastSlashIndex := strings.LastIndex(string(*v), "/")
			if lastSlashIndex != -1 {
				result.Name = strings.Trim(string(*v)[lastSlashIndex+1:], ".json")
			}
		case *jsonschema.Description:
			result.Description = string(*v)
		case *jsonschema.Type:
			if len(v.Vals) == 1 {
				switch v.Vals[0] {
				case "object":
					result.IsSchema = true
					continue
				case "array":
					isArray = true
					continue
				}
			}
			typeList.Items = append(typeList.Items, jsonTypesToKclTypes(v.Vals))
		case *jsonschema.Items:
			if !v.Single {
				logger.GetLogger().Warningf("unsupported multiple items: %#v", v)
				break
			}
			for _, i := range v.Schemas {
				item := convertSchemaFromJsonSchema(ctx, i, "items")
				if item.IsSchema {
					typeList.Items = append(typeList.Items, typeCustom{Name: item.Name})
				} else {
					typeList.Items = append(typeList.Items, item.Type)
				}
			}
		case *jsonschema.Required:
			for _, key := range []string(*v) {
				required[key] = struct{}{}
			}
		case *jsonschema.Properties:
			for _, prop := range *v {
				key := prop.Key
				val := prop.Value
				propSch := convertSchemaFromJsonSchema(ctx, val, key)
				_, propSch.Required = required[key]
				if propSch.IsSchema {
					propSch.Name = strcase.ToCamel(key)
					ctx.resultMap[propSch.Name] = propSch
				}
				propSch.Name = strcase.ToSnake(key)
				result.Properties = append(result.Properties, propSch.property)
				if !propSch.IsSchema {
					for _, validate := range propSch.Validations {
						validate.Name = propSch.Name
						result.Validations = append(result.Validations, validate)
					}
				}
			}
		case *jsonschema.Default:
			result.HasDefault = true
			result.DefaultValue = v.Data
		case *jsonschema.Enum:
			for _, val := range *v {
				unmarshalledVal := interface{}(nil)
				err := json.Unmarshal(val, &unmarshalledVal)
				if err != nil {
					logger.GetLogger().Warningf("failed to unmarshal enum value: %s", err)
					continue
				}
				typeList.Items = append(typeList.Items, typeValue{
					Value: unmarshalledVal,
				})
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
		default:
			logger.GetLogger().Warningf("unknown Keyword: %s", k)
		}
	}

	if result.IsSchema {
		result.Type = typeCustom{Name: strcase.ToCamel(name)}
	} else {
		if isArray {
			result.Type = typeArray{Items: typeList}
		} else {
			result.Type = typeList
		}
	}
	result.schema.Name = strcase.ToCamel(result.Name)
	result.schema.Description = result.Description
	result.property.Name = strcase.ToSnake(result.Name)
	result.property.Description = result.Description
	return result
}

func jsonTypesToKclTypes(t []string) typeInterface {
	var kclTypes typeUnion
	for _, v := range t {
		kclTypes.Items = append(kclTypes.Items, jsonTypeToKclType(v))
	}
	return kclTypes
}

func jsonTypeToKclType(t string) typeInterface {
	switch t {
	case "string":
		return typePrimitive(typStr)
	case "boolean":
		return typePrimitive(typBool)
	case "integer":
		return typePrimitive(typInt)
	case "number":
		return typePrimitive(typFloat)
	default:
		logger.GetLogger().Warningf("unknown type: %s", t)
		return typePrimitive(typStr)
	}
}
