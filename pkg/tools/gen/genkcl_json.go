package gen

import (
	"encoding/json"
	"io"
	"sort"
)

func (k *kclGenerator) genKclFromJsonData(w io.Writer, filename string, src interface{}) error {
	code, err := readSource(filename, src)
	if err != nil {
		return err
	}
	jsonData := &map[string]interface{}{}
	if err = json.Unmarshal(code, jsonData); err != nil {
		return err
	}

	// convert json data to kcl
	result := convertKclFromJson(jsonData)

	// generate kcl code
	return k.genKcl(w, kclFile{Data: result})
}

func convertKclFromJson(jsonData *map[string]interface{}) []data {
	var result []data
	for key, value := range *jsonData {
		switch value := value.(type) {
		case map[string]interface{}:
			result = append(result, data{
				Key:   key,
				Value: convertKclFromJson(&value),
			})
		case []interface{}:
			var vals []interface{}
			for _, v := range value {
				switch v := v.(type) {
				case map[string]interface{}:
					vals = append(vals, convertKclFromJson(&v))
				default:
					vals = append(vals, v)
				}
			}
			result = append(result, data{Key: key, Value: vals})
		default:
			result = append(result, data{Key: key, Value: value})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result
}
