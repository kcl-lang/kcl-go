// Copyright 2021 The KCL Authors. All rights reserved.

/*
KCLVM binding for Go

	┌─────────────────┐         ┌─────────────────┐           ┌─────────────────┐
	│     kcl files   │         │   KCLVM-Go-API  │           │  KCLResultList  │
	│  ┌───────────┐  │         │                 │           │                 │
	│  │    1.k    │  │         │                 │           │                 │
	│  └───────────┘  │         │                 │           │  ┌───────────┐  │         ┌───────────────┐
	│  ┌───────────┐  │         │  ┌───────────┐  │           │  │ KCLResult │──┼────────▶│x.Get("a.b.c") │
	│  │    2.k    │  │         │  │ Run(path) │  │           │  └───────────┘  │         └───────────────┘
	│  └───────────┘  │────┐    │  └───────────┘  │           │                 │
	│  ┌───────────┐  │    │    │                 │           │  ┌───────────┐  │         ┌───────────────┐
	│  │    3.k    │  │    │    │                 │           │  │ KCLResult │──┼────────▶│x.Get("k", &v) │
	│  └───────────┘  │    │    │                 │           │  └───────────┘  │         └───────────────┘
	│  ┌───────────┐  │    ├───▶│  ┌───────────┐  │──────────▶│                 │
	│  │setting.yml│  │    │    │  │RunFiles() │  │           │  ┌───────────┐  │         ┌───────────────┐
	│  └───────────┘  │    │    │  └───────────┘  │           │  │ KCLResult │──┼────────▶│x.JSONString() │
	└─────────────────┘    │    │                 │           │  └───────────┘  │         └───────────────┘
	                       │    │                 │           │                 │
	┌─────────────────┐    │    │                 │           │  ┌───────────┐  │         ┌───────────────┐
	│     Options     │    │    │  ┌───────────┐  │           │  │ KCLResult │──┼────────▶│x.YAMLString() │
	│WithOptions      │    │    │  │MustRun()  │  │           │  └───────────┘  │         └───────────────┘
	│WithOverrides    │────┘    │  └───────────┘  │           │                 │
	│WithWorkDir      │         │                 │           │                 │
	│WithDisableNone  │         │                 │           │                 │
	└─────────────────┘         └─────────────────┘           └─────────────────┘
*/
package kclvm

import (
	"kusionstack.io/kclvm-go/pkg/kcl"
	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/tools/format"
	"kusionstack.io/kclvm-go/pkg/tools/lint"
	"kusionstack.io/kclvm-go/pkg/tools/list"
	"kusionstack.io/kclvm-go/pkg/tools/override"
	"kusionstack.io/kclvm-go/pkg/tools/validate"
)

type (
	Option             = kcl.Option
	ListDepFilesOption = list.Option
	ValidateOptions    = validate.ValidateOptions
	KCLResult          = kcl.KCLResult
	KCLResultList      = kcl.KCLResultList

	KclType = kcl.KclType
)

// InitKclvmPath init kclvm path.
func InitKclvmPath(kclvmRoot string) {
	kclvm_runtime.InitKclvmPath(kclvmRoot)
}

// InitKclvmRuntime init kclvm process.
func InitKclvmRuntime(n int) {
	kclvm_runtime.InitRuntime(n)
}

// MustRun is like Run but panics if return any error.
func MustRun(path string, opts ...Option) *KCLResultList {
	return kcl.MustRun(path, opts...)
}

// Run evaluates the KCL program with path and opts, then returns the object list.
func Run(path string, opts ...Option) (*KCLResultList, error) {
	return kcl.Run(path, opts...)
}

// RunFiles evaluates the KCL program with multi file path and opts, then returns the object list.
func RunFiles(paths []string, opts ...Option) (*KCLResultList, error) {
	return kcl.RunFiles(paths, opts...)
}

// WithCode returns a Option which hold a kcl source code list.
func WithCode(codes ...string) Option { return kcl.WithCode(codes...) }

// WithKFilenames returns a Option which hold a filenames list.
func WithKFilenames(filenames ...string) Option { return kcl.WithKFilenames(filenames...) }

// WithOptions returns a Option which hold a key=value pair list for option function.
func WithOptions(key_value_list ...string) Option { return kcl.WithOptions(key_value_list...) }

// WithOverrides returns a Option which hold a override list.
func WithOverrides(override_list ...string) Option { return kcl.WithOverrides(override_list...) }

// WithSettings returns a Option which hold a settings file.
func WithSettings(filename string) Option { return kcl.WithSettings(filename) }

// WithWorkDir returns a Option which hold a work dir.
func WithWorkDir(workDir string) Option { return kcl.WithWorkDir(workDir) }

// WithDisableNone returns a Option which hold a disable none switch.
func WithDisableNone(disableNone bool) Option { return kcl.WithDisableNone(disableNone) }

// WithPrintOverridesAST returns a Option which hold a printOverridesAST switch.
func WithPrintOverridesAST(printOverridesAST bool) Option {
	return kcl.WithPrintOverridesAST(printOverridesAST)
}

// WithSortKeys returns a Option which hold a sortKeys switch.
func WithSortKeys(sortKeys bool) Option {
	return kcl.WithSortKeys(sortKeys)
}

// WithIncludeSchemaTypePath returns a Option which hold a includeSchemaTypePath switch.
func WithIncludeSchemaTypePath(includeSchemaTypePath bool) Option {
	return kcl.WithIncludeSchemaTypePath(includeSchemaTypePath)
}

// FormatCode returns the formatted code.
func FormatCode(code interface{}) ([]byte, error) {
	return format.FormatCode(code)
}

// FormatPath formats files from the given path
// path:
// if path is `.` or empty string, all KCL files in current directory will be formatted, not recursively
// if path is `path/file.k`, the specified KCL file will be formatted
// if path is `path/to/dir`, all KCL files in the specified dir will be formatted, not recursively
// if path is `path/to/dir/...`, all KCL files in the specified dir will be formatted recursively
//
// the returned changedPaths are the changed file paths (relative path)
func FormatPath(path string) (changedPaths []string, err error) {
	return format.FormatPath(path)
}

// ListDepFiles return the depend files from the given path
func ListDepFiles(workDir string, opt *ListDepFilesOption) (files []string, err error) {
	return list.ListDepFiles(workDir, opt)
}

// LintPath lint files from the given path
func LintPath(path string) (results []string, err error) {
	return lint.LintPath(path)
}

// OverrideFile rewrites a file with override spec
// file: string. The File that need to be overridden
// specs: []string. List of specs that need to be overridden.
//     Each spec string satisfies the form: <pkgpath>:<field_path>=<filed_value> or <pkgpath>:<field_path>-
//     When the pkgpath is '__main__', it can be omitted.
// importPaths. List of import statements that need to be added
func OverrideFile(file string, specs, importPaths []string) (bool, error) {
	return override.OverrideFile(file, specs, importPaths)
}

// ValidateCode validate data match code
func ValidateCode(data, code string, opt *ValidateOptions) (ok bool, err error) {
	return validate.ValidateCode(data, code, opt)
}

func EvalCode(code string) (*KCLResult, error) {
	return kcl.EvalCode(code)
}

// GetSchemaType returns schema types from a kcl file or code.
//
// file: string
//     The kcl filename
// code: string
//     The kcl code string
// schema_name: string
//    The schema name got, when the schema name is empty, all schemas are returned.
func GetSchemaType(file, code, schemaName string) ([]*KclType, error) {
	return kcl.GetSchemaType(file, code, schemaName)
}
