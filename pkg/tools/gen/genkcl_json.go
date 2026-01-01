package gen

import (
	"io"

	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/source"
)

func (k *kclGenerator) genKclFromJsonData(w io.Writer, filename string, src any) error {
	code, err := source.ReadSource(filename, src)
	if err != nil {
		return err
	}

	// as yaml can be viewed as a superset of json,
	// we can handle json data like yaml.
	yamlData := &yaml.MapSlice{}
	if err = yaml.UnmarshalWithOptions(code, yamlData, yaml.UseOrderedMap(), yaml.UseJSONUnmarshaler()); err != nil {
		return err
	}

	// convert to kcl
	result := convertKclFromYaml(yamlData)

	// generate kcl code
	return k.genKcl(w, kclFile{Config: []config{
		{Data: result},
	}})
}
