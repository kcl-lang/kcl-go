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
	Key   string      `yaml:"key"`
	Value interface{} `yaml:"value"`
}

func LoadFile(filename string, src interface{}) (f *SettingsFile, err error) {
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

	var settings SettingsFile
	if err := yaml.Unmarshal([]byte(code), &settings); err != nil {
		return nil, err
	}

	settings.Filename = filename
	return &settings, nil
}

func (settings *SettingsFile) To_ExecProgram_Args() *gpyrpc.ExecProgram_Args {
	args := &gpyrpc.ExecProgram_Args{
		WorkDir: filepath.Dir(settings.Filename),

		KFilenameList: []string{},
		KCodeList:     []string{},

		Args:      []*gpyrpc.CmdArgSpec{},
		Overrides: []*gpyrpc.CmdOverrideSpec{},

		DisableYamlResult: false,
		PrintOverrideAst:  false,

		StrictRangeCheck:      settings.Config.StrictRangeCheck,
		DisableNone:           settings.Config.DisableNone,
		Verbose:               int32(settings.Config.Verbose),
		Debug:                 0,
		SortKeys:              settings.Config.SortKeys,
		ShowHidden:            settings.Config.ShowHidden,
		IncludeSchemaTypePath: settings.Config.IncludeSchemaTypePath,
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

	// kcl -O pkgpath:path.to.field=field_value
	for _, kv := range settings.Config.Overrides {
		idx0 := strings.Index(kv, ":")
		idx1 := strings.Index(kv, "=")
		if idx0 >= 0 && idx1 >= 0 && idx0 < idx1 {
			var pkgpath = kv[:idx0]
			var field_path = kv[idx0+1 : idx1]
			var field_value = kv[idx1+1:]

			args.Overrides = append(args.Overrides, &gpyrpc.CmdOverrideSpec{
				Pkgpath:    pkgpath,
				FieldPath:  field_path,
				FieldValue: field_value,
			})
		}
	}

	// kcl -D aa=11 -D bb=22 main.k
	for _, t := range settings.Options {
		var key string = t.Key
		var val string

		switch v := t.Value.(type) {
		case map[string]interface{}:
			if s, err := json.Marshal(v); err == nil {
				val = string(s)
			} else {
				val = fmt.Sprint(v)
			}
		case []interface{}:
			if s, err := json.Marshal(v); err == nil {
				val = string(s)
			} else {
				val = fmt.Sprint(v)
			}
		default:
			val = fmt.Sprint(v)
		}

		args.Args = append(args.Args, &gpyrpc.CmdArgSpec{
			Name:  key,
			Value: val,
		})
	}

	// kcl -E k8s=../vendor/k8s
	for name, path := range settings.Config.PackageMaps {
		externalPkg := gpyrpc.CmdExternalPkgSpec{
			PkgName: name,
			PkgPath: path,
		}
		args.ExternalPkgs = append(args.ExternalPkgs, &externalPkg)
	}

	return args
}
