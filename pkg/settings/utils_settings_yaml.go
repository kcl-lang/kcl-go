// Copyright The KCL Authors. All rights reserved.

package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
	tools_list "kcl-lang.io/kcl-go/pkg/tools/list"
)

type SettingsFile struct {
	Filename string           `yaml:"-"`
	Config   ConfigStruct     `yaml:"kcl_cli_configs"`
	Options  []KeyValueStruct `yaml:"kcl_options"`
}

type ConfigStruct struct {
	InputFile  []string `yaml:"file"`
	InputFiles []string `yaml:"files"`
	Output     string   `yaml:"output"`

	Overrides    []string `yaml:"overrides"`
	PathSelector []string `yaml:"path_selector"`

	StrictRangeCheck      bool              `yaml:"strict_range_check"`
	DisableNone           bool              `yaml:"disable_none"`
	Verbose               int               `yaml:"verbose"`
	Debug                 bool              `yaml:"debug"`
	PackageMaps           map[string]string `yaml:"package_maps"`
	SortKeys              bool              `yaml:"sort_keys"`
	ShowHidden            bool              `yaml:"show_hidden"`
	IncludeSchemaTypePath bool              `yaml:"include_schema_type_path"`
}

type KeyValueStruct struct {
	Key   string `yaml:"key"`
	Value any    `yaml:"value"`
	// Store the original YAML value node to preserve order (the node under key: "value")
	originalValueNode *yaml.Node `yaml:"-"`
}

// enhanceWithOrderInfo adds order-preserving information to settings from YAML nodes
func enhanceWithOrderInfo(settings *SettingsFile, rootNode *yaml.Node) error {
	if len(rootNode.Content) == 0 {
		return nil
	}

	// Find the kcl_options section in the node tree
	optionsNode := findOptionsNode(rootNode.Content[0])
	if optionsNode == nil {
		return nil
	}

	// Map the original value nodes to the parsed options
	if optionsNode.Kind == yaml.SequenceNode {
		for i, optionNode := range optionsNode.Content {
			if i < len(settings.Options) {
				if optionNode != nil && optionNode.Kind == yaml.MappingNode {
					if val := getMappingValueNode(optionNode, "value"); val != nil {
						settings.Options[i].originalValueNode = val
					}
				}
			}
		}
	}

	return nil
}

// findOptionsNode finds the kcl_options node in the YAML tree
func findOptionsNode(node *yaml.Node) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Value == "kcl_options" {
			return valueNode
		}
	}

	return nil
}

// getMappingValueNode returns the value node for the provided key in a mapping node.
func getMappingValueNode(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]
		if k != nil && k.Value == key {
			return v
		}
	}
	return nil
}

// nodeToOrderedJSON converts any YAML node to JSON while preserving order
func nodeToOrderedJSON(node *yaml.Node) (string, error) {
	switch node.Kind {
	case yaml.MappingNode:
		var pairs []string
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			key := keyNode.Value
			valueJSON, err := nodeToOrderedJSON(valueNode)
			if err != nil {
				continue
			}

			pairs = append(pairs, fmt.Sprintf("%q:%s", key, valueJSON))
		}
		return "{" + strings.Join(pairs, ",") + "}", nil

	case yaml.SequenceNode:
		var items []string
		for _, itemNode := range node.Content {
			itemJSON, err := nodeToOrderedJSON(itemNode)
			if err != nil {
				continue
			}
			items = append(items, itemJSON)
		}
		return "[" + strings.Join(items, ",") + "]", nil

	case yaml.ScalarNode:
		var value any
		if err := node.Decode(&value); err != nil {
			return fmt.Sprintf("%q", node.Value), nil
		}
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Sprintf("%q", node.Value), nil
		}
		return string(jsonBytes), nil

	default:
		return "null", nil
	}
}

func LoadFile(filename string, src any) (f *SettingsFile, err error) {
	if !filepath.IsAbs(filename) {
		if s, _ := filepath.Abs(filename); s != "" {
			filename = s
		}
	}
	if src == nil {
		src, err = os.ReadFile(filename)
		if err != nil {
			return
		}
	}

	var code string
	switch src := src.(type) {
	case []byte:
		code = string(src)
	case string:
		code = src
	case io.Reader:
		d, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		code = string(d)
	default:
		return nil, fmt.Errorf("unsupported src type: %T", src)
	}
	if code == "" {
		return &SettingsFile{Filename: filename}, nil
	}

	// First parse with Node API to preserve order
	var rootNode yaml.Node
	if err := yaml.Unmarshal([]byte(code), &rootNode); err != nil {
		return nil, err
	}

	// Also parse normally for basic structure
	var settings SettingsFile
	if err := yaml.Unmarshal([]byte(code), &settings); err != nil {
		return nil, err
	}

	// Enhance the settings with order-preserving information
	if err := enhanceWithOrderInfo(&settings, &rootNode); err != nil {
		// If enhancement fails, continue with regular settings
		// The order preservation is a best-effort feature
	}

	settings.Filename = filename
	return &settings, nil
}

func (settings *SettingsFile) To_ExecProgramArgs() *gpyrpc.ExecProgramArgs {
	args := &gpyrpc.ExecProgramArgs{
		KFilenameList: []string{},
		KCodeList:     []string{},

		Args:      []*gpyrpc.Argument{},
		Overrides: []string{},

		DisableYamlResult: false,
		PrintOverrideAst:  false,

		StrictRangeCheck:      settings.Config.StrictRangeCheck,
		DisableNone:           settings.Config.DisableNone,
		Verbose:               int32(settings.Config.Verbose),
		Debug:                 0,
		SortKeys:              settings.Config.SortKeys,
		ShowHidden:            settings.Config.ShowHidden,
		IncludeSchemaTypePath: settings.Config.IncludeSchemaTypePath,
		PathSelector:          settings.Config.PathSelector,
	}
	if settings.Config.Debug {
		args.Debug = 1
	}

	pkgroot, _, _ := tools_list.FindPkgInfo(args.WorkDir)

	// Input files may be a KCL file folder or a single KCL file.

	var files []string
	files = append(files, settings.Config.InputFile...)
	files = append(files, settings.Config.InputFiles...)

	for _, s := range files {
		if strings.Contains(s, "${PWD}") {
			s = strings.ReplaceAll(s, "${PWD}", args.WorkDir)
		}
		if strings.Contains(s, "${KCL_MOD}") {
			if pkgroot != "" {
				s = strings.ReplaceAll(s, "${KCL_MOD}", pkgroot)
			}
		}

		if strings.HasPrefix(s, ".") {
			args.KFilenameList = append(args.KFilenameList, filepath.Join(args.WorkDir, s))
		} else {
			// ${KCL_MOD}/...
			if !strings.HasPrefix(s, "${") && !filepath.IsAbs(s) {
				args.KFilenameList = append(args.KFilenameList, filepath.Join(args.WorkDir, s))
			} else {
				args.KFilenameList = append(args.KFilenameList, s)
			}
		}
	}

	// kcl -O path.to.field=field_value
	for _, override := range settings.Config.Overrides {
		args.Overrides = append(args.Overrides, override)
	}

	// kcl -D aa=11 -D bb=22 main.k
	for _, t := range settings.Options {
		var key string = t.Key
		var val string

		switch v := t.Value.(type) {
		case map[string]any:
			// Check if we should preserve order or sort keys
			if !settings.Config.SortKeys && t.originalValueNode != nil {
				// Preserve original YAML order when sort_keys is false
				if s, err := nodeToOrderedJSON(t.originalValueNode); err == nil {
					val = s
				} else {
					// Fallback to regular JSON marshaling (which sorts)
					if s, err := json.Marshal(v); err == nil {
						val = string(s)
					} else {
						val = fmt.Sprint(v)
					}
				}
			} else {
				// Use regular JSON marshaling (which sorts) when sort_keys is true
				if s, err := json.Marshal(v); err == nil {
					val = string(s)
				} else {
					val = fmt.Sprint(v)
				}
			}
		case []any:
			if s, err := json.Marshal(v); err == nil {
				val = string(s)
			} else {
				val = fmt.Sprint(v)
			}
		default:
			val = fmt.Sprint(v)
		}

		args.Args = append(args.Args, &gpyrpc.Argument{
			Name:  key,
			Value: val,
		})
	}

	// kcl -E k8s=../vendor/k8s
	for name, path := range settings.Config.PackageMaps {
		externalPkg := gpyrpc.ExternalPkg{
			PkgName: name,
			PkgPath: path,
		}
		args.ExternalPkgs = append(args.ExternalPkgs, &externalPkg)
	}

	return args
}
