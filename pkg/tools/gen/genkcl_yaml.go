package gen

import (
	"bytes"
	"io"

	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/kcl"
)

func (k *kclGenerator) genKclFromYaml(w io.Writer, filename string, src interface{}) error {
	code, err := readSource(filename, src)
	if err != nil {
		return err
	}
	// convert yaml data to kcl
	result, err := convertKclFromYamlString(code)
	if err != nil {
		return err
	}
	// generate kcl code
	return k.genKcl(w, kclFile{Config: []config{
		{Data: result},
	}})
}

func convertKclFromYaml(yamlData *yaml.MapSlice) []data {
	var result []data
	for _, item := range *yamlData {
		key, ok := item.Key.(string)
		if !ok {
			continue
		}
		switch value := item.Value.(type) {
		case yaml.MapSlice:
			result = append(result, data{
				Key:   key,
				Value: convertKclFromYaml(&value),
			})
		case []interface{}:
			var vals []interface{}
			for _, v := range value {
				switch v := v.(type) {
				case yaml.MapSlice:
					vals = append(vals, convertKclFromYaml(&v))
				default:
					vals = append(vals, v)
				}
			}
			result = append(result, data{Key: key, Value: vals})
		default:
			result = append(result, data{Key: key, Value: value})
		}
	}
	return result
}

func convertKclFromYamlString(byteData []byte) ([]data, error) {
	byteData = bytes.ReplaceAll(byteData, []byte("\r\n"), []byte("\n"))
	var result []data
	// split yaml with ‘---’
	items, err := kcl.SplitDocuments(string(byteData))
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		yamlData := &yaml.MapSlice{}
		if err := yaml.UnmarshalWithOptions([]byte(item), yamlData, yaml.UseOrderedMap()); err != nil {
			return nil, err
		}
		// convert yaml data to kcl
		d := convertKclFromYaml(yamlData)
		result = append(result, d...)
	}
	return result, nil
}
