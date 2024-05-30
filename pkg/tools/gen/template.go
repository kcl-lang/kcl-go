package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

var (
	//go:embed templates/kcl/config.gotmpl
	configTmpl string
	//go:embed templates/kcl/data.gotmpl
	dataTmpl string
	//go:embed templates/kcl/document.gotmpl
	documentTmpl string
	//go:embed templates/kcl/header.gotmpl
	headerTmpl string
	//go:embed templates/kcl/validator.gotmpl
	validatorTmpl string
	//go:embed templates/kcl/schema.gotmpl
	schemaTmpl string
	//go:embed templates/kcl/index.gotmpl
	indexTmpl string
)

var funcs = template.FuncMap{
	"formatType":  formatType,
	"formatValue": formatValue,
	"formatName":  formatName,
	"indentLines": indentLines,
	"isKclData": func(v interface{}) bool {
		_, ok := v.([]data)
		return ok
	},
	"isKclConfig": func(v interface{}) bool {
		_, ok := v.(config)
		return ok
	},
	"isArray": func(v interface{}) bool {
		switch v.(type) {
		case []data:
			return true
		case []config:
			return true
		case []interface{}:
			return true
		default:
			return false
		}
	},
}
var tmpl *template.Template = &template.Template{}

func init() {
	// add "include" function. It works like "template" but can be used in pipeline.
	funcs["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	tmpl = addTemplate(tmpl, "config", configTmpl)
	tmpl = addTemplate(tmpl, "data", dataTmpl)
	tmpl = addTemplate(tmpl, "document", documentTmpl)
	tmpl = addTemplate(tmpl, "header", headerTmpl)
	tmpl = addTemplate(tmpl, "validator", validatorTmpl)
	tmpl = addTemplate(tmpl, "schema", schemaTmpl)
	tmpl = addTemplate(tmpl, "index", indexTmpl)
	tmpl = tmpl.Funcs(funcs)
}

func (k *kclGenerator) genKcl(w io.Writer, s kclFile) error {
	return tmpl.Execute(w, s)
}

func addTemplate(tmpl *template.Template, name, data string) *template.Template {
	newTmpl := template.Must(template.New(name).Funcs(funcs).Parse(data))
	return template.Must(tmpl.AddParseTree(name, newTmpl.Tree))
}

func formatType(t typeInterface) string {
	if t != nil {
		return t.Format()
	}
	return typAny
}

func formatValue(v interface{}) string {
	if v == nil {
		return "None"
	}
	switch value := v.(type) {
	case string:
		if isStringEscaped(value) {
			if value[len(value)-1] == '"' {
				// if the string ends with '"' then we need to add a space after the closing triple quote
				return fmt.Sprintf(`r"""%s """`, value)
			}
			return fmt.Sprintf(`r"""%s"""`, value)
		}

		return fmt.Sprintf(`"%s"`, value)
	case bool:
		if value {
			return "True"
		}
		return "False"
	case map[string]bool:
		return formatMap(value)
	case map[string]float32:
		return formatMap(value)
	case map[string]float64:
		return formatMap(value)
	case map[string]int:
		return formatMap(value)
	case map[string]string:
		return formatMap(value)
	case map[string]interface{}:
		return formatMap(value)
	case []interface{}:
		var s strings.Builder
		for i, item := range value {
			if i != 0 {
				s.WriteString(", ")
			}
			s.WriteString(formatValue(item))
		}
		return "[" + s.String() + "]"
	default:
		return fmt.Sprintf("%v", value)
	}
}

func formatMap[V any](value map[string]V) string {
	var s strings.Builder
	for _, key := range getSortedKeys(value) {
		if s.Len() != 0 {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%s: %s", formatValue(key), formatValue(value[key])))
	}
	return "{" + s.String() + "}"
}

var kclKeywords = map[string]struct{}{
	"True":      {},
	"False":     {},
	"None":      {},
	"Undefined": {},
	"import":    {},
	"and":       {},
	"or":        {},
	"in":        {},
	"is":        {},
	"not":       {},
	"as":        {},
	"if":        {},
	"else":      {},
	"elif":      {},
	"for":       {},
	"schema":    {},
	"mixin":     {},
	"protocol":  {},
	"check":     {},
	"assert":    {},
	"all":       {},
	"any":       {},
	"map":       {},
	"filter":    {},
	"lambda":    {},
	"rule":      {},
}

var validNameRegexp = regexp.MustCompile(`\$?^[a-zA-Z_][a-zA-Z0-9_]*$`)

func formatName(name string) string {
	if _, ok := kclKeywords[name]; ok {
		return fmt.Sprintf("$%s", name)
	}

	if !validNameRegexp.MatchString(name) {
		return fmt.Sprintf(`"%s"`, name)
	}

	return name
}

func indentLines(s, indent string) string {
	s = strings.Replace(s, "\r\n", "\n", -1)
	var b strings.Builder
	raw := false
	for i, line := range strings.Split(s, "\n") {
		if i != 0 {
			b.WriteString("\n")
		}
		if line == "" {
			continue
		}

		if raw {
			if strings.HasSuffix(line, `"""`) {
				raw = false
			}
			b.WriteString(line)
			continue
		}

		if strings.Contains(line, `r"""`) && !strings.HasSuffix(line, `"""`) {
			raw = true
		}

		b.WriteString(indent)
		b.WriteString(line)
	}

	return b.String()
}

func isStringEscaped(s string) bool {
	_, err := strconv.Unquote(`"` + s + `"`)
	return err != nil || strings.Contains(s, "$")
}
