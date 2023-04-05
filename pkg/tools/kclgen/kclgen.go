package kclgen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

type KclField struct {
	Name       string
	SchemaType string
}

type GenKclSchema map[SchemaName][]*KclField

type SchemaName string

func parseGoStruct(s interface{}) GenKclSchema {
	var schemaName SchemaName
	genKclSchema := make(GenKclSchema, 0)
	kclFieldList := make([]*KclField, 0)
	structInst := structs.New(s)
	schemaName = SchemaName(structInst.Name())
	for _, f := range structs.Fields(s) {
		structTag := f.Tag("kcl")
		tagMap := make(map[string]string, 0)
		s1 := strings.Split(structTag, ",")
		for _, s := range s1 {
			s2 := strings.Split(s, "=")
			tagMap[s2[0]] = s2[1]
		}
		kclFieldList = append(kclFieldList, &KclField{
			Name:       tagMap["name"],
			SchemaType: tagMap["type"],
		})
	}
	genKclSchema[schemaName] = kclFieldList
	return genKclSchema
}

func GenKclSchemaCode(s interface{}) string {
	genKclSchema := parseGoStruct(s)
	var buf bytes.Buffer
	for k, v := range genKclSchema {
		fmt.Fprintf(&buf, "schema %s:\n", k)
		for _, field := range v {
			fmt.Fprintf(&buf, "    %s: %s\n", field.Name, field.SchemaType)
		}
		fmt.Fprintf(&buf, "\n")
	}
	return buf.String()
}
