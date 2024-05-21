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
)

var defaultFieldTypeMap = map[string]string{
	"uint32":              "int",
	"uint64":              "int",
	"int32":               "int",
	"int64":               "int",
	"sint32":              "int",
	"sint64":              "int",
	"string":              "str",
	"google.protobuf.Any": "any",
	"bool":                "bool",
	"float":               "float",
	"double":              "float",
}

// genKclFromProtoData
func (k *kclGenerator) genKclFromProtoData(w io.Writer, filename string, src interface{}) error {
	lineBreak := "\n"
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}

	code, err := readSource(filename, src)
	if err != nil {
		return err
	}

	parser := proto.NewParser(bytes.NewBuffer(code))
	definitions, err := parser.Parse()
	if err != nil {
		return fmt.Errorf(`error parsing proto file %v: %v`, filename, err)
	}

	fieldTypeMap := k.genFieldTypeMap(definitions)
	builder := bufio.NewWriter(w)
	for _, definition := range definitions.Elements {
		message, ok := definition.(*proto.Message)
		if !ok {
			continue
		}

		builder.WriteString("schema ")
		builder.WriteString(message.Name)
		builder.WriteString(":")
		builder.WriteString(lineBreak)

		for _, element := range message.Elements {
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
	}

	if err = builder.Flush(); err != nil {
		return err
	}

	return nil
}

// GenFieldTypeMap
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
