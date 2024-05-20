package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"runtime"

	"github.com/emicklei/proto"
)

var fieldTypeMap = map[string]string{
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
		return fmt.Errorf(`error parsing proto file:%v`, err)
	}

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
				builder.WriteString(getFieldType(field.Type))
				if field.Repeated {
					builder.WriteString("]")
				}
				builder.WriteString(lineBreak)

			case *proto.MapField:
				builder.WriteString("    ")
				builder.WriteString(field.Name)
				builder.WriteString(": {")
				builder.WriteString(getFieldType(field.KeyType))
				builder.WriteString(":")
				builder.WriteString(getFieldType(field.Type))
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

func getFieldType(fieldType string) string {
	if mappedType, ok := fieldTypeMap[fieldType]; ok {
		return mappedType
	}
	return fieldType
}
