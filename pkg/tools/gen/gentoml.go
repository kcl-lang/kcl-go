package gen

import (
	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/3rdparty/toml"
)

// Marshal returns a TOML representation of the Go value.
func MarshalTOML(data *yaml.MapSlice) ([]byte, error) {
	return toml.Marshal(data)
}

// MarshalYamlMapSliceToTOML convert an ordered yaml data to an ordered toml
func MarshalYamlMapSliceToTOML(data *yaml.MapSlice) ([]byte, error) {
	return toml.Marshal(data)
}
