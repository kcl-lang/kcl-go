package gen

import (
	"bytes"
	"io"

	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/loader"
)

const (
	manifestPkgPath      = "manifests"
	yamlStreamOutputFunc = "manifests.yaml_stream(items)\n"
)

func (k *kclGenerator) genKclFromYaml(w io.Writer, filename string, src interface{}) error {
	code, err := loader.ReadSource(filename, src)
	if err != nil {
		return err
	}
	// convert yaml data to kcl
	result, err := convertKclFromYamlStreamString(code)
	if err != nil {
		return err
	}
	// generate kcl code
	if len(result) == 0 {
		return k.genKcl(w, kclFile{Config: []config{
			{Data: []data{}},
		}})
	}
	if len(result) == 1 {
		return k.genKcl(w, kclFile{Config: []config{
			{Data: result[0]},
		}})
	} else {
		var value []config
		for _, r := range result {
			value = append(value, config{
				Data: r,
			})
		}
		return k.genKcl(
			w,
			kclFile{
				Imports: []kImport{{PkgPath: manifestPkgPath}},
				Data: []data{{
					Key:   "items",
					Value: value,
				}},
				ExtraCode: yamlStreamOutputFunc,
			},
		)
	}
}

func convertKclFromYaml(yamlData *yaml.MapSlice) []data {
	var result []data
	for _, item := range *yamlData {
		key, ok := item.Key.(string)
		if !ok {
			continue
		}
		switch value := item.Value.(type) {
		case *yaml.MapSlice:
			result = append(result, data{
				Key:   key,
				Value: convertKclFromYaml(value),
			})
		case yaml.MapSlice:
			result = append(result, data{
				Key:   key,
				Value: convertKclFromYaml(&value),
			})
		case []interface{}:
			var vals []interface{}
			for _, v := range value {
				switch v := v.(type) {
				case *yaml.MapSlice:
					vals = append(vals, convertKclFromYaml(v))
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
	result, err := convertKclFromYamlStreamString(byteData)
	if err != nil {
		return nil, err
	}
	if len(result) >= 1 {
		return result[0], err
	}
	return nil, nil
}

func convertKclFromYamlStreamString(byteData []byte) ([][]data, error) {
	byteData = bytes.ReplaceAll(byteData, []byte("\r\n"), []byte("\n"))
	var result [][]data
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
		result = append(result, d)
	}
	return result, nil
}
