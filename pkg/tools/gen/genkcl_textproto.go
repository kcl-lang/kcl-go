package gen

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"kcl-lang.io/kcl-go/pkg/ast"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/loader"

	pbast "github.com/protocolbuffers/txtpbfmt/ast"
	"github.com/protocolbuffers/txtpbfmt/parser"
	"github.com/protocolbuffers/txtpbfmt/unquote"
)

var (
	ErrNoSchemaFound = errors.New("no expected schema found")
)

type TextProtoGenerator struct {
	file string
}

// Parse parses the given textproto bytes and converts them to KCL configs.
// Note fields in the textproto that have no corresponding field in schema
// are ignored.
func (d *TextProtoGenerator) Gen(filename string, src any, schema *kcl.KclType) (*config, error) {
	source, err := loader.ReadSource(filename, src)
	if runtime.GOOS == "windows" {
		source = []byte(strings.Replace(string(source), "\r\n", "\n", -1))
	}
	if err != nil {
		return nil, err
	}
	cfg := parser.Config{}
	d.file = filename
	nodes, err := parser.ParseWithConfig(source, cfg)
	if err != nil {
		return nil, err
	}
	return d.genProperties(schema, nodes)
}

// ParseFromSchemaFile parses the given textproto bytes and converts them
// to KCL configs with the schema file. Note fields in the textproto that
// have no corresponding field in schema are ignored.
func (d *TextProtoGenerator) GenFromSchemaFile(filename string, src any, schemaFile string, schemaSrc any, schemaName string) (*config, error) {
	types, err := kcl.GetSchemaType(schemaFile, schemaSrc, schemaName)
	if err != nil {
		return nil, err
	}
	if len(types) == 0 {
		return nil, ErrNoSchemaFound
	}
	return d.Gen(filename, src, types[0])
}

func (d *TextProtoGenerator) genProperties(ty *kcl.KclType, nodes []*pbast.Node) (*config, error) {
	var values []data
	for _, n := range nodes {
		var comments []*ast.Comment
		if n.Values == nil && n.Children == nil {
			if comments = addComments(n.PreComments...); comments != nil {
				continue
			}
		}
		if ty == nil || ty.Properties == nil {
			continue
		}
		ty, ok := ty.Properties[n.Name]
		// Ignore unknown attributes that not defined in the schema
		if !ok {
			continue
		}
		value, err := d.genValue(ty, n)
		if err != nil {
			return nil, err
		}
		values = append(values, data{
			Key:      n.Name,
			Value:    value,
			Comments: comments,
		})
	}
	return &config{
		Name: ty.SchemaName,
		Data: values,
	}, nil
}

func (d *TextProtoGenerator) genValue(ty *kcl.KclType, n *pbast.Node) (any, error) {
	if n == nil {
		return nil, nil
	}
	tyStr := typAny
	if ty != nil {
		tyStr = ty.Type
	}
	switch tyStr {
	case typSchema:
		if k := len(n.Values); k > 0 {
			return nil, d.errorf(n, "not allowed for the message type; found %d", k)
		}
		return d.genProperties(ty, n.Children)
	case typDict:
		if k := len(n.Values); k > 0 {
			return nil, d.errorf(n, "not allowed for the message type; found %d", k)
		}
		var values []data
		var key string
		var value any
		var comments []*ast.Comment
		for _, c := range n.Children {
			if len(c.Values) != 1 {
				return nil, d.errorf(n, "expected 1 value, found %d", len(c.Values))
			}
			switch c.Name {
			case "key":
				s, err := d.genValue(ty.Key, c)
				if err != nil {
					return nil, err
				}
				key = s.(string)
			case "value":
				s, err := d.genValue(ty.Item, c)
				if err != nil {
					return nil, err
				}
				value = s
				comments = addComments(n.ClosingBraceComment)
			default:
				return nil, d.errorf(c, "unsupported key name %q in map", c.Name)
			}
		}
		if key != "" {
			values = append(values, data{
				Key:      key,
				Value:    value,
				Comments: comments,
			})
		}
		return values, nil
	case typList:
		var values []any
		for _, v := range n.Values {
			if comments := addComments(n.PreComments...); comments != nil {
				continue
			}
			y := *n
			y.Values = []*pbast.Value{v}
			genV, err := d.genValue(ty.Item, &y)
			if err != nil {
				return nil, err
			}
			values = append(values, genV)
		}
		return values, nil
	case typInt:
		if len(n.Values) != 1 {
			return nil, d.errorf(n, "expected 1 value, found %d", len(n.Values))
		}
		s := n.Values[0].Value
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, d.errorf(n, "invalid number %s", s)
		}
		return v, nil
	case typFloat:
		if len(n.Values) != 1 {
			return nil, d.errorf(n, "expected 1 value, found %d", len(n.Values))
		}
		s := n.Values[0].Value
		switch s {
		case "inf", "nan":
			return nil, d.errorf(n, "unexpected float value %s", s)
		}
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, d.errorf(n, "invalid number %s", s)
		}
		return v, nil
	case typBool:
		if len(n.Values) != 1 {
			return nil, d.errorf(n, "expected 1 value, found %d", len(n.Values))
		}
		s := n.Values[0].Value
		switch s {
		case "true":
			return true, nil
		default:
			return false, nil
		}
	case typStr, typAny, typUnion:
		s, _, err := unquote.Unquote(n)
		if err != nil {
			return nil, d.errorf(n, "invalid value to string %s", err.Error())
		}
		return s, nil
	default:
		return nil, fmt.Errorf("unsupported type '%v'", ty.Type)
	}
}

func (d *TextProtoGenerator) errorf(n *pbast.Node, format string, a ...any) error {
	return errors.New(d.locationFormat(n) + ": " + fmt.Sprintf(format, a...))
}

func (d *TextProtoGenerator) locationFormat(n *pbast.Node) string {
	return fmt.Sprintf("%s:%d:%d", d.file, n.Start.Line, n.Start.Column)
}

func addComments(lines ...string) []*ast.Comment {
	var comments []*ast.Comment
	for _, c := range lines {
		if !strings.HasPrefix(c, "#") {
			continue
		}
		comments = append(comments, &ast.Comment{Text: c})
	}
	return comments
}
