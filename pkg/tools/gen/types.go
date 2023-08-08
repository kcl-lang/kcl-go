package gen

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	pb "kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
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

// kclSchema is the top-level structure for a kcl schema file.
// It contains all the imports and schemas in this file.
type kclSchema struct {
	Imports []string
	Schemas []schema
}

// schema is a kcl schema definition.
type schema struct {
	Name              string
	Description       string
	Properties        []property
	Validations       []validation
	HasIndexSignature bool
	IndexSignature    indexSignature
}

// property is a kcl schema property definition.
type property struct {
	Name         string
	Description  string
	Type         typeInterface
	Required     bool
	HasDefault   bool
	DefaultValue interface{}
}

// validation is a kcl schema validation definition.
type validation struct {
	Name             string
	Minimum          *float64
	ExclusiveMinimum bool
	Maximum          *float64
	ExclusiveMaximum bool
	MinLength        *int
	MaxLength        *int
	Regex            *regexp.Regexp
}

// indexSignature is a kcl schema index signature definition.
// It can be used to construct a dict with type.
type indexSignature struct {
	Alias string
	Type  typeInterface
}

type typeInterface interface {
	Format() string
}

type typePrimitive string

func (t typePrimitive) Format() string {
	return string(t)
}

type typeArray struct {
	Items typeInterface
}

func (t typeArray) Format() string {
	return "[" + t.Items.Format() + "]"
}

type typeUnion struct {
	Items []typeInterface
}

func (t typeUnion) Format() string {
	var items []string
	for _, v := range t.Items {
		items = append(items, v.Format())
	}
	return strings.Join(items, " | ")
}

type typeDict struct {
	Key   typeInterface
	Value typeInterface
}

func (t typeDict) Format() string {
	return "{" + t.Key.Format() + ":" + t.Value.Format() + "}"
}

type typeCustom struct {
	Name string
}

func (t typeCustom) Format() string {
	return t.Name
}

type typeValue struct {
	Value interface{}
}

func (t typeValue) Format() string {
	return formatValue(t.Value)
}
