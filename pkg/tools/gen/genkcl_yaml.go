package gen

import (
	"bytes"
	"io"

	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/source"
)

const (
	manifestPkgPath      = "manifests"
	yamlStreamOutputFunc = "manifests.yaml_stream(items)\n"
)

func (k *kclGenerator) genKclFromYaml(w io.Writer, filename string, src any) error {
	code, err := source.ReadSource(filename, src)
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
		result = append(result, data{
			Key:   key,
			Value: convertKclValue(item.Value),
		})
	}
	return result
}

// convertKclValue converts an arbitrary YAML/JSON value (as produced by
// goccy/go-yaml's MapSlice decoder) into a value ready for walkValue.
//
// MapSlice values are unwrapped into their underlying []data representation
// so the renderer can produce `{ key = value, ... }` blocks. Lists are walked
// element-by-element so nested lists (e.g. `[[{...}]]`) are recursively
// flattened in the same way instead of being kept as raw `[]any` slices.
func convertKclValue(value any) any {
	switch v := value.(type) {
	case *yaml.MapSlice:
		return convertKclFromYaml(v)
	case yaml.MapSlice:
		vs := v
		return convertKclFromYaml(&vs)
	case []any:
		vals := make([]any, 0, len(v))
		for _, item := range v {
			vals = append(vals, convertKclValue(item))
		}
		return vals
	case nil:
		return nil
	default:
		return value
	}
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
