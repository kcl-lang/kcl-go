package gen

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

var (
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
}

func (k *kclGenerator) genKclSchema(w io.Writer, s kclSchema) error {
	tmpl := &template.Template{}
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
