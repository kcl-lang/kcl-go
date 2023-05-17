// Copyright 2021 The KCL Authors. All rights reserved.

package kcl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kusionstack.io/kclvm-go/pkg/settings"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

type Option struct {
	*gpyrpc.ExecProgram_Args
	Err error
}

func newOption() *Option {
	return &Option{
		ExecProgram_Args: new(gpyrpc.ExecProgram_Args),
	}
}

func (p *Option) JSONString() string {
	x, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return ""
	}
	return string(x)
}

func ParseArgs(pathList []string, opts ...Option) (Option, error) {
	var tmpOptList []Option
	for _, s := range pathList {
		switch {
		case strings.HasSuffix(s, ".k"):
			tmpOptList = append(tmpOptList, WithKFilenames(s))
		case strings.HasSuffix(s, ".yaml") || strings.HasSuffix(s, ".yml"):
			tmpOptList = append(tmpOptList, WithSettings(s))
		case isDir(s):
			tmpOptList = append(tmpOptList, WithWorkDir(s))
		default:
			tmpOptList = append(tmpOptList, WithKFilenames(s))
		}
	}

	args := newOption().merge(opts...).merge(tmpOptList...)
	if err := args.Err; err != nil {
		return Option{}, err
	}

	if args.WorkDir == "" {
		if len(args.KCodeList) == 0 {
			if len(args.KFilenameList) > 0 {
				if filepath.IsAbs(args.KFilenameList[0]) {
					args.WorkDir = filepath.Dir(args.KFilenameList[0])
				} else {
					args.WorkDir, _ = os.Getwd()
				}
			}
		}
	}

	if len(args.KFilenameList) == 0 {
		return Option{}, fmt.Errorf("kcl.Run: no kcl file")
	}

	return *args, nil
}

func WithWorkDir(s string) Option {
	var opt = newOption()
	opt.WorkDir = s
	return *opt
}

func WithKFilenames(filenames ...string) Option {
	var opt = newOption()
	opt.KFilenameList = filenames
	return *opt
}

func WithCode(codes ...string) Option {
	var opt = newOption()
	opt.KCodeList = codes
	return *opt
}

// kcl -D aa=11 -D bb=22 main.k
func WithOptions(key_value_list ...string) Option {
	var args []*gpyrpc.CmdArgSpec
	for _, kv := range key_value_list {
		if idx := strings.Index(kv, "="); idx > 0 {
			name, value := kv[:idx], kv[idx+1:]
			args = append(args, &gpyrpc.CmdArgSpec{
				Name:  name,
				Value: value,
			})
		}
	}
	var opt = newOption()
	opt.Args = args
	return *opt
}

// kcl -O pkgpath:path.to.field=field_value
func WithOverrides(override_list ...string) Option {
	var overrides []*gpyrpc.CmdOverrideSpec
	for _, kv := range override_list {
		idx0 := strings.Index(kv, ":")
		idx1 := strings.Index(kv, "=")
		if idx0 >= 0 && idx1 >= 0 && idx0 < idx1 {
			var pkgpath = kv[:idx0]
			var field_path = kv[idx0+1 : idx1]
			var field_value = kv[idx1+1:]
			overrides = append(overrides, &gpyrpc.CmdOverrideSpec{
				Pkgpath:    pkgpath,
				FieldPath:  field_path,
				FieldValue: field_value,
			})
		}
	}
	var opt = newOption()
	opt.Overrides = overrides
	return *opt
}

func WithPrintOverridesAST(printOverrideAst bool) Option {
	var opt = newOption()
	opt.PrintOverrideAst = printOverrideAst
	return *opt
}

// kcl -Y settings.yaml
func WithSettings(filename string) Option {
	if filename == "" {
		return Option{}
	}
	f, err := settings.LoadFile(filename, nil)
	if err != nil {
		return Option{Err: fmt.Errorf("kcl.WithSettings(%q): %v", filename, err)}
	}
	var opt = newOption()
	opt.ExecProgram_Args = f.To_ExecProgram_Args()
	return *opt
}

// kcl -n
func WithDisableNone(disableNone bool) Option {
	var opt = newOption()
	opt.DisableNone = disableNone
	return *opt
}

func WithSortKeys(sortKeys bool) Option {
	var opt = newOption()
	opt.SortKeys = sortKeys
	return *opt
}

func (p *Option) merge(opts ...Option) *Option {
	for _, opt := range opts {
		if opt.ExecProgram_Args == nil {
			continue
		}

		if opt.Err != nil {
			p.Err = opt.Err
		}

		if opt.WorkDir != "" {
			p.WorkDir = opt.WorkDir
		}

		if len(opt.KFilenameList) > 0 {
			p.KFilenameList = append(p.KFilenameList, opt.KFilenameList...)
		}
		if len(opt.KCodeList) > 0 {
			p.KCodeList = append(p.KCodeList, opt.KCodeList...)
		}

		if len(opt.Args) > 0 {
			p.Args = append(p.Args, opt.Args...)
		}
		if len(opt.Overrides) > 0 {
			p.Overrides = append(p.Overrides, opt.Overrides...)
		}

		if opt.DisableYamlResult {
			p.DisableYamlResult = opt.DisableYamlResult
		}

		if opt.PrintOverrideAst {
			p.PrintOverrideAst = opt.PrintOverrideAst
		}

		if opt.StrictRangeCheck {
			p.StrictRangeCheck = opt.StrictRangeCheck
		}
		if opt.DisableNone {
			p.DisableNone = opt.DisableNone
		}
		if opt.Verbose > 0 {
			p.Verbose = opt.Verbose
		}

		if opt.Debug != 0 {
			p.Debug = opt.Debug
		}

		if opt.SortKeys {
			p.SortKeys = opt.SortKeys
		}
		if opt.IncludeSchemaTypePath {
			p.IncludeSchemaTypePath = opt.IncludeSchemaTypePath
		}
	}
	return p
}
