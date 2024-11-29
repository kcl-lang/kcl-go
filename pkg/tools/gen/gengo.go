package gen

import (
	"fmt"
	"io"
	"strings"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/source"
	pb "kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const goAnyType = "interface{}"

var _ Generator = &goGenerator{}

type GenGoOptions struct {
	Package    string
	AnyType    string
	UseValue   bool
	RenameFunc func(string) string
}

// GenGo translate kcl schema type to go struct.
func GenGo(w io.Writer, filename string, src interface{}, opts *GenGoOptions) error {
	return newGoGenerator(opts).GenFromSource(w, filename, src)
}

type goGenerator struct {
	opts       *GenGoOptions
	renameFunc func(string) string
}

func newGoGenerator(opts *GenGoOptions) *goGenerator {
	if opts == nil {
		opts = &GenGoOptions{
			AnyType: goAnyType,
		}
	}
	generator := &goGenerator{
		opts:       opts,
		renameFunc: opts.RenameFunc,
	}
	if generator.renameFunc == nil {
		generator.renameFunc = func(name string) string {
			return name
		}
	}
	return generator
}

func (g *goGenerator) GenFromSource(w io.Writer, filename string, src interface{}) error {
	code, err := source.ReadSource(filename, src)
	if err != nil {
		return err
	}

	types, err := kcl.GetSchemaType(filename, string(code), "")
	if err != nil {
		return err
	}

	g.GenFromTypes(w, types...)

	return nil
}

func (g *goGenerator) GenFromTypes(w io.Writer, types ...*pb.KclType) {
	for _, typ := range types {
		switch typ.Type {
		case typSchema:
			g.GenSchema(w, typ)
		}
	}
}

func (g *goGenerator) GenSchema(w io.Writer, typ *pb.KclType) {
	assert(typ.Type == typSchema)

	fmt.Fprintln(w)

	if doc := getSchemaDoc(typ); doc != "" {
		fmt.Fprint(w, doc)
	}

	fmt.Fprintf(w, "type %s struct {\n", typ.SchemaName)
	defer fmt.Fprintf(w, "}\n")

	var (
		sortedFieldNames = getSortedFieldNames(typ.Properties)

		goFieldDefines []string
		goFieldDocs    []string

		maxFieldDefineLen int
	)

	for i, fieldName := range sortedFieldNames {
		fieldType := typ.Properties[fieldName]

		goFieldType := g.GetTypeName(fieldType)
		kclFieldType := getKclTypeName(fieldType)

		goTagInfo := fmt.Sprintf(`kcl:"name=%s,type=%s"`, fieldName, g.GetFieldTag(fieldType))
		goFieldDefines = append(goFieldDefines,
			fmt.Sprintf("%s %s %s", g.renameFunc(fieldName), goFieldType, "`"+goTagInfo+"`"),
		)
		goFieldDocs = append(goFieldDocs,
			fmt.Sprintf("// kcl-type: %s", kclFieldType),
		)
		if n := len(goFieldDefines[i]); n > maxFieldDefineLen {
			maxFieldDefineLen = n
		}
	}

	for i := range sortedFieldNames {
		fmt.Fprintf(w, "    %-*s %s\n", maxFieldDefineLen, goFieldDefines[i], goFieldDocs[i])
	}
}

func (g *goGenerator) GetTypeName(typ *pb.KclType) string {
	switch typ.Type {
	case typSchema:
		{
			name := typ.SchemaName
			if !g.opts.UseValue {
				// Use pointer value
				name = "*" + name
			}
			return name
		}
	case typDict:
		return fmt.Sprintf("map[%s]%s", g.GetTypeName(typ.Key), g.GetTypeName(typ.Item))
	case typList:
		return fmt.Sprintf("[]%s", g.GetTypeName(typ.Item))
	case typStr:
		return "string"
	case typInt:
		return "int"
	case typFloat:
		return "float64"
	case typBool:
		return "bool"
	case typAny:
		return g.opts.AnyType

	case typUnion:
		var m = make(map[string]bool)
		for _, t := range typ.UnionTypes {
			m[g.GetTypeName(t)] = true
		}
		if len(m) == 1 {
			for k := range m {
				return k
			}
		}
		return g.opts.AnyType

	case typNumberMultiplier:
		return "int"

	default:
		if isLit, basicTyp, _ := IsLitType(typ); isLit {
			switch basicTyp {
			case typBool:
				return "bool"
			case typInt:
				return "int"
			case typFloat:
				return "float64"
			case typStr:
				return "string"
			}
		}
		panic(fmt.Sprintf("ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ)))
	}
}

func (g *goGenerator) GetFieldTag(typ *pb.KclType) string {
	switch typ.Type {
	case typSchema:
		return typ.SchemaName
	case typDict:
		return fmt.Sprintf("{%s:%s}", g.GetFieldTag(typ.Key), g.GetFieldTag(typ.Item))
	case typList:
		return fmt.Sprintf("[%s]", g.GetFieldTag(typ.Item))
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
			ss = append(ss, t.Type)
		}
		return strings.Join(ss, "|")

	case typNumberMultiplier:
		return "units.NumberMultiplier"

	default:
		panic(fmt.Sprintf("ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ)))
	}
}
