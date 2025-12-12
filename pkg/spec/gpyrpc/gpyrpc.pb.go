// Copyright The KCL Authors. All rights reserved.

package gpyrpc

import (
	"kcl-lang.io/lib/go/api"
)

// Message representing an external package for KCL.
// kcl main.k -E pkg_name=pkg_path
type ExternalPkg = api.ExternalPkg

// Message representing a key-value argument for KCL.
// kcl main.k -D name=value
type Argument = api.Argument

// Message representing an error.
type Error = api.Error

// Message representing a detailed error message with a position.
type Message = api.Message

// Message for ping request arguments.
type PingArgs = api.PingArgs

// Message for ping response.
type PingResult = api.PingResult

// Message for version request arguments. Empty message.
type GetVersionArgs = api.GetVersionArgs

// Message for version response.
type GetVersionResult = api.GetVersionResult

// Message for list method request arguments. Empty message.
type ListMethodArgs = api.ListMethodArgs

// Message for list method response.
type ListMethodResult = api.ListMethodResult

// Message for parse file request arguments.
type ParseFileArgs = api.ParseFileArgs

// Message for parse file response.
type ParseFileResult = api.ParseFileResult

// Message for parse program request arguments.
type ParseProgramArgs = api.ParseProgramArgs

// Message for parse program response.
type ParseProgramResult = api.ParseProgramResult

// Message for load package request arguments.
type LoadPackageArgs = api.LoadPackageArgs

// Message for load package response.
type LoadPackageResult = api.LoadPackageResult

// Message for list options response.
type ListOptionsResult = api.ListOptionsResult

// Message representing a help option.
type OptionHelp = api.OptionHelp

// Message representing a symbol in KCL.
type Symbol = api.Symbol

// Message representing a scope in KCL.
type Scope = api.Scope

// Message representing a symbol index.
type SymbolIndex = api.SymbolIndex

// Message representing a scope index.
type ScopeIndex = api.ScopeIndex

// Message for execute program request arguments.
type ExecProgramArgs = api.ExecProgramArgs

// Message for execute program response.
type ExecProgramResult = api.ExecProgramResult

// Message for build program request arguments.
type BuildProgramArgs = api.BuildProgramArgs

// Message for build program response.
type BuildProgramResult = api.BuildProgramResult

// Message for execute artifact request arguments.
type ExecArtifactArgs = api.ExecArtifactArgs

// Message for format code request arguments.
type FormatCodeArgs = api.FormatCodeArgs

// Message for format code response.
type FormatCodeResult = api.FormatCodeResult

// Message for format file path request arguments.
type FormatPathArgs = api.FormatPathArgs

// Message for format file path response.
type FormatPathResult = api.FormatPathResult

// Message for lint file path request arguments.
type LintPathArgs = api.LintPathArgs

// Message for lint file path response.
type LintPathResult = api.LintPathResult

// Message for override file request arguments.
type OverrideFileArgs = api.OverrideFileArgs

// Message for override file response.
type OverrideFileResult = api.OverrideFileResult

// Message for list variables options.
type ListVariablesOptions = api.ListVariablesOptions

// Message representing a list of variables.
type VariableList = api.VariableList

// Message for list variables request arguments.
type ListVariablesArgs = api.ListVariablesArgs

// Message for list variables response.
type ListVariablesResult = api.ListVariablesResult

// Message representing a variable.
type Variable = api.Variable

// Message representing a map entry.
type MapEntry = api.MapEntry

// Message for get schema type mapping request arguments.
type GetSchemaTypeMappingArgs = api.GetSchemaTypeMappingArgs

// Message for get schema type mapping response.
type GetSchemaTypeMappingResult = api.GetSchemaTypeMappingResult

// Message for validate code request arguments.
type ValidateCodeArgs = api.ValidateCodeArgs

// Message for validate code response.
type ValidateCodeResult = api.ValidateCodeResult

// Message representing a position in the source code.
type Position = api.Position

// Message for list dependency files request arguments.
type ListDepFilesArgs = api.ListDepFilesArgs

// Message for list dependency files response.
type ListDepFilesResult = api.ListDepFilesResult

// Message for load settings files request arguments.
type LoadSettingsFilesArgs = api.LoadSettingsFilesArgs

// Message for load settings files response.
type LoadSettingsFilesResult = api.LoadSettingsFilesResult

// Message representing KCL CLI configuration.
type CliConfig = api.CliConfig

// Message representing a key-value pair.
type KeyValuePair = api.KeyValuePair

// Message for rename request arguments.
type RenameArgs = api.RenameArgs

// Message for rename response.
type RenameResult = api.RenameResult

// Message for rename code request arguments.
type RenameCodeArgs = api.RenameCodeArgs

// Message for rename code response.
type RenameCodeResult = api.RenameCodeResult

// Message for test request arguments.
type TestArgs = api.TestArgs

// Message for test response.
type TestResult = api.TestResult

// Message representing information about a single test case.
type TestCaseInfo = api.TestCaseInfo

// Message for update dependencies request arguments.
type UpdateDependenciesArgs = api.UpdateDependenciesArgs

// Message for update dependencies response.
type UpdateDependenciesResult = api.UpdateDependenciesResult

// Message representing a KCL type.
type KclType = api.KclType

// Message representing a KCL function type.
type FunctionType = api.FunctionType

// Message representing a KCL function parameter type.
type Parameter = api.Parameter

// Message representing a KCL schema index signature.
type IndexSignature = api.IndexSignature

// Message representing a decorator in KCL.
type Decorator = api.Decorator

// Message representing an example in KCL.
type Example = api.Example
