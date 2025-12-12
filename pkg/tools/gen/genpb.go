// Copyright 2023 The KCL Authors. All rights reserved.

package gen

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/source"
	pb "kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const (
	pbTypAny        = "google.protobuf.Any"
	pbTypAnyPkgPath = "google/protobuf/any.proto"
)

type Options struct {
	GoPackage string
	PbPackage string
}

// GenProto translate kcl schema type to protobuf message.
func GenProto(filename string, src interface{}, opt *Options) (string, error) {
	return newPbGenerator(opt).GenProto(filename, src)
}

type pbGenerator struct {
	opt Options

	needImportAny bool
}

func newPbGenerator(opt *Options) *pbGenerator {
	if opt == nil {
		opt = new(Options)
	}
	return &pbGenerator{
		opt: *opt,
	}
}

func (p *pbGenerator) GenProto(filename string, src interface{}) (string, error) {
	code, err := source.ReadSource(filename, src)
	if err != nil {
		return "", err
	}

	if p.opt.GoPackage == "" {
		p.opt.GoPackage = p.getOptionGoPackage(string(code))
	}
	if p.opt.PbPackage == "" {
		p.opt.PbPackage = p.getOptionPbPackage(string(code))
	}

	types, err := kcl.GetSchemaType(filename, string(code), "")
	if err != nil {
		return "", err
	}

	if p.opt.PbPackage == "" {
		return "", fmt.Errorf("opt.PbPackage missing")
	}
	if p.opt.GoPackage == "" {
		return "", fmt.Errorf("opt.GoPackage missing")
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "syntax = \"proto3\";\n\n")
	fmt.Fprintf(&buf, "package %s;\n\n", p.opt.PbPackage)

	fmt.Fprintf(&buf, "option go_package = \"%s\";\n", p.opt.GoPackage)

	var messageBody = p.genProtoMessages(types...)

	if p.needImportAny {
		fmt.Fprintln(&buf)
		fmt.Fprintf(&buf, "import \"google/protobuf/any.proto\";\n")
	}

	fmt.Fprint(&buf, messageBody)

	return buf.String(), nil
}

func (p *pbGenerator) genProtoMessages(types ...*pb.KclType) string {
	var buf bytes.Buffer
	for _, typ := range types {
		switch typ.Type {
		case typSchema:
			p.genProtoSchema(&buf, typ)
		default:
			fmt.Fprintf(&buf, "ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ))
		}
	}
	return buf.String()
}

func (p *pbGenerator) genProtoSchema(w io.Writer, typ *pb.KclType) {
	assert(typ.Type == typSchema)

	fmt.Fprintln(w)

	if doc := getSchemaDoc(typ); doc != "" {
		fmt.Fprint(w, doc)
	}

	fmt.Fprintf(w, "message %s {\n", typ.SchemaName)
	defer fmt.Fprintf(w, "}\n")

	var (
		sortedFieldNames = getSortedFieldNames(typ.Properties)

		pbFieldDefines []string
		pbFieldDocs    []string

		maxFieldDefineLen int
	)

	for i, fieldName := range sortedFieldNames {
		fieldType := typ.Properties[fieldName]

		pbFieldType := getPbTypeName(fieldType)
		kclFieldType := getKclTypeName(fieldType)

		if pbFieldType == pbTypAny {
			p.needImportAny = true
		}

		pbFieldDefines = append(pbFieldDefines,
			fmt.Sprintf("%s %s = %d;", pbFieldType, fieldName, i+1),
		)
		pbFieldDocs = append(pbFieldDocs,
			fmt.Sprintf("// kcl-type: %s", kclFieldType),
		)
		if n := len(pbFieldDefines[i]); n > maxFieldDefineLen {
			maxFieldDefineLen = n
		}
	}

	for i := range sortedFieldNames {
		fmt.Fprintf(w, "    %-*s %s\n", maxFieldDefineLen, pbFieldDefines[i], pbFieldDocs[i])
	}
}

// #kcl/genpb: option go_package = kcl_gen/_/hello_k
func (p *pbGenerator) getOptionGoPackage(code string) string {
	if !strings.Contains(code, `#kcl/genpb:`) {
		return ""
	}
	const prefix = `#kcl/genpb: option go_package =`
	for _, line := range strings.Split(code, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

// #kcl/genpb: option pb_package = kcl_gen._.hello_k
func (p *pbGenerator) getOptionPbPackage(code string) string {
	if !strings.Contains(code, `#kcl/genpb:`) {
		return ""
	}
	const prefix = `#kcl/genpb: option pb_package =`
	for _, line := range strings.Split(code, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

func getPbTypeName(typ *pb.KclType) string {
	switch typ.Type {
	case typSchema:
		return typ.SchemaName
	case typDict:
		return fmt.Sprintf("map<%s, %s>", getPbTypeName(typ.Key), getPbTypeName(typ.Item))
	case typList:
		return fmt.Sprintf("repeated %s", getPbTypeName(typ.Item))
	case typStr:
		return "string"
	case typInt:
		return "int64"
	case typFloat:
		return "double"
	case typBool:
		return "bool"
	case typAny:
		return pbTypAny
	case typUnion:
		var m = make(map[string]bool)
		for _, t := range typ.UnionTypes {
			m[getPbTypeName(t)] = true
		}
		if len(m) == 1 {
			for k := range m {
				return k
			}
		}
		return pbTypAny
	case typFunction:
		return pbTypAny
	case typNumberMultiplier:
		return "int64"
	default:
		if isLit, basicTyp, _ := IsLitType(typ); isLit {
			switch basicTyp {
			case typBool:
				return "bool"
			case typInt:
				return "int64"
			case typFloat:
				return "double"
			case typStr:
				return "string"
			}
		}
		panic(fmt.Sprintf("ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ)))
	}
}
