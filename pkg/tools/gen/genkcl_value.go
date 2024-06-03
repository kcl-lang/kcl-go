package gen

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type printer struct {
	indent          uint
	writer          io.Writer
	listInOneLine   bool
	configInOneLine bool
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

func (p *printer) walkValue(v any) error {
	if v == nil {
		p.write("None")
		return nil
	}
	ty := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	switch ty.Kind() {
	case reflect.Bool:
		if val.Bool() {
			p.write("True")
		} else {
			p.write("False")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.write(strconv.FormatInt(val.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.write(strconv.FormatUint(val.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		p.write(strconv.FormatFloat(val.Float(), 'f', -1, ty.Bits()))
	case reflect.Array, reflect.Slice:
		p.write("[")
		if !p.listInOneLine {
			p.writeIndentWithNewLine()
		}
		for i := 0; i < val.Len(); i++ {
			if i > 0 {
				if !p.listInOneLine {
					p.writeNewLine()
				} else {
					p.write(", ")
				}
			}
			if err := p.walkValue(val.Index(i).Interface()); err != nil {
				return err
			}
		}
		if !p.listInOneLine {
			p.writeDedentWithNewLine()
		}
		p.write("]")
	case reflect.Map:
		p.write("{")
		if !p.configInOneLine {
			p.writeIndentWithNewLine()
		}
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
		for i, key := range keys {
			if i > 0 {
				if !p.configInOneLine {
					p.writeNewLine()
				} else {
					p.write(", ")
				}
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
		if !p.configInOneLine {
			p.writeDedentWithNewLine()
		}
		p.write("}")
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
		p.write("{")
		p.writeIndentWithNewLine()
		for i := 0; i < ty.NumField(); i++ {
			if i > 0 {
				if !p.configInOneLine {
					p.writeNewLine()
				} else {
					p.write(", ")
				}
			}
			field := ty.Field(i)
			fieldValue := val.Field(i).Interface()
			p.write(formatName(field.Name) + " = ")
			if err := p.walkValue(fieldValue); err != nil {
				return err
			}
		}
		p.writeDedentWithNewLine()
		p.write("}")
	default:
		return fmt.Errorf("invalid value type to kcl: %v", ty)
	}
	return nil
}

// Generate KCL code from go value
func GenKclFromValue(w io.Writer, v any) error {
	p := &printer{
		writer: w,
	}
	return p.walkValue(v)
}

// Generate KCL code from go value
func GenKclFromValueWithIndent(w io.Writer, v any, indent uint) error {
	p := &printer{
		indent: indent,
		writer: w,
	}
	return p.walkValue(v)
}
