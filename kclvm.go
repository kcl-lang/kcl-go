// Copyright 2023 The KCL Authors. All rights reserved.

// Package kclvm
/*
KCL Go SDK

	┌─────────────────┐         ┌─────────────────┐           ┌─────────────────┐
	│     kcl files   │         │    KCL-Go-API   │           │  KCLResultList  │
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
	"io"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
	"kcl-lang.io/kcl-go/pkg/tools/format"
	"kcl-lang.io/kcl-go/pkg/tools/lint"
	"kcl-lang.io/kcl-go/pkg/tools/list"
	"kcl-lang.io/kcl-go/pkg/tools/override"
	"kcl-lang.io/kcl-go/pkg/tools/testing"
	"kcl-lang.io/kcl-go/pkg/tools/validate"
)

type (
	Option             = kcl.Option
	ListDepsOptions    = list.DepOptions
	ListDepFilesOption = list.Option
	ValidateOptions    = validate.ValidateOptions
	TestOptions        = testing.TestOptions
	TestCaseInfo       = testing.TestCaseInfo
	TestResult         = testing.TestResult
	KCLResult          = kcl.KCLResult
	KCLResultList      = kcl.KCLResultList

	KclType = kcl.KclType
)

// InitKclvmPath init kclvm path.
func InitKclvmPath(kclvmRoot string) {
	kclvm_runtime.InitKclvmRoot(kclvmRoot)
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

// WithExternalPkgs returns a Option which hold a external package list.
func WithExternalPkgs(externalPkgs ...string) Option { return kcl.WithExternalPkgs(externalPkgs...) }

// WithKFilenames returns a Option which hold a filenames list.
func WithKFilenames(filenames ...string) Option { return kcl.WithKFilenames(filenames...) }

// WithOptions returns a Option which hold a key=value pair list for option function.
func WithOptions(key_value_list ...string) Option { return kcl.WithOptions(key_value_list...) }

// WithOverrides returns a Option which hold a override list.
func WithOverrides(override_list ...string) Option { return kcl.WithOverrides(override_list...) }

// WithSelectors returns a Option which hold a path selector list.
func WithSelectors(selectors ...string) Option { return kcl.WithSelectors(selectors...) }

// WithSettings returns a Option which hold a settings file.
func WithSettings(filename string) Option { return kcl.WithSettings(filename) }

// WithWorkDir returns a Option which hold a work dir.
func WithWorkDir(workDir string) Option { return kcl.WithWorkDir(workDir) }

// WithDisableNone returns a Option which hold a disable none switch.
func WithDisableNone(disableNone bool) Option { return kcl.WithDisableNone(disableNone) }

// WithIncludeSchemaTypePath returns a Option which hold a include schema type path switch.
func WithIncludeSchemaTypePath(includeSchemaTypePath bool) Option {
	return kcl.WithIncludeSchemaTypePath(includeSchemaTypePath)
}

// WithPrintOverridesAST returns a Option which hold a printOverridesAST switch.
func WithPrintOverridesAST(printOverridesAST bool) Option {
	return kcl.WithPrintOverridesAST(printOverridesAST)
}

// WithSortKeys returns a Option which hold a sortKeys switch.
func WithSortKeys(sortKeys bool) Option {
	return kcl.WithSortKeys(sortKeys)
}

// WithLogger returns a Option which hold a logger.
func WithLogger(l io.Writer) Option {
	return kcl.WithLogger(l)
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

// ListUpStreamFiles return a list of upstream depend files from the given path list
func ListUpStreamFiles(workDir string, opt *ListDepsOptions) (deps []string, err error) {
	return list.ListUpStreamFiles(workDir, opt)
}

// ListDownStreamFiles return a list of downstream depend files from the given changed path list.
func ListDownStreamFiles(workDir string, opt *ListDepsOptions) ([]string, error) {
	return list.ListDownStreamFiles(workDir, opt)
}

// LintPath lint files from the given path
func LintPath(paths []string) (results []string, err error) {
	return lint.LintPath(paths)
}

// OverrideFile rewrites a file with override spec
// file: string. The File that need to be overridden
// specs: []string. List of specs that need to be overridden.
//
//	Each spec string satisfies the form: <pkgpath>:<field_path>=<filed_value> or <pkgpath>:<field_path>-
//	When the pkgpath is '__main__', it can be omitted.
//
// importPaths. List of import statements that need to be added
func OverrideFile(file string, specs, importPaths []string) (bool, error) {
	return override.OverrideFile(file, specs, importPaths)
}

// ValidateCode validate data string match code string
func ValidateCode(data, code string, opts *ValidateOptions) (ok bool, err error) {
	return validate.ValidateCode(data, code, opts)
}

// Validate validates the given data file against the specified
// schema file with the provided options.
func Validate(dataFile, schemaFile string, opts *ValidateOptions) (ok bool, err error) {
	return validate.Validate(dataFile, schemaFile, opts)
}

// Test calls the test tool to run uni tests in packages.
func Test(testOpts *TestOptions, opts ...Option) (TestResult, error) {
	return testing.Test(testOpts, opts...)
}

// GetSchemaType returns schema types from a kcl file or code.
//
// file: string
//
//	The kcl filename
//
// code: string
//
//	The kcl code string
//
// schema_name: string
//
//	The schema name got, when the schema name is empty, all schemas are returned.
func GetSchemaType(file, code, schemaName string) ([]*KclType, error) {
	return kcl.GetSchemaType(file, code, schemaName)
}

// GetSchemaTypeMapping returns a <schemaName>:<schemaType> mapping of schema types from a kcl file or code.
//
// file: string
//
//	The kcl filename
//
// code: string
//
//	The kcl code string
//
// schema_name: string
//
//	The schema name got, when the schema name is empty, all schemas are returned.
func GetSchemaTypeMapping(file, code, schemaName string) (map[string]*KclType, error) {
	return kcl.GetSchemaTypeMapping(file, code, schemaName)
}
