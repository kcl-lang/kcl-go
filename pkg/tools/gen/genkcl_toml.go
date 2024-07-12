package gen

import (
	"io"
	"reflect"
	"strings"

	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/3rdparty/toml"
	"kcl-lang.io/kcl-go/pkg/source"
)

func (k *kclGenerator) genKclFromToml(w io.Writer, filename string, src interface{}) error {
	code, err := source.ReadSource(filename, src)
	if err != nil {
		return err
	}

	// as yaml can be viewed as a superset of json,
	// we can handle json data like yaml.
	yamlData := &yaml.MapSlice{}
	mappingData := make(map[string]any)
	meta, err := toml.Decode(string(code), &mappingData)
	if err != nil {
		return err
	}
	for _, key := range meta.Keys() {
		key := key.String()
		mapSliceSet(yamlData, key, mapGet(mappingData, key))
	}
	// convert to kcl
	result := convertKclFromYaml(yamlData)

	// generate kcl code
	return k.genKcl(w, kclFile{Config: []config{
		{Data: result},
	}})
}

func mapGet(m interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	value := reflect.ValueOf(m)
	if value.Kind() != reflect.Map {
		return nil
	}
	for _, k := range keys {
		elem := value.MapIndex(reflect.ValueOf(k))
		if !elem.IsValid() {
			return nil
		}
		value = elem.Elem()
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}
	}
	return value.Interface()
}

func mapSliceSet(m *yaml.MapSlice, key string, value interface{}) {
	keys := strings.Split(key, ".")
	currentMap := m
	for i, k := range keys {
		found := false
		for j, item := range *currentMap {
			if kkey, ok := item.Key.(string); ok && kkey == k {
				found = true
				if i == len(keys)-1 {
					(*currentMap)[j] = yaml.MapItem{Key: k, Value: value}
					return
				} else {
					if mp, ok := item.Value.(yaml.MapSlice); ok {
						currentMap = &mp
						break
					} else if mp, ok := item.Value.(*yaml.MapSlice); ok {
						currentMap = mp
						break
					}
				}
			}
		}
		if !found {
			if i == len(keys)-1 {
				if reflect.TypeOf(value).Kind() == reflect.Map {
					newMap := make(yaml.MapSlice, 0)
					*currentMap = append(*currentMap, yaml.MapItem{Key: k, Value: &newMap})
					currentMap = &newMap
				} else {
					*currentMap = append(*currentMap, yaml.MapItem{Key: k, Value: value})
				}
			} else {
				newMap := make(yaml.MapSlice, 0)
				*currentMap = append(*currentMap, yaml.MapItem{Key: k, Value: &newMap})
				currentMap = &newMap
			}
		}
	}
}
