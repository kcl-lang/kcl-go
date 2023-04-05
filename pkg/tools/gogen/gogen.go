package gogen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"kusionstack.io/kclvm-go/pkg/logger"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

type GenStruct map[StructName][]*Field

type StructName string

type Field struct {
	Name       string
	SchemaType string
	OmitEmpty  bool
	Tag        string
}

func parseKclType(ktList []*gpyrpc.KclType) *GenStruct {
	var schemaName StructName
	genStruct := make(GenStruct, 0)
	for _, kt := range ktList {
		if kt.Type == "schema" {
			schemaName = StructName(kt.SchemaName)
			var schemaList []*Field
			for k, p := range kt.Properties {
				var field Field
				field.Name = strings.ToUpper(k[:1]) + k[1:]
				switch p.Type {
				case "schema":
					field.SchemaType = "*" + p.SchemaName
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "schema")
				case "dict":
					kType := p.Key.Type
					if kType == "str" {
						kType = "string"
					} else if kType == "int" {
						kType = "int"
					}
					vType := p.Item.Type
					if vType == "schema" {
						vType = p.Item.SchemaName
					}
					field.SchemaType = "map[" + kType + "]" + vType
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "{"+p.Key.Type+":"+vType+"}")
				case "list":
					vType := p.Item.Type
					if vType == "schema" {
						vType = p.Item.SchemaName
					}
					field.SchemaType = "[]" + "*" + vType
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "["+vType+"]")
				case "str":
					field.SchemaType = "string"
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "str")
				case "int":
					field.SchemaType = "int"
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "int")
				case "float":
					field.SchemaType = "float32"
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "float")
				case "bool":
					field.SchemaType = "bool"
					field.Tag = fmt.Sprintf(`kcl:"name=%s,type=%s"`, k, "bool")
				case "null":
					field.SchemaType = "nil"
				}
				schemaList = append(schemaList, &field)
			}
			genStruct[schemaName] = schemaList
		}
	}
	return &genStruct
}

func GenGoCodeFromKclType(ktList []*gpyrpc.KclType) string {
	genStruct := parseKclType(ktList)
	var buf bytes.Buffer
	for k, v := range *genStruct {
		fmt.Fprintf(&buf, "type %s struct {\n", k)
		for _, field := range v {
			if field.Tag != "" {
				fmt.Fprintf(&buf, " %s %s %s\n", field.Name, field.SchemaType, "`"+field.Tag+"`")
			} else {
				fmt.Fprintf(&buf, " %s %s\n", field.Name, field.SchemaType)
			}
		}
		fmt.Fprintf(&buf, "}\n")
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		logger.GetLogger().Errorf("Failed to format kclvm source code: %s", err.Error())
	}

	return string(source)
}
