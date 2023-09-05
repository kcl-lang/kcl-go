package gen

import (
	"github.com/goccy/go-yaml"
	"io"
)

func (k *kclGenerator) genKclFromJsonData(w io.Writer, filename string, src interface{}) error {
	code, err := readSource(filename, src)
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
	return k.genKcl(w, kclFile{Data: result})
}
