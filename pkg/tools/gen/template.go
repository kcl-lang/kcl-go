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
	//go:embed templates/kcl/globals.gotmpl
	globalsTmpl string
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
	"condExpr":              condExpr,
	"indentLines":           indentLines,
	"isKclData": func(v any) bool {
		_, ok := v.([]data)
		return ok
	},
	"isKclConfig": func(v any) bool {
		_, ok := v.(config)
		return ok
	},
	"isArray": func(v any) bool {
		switch v.(type) {
		case []data:
			return true
		case []config:
			return true
		case []any:
			return true
		default:
			return false
		}
	},
}
var tmpl *template.Template = &template.Template{}

func init() {
	// add "include" function. It works like "template" but can be used in pipeline.
	funcs["include"] = func(name string, data any) (string, error) {
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
	tmpl = addTemplate(tmpl, "globals", globalsTmpl)
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

func formatValueWithEscape(v any, escape bool) string {
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

func formatValue(v any) string {
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

// isKclKeyword reports whether the supplied identifier is reserved by the
// KCL grammar and therefore cannot be used as a bare KCL identifier. The
// keyword set matches the 27 symbols declared in
// crates/span/src/symbol.rs at HEAD.
func isKclKeyword(name string) bool {
	_, ok := kclKeywords[name]
	return ok
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

// condExpr renders a validation as a single inline KCL boolean expression,
// joining every constraint it carries with `and`. It exists because JSON
// Schema `if`/`then`/`else` needs its operands embedded inside a conditional
// check statement, where the multi-line statement form used by the rest of
// the validator template does not fit.
//
// The returned condResult has two pieces:
//   - Expr holds the rendered KCL boolean expression.
//   - OptionalGuards lists the names of any optional fields the expression
//     references. The template ANDs them into the surrounding condition so
//     the check is short-circuited when the field is undefined, mirroring
//     how the regular validator template appends `if <field>` for optional
//     fields.
//
// Both fields are empty when the validation holds no constraint that maps
// onto a KCL expression.
type condResult struct {
	Expr           string
	OptionalGuards []string
}

func condExpr(v *validation) condResult {
	if v == nil {
		return condResult{}
	}
	name := formatName(v.Name)
	var parts []string
	var guards []string
	// Only emit a top-level guard when the validation targets a single,
	// specific field (i.e. SubConstraints is empty AND a non-empty Name).
	// When this is a top-level if/then/else expression whose actual
	// references are carried by SubConstraints, the per-field guards
	// collected below are what matter.
	if v.Name != "" && len(v.SubConstraints) == 0 && !v.Required {
		guards = append(guards, formatName(v.Name))
	}
	if v.TypeName != "" {
		parts = append(parts, fmt.Sprintf("typeof(%s) == %q", name, v.TypeName))
	}
	if v.ConstValue != nil {
		parts = append(parts, fmt.Sprintf("%s == %s", name, formatValue(v.ConstValue)))
	}
	if v.Minimum != nil {
		op := ">="
		if v.ExclusiveMinimum {
			op = ">"
		}
		parts = append(parts, fmt.Sprintf("%s %s %v", name, op, *v.Minimum))
	}
	if v.Maximum != nil {
		op := "<="
		if v.ExclusiveMaximum {
			op = "<"
		}
		parts = append(parts, fmt.Sprintf("%s %s %v", name, op, *v.Maximum))
	}
	if v.MinLength != nil {
		parts = append(parts, fmt.Sprintf("len(%s) >= %d", name, *v.MinLength))
	}
	if v.MaxLength != nil {
		parts = append(parts, fmt.Sprintf("len(%s) <= %d", name, *v.MaxLength))
	}
	if v.Regex != nil {
		parts = append(parts, fmt.Sprintf(`regex.match(%s, r"%s")`, name, v.Regex))
	}
	if v.MultiplyOf != nil {
		parts = append(parts, fmt.Sprintf("multiplyof(%s, %d)", name, *v.MultiplyOf))
	}
	for _, sub := range v.SubConstraints {
		if sub == nil {
			continue
		}
		subRes := condExpr(sub)
		if subRes.Expr == "" {
			continue
		}
		parts = append(parts, subRes.Expr)
		guards = append(guards, subRes.OptionalGuards...)
	}
	return condResult{
		Expr:           strings.Join(parts, " and "),
		OptionalGuards: guards,
	}
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
