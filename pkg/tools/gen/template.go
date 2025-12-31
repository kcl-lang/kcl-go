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
	"formatType":            formatType,
	"formatValue":           formatValue,
	"formatValueWithEscape": formatValueWithEscape,
	"formatName":            formatName,
	"indentLines":           indentLines,
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
		result := buf.String()
		// Fix spacing issues in anyOf expressions
		result = strings.ReplaceAll(result, ")or ", ") or ")
		result = strings.ReplaceAll(result, ")or  ", ") or ")
		result = strings.ReplaceAll(result, ")if ", ") if ")
		return result, nil
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
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return err
	}
	result := buf.String()
	// Fix spacing issues in anyOf expressions
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		// Only process validation lines (lines that start with spaces and contain validation patterns)
		// Check if line looks like a validation (starts with spaces, contains "or" between expressions)
		if strings.HasPrefix(line, "        ") && (strings.Contains(line, ")or ") || strings.Contains(line, "match(")) {
			// Fix spacing around "or" in anyOf expressions
			// Replace "or" followed by multiple spaces with "or " followed by single space
			re := regexp.MustCompile(`or +`)
			line = re.ReplaceAllString(line, "or ")
			// Also fix )or case
			line = strings.ReplaceAll(line, ")or ", ") or ")
			lines[i] = line
		}
	}
	result = strings.Join(lines, "\n")
	result = strings.ReplaceAll(result, ")if ", ") if ")
	_, err := io.WriteString(w, result)
	return err
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

func formatValueWithEscape(v interface{}, escape bool) string {
	var buf bytes.Buffer
	p := &printer{
		listInline:   true,
		configInline: true,
		writer:       &buf,
		escape:       escape,
	}
	err := p.walkValue(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	} else {
		return buf.String()
	}
}

func formatValue(v interface{}) string {
	return formatValueWithEscape(v, true)
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
	"type":      {},
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
