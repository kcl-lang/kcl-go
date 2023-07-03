package gen

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"kcl-lang.io/kcl-go/pkg/logger"
)

type GenKclOptions struct {
	ParseFromTag bool
}

type kclGenerator struct {
	opts *GenKclOptions
}

// GenKcl translate go struct to kcl schema code.
func GenKcl(w io.Writer, filename string, src interface{}, opts *GenKclOptions) error {
	return newKclGenerator(opts).GenSchema(w, filename, src)
}

func newKclGenerator(opts *GenKclOptions) *kclGenerator {
	if opts == nil {
		opts = new(GenKclOptions)
	}
	return &kclGenerator{
		opts: opts,
	}
}

func (k *kclGenerator) GenSchema(w io.Writer, filename string, src interface{}) error {
	fmt.Fprintln(w)
	goStructs, err := ParseGoSourceCode(filename, src)
	if err != nil {
		return err
	}
	for _, goStruct := range goStructs {
		fmt.Fprintf(w, "schema %s:\n", goStruct.Name)
		if goStruct.StructComment != "" {
			fmt.Fprintf(w, "    \"\"\"%s\"\"\"\n", goStruct.StructComment)
		}
		for _, field := range goStruct.Fields {
			kclFieldName, kclFieldType, err := k.GetTypeName(field)
			if err != nil {
				logger.GetLogger().Warningf("get struct tag key kcl info err: %s, will generate kcl schema from the struct field metadata data, field info %#v", err.Error(), field)
				kclFieldName, kclFieldType = k.GetKclTypeFromStructField(field)
			}
			fmt.Fprintf(w, "    %s: %s\n", kclFieldName, kclFieldType)
		}
		fmt.Fprintf(w, "\n")
	}
	return nil
}

func (k *kclGenerator) GetTypeName(f *GoStructField) (string, string, error) {
	if k.opts.ParseFromTag {
		return k.parserGoStructFieldTag(f.FieldTag)
	}
	fieldName, fieldType := k.GetKclTypeFromStructField(f)
	return fieldName, fieldType, nil
}

func (k *kclGenerator) parserGoStructFieldTag(tag string) (string, string, error) {
	tagMap := make(map[string]string, 0)
	sp := strings.Split(tag, "`")
	if len(sp) == 1 {
		return "", "", errors.New("this field not found tag string like `` !")
	}
	value, ok := k.Lookup(sp[1], "kcl")
	if !ok {
		return "", "", errors.New("not found tag key named kcl")
	}
	reg := "name=.*,type=.*"
	match, err := regexp.Match(reg, []byte(value))
	if err != nil {
		return "", "", err
	}
	if !match {
		return "", "", errors.New("don't match the kcl tag info, the tag info style is name=NAME,type=TYPE")
	}
	tagInfo := strings.Split(value, ",")
	for _, s := range tagInfo {
		t := strings.Split(s, "=")
		tagMap[t[0]] = t[1]
	}
	fieldType := tagMap["type"]
	if strings.Contains(tagMap["type"], ")|") {
		typeUnionList := strings.Split(tagMap["type"], "|")
		var ss []string
		for _, u := range typeUnionList {
			_, _, litValue := k.isLitType(u)
			ss = append(ss, litValue)
		}
		fieldType = strings.Join(ss, "|")
	}
	return tagMap["name"], fieldType, nil
}

func (k *kclGenerator) GetKclTypeFromStructField(f *GoStructField) (string, string) {
	return f.FieldName, k.isLitGoType(f.FieldType)
}

func (k *kclGenerator) isLitType(fieldType string) (ok bool, basicTyp, litValue string) {
	if !strings.HasSuffix(fieldType, ")") {
		return
	}

	i := strings.Index(fieldType, "(") + 1
	j := strings.LastIndex(fieldType, ")")

	switch {
	case strings.HasPrefix(fieldType, "bool("):
		return true, "bool", fieldType[i:j]
	case strings.HasPrefix(fieldType, "int("):
		return true, "int", fieldType[i:j]
	case strings.HasPrefix(fieldType, "float("):
		return true, "float", fieldType[i:j]
	case strings.HasPrefix(fieldType, "str("):
		return true, "str", strconv.Quote(fieldType[i:j])
	}
	return
}

func (k *kclGenerator) isLitGoType(fieldType string) string {
	switch fieldType {
	case "int", "int32", "int64":
		return "int"
	case "float", "float64":
		return "float"
	case "string":
		return "str"
	case "bool":
		return "bool"
	case "interface{}":
		return "any"
	default:
		if strings.HasPrefix(fieldType, "*") {
			i := strings.Index(fieldType, "*") + 1
			return k.isLitGoType(fieldType[i:])
		}
		if strings.HasPrefix(fieldType, "map") {
			i := strings.Index(fieldType, "[") + 1
			j := strings.Index(fieldType, "]")
			return fmt.Sprintf("{%s:%s}", k.isLitGoType(fieldType[i:j]), k.isLitGoType(fieldType[j+1:]))
		}
		if strings.HasPrefix(fieldType, "[]") {
			i := strings.Index(fieldType, "]") + 1
			return fmt.Sprintf("[%s]", k.isLitGoType(fieldType[i:]))
		}
		return fieldType
	}
}

func (k *kclGenerator) Lookup(tag, key string) (value string, ok bool) {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
