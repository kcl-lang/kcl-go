package gen

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

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

func getKclTypeName(typ *pb.KclType) string {
	if isLit, _, litValue := isLitType(typ); isLit {
		return litValue
	}

	switch typ.Type {
	case typSchema:
		return typ.SchemaName
	case typDict:
		return fmt.Sprintf("{%s:%s}", getKclTypeName(typ.Key), getKclTypeName(typ.Item))
	case typList:
		return fmt.Sprintf("[%s]", getKclTypeName(typ.Item))
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
			ss = append(ss, getKclTypeName(t))
		}
		return strings.Join(ss, "|")

	case typNumberMultiplier:
		return "units.NumberMultiplier"

	default:
		panic(fmt.Sprintf("ERR: unknown '%v', json = %v\n", typ.Type, jsonString(typ)))
	}
}

func isLitType(typ *pb.KclType) (ok bool, basicTyp, litValue string) {
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

func getSchemaDoc(typ *pb.KclType) (doc string) {
	var w = new(bytes.Buffer)
	if doc := strings.TrimSpace(typ.SchemaDoc); doc != "" {
		for _, s := range strings.Split(doc, "\n") {
			fmt.Fprintf(w, "// %s\n", s)
		}
	}
	doc = w.String()
	return
}

func getSortedFieldNames(fields map[string]*pb.KclType) []string {
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
