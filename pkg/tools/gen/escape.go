package gen

import "strings"

func DeepEscapePipesInPackage(pkg *KclPackage) *KclPackage {
	if pkg == nil {
		return nil
	}
	escaped := &KclPackage{
		Name:           pkg.Name,
		Version:        pkg.Version,
		Description:    pkg.Description,
		SchemaList:     nil,
		SubPackageList: nil,
	}

	if pkg.SchemaList != nil {
		escaped.SchemaList = make([]*KclOpenAPIType, len(pkg.SchemaList))
		for i, schema := range pkg.SchemaList {
			escaped.SchemaList[i] = deepEscapePipesInSchema(schema)
		}
	}
	if pkg.SubPackageList != nil {
		escaped.SubPackageList = make([]*KclPackage, len(pkg.SubPackageList))
		for i, sub := range pkg.SubPackageList {
			escaped.SubPackageList[i] = DeepEscapePipesInPackage(sub)
		}
	}
	return escaped
}

func deepEscapePipesInSchema(schema *KclOpenAPIType) *KclOpenAPIType {
	if schema == nil {
		return nil
	}
	escaped := *schema
	if schema.Enum != nil {
		escaped.Enum = make([]string, len(schema.Enum))
		for i, v := range schema.Enum {
			escaped.Enum[i] = strings.ReplaceAll(v, "|", "\\|")
		}
	}
	return &escaped
}
