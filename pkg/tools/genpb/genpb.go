// Copyright 2022 The KCL Authors. All rights reserved.

package genpb

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"kusionstack.io/kclvm-go/pkg/kcl"
	pb "kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

const (
	typSchema = "schema"
	typDict   = "dict"
	typList   = "list"
	typStr    = "str"
	typInt    = "int"
	typFloat  = "float"
	typBool   = "bool"

	typAny              = "any"
	typUnion            = "union"
	typNumberMultiplier = "number_multiplier"
)

const (
	pbTypAny = "google.protobuf.Any"
)

type Options struct {
	GoPackage string
	PbPackage string
}

// GenProto translate kcl schema type to protobuf message.
func GenProto(filename string, src interface{}, opt *Options) (string, error) {
	return newPbGenerator(opt).GenProto(filename, src)
}

type _PbGenerator struct {
	opt Options

	needImportAny bool
}

func newPbGenerator(opt *Options) *_PbGenerator {
	if opt == nil {
		opt = new(Options)
	}
	return &_PbGenerator{
		opt: *opt,
	}
}

func (p *_PbGenerator) GenProto(filename string, src interface{}) (string, error) {
	code, err := readSource(filename, src)
	if err != nil {
		return "", err
	}

	if p.opt.GoPackage == "" {
		p.opt.GoPackage = p.getOptopn_go_package(string(code))
	}
	if p.opt.PbPackage == "" {
		p.opt.PbPackage = p.getOptopn_pb_package(string(code))
	}

	typs, err := kcl.GetSchemaType(filename, string(code), "")
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

	var messageBody = p.genProto_messages(typs...)

	if p.needImportAny {
		fmt.Fprintln(&buf)
		fmt.Fprintf(&buf, "import \"google/protobuf/any.proto\";\n")
	}

	fmt.Fprint(&buf, messageBody)

	return buf.String(), nil
}

func (p *_PbGenerator) genProto_messages(typs ...*pb.KclType) string {
	var buf bytes.Buffer
	for _, typ := range typs {
		switch typ.Type {
		case typSchema:
			p.genProto_shema(&buf, typ)
		default:
			fmt.Fprintf(&buf, "ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ))
		}
	}
	return buf.String()
}

func (p *_PbGenerator) genProto_shema(w io.Writer, typ *pb.KclType) {
	assert(typ.Type == typSchema)

	fmt.Fprintln(w)

	if doc := p.getSchemaDoc(typ); doc != "" {
		fmt.Fprint(w, doc)
	}

	fmt.Fprintf(w, "message %s {\n", typ.SchemaName)
	defer fmt.Fprintf(w, "}\n")

	var (
		sortedFieldNames = p.getSortedFieldNames(typ.Properties)

		pbFieldDefines []string
		pbFieldDocs    []string

		maxFieldDefineLen int
	)

	for i, fieldName := range sortedFieldNames {
		fieldType := typ.Properties[fieldName]

		pbFieldType := p.getPbTypeName(fieldType)
		kclFieldType := p.getKclTypeName(fieldType)

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

func (p *_PbGenerator) getKclTypeName(typ *pb.KclType) string {
	if isLit, _, litValue := p.isLitType(typ); isLit {
		return litValue
	}

	switch typ.Type {
	case typSchema:
		return typ.SchemaName
	case typDict:
		return fmt.Sprintf("{%s:%s}", p.getKclTypeName(typ.Key), p.getKclTypeName(typ.Item))
	case typList:
		return fmt.Sprintf("[%s]", p.getKclTypeName(typ.Item))
	case typStr:
		return "str"
	case typInt:
		return "int"
	case typFloat:
		return "float"
	case typBool:
		return "bool"

	case typAny:
		return "any"
	case typUnion:
		var ss []string
		for _, t := range typ.UnionTypes {
			ss = append(ss, p.getKclTypeName(t))
		}
		return strings.Join(ss, "|")

	case typNumberMultiplier:
		return "units.NumberMultiplier"

	default:
		panic(fmt.Sprintf("ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ)))
	}
}

func (p *_PbGenerator) getPbTypeName(typ *pb.KclType) string {
	switch typ.Type {
	case typSchema:
		return typ.SchemaName
	case typDict:
		return fmt.Sprintf("map<%s, %s>", p.getPbTypeName(typ.Key), p.getPbTypeName(typ.Item))
	case typList:
		return fmt.Sprintf("repeated %s", p.getPbTypeName(typ.Item))
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
			m[p.getPbTypeName(t)] = true
		}
		if len(m) == 1 {
			for k := range m {
				return k
			}
		}
		return pbTypAny

	case typNumberMultiplier:
		return "int64"

	default:
		if isLit, basicTyp, _ := p.isLitType(typ); isLit {
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

func (p *_PbGenerator) isLitType(typ *pb.KclType) (ok bool, basicTyp, litValue string) {
	if !strings.HasSuffix(typ.Type, ")") {
		return
	}

	i := strings.Index(typ.Type, "(") + 1
	j := strings.LastIndex(typ.Type, ")")

	switch {
	case strings.HasPrefix(typ.Type, "bool("):
		return true, "bool", typ.Type[i:j]
	case strings.HasPrefix(typ.Type, "int("):
		return true, "int", typ.Type[i:j]
	case strings.HasPrefix(typ.Type, "float("):
		return true, "float", typ.Type[i:j]
	case strings.HasPrefix(typ.Type, "str("):
		return true, "str", strconv.Quote(typ.Type[i:j])
	}
	return
}

func (p *_PbGenerator) getSchemaDoc(typ *pb.KclType) (doc string) {
	var w = new(bytes.Buffer)
	if doc := strings.TrimSpace(typ.SchemaDoc); doc != "" {
		for _, s := range strings.Split(doc, "\n") {
			fmt.Fprintf(w, "// %s\n", s)
		}
	}
	doc = w.String()
	return
}

func (p *_PbGenerator) getSortedFieldNames(fields map[string]*pb.KclType) []string {
	type FieldInfo struct {
		Name string
		Type *pb.KclType
	}

	var infos []FieldInfo
	for name, typ := range fields {
		infos = append(infos, FieldInfo{
			Name: name,
			Type: typ,
		})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Type.Line < infos[j].Type.Line
	})

	var ss []string
	for _, x := range infos {
		ss = append(ss, x.Name)
	}
	return ss
}

// #kclvm/genpb: option go_package = kcl_gen/_/hello_k
func (p *_PbGenerator) getOptopn_go_package(code string) string {
	if !strings.Contains(code, `#kclvm/genpb:`) {
		return ""
	}
	const prefix = `#kclvm/genpb: option go_package =`
	for _, line := range strings.Split(code, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

// #kclvm/genpb: option pb_package = kcl_gen._.hello_k
func (p *_PbGenerator) getOptopn_pb_package(code string) string {
	if !strings.Contains(code, `#kclvm/genpb:`) {
		return ""
	}
	const prefix = `#kclvm/genpb: option pb_package =`
	for _, line := range strings.Split(code, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}
