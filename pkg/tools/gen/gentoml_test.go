package gen

import (
	"testing"

	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go/pkg/3rdparty/toml"
)

func TestMarshalYamlMapSliceToTOML(t *testing.T) {
	tests := []struct {
		name         string
		data         any
		expectedTOML string
		expectErr    error
	}{
		{
			name: "Simple MapSlice",
			data: &yaml.MapSlice{
				{Key: "b_key", Value: "value1"},
				{Key: "a_key", Value: "value2"},
				{Key: "c_key", Value: "value3"},
			},
			expectedTOML: `b_key = "value1"
a_key = "value2"
c_key = "value3"
`,
			expectErr: nil,
		},
		{
			name: "Nested MapSlice",
			data: &yaml.MapSlice{
				{Key: "outer_key1", Value: yaml.MapSlice{
					{Key: "inner_key1", Value: "inner_value1"},
					{Key: "inner_key2", Value: "inner_value2"},
				}},
				{Key: "outer_key2", Value: "outer_value2"},
			},
			expectedTOML: `[outer_key1]
  inner_key1 = "inner_value1"
  inner_key2 = "inner_value2"
outer_key2 = "outer_value2"
`,
			expectErr: nil,
		},
		{
			name: "Nested MapSlice with Slice",
			data: &yaml.MapSlice{
				{Key: "key_with_slices", Value: []yaml.MapSlice{
					{
						{Key: "inner_key1", Value: "value1"},
					},
					{
						{Key: "inner_key2", Value: "value2"},
					},
				}},
				{Key: "simple_key", Value: "simple_value"},
			},
			expectedTOML: `key_with_slices = [{inner_key1 = "value1"}, {inner_key2 = "value2"}]
simple_key = "simple_value"
`,
			expectErr: nil,
		},
		{
			name: "Nested Map, MapSlice with Slice",
			data: map[string]any{
				"map": yaml.MapSlice{
					{Key: "key_with_slices", Value: []yaml.MapSlice{
						{
							{Key: "inner_key1", Value: "value1"},
						},
						{
							{Key: "inner_key2", Value: "value2"},
						},
					}},
					{Key: "simple_key", Value: "simple_value"},
				},
			},
			expectedTOML: `[map]
  key_with_slices = [{inner_key1 = "value1"}, {inner_key2 = "value2"}]
  simple_key = "simple_value"
`,
			expectErr: nil,
		},
		{
			name: "Simple MapSlice",
			data: &yaml.MapSlice{
				{Key: "b_key", Value: "value1"},
				{Key: "a_key", Value: map[string]string{
					"a_a_key": "value2",
				}},
				{Key: "c_key", Value: "value3"},
			},
			expectedTOML: `b_key = "value1"

[a_key]
  a_a_key = "value2"
c_key = "value3"
`,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tomlData, err := toml.Marshal(tt.data)
			if err != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if got := string(tomlData); got != tt.expectedTOML {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expectedTOML, got)
			}
		})
	}
}
