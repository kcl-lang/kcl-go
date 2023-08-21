package gen

import (
	"encoding/json"
	"github.com/iancoleman/strcase"
	"io"
	"kcl-lang.io/kcl-go/pkg/logger"
	"sort"
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

func (k *kclGenerator) genSchemaFromTerraformSchema(w io.Writer, filename string, src interface{}) error {
	code, err := readSource(filename, src)
	if err != nil {
		return err
	}
	tfSch := &tfSchema{}
	if err = json.Unmarshal(code, tfSch); err != nil {
		return err
	}

	// convert terraform schema to kcl schema
	var result []schema
	for _, providerSchema := range tfSch.ProviderSchemas {
		for resKey, resourceSchema := range providerSchema.ResourceSchemas {
			sch := schema{
				Name:        strcase.ToCamel(resKey),
				Description: resourceSchema.Block.Description,
			}
			for attrKey, attr := range resourceSchema.Block.Attributes {
				sch.Properties = append(sch.Properties, property{
					Name:        attrKey,
					Description: attr.Description,
					Type:        tfTypeToKclType(attr.Type),
					Required:    attr.Required,
				})
				if t, ok := attr.Type.([]interface{}); ok && t[0] == "set" {
					sch.Validations = append(sch.Validations, validation{
						Name:   attrKey,
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
			result = append(result, sch)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	// generate kcl schema code
	kclSch := kclSchema{
		Imports: []string{},
		Schemas: result,
	}
	return k.genKclSchema(w, kclSch)
}

func tfTypeToKclType(t interface{}) typeInterface {
	switch t := t.(type) {
	case string:
		return jsonTypeToKclType(t)
	case []interface{}:
		switch t[0] {
		case "list":
			return typeArray{Items: tfTypeToKclType(t[1])}
		case "map":
			return typeDict{Key: typePrimitive(typStr), Value: tfTypeToKclType(t[1])}
		case "set":
			return typeArray{Items: tfTypeToKclType(t[1])}
		case "object", "tuple":
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
