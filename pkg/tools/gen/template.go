package gen

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"
)

var (
	//go:embed templates/document.gotmpl
	documentTmpl string
	//go:embed templates/header.gotmpl
	headerTmpl string
	//go:embed templates/validator.gotmpl
	validatorTmpl string
	//go:embed templates/schema.gotmpl
	schemaTmpl string
)

var funcs = template.FuncMap{
	"formatType":  formatType,
	"formatValue": formatValue,
	"formatName":  formatName,
	"formatDoc":   formatDoc,
}

func (k *kclGenerator) genKclSchema(w io.Writer, s kclSchema) error {
	tmpl := &template.Template{}
	tmpl = addTemplate(tmpl, "document", documentTmpl)
	tmpl = addTemplate(tmpl, "header", headerTmpl)
	tmpl = addTemplate(tmpl, "validator", validatorTmpl)
	tmpl = addTemplate(tmpl, "schema", schemaTmpl)
	return tmpl.Funcs(funcs).Execute(w, s)
}

func addTemplate(tmpl *template.Template, name, data string) *template.Template {
	newTmpl := template.Must(template.New(name).Funcs(funcs).Parse(data))
	return template.Must(tmpl.AddParseTree(name, newTmpl.Tree))
}

func formatType(t typeInterface) string {
	return t.Format()
}

func formatValue(v interface{}) string {
	if v == nil {
		return "None"
	}
	switch value := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", value)
	case bool:
		if value {
			return "True"
		}
		return "False"
	default:
		return fmt.Sprintf("%v", value)
	}
}

func formatName(name string) string {
	return name
}

func formatDoc(doc, indent string) string {
	doc = strings.Replace(doc, "\r\n", "\n", -1)
	return indent + strings.Replace(doc, "\n", "\n"+indent, -1)
}
