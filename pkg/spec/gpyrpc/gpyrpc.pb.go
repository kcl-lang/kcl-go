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
type Ping_Args = api.Ping_Args

// Message for ping response.
type Ping_Result = api.Ping_Result

// Message for version request arguments. Empty message.
type GetVersion_Args = api.GetVersion_Args

// Message for version response.
type GetVersion_Result = api.GetVersion_Result

// Message for list method request arguments. Empty message.
type ListMethod_Args = api.ListMethod_Args

// Message for list method response.
type ListMethod_Result = api.ListMethod_Result

// Message for parse file request arguments.
type ParseFile_Args = api.ParseFile_Args

// Message for parse file response.
type ParseFile_Result = api.ParseFile_Result

// Message for parse program request arguments.
type ParseProgram_Args = api.ParseProgram_Args

// Message for parse program response.
type ParseProgram_Result = api.ParseProgram_Result

// Message for load package request arguments.
type LoadPackage_Args = api.LoadPackage_Args

// Message for load package response.
type LoadPackage_Result = api.LoadPackage_Result

// Message for list options response.
type ListOptions_Result = api.ListOptions_Result

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
type ExecProgram_Args = api.ExecProgram_Args

// Message for execute program response.
type ExecProgram_Result = api.ExecProgram_Result

// Message for build program request arguments.
type BuildProgram_Args = api.BuildProgram_Args

// Message for build program response.
type BuildProgram_Result = api.BuildProgram_Result

// Message for execute artifact request arguments.
type ExecArtifact_Args = api.ExecArtifact_Args

// Message for format code request arguments.
type FormatCode_Args = api.FormatCode_Args

// Message for format code response.
type FormatCode_Result = api.FormatCode_Result

// Message for format file path request arguments.
type FormatPath_Args = api.FormatPath_Args

// Message for format file path response.
type FormatPath_Result = api.FormatPath_Result

// Message for lint file path request arguments.
type LintPath_Args = api.LintPath_Args

// Message for lint file path response.
type LintPath_Result = api.LintPath_Result

// Message for override file request arguments.
type OverrideFile_Args = api.OverrideFile_Args

// Message for override file response.
type OverrideFile_Result = api.OverrideFile_Result

// Message for list variables options.
type ListVariables_Options = api.ListVariables_Options

// Message representing a list of variables.
type VariableList = api.VariableList

// Message for list variables request arguments.
type ListVariables_Args = api.ListVariables_Args

// Message for list variables response.
type ListVariables_Result = api.ListVariables_Result

// Message representing a variable.
type Variable = api.Variable

// Message representing a map entry.
type MapEntry = api.MapEntry

// Message for get schema type mapping request arguments.
type GetSchemaTypeMapping_Args = api.GetSchemaTypeMapping_Args

// Message for get schema type mapping response.
type GetSchemaTypeMapping_Result = api.GetSchemaTypeMapping_Result

// Message for validate code request arguments.
type ValidateCode_Args = api.ValidateCode_Args

// Message for validate code response.
type ValidateCode_Result = api.ValidateCode_Result

// Message representing a position in the source code.
type Position = api.Position

// Message for list dependency files request arguments.
type ListDepFiles_Args = api.ListDepFiles_Args

// Message for list dependency files response.
type ListDepFiles_Result = api.ListDepFiles_Result

// Message for load settings files request arguments.
type LoadSettingsFiles_Args = api.LoadSettingsFiles_Args

// Message for load settings files response.
type LoadSettingsFiles_Result = api.LoadSettingsFiles_Result

// Message representing KCL CLI configuration.
type CliConfig = api.CliConfig

// Message representing a key-value pair.
type KeyValuePair = api.KeyValuePair

// Message for rename request arguments.
type Rename_Args = api.Rename_Args

// Message for rename response.
type Rename_Result = api.Rename_Result

// Message for rename code request arguments.
type RenameCode_Args = api.RenameCode_Args

// Message for rename code response.
type RenameCode_Result = api.RenameCode_Result

// Message for test request arguments.
type Test_Args = api.Test_Args

// Message for test response.
type Test_Result = api.Test_Result

// Message representing information about a single test case.
type TestCaseInfo = api.TestCaseInfo

// Message for update dependencies request arguments.
type UpdateDependencies_Args = api.UpdateDependencies_Args

// Message for update dependencies response.
type UpdateDependencies_Result = api.UpdateDependencies_Result

// Message representing a KCL type.
type KclType = api.KclType

// Message representing a decorator in KCL.
type Decorator = api.Decorator

// Message representing an example in KCL.
type Example = api.Example
