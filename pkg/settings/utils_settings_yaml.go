// Copyright 2021 The KCL Authors. All rights reserved.

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
	InputFiles []string `yaml:"file"`
	Output     string   `yaml:"output"`

	Overrides    []string `yaml:"overrides"`
	PathSelector []string `yaml:"path_selector"`

	StrictRangeCheck bool `yaml:"strict_range_check"`
	DisableNone      bool `yaml:"disable_none"`
	Verbose          int  `yaml:"verbose"`
	Debug            bool `yaml:"debug"`
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
		code = string(src)
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
		if err2 := decodeTestFormatSettingsFile(code, &settings); err2 != nil {
			return nil, err
		}
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

		StrictRangeCheck: settings.Config.StrictRangeCheck,
		DisableNone:      settings.Config.DisableNone,
		Verbose:          int32(settings.Config.Verbose),
		Debug:            0,
	}
	if settings.Config.Debug {
		args.Debug = 1
	}

	pkgroot, _, _ := tools_list.FindPkgInfo(args.WorkDir)

	// Input files may be a KCL file folder or a single KCL file.
	for _, s := range settings.Config.InputFiles {
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

	return args
}

// kcl_options: -Y config1.yaml config2.yaml
// kcl_options: -D key0=value0 -D key1=value1 -Y temp.yaml
// kcl_options: -O :alice.labels.skin=white
// kcl_options: -S :JohnDoe.* -S :list_data.* -S :dict_data.*
// kcl_options: -d
// kcl_options: -r
// kcl_options: -n
// kcl_options: pkg.k
func decodeTestFormatSettingsFile(src string, settings *SettingsFile) error {
	// kcl_options: -Y config1.yaml config2.yaml
	var test_settings_yaml struct {
		CmdArgs string `yaml:"kcl_options"`
	}
	err := yaml.Unmarshal([]byte(src), &test_settings_yaml)
	if err != nil {
		return fmt.Errorf("invalid settings file: %v, %v", err, src)
	}

	cmdArgs := strings.Fields(test_settings_yaml.CmdArgs)
	*settings = SettingsFile{}

	for len(cmdArgs) > 0 {
		switch cmdArgs[0] {
		case "-Y":
			cmdArgs = cmdArgs[1:]
			for len(cmdArgs) > 0 {
				if strings.HasSuffix(cmdArgs[0], ".yaml") || strings.HasSuffix(cmdArgs[0], ".yml") {
					settings.Config.InputFiles = append(settings.Config.InputFiles, cmdArgs[0])
					cmdArgs = cmdArgs[1:]
					continue
				}
				break
			}
		case "-D":
			cmdArgs = cmdArgs[1:]
			if len(cmdArgs) == 0 {
				return fmt.Errorf("invalid kcl_options: %v %v", "-D", cmdArgs)
			}

			_D_arg := cmdArgs[0]
			cmdArgs = cmdArgs[1:]

			kv := strings.Split(_D_arg, "=")
			if len(kv) == 0 {
				return fmt.Errorf("invalid kcl_options: %v %v", "-D", cmdArgs)
			}
			if len(kv) < 2 {
				kv = append(kv, "")
			}

			settings.Options = append(settings.Options, KeyValueStruct{
				Key:   kv[0],
				Value: kv[1],
			})

		case "-O":
			// "--overrides",
			cmdArgs = cmdArgs[1:]
			if len(cmdArgs) == 0 {
				return fmt.Errorf("invalid kcl_options: %v %v", "-O", cmdArgs)
			}

			_O_arg := cmdArgs[0]
			cmdArgs = cmdArgs[1:]

			settings.Config.Overrides = append(settings.Config.Overrides, _O_arg)

		case "-S":
			// "--path-selector",
			cmdArgs = cmdArgs[1:]
			if len(cmdArgs) == 0 {
				return fmt.Errorf("invalid kcl_options: %v %v", "-S", cmdArgs)
			}

			_S_arg := cmdArgs[0]
			cmdArgs = cmdArgs[1:]

			settings.Config.PathSelector = append(settings.Config.PathSelector, _S_arg)

		case "-n":
			// --disable-none
			cmdArgs = cmdArgs[1:]
			settings.Config.DisableNone = true
		case "-r":
			// --strict-range-check
			cmdArgs = cmdArgs[1:]
			settings.Config.StrictRangeCheck = true
		case "-d":
			// --debug
			cmdArgs = cmdArgs[1:]
			settings.Config.Debug = true
		default:
			if strings.HasSuffix(cmdArgs[0], ".k") {
				settings.Config.InputFiles = append(settings.Config.InputFiles, cmdArgs[0])
				cmdArgs = cmdArgs[1:]
			} else {
				return fmt.Errorf("invalid kcl_options: %v", test_settings_yaml.CmdArgs)
			}
		}
	}

	return nil
}
