package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/emicklei/proto"
	"github.com/iancoleman/strcase"
	"kcl-lang.io/kcl-go/pkg/logger"
	"kcl-lang.io/kcl-go/pkg/source"
)

var defaultFieldTypeMap = map[string]string{
	"uint32":              "int",
	"uint64":              "int",
	"int32":               "int",
	"int64":               "int",
	"sint32":              "int",
	"sint64":              "int",
	"fixed32":             "int",
	"fixed64":             "int",
	"sfixed32":            "int",
	"sfixed64":            "int",
	"string":              "str",
	"bytes":               "str",
	"google.protobuf.Any": "any",
	"bool":                "bool",
	"float":               "float",
	"double":              "float",
}

// genKclFromProto converts the .proto config to KCL schema.
func (k *kclGenerator) genKclFromProto(w io.Writer, filename string, src any) error {
	code, err := source.ReadSource(filename, src)
	if err != nil {
		return err
	}
	parser := proto.NewParser(bytes.NewBuffer(code))
	definitions, err := parser.Parse()
	if err != nil {
		return fmt.Errorf(`error parsing proto file %v: %v`, filename, err)
	}
	builder := bufio.NewWriter(w)
	k.genKclFromProtoDef(builder, definitions)
	if err = builder.Flush(); err != nil {
		return err
	}
	return nil
}

// genKclFromProto converts the .proto config to KCL schema.
func (k *kclGenerator) genKclFromProtoDef(builder *bufio.Writer, definitions *proto.Proto) error {
	lineBreak := "\n"
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}
	fieldTypeMap := k.genFieldTypeMap(definitions)
	var oneOfSchemas []proto.Visitee
	for _, definition := range definitions.Elements {
		switch def := definition.(type) {
		// Convert proto message to kcl schema
		case *proto.Message:
			builder.WriteString("schema ")
			builder.WriteString(def.Name)
			builder.WriteString(":")
			builder.WriteString(lineBreak)

			for _, element := range def.Elements {
				switch field := element.(type) {
				case *proto.NormalField:
					builder.WriteString("    ")
					builder.WriteString(field.Name)
					if field.Optional {
						builder.WriteString("?")
					}
					builder.WriteString(": ")

					if field.Repeated {
						builder.WriteString("[")
					}

					fieldType, err := getFieldType(fieldTypeMap, field.Type)
					if err != nil {
						return err
					}
					builder.WriteString(fieldType)

					if field.Repeated {
						builder.WriteString("]")
					}
					builder.WriteString(lineBreak)
				case *proto.Oneof:
					builder.WriteString("    ")
					builder.WriteString(field.Name)
					builder.WriteString(": ")
					elementsLen := len(field.Elements) - 1
					for i, element := range field.Elements {
						switch v := element.(type) {
						case *proto.OneOfField:
							oneOfSchemaName := fmt.Sprintf("%s%sOneOf%v", def.Name, strcase.ToCamel(field.Name), i)
							builder.WriteString(oneOfSchemaName)
							if elementsLen > i {
								builder.WriteString(` | `)
							}
							oneOfSchemas = append(oneOfSchemas, &proto.Message{
								Name:     oneOfSchemaName,
								Elements: []proto.Visitee{&proto.NormalField{Field: v.Field}},
							})
						}
					}
					builder.WriteString(lineBreak)
				case *proto.MapField:
					builder.WriteString("    ")
					builder.WriteString(field.Name)
					builder.WriteString(": {")
					keyType, err := getFieldType(fieldTypeMap, field.KeyType)
					if err != nil {
						return err
					}
					builder.WriteString(keyType)
					builder.WriteString(":")
					fieldType, err := getFieldType(fieldTypeMap, field.Type)
					if err != nil {
						return err
					}
					builder.WriteString(fieldType)
					builder.WriteString("}")
					builder.WriteString(lineBreak)
				}
			}

			builder.WriteString(lineBreak)
		// Convert proto enum to kcl type alias
		case *proto.Enum:
			elementsLen := len(def.Elements) - 1
			builder.WriteString("type ")
			builder.WriteString(def.Name)
			builder.WriteString(" = ")
			for i, element := range def.Elements {
				switch v := element.(type) {
				case *proto.EnumField:
					value := fmt.Sprintf(`"%v"`, v.Name)
					if k.opts.UseIntegersForNumbers {
						value = strconv.Itoa(v.Integer)
					}

					builder.WriteString(value)
					if elementsLen > i {
						builder.WriteString(` | `)
					}

					fieldTypeMap[v.Name] = v.Name
				}
			}
			builder.WriteString(lineBreak)
		// TODO: Import node on multi proto files
		case *proto.Import:
			if def.Filename != pbTypAnyPkgPath {
				logger.GetLogger().Warningf("unsupported import statement for %v", def.Filename)
			}
		}
	}
	if len(oneOfSchemas) > 0 {
		k.genKclFromProtoDef(builder, &proto.Proto{
			Elements: oneOfSchemas,
		})
	}
	return nil
}

func (k *kclGenerator) genFieldTypeMap(definitions *proto.Proto) map[string]string {
	fieldTypeMap := make(map[string]string)
	for key, value := range defaultFieldTypeMap {
		fieldTypeMap[key] = value
	}

	for _, definition := range definitions.Elements {
		switch visitee := definition.(type) {
		case *proto.Message:
			fieldTypeMap[visitee.Name] = visitee.Name
		case *proto.Enum:
			var builder strings.Builder
			elementsLen := len(visitee.Elements) - 1
			for i, e := range visitee.Elements {
				v, ok := e.(*proto.EnumField)
				if !ok {
					continue
				}

				value := fmt.Sprintf(`"%v"`, v.Name)
				if k.opts.UseIntegersForNumbers {
					value = strconv.Itoa(v.Integer)
				}

				builder.WriteString(value)
				if elementsLen > i {
					builder.WriteString(` | `)
				}

				fieldTypeMap[v.Name] = v.Name
			}
			fieldTypeMap[visitee.Name] = builder.String()
		}
	}

	return fieldTypeMap
}

func getFieldType(fieldTypeMap map[string]string, fieldType string) (string, error) {
	value, ok := fieldTypeMap[fieldType]
	if !ok {
		return "", fmt.Errorf(`this "%v" is not currently supported`, fieldType)
	}

	return value, nil
}
