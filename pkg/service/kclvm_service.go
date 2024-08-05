package service

import "kcl-lang.io/kcl-go/pkg/spec/gpyrpc"

type KclvmService interface {
	// Ping KclvmService, return the same value as the parameter
	Ping(in *gpyrpc.Ping_Args) (out *gpyrpc.Ping_Result, err error)
	// Execute KCL file with arguments and return the JSON/YAML result.
	ExecProgram(in *gpyrpc.ExecProgram_Args) (out *gpyrpc.ExecProgram_Result, err error)
	// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecutProgram method and will be removed in v0.11.0.
	BuildProgram(in *gpyrpc.BuildProgram_Args) (out *gpyrpc.BuildProgram_Result, err error)
	// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecutProgram method and will be removed in v0.11.0.
	ExecArtifact(in *gpyrpc.ExecArtifact_Args) (out *gpyrpc.ExecProgram_Result, err error)
	// Parse KCL single file to Module AST JSON string with import dependencies and parse errors.
	ParseFile(in *gpyrpc.ParseFile_Args) (out *gpyrpc.ParseFile_Result, err error)
	// Parse KCL program with entry files and return the AST JSON string.
	ParseProgram(in *gpyrpc.ParseProgram_Args) (out *gpyrpc.ParseProgram_Result, err error)
	// ListOptions provides users with the ability to parse KCL program and get all option information.
	ListOptions(in *gpyrpc.ParseProgram_Args) (out *gpyrpc.ListOptions_Result, err error)
	// ListVariables provides users with the ability to parse KCL program and get all variables by specs.
	ListVariables(in *gpyrpc.ListVariables_Args) (out *gpyrpc.ListVariables_Result, err error)
	// LoadPackage provides users with the ability to parse KCL program and semantic model information including symbols, types, definitions, etc.
	LoadPackage(in *gpyrpc.LoadPackage_Args) (out *gpyrpc.LoadPackage_Result, err error)
	// Format the code source.
	FormatCode(in *gpyrpc.FormatCode_Args) (out *gpyrpc.FormatCode_Result, err error)
	// Format KCL file or directory path contains KCL files and returns the changed file paths.
	FormatPath(in *gpyrpc.FormatPath_Args) (out *gpyrpc.FormatPath_Result, err error)
	// Lint files and return error messages including errors and warnings.
	LintPath(in *gpyrpc.LintPath_Args) (out *gpyrpc.LintPath_Result, err error)
	// Override KCL file with arguments. See [https://www.kcl-lang.io/docs/user_docs/guides/automation](https://www.kcl-lang.io/docs/user_docs/guides/automation) for more override spec guide.
	OverrideFile(in *gpyrpc.OverrideFile_Args) (out *gpyrpc.OverrideFile_Result, err error)
	// Get schema type mapping defined in the program.
	GetSchemaTypeMapping(in *gpyrpc.GetSchemaTypeMapping_Args) (out *gpyrpc.GetSchemaTypeMapping_Result, err error)
	// Validate code using schema and JSON/YAML data strings.
	ValidateCode(in *gpyrpc.ValidateCode_Args) (out *gpyrpc.ValidateCode_Result, err error)
	// List dependencies files of input paths.
	ListDepFiles(in *gpyrpc.ListDepFiles_Args) (out *gpyrpc.ListDepFiles_Result, err error)
	// Load the setting file config defined in `kcl.yaml`.
	LoadSettingsFiles(in *gpyrpc.LoadSettingsFiles_Args) (out *gpyrpc.LoadSettingsFiles_Result, err error)
	// Rename all the occurrences of the target symbol in the files. This API will rewrite files if they contain symbols to be renamed. Return the file paths that got changed.
	Rename(in *gpyrpc.Rename_Args) (out *gpyrpc.Rename_Result, err error)
	// Rename all the occurrences of the target symbol and return the modified code if any code has been changed. This API won't rewrite files but return the changed code.
	RenameCode(in *gpyrpc.RenameCode_Args) (out *gpyrpc.RenameCode_Result, err error)
	// Test KCL packages with test arguments.
	Test(in *gpyrpc.Test_Args) (out *gpyrpc.Test_Result, err error)
	// Download and update dependencies defined in the `kcl.mod` file and return the external package name and location list.
	UpdateDependencies(in *gpyrpc.UpdateDependencies_Args) (out *gpyrpc.UpdateDependencies_Result, err error)
	// GetVersion KclvmService, return the kclvm service version information
	GetVersion(in *gpyrpc.GetVersion_Args) (out *gpyrpc.GetVersion_Result, err error)
}
