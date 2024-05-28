package gen

import (
	"encoding/json"
	"io"
	"reflect"
	"sort"
	"strconv"

	"github.com/iancoleman/strcase"
	"kcl-lang.io/kcl-go/pkg/loader"
	"kcl-lang.io/kcl-go/pkg/logger"
)

type tfSchema struct {
	FormatVersion   string                      `json:"format_version"`
	ProviderSchemas map[string]tfProviderSchema `json:"provider_schemas"`
}

type tfProviderSchema struct {
	Provider          map[string]interface{}      `json:"provider"`
	ResourceSchemas   map[string]tfResourceSchema `json:"resource_schemas"`
	DataSourceSchemas map[string]interface{}      `json:"data_source_schemas"`
}

type tfResourceSchema struct {
	Version int     `json:"version"`
	Block   tfBlock `json:"block"`
}

type tfBlock struct {
	Attributes      map[string]tfAttribute `json:"attributes"`
	BlockTypes      map[string]interface{} `json:"block_types"`
	Description     string                 `json:"description"`
	DescriptionKind string                 `json:"description_kind"`
}

type tfAttribute struct {
	Type            interface{} `json:"type"`
	Description     string      `json:"description"`
	DescriptionKind string      `json:"description_kind"`
	Required        bool        `json:"required"`
	Optional        bool        `json:"optional"`
	Computed        bool        `json:"computed"`
}

type tfConvertContext struct {
	resultMap  map[string]schema
	attrKeyNow string
}

func (k *kclGenerator) genSchemaFromTerraformSchema(w io.Writer, filename string, src interface{}) error {
	code, err := loader.ReadSource(filename, src)
	if err != nil {
		return err
	}
	tfSch := &tfSchema{}
	if err = json.Unmarshal(code, tfSch); err != nil {
		return err
	}

	// convert terraform schema to kcl schema
	ctx := &tfConvertContext{
		resultMap: make(map[string]schema),
	}
	for _, providerSchema := range tfSch.ProviderSchemas {
		convertSchemaFromTFSchema(ctx, providerSchema)
	}
	result := make([]schema, 0, len(ctx.resultMap))
	for _, sch := range ctx.resultMap {
		result = append(result, sch)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	// generate kcl schema code
	kclSch := kclFile{
		Schemas: result,
	}
	return k.genKcl(w, kclSch)
}

// convertSchemaFromTFSchema converts terraform provider schema to kcl schema and save to ctx.resultMap
func convertSchemaFromTFSchema(ctx *tfConvertContext, tfSch tfProviderSchema) {
	for resKey, resourceSchema := range tfSch.ResourceSchemas {
		sch := schema{
			Name:        strcase.ToCamel(resKey),
			Description: resourceSchema.Block.Description,
		}
		for attrKey, attr := range resourceSchema.Block.Attributes {
			ctx.attrKeyNow = attrKey
			sch.Properties = append(sch.Properties, property{
				Name:        attrKey,
				Description: attr.Description,
				Type:        tfTypeToKclType(ctx, attr.Type),
				Required:    attr.Required,
			})
			if t, ok := attr.Type.([]interface{}); ok && t[0] == "set" {
				sch.Validations = append(sch.Validations, validation{
					Required: attr.Required,
					Name:     attrKey,
					Unique:   true,
				})
			}
		}
		sort.Slice(sch.Properties, func(i, j int) bool {
			return sch.Properties[i].Name < sch.Properties[j].Name
		})
		sort.Slice(sch.Validations, func(i, j int) bool {
			return sch.Validations[i].Name < sch.Validations[j].Name
		})
		ctx.resultMap[sch.Name] = sch
	}
}

// convertTFNestedSchema converts nested object schema to kcl schema, save to ctx.resultMap and return schema name
func convertTFNestedSchema(ctx *tfConvertContext, tfSchema map[string]interface{}) string {
	resultSchemaName := strcase.ToCamel(ctx.attrKeyNow + "Item")
	sch := schema{}
	for key, typ := range tfSchema {
		ctx.attrKeyNow = key
		sch.Properties = append(sch.Properties, property{
			Name: key,
			Type: tfTypeToKclType(ctx, typ),
		})
		if t, ok := typ.([]interface{}); ok && t[0] == "set" {
			sch.Validations = append(sch.Validations, validation{
				Name:   key,
				Unique: true,
			})
		}
	}
	sort.Slice(sch.Properties, func(i, j int) bool {
		return sch.Properties[i].Name < sch.Properties[j].Name
	})
	sort.Slice(sch.Validations, func(i, j int) bool {
		return sch.Validations[i].Name < sch.Validations[j].Name
	})

	// for the name of the schema, we will try xxxItem first
	// if it is already used and not equal to the schema, we will try xxxItem1, xxxItem2, ...
	for i := 0; true; i++ {
		sch.Name = resultSchemaName
		if i != 0 {
			sch.Name += strconv.Itoa(i)
		}
		if _, ok := ctx.resultMap[sch.Name]; !ok || reflect.DeepEqual(ctx.resultMap[sch.Name], sch) {
			break
		}
	}
	ctx.resultMap[sch.Name] = sch
	return sch.Name
}

func tfTypeToKclType(ctx *tfConvertContext, t interface{}) typeInterface {
	switch t := t.(type) {
	case string:
		return jsonTypeToKclType(t)
	case []interface{}:
		switch t[0] {
		case "list":
			return typeArray{Items: tfTypeToKclType(ctx, t[1])}
		case "map":
			return typeDict{Key: typePrimitive(typStr), Value: tfTypeToKclType(ctx, t[1])}
		case "set":
			return typeArray{Items: tfTypeToKclType(ctx, t[1])}
		case "object":
			return typeCustom{Name: convertTFNestedSchema(ctx, t[1].(map[string]interface{}))}
		case "tuple":
			// todo
			return typePrimitive(typAny)
		default:
			logger.GetLogger().Warningf("unknown type: %#v", t)
			return typePrimitive(typAny)
		}
	default:
		logger.GetLogger().Warningf("unknown type: %#v", t)
		return typePrimitive(typAny)
	}
}
