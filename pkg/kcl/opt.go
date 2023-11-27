// Copyright 2021 The KCL Authors. All rights reserved.

package kcl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kcl-lang.io/kcl-go/pkg/logger"
	"kcl-lang.io/kcl-go/pkg/settings"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
	"kcl-lang.io/kcl-go/pkg/tools/override"
)

type Option struct {
	*gpyrpc.ExecProgram_Args
	logger logger.Logger
	Err    error
}

// NewOption returns a new Option.
func NewOption() *Option {
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

func (p *Option) GetLogger() logger.Logger {
	return p.logger
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
			tmpOptList = append(tmpOptList, WithKFilenames(s))
		default:
			tmpOptList = append(tmpOptList, WithKFilenames(s))
		}
	}

	args := NewOption().Merge(opts...).Merge(tmpOptList...)
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

	return *args, nil
}

func WithLogger(l logger.Logger) Option {
	var opt = NewOption()
	opt.logger = l
	return *opt
}

func WithWorkDir(s string) Option {
	var opt = NewOption()
	opt.WorkDir = s
	return *opt
}

func WithKFilenames(filenames ...string) Option {
	var opt = NewOption()
	opt.KFilenameList = filenames
	return *opt
}

func WithCode(codes ...string) Option {
	var opt = NewOption()
	opt.KCodeList = codes
	return *opt
}

// kcl -E aaa=/xx/xxx/aaa main.k
func WithExternalPkgs(key_value_list ...string) Option {
	var args []*gpyrpc.CmdExternalPkgSpec
	for _, kv := range key_value_list {
		if idx := strings.Index(kv, "="); idx > 0 {
			name, value := kv[:idx], kv[idx+1:]
			args = append(args, &gpyrpc.CmdExternalPkgSpec{
				PkgName: name,
				PkgPath: value,
			})
		}
	}
	var opt = NewOption()
	opt.ExternalPkgs = args
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
	var opt = NewOption()
	opt.Args = args
	return *opt
}

// kcl -O pkgpath:path.to.field=field_value
// kcl -O pkgpath.path.to.field-
func WithOverrides(override_list ...string) Option {
	var overrides []*gpyrpc.CmdOverrideSpec
	for _, spec := range override_list {
		o, _ := override.ParseOverrideSpec(spec)
		overrides = append(overrides, o)
	}
	var opt = NewOption()
	opt.Overrides = overrides
	return *opt
}

// kcl -O pkgpath:path.to.field=field_value
// kcl -O pkgpath.path.to.field-
func WithOverridesError(override_list ...string) (Option, error) {
	var overrides []*gpyrpc.CmdOverrideSpec
	var opt = NewOption()
	for _, spec := range override_list {
		o, err := override.ParseOverrideSpec(spec)
		if err != nil {
			return *opt, err
		}
		overrides = append(overrides, o)
	}
	opt.Overrides = overrides
	return *opt, nil
}

// kcl -S path.to.field
func WithSelectors(selectors ...string) Option {
	var opt = NewOption()
	opt.PathSelector = selectors
	return *opt
}

func WithPrintOverridesAST(printOverrideAst bool) Option {
	var opt = NewOption()
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
	var opt = NewOption()
	opt.ExecProgram_Args = f.To_ExecProgram_Args()
	return *opt
}

// kcl -n --disable_none
func WithDisableNone(disableNone bool) Option {
	var opt = NewOption()
	opt.DisableNone = disableNone
	return *opt
}

// WithIncludeSchemaTypePath returns a Option which hold a include schema type path switch.
func WithIncludeSchemaTypePath(includeSchemaTypePath bool) Option {
	var opt = NewOption()
	opt.IncludeSchemaTypePath = includeSchemaTypePath
	return *opt
}

// kcl -k --sort_keys
func WithSortKeys(sortKeys bool) Option {
	var opt = NewOption()
	opt.SortKeys = sortKeys
	return *opt
}

// Merge will merge all options into one.
func (p *Option) Merge(opts ...Option) *Option {
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

		if len(opt.PathSelector) > 0 {
			p.PathSelector = append(p.PathSelector, opt.PathSelector...)
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
		if opt.CompileOnly {
			p.CompileOnly = opt.CompileOnly
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
		if opt.ExternalPkgs != nil {
			p.ExternalPkgs = append(p.ExternalPkgs, opt.ExternalPkgs...)
		}
	}
	return p
}
