package gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

const (
	KclNoneValue  = "None"
	KclTrueValue  = "True"
	KclFalseValue = "False"
)

// Marshaler is the interface implemented by types that can marshal
// themselves into valid KCL.
type Marshaler interface {
	MarshalKcl() ([]byte, error)
}

// Marshal returns a KCL representation of the Go value.
func Marshal(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	p := &printer{
		writer: buf,
	}
	if err := p.walkValue(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var (
	marshalTy = reflect.TypeOf((*Marshaler)(nil)).Elem()
)

func isMarshalTy(rv reflect.Value) bool {
	return rv.Type().Implements(marshalTy)
}

type printer struct {
	indent       uint
	writer       io.Writer
	listInline   bool
	configInline bool
}

func (p *printer) writeIndent() {
	p.enter()
}

func (p *printer) writeDedent() {
	p.leave()
}

func (p *printer) writeNewLine() {
	p.writeln("")
}

func (p *printer) writeIndentWithNewLine() {
	p.writeIndent()
	p.writeln("")
}

func (p *printer) writeDedentWithNewLine() {
	p.writeDedent()
	p.writeln("")
}

func (p *printer) write(text string) {
	p.writer.Write([]byte(text))
}

func (p *printer) writeln(text string) {
	p.write(text)
	p.write("\n")
	p.fill("")
}

func (p *printer) fill(text string) {
	p.write(strings.Repeat("    ", int(p.indent)))
	p.write(text)
}

func (p *printer) enter() {
	p.indent += 1
}

func (p *printer) leave() {
	if p.indent > 0 {
		p.indent -= 1
	}
}

func (p *printer) writeListBegin() {
	p.write("[")
	if !p.listInline {
		p.writeIndentWithNewLine()
	}
}

func (p *printer) writeListEnd() {
	if !p.listInline {
		p.writeDedentWithNewLine()
	}
	p.write("]")
}

func (p *printer) writeListSep() {
	if !p.listInline {
		p.writeNewLine()
	} else {
		p.write(", ")
	}
}

func (p *printer) writeConfigBegin() {
	p.write("{")
	if !p.configInline {
		p.writeIndentWithNewLine()
	}
}

func (p *printer) writeConfigEnd() {
	if !p.configInline {
		p.writeDedentWithNewLine()
	}
	p.write("}")
}

func (p *printer) writeConfigSep() {
	if !p.configInline {
		p.writeNewLine()
	} else {
		p.write(", ")
	}
}

func (p *printer) walkValue(v any) error {
	if v == nil {
		p.write(KclNoneValue)
		return nil
	}
	ty := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	switch v := val.Interface().(type) {
	case Marshaler:
		s, err := v.MarshalKcl()
		if err != nil {
			return err
		}
		if s == nil {
			return errors.New("MarshalKcl returned nil and no error")
		}
		p.writer.Write(s)
		return nil
	case time.Duration:
		p.write(v.String())
		return nil
	case json.Number:
		n, _ := val.Interface().(json.Number)
		if n == "" {
			p.write("0")
			return nil
		} else if v, err := n.Int64(); err == nil {
			p.walkValue(v)
			return nil
		} else if v, err := n.Float64(); err == nil {
			p.walkValue(v)
			return nil
		}
		return fmt.Errorf("unable to convert %q to int64 or float64", n)
	case yaml.MapSlice:
		n, _ := val.Interface().(yaml.MapSlice)
		p.writeConfigBegin()
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
		for i, item := range n {
			if i > 0 {
				p.writeConfigSep()
			}
			p.walkValue(item.Key)
			p.write(" = ")
			p.walkValue(item.Value)
		}
		p.writeConfigEnd()
	}

	switch ty.Kind() {
	case reflect.Bool:
		if val.Bool() {
			p.write(KclTrueValue)
		} else {
			p.write(KclFalseValue)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.write(strconv.FormatInt(val.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.write(strconv.FormatUint(val.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		p.write(strconv.FormatFloat(val.Float(), 'f', -1, ty.Bits()))
	case reflect.Array, reflect.Slice:
		p.writeListBegin()
		for i := 0; i < val.Len(); i++ {
			if i > 0 {
				p.writeListSep()
			}
			if err := p.walkValue(val.Index(i).Interface()); err != nil {
				return err
			}
		}
		p.writeListEnd()
	case reflect.Map:
		p.writeConfigBegin()
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
		for i, key := range keys {
			if i > 0 {
				p.writeConfigSep()
			}
			if key.Kind() == reflect.String {
				p.write(formatName(key.String()))
			} else {
				if err := p.walkValue(key.Interface()); err != nil {
					return err
				}
			}
			p.write(" = ")
			if err := p.walkValue(val.MapIndex(key).Interface()); err != nil {
				return err
			}
		}
		p.writeConfigEnd()
	case reflect.String:
		value := val.String()
		if isStringEscaped(value) {
			if value[len(value)-1] == '"' {
				// if the string ends with '"' then we need to add a space after the closing triple quote
				p.write(fmt.Sprintf(`r"""%s """`, value))
			} else {
				p.write(fmt.Sprintf(`r"""%s"""`, value))
			}
		} else {
			p.write(strconv.Quote(value))
		}
	case reflect.Struct:
		p.writeConfigBegin()
		fields := TypeFields(ty, val)
		for i, field := range fields {
			if i > 0 {
				if !p.configInline {
					p.writeNewLine()
				} else {
					p.write(", ")
				}
			}
			keyName := field.name
			opts := field.opts
			if opts != nil {
				if field.opts.skip {
					continue
				}
				if opts.omitempty && IsEmpty(field.val) {
					continue
				}
				if opts.omitzero && IsZero(field.val) {
					continue
				}
				if opts.name != "" {
					keyName = opts.name
				}
			}
			p.write(formatName(keyName) + " = ")
			if err := p.walkValue(field.val.Interface()); err != nil {
				return err
			}
		}
		p.writeConfigEnd()
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			p.write(KclNoneValue)
		} else {
			return p.walkValue(val.Elem())
		}
	default:
		return fmt.Errorf("invalid value type to kcl: %v", ty)
	}
	return nil
}
