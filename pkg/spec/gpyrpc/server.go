package gpyrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KclServiceServer is the server API for KclService service.
type KclServiceServer interface {
	// / Ping KclService, return the same value as the parameter
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "Ping",
	// /     "params": {
	// /         "value": "hello"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "value": "hello"
	// /     },
	// /     "id": 1
	// / }
	// / ```
	Ping(context.Context, *PingArgs) (*PingResult, error)
	// / GetVersion KclService, return the kcl service version information
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "GetVersion",
	// /     "params": {},
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "version": "0.9.1",
	// /         "checksum": "c020ab3eb4b9179219d6837a57f5d323",
	// /         "git_sha": "1a9a72942fffc9f62cb8f1ae4e1d5ca32aa1f399",
	// /         "version_info": "Version: 0.9.1-c020ab3eb4b9179219d6837a57f5d323\nPlatform: aarch64-apple-darwin\nGitCommit: 1a9a72942fffc9f62cb8f1ae4e1d5ca32aa1f399"
	// /     },
	// /     "id": 1
	// / }
	// / ```
	GetVersion(context.Context, *GetVersionArgs) (*GetVersionResult, error)
	// / Parse KCL program with entry files.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ParseProgram",
	// /     "params": {
	// /         "paths": ["./src/testdata/test.k"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "ast_json": "{...}",
	// /         "paths": ["./src/testdata/test.k"],
	// /         "errors": []
	// /     },
	// /     "id": 1
	// / }
	// / ```
	ParseProgram(context.Context, *ParseProgramArgs) (*ParseProgramResult, error)
	// / Parse KCL single file to Module AST JSON string with import dependencies
	// / and parse errors.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ParseFile",
	// /     "params": {
	// /         "path": "./src/testdata/parse/main.k"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "ast_json": "{...}",
	// /         "deps": ["./dep1", "./dep2"],
	// /         "errors": []
	// /     },
	// /     "id": 1
	// / }
	// / ```
	ParseFile(context.Context, *ParseFileArgs) (*ParseFileResult, error)
	// / load_package provides users with the ability to parse kcl program and semantic model
	// / information including symbols, types, definitions, etc.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "LoadPackage",
	// /     "params": {
	// /         "parse_args": {
	// /             "paths": ["./src/testdata/parse/main.k"]
	// /         },
	// /         "resolve_ast": true
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "program": "{...}",
	// /         "paths": ["./src/testdata/parse/main.k"],
	// /         "parse_errors": [],
	// /         "type_errors": [],
	// /         "symbols": { ... },
	// /         "scopes": { ... },
	// /         "node_symbol_map": { ... },
	// /         "symbol_node_map": { ... },
	// /         "fully_qualified_name_map": { ... },
	// /         "pkg_scope_map": { ... }
	// /     },
	// /     "id": 1
	// / }
	// / ```
	LoadPackage(context.Context, *LoadPackageArgs) (*LoadPackageResult, error)
	// / list_options provides users with the ability to parse kcl program and get all option information.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ListOptions",
	// /     "params": {
	// /         "paths": ["./src/testdata/option/main.k"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "options": [
	// /             { "name": "option1", "type": "str", "required": true, "default_value": "", "help": "option 1 help" },
	// /             { "name": "option2", "type": "int", "required": false, "default_value": "0", "help": "option 2 help" },
	// /             { "name": "option3", "type": "bool", "required": false, "default_value": "false", "help": "option 3 help" }
	// /         ]
	// /     },
	// /     "id": 1
	// / }
	// / ```
	ListOptions(context.Context, *ParseProgramArgs) (*ListOptionsResult, error)
	// / list_variables provides users with the ability to parse kcl program and get all variables by specs.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ListVariables",
	// /     "params": {
	// /         "files": ["./src/testdata/variables/main.k"],
	// /         "specs": ["a"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "variables": {
	// /             "a": {
	// /                 "variables": [
	// /                     { "value": "1", "type_name": "int", "op_sym": "", "list_items": [], "dict_entries": [] }
	// /                 ]
	// /             }
	// /         },
	// /         "unsupported_codes": [],
	// /         "parse_errors": []
	// /     },
	// /     "id": 1
	// / }
	// / ```
	ListVariables(context.Context, *ListVariablesArgs) (*ListVariablesResult, error)
	// / Execute KCL file with args. **Note that it is not thread safe.**
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ExecProgram",
	// /     "params": {
	// /         "work_dir": "./src/testdata",
	// /         "k_filename_list": ["test.k"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "json_result": "{\"alice\": {\"age\": 18}}",
	// /         "yaml_result": "alice:\n  age: 18",
	// /         "log_message": "",
	// /         "err_message": ""
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Request with code
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ExecProgram",
	// /     "params": {
	// /         "k_filename_list": ["file.k"],
	// /         "k_code_list": ["alice = {age = 18}"]
	// /     },
	// /     "id": 2
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "json_result": "{\"alice\": {\"age\": 18}}",
	// /         "yaml_result": "alice:\n  age: 18",
	// /         "log_message": "",
	// /         "err_message": ""
	// /     },
	// /     "id": 2
	// / }
	// /
	// / // Error case - cannot find file
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ExecProgram",
	// /     "params": {
	// /         "k_filename_list": ["invalid_file.k"]
	// /     },
	// /     "id": 3
	// / }
	// /
	// / // Error Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "error": {
	// /         "code": -32602,
	// /         "message": "Cannot find the kcl file"
	// /     },
	// /     "id": 3
	// / }
	// /
	// / // Error case - no input files
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ExecProgram",
	// /     "params": {
	// /         "k_filename_list": []
	// /     },
	// /     "id": 4
	// / }
	// /
	// / // Error Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "error": {
	// /         "code": -32602,
	// /         "message": "No input KCL files or paths"
	// /     },
	// /     "id": 4
	// / }
	// / ```
	ExecProgram(context.Context, *ExecProgramArgs) (*ExecProgramResult, error)
	// / Build the KCL program to an artifact.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "BuildProgram",
	// /     "params": {
	// /         "exec_args": {
	// /             "work_dir": "./src/testdata",
	// /             "k_filename_list": ["test.k"]
	// /         },
	// /         "output": "./build"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "path": "./build/test.k"
	// /     },
	// /     "id": 1
	// / }
	// / ```
	// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
	BuildProgram(context.Context, *BuildProgramArgs) (*BuildProgramResult, error)
	// / Execute the KCL artifact with args. **Note that it is not thread safe.**
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ExecArtifact",
	// /     "params": {
	// /         "path": "./artifact_path",
	// /         "exec_args": {
	// /             "work_dir": "./src/testdata",
	// /             "k_filename_list": ["test.k"]
	// /         }
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "json_result": "{\"alice\": {\"age\": 18}}",
	// /         "yaml_result": "alice:\n  age: 18",
	// /         "log_message": "",
	// /         "err_message": ""
	// /     },
	// /     "id": 1
	// / }
	// / ```
	// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
	ExecArtifact(context.Context, *ExecArtifactArgs) (*ExecProgramResult, error)
	// / Override KCL file with args.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "OverrideFile",
	// /     "params": {
	// /         "file": "./src/testdata/test.k",
	// /         "specs": ["alice.age=18"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "result": true,
	// /         "parse_errors": []
	// /     },
	// /     "id": 1
	// / }
	// / ```
	OverrideFile(context.Context, *OverrideFileArgs) (*OverrideFileResult, error)
	// / Get schema type mapping.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "GetSchemaTypeMapping",
	// /     "params": {
	// /         "exec_args": {
	// /             "work_dir": "./src/testdata",
	// /             "k_filename_list": ["main.k"],
	// /             "external_pkgs": [
	// /                 {
	// /                     "pkg_name":"pkg",
	// /                     "pkg_path": "./src/testdata/pkg"
	// /                 }
	// /             ]
	// /         },
	// /         "schema_name": "Person"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "schema_type_mapping": {
	// /             "Person": {
	// /                 "type": "schema",
	// /                 "schema_name": "Person",
	// /                 "properties": {
	// /                     "name": { "type": "str" },
	// /                     "age": { "type": "int" }
	// /                 },
	// /                 "required": ["name", "age"],
	// /                 "decorators": []
	// /             }
	// /         }
	// /     },
	// /     "id": 1
	// / }
	// / ```
	GetSchemaTypeMapping(context.Context, *GetSchemaTypeMappingArgs) (*GetSchemaTypeMappingResult, error)
	// / Format code source.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "FormatCode",
	// /     "params": {
	// /         "source": "schema Person {\n    name: str\n    age: int\n}\nperson = Person {\n    name = \"Alice\"\n    age = 18\n}\n"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "formatted": "schema Person {\n    name: str\n    age: int\n}\nperson = Person {\n    name = \"Alice\"\n    age = 18\n}\n"
	// /     },
	// /     "id": 1
	// / }
	// / ```
	FormatCode(context.Context, *FormatCodeArgs) (*FormatCodeResult, error)
	// / Format KCL file or directory path contains KCL files and returns the changed file paths.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "FormatPath",
	// /     "params": {
	// /         "path": "./src/testdata/test.k"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "changed_paths": []
	// /     },
	// /     "id": 1
	// / }
	// / ```
	FormatPath(context.Context, *FormatPathArgs) (*FormatPathResult, error)
	// / Lint files and return error messages including errors and warnings.
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "LintPath",
	// /     "params": {
	// /         "paths": ["./src/testdata/test-lint.k"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "results": ["Module 'math' imported but unused"]
	// /     },
	// /     "id": 1
	// / }
	// / ```
	LintPath(context.Context, *LintPathArgs) (*LintPathResult, error)
	// / Validate code using schema and data strings.
	// /
	// / **Note that it is not thread safe.**
	// /
	// / # Examples
	// /
	// / ```jsonrpc
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "ValidateCode",
	// /     "params": {
	// /         "code": "schema Person {\n    name: str\n    age: int\n    check: 0 < age < 120\n}",
	// /         "data": "{\"name\": \"Alice\", \"age\": 10}"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "success": true,
	// /         "err_message": ""
	// /     },
	// /     "id": 1
	// / }
	// / ```
	ValidateCode(context.Context, *ValidateCodeArgs) (*ValidateCodeResult, error)
	ListDepFiles(context.Context, *ListDepFilesArgs) (*ListDepFilesResult, error)
	// / Build setting file config from args.
	// /
	// / # Examples
	// /
	// /
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "LoadSettingsFiles",
	// /     "params": {
	// /         "work_dir": "./src/testdata/settings",
	// /         "files": ["./src/testdata/settings/kcl.yaml"]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "kcl_cli_configs": {
	// /             "files": ["./src/testdata/settings/kcl.yaml"],
	// /             "output": "",
	// /             "overrides": [],
	// /             "path_selector": [],
	// /             "strict_range_check": false,
	// /             "disable_none": false,
	// /             "verbose": 0,
	// /             "debug": false,
	// /             "sort_keys": false,
	// /             "show_hidden": false,
	// /             "include_schema_type_path": false,
	// /             "fast_eval": false
	// /         },
	// /         "kcl_options": []
	// /     },
	// /     "id": 1
	// / }
	// / ```
	LoadSettingsFiles(context.Context, *LoadSettingsFilesArgs) (*LoadSettingsFilesResult, error)
	// / Rename all the occurrences of the target symbol in the files. This API will rewrite files if they contain symbols to be renamed.
	// / Return the file paths that got changed.
	// /
	// / # Examples
	// /
	// /
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "Rename",
	// /     "params": {
	// /         "package_root": "./src/testdata/rename_doc",
	// /         "symbol_path": "a",
	// /         "file_paths": ["./src/testdata/rename_doc/main.k"],
	// /         "new_name": "a2"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "changed_files": ["./src/testdata/rename_doc/main.k"]
	// /     },
	// /     "id": 1
	// / }
	// / ```
	Rename(context.Context, *RenameArgs) (*RenameResult, error)
	// / Rename all the occurrences of the target symbol and return the modified code if any code has been changed. This API won't rewrite files but return the changed code.
	// /
	// / # Examples
	// /
	// /
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "RenameCode",
	// /     "params": {
	// /         "package_root": "/mock/path",
	// /         "symbol_path": "a",
	// /         "source_codes": {
	// /             "/mock/path/main.k": "a = 1\nb = a"
	// /         },
	// /         "new_name": "a2"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "changed_codes": {
	// /             "/mock/path/main.k": "a2 = 1\nb = a2"
	// /         }
	// /     },
	// /     "id": 1
	// / }
	// / ```
	RenameCode(context.Context, *RenameCodeArgs) (*RenameCodeResult, error)
	// / Test KCL packages with test arguments.
	// /
	// / # Examples
	// /
	// /
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "Test",
	// /     "params": {
	// /         "exec_args": {
	// /             "work_dir": "./src/testdata/testing/module",
	// /             "k_filename_list": ["main.k"]
	// /         },
	// /         "pkg_list": ["./src/testdata/testing/module/..."]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "info": [
	// /             {"name": "test_case_1", "error": "", "duration": 1000, "log_message": ""},
	// /             {"name": "test_case_2", "error": "some error", "duration": 2000, "log_message": ""}
	// /         ]
	// /     },
	// /     "id": 1
	// / }
	// / ```
	Test(context.Context, *TestArgs) (*TestResult, error)
	// / Download and update dependencies defined in the kcl.mod file.
	// /
	// / # Examples
	// /
	// /
	// / // Request
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "UpdateDependencies",
	// /     "params": {
	// /         "manifest_path": "./src/testdata/update_dependencies"
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "external_pkgs": [
	// /             {"pkg_name": "pkg1", "pkg_path": "./src/testdata/update_dependencies/pkg1"}
	// /         ]
	// /     },
	// /     "id": 1
	// / }
	// /
	// / // Request with vendor flag
	// / {
	// /     "jsonrpc": "2.0",
	// /     "method": "UpdateDependencies",
	// /     "params": {
	// /         "manifest_path": "./src/testdata/update_dependencies",
	// /         "vendor": true
	// /     },
	// /     "id": 2
	// / }
	// /
	// / // Response
	// / {
	// /     "jsonrpc": "2.0",
	// /     "result": {
	// /         "external_pkgs": [
	// /             {"pkg_name": "pkg1", "pkg_path": "./src/testdata/update_dependencies/pkg1"}
	// /         ]
	// /     },
	// /     "id": 2
	// / }
	// / ```
	UpdateDependencies(context.Context, *UpdateDependenciesArgs) (*UpdateDependenciesResult, error)
}

// UnimplementedKclServiceServer can be embedded to have forward compatible implementations.
type UnimplementedKclServiceServer struct {
}

func (*UnimplementedKclServiceServer) Ping(context.Context, *PingArgs) (*PingResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedKclServiceServer) GetVersion(context.Context, *GetVersionArgs) (*GetVersionResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (*UnimplementedKclServiceServer) ParseProgram(context.Context, *ParseProgramArgs) (*ParseProgramResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseProgram not implemented")
}
func (*UnimplementedKclServiceServer) ParseFile(context.Context, *ParseFileArgs) (*ParseFileResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseFile not implemented")
}
func (*UnimplementedKclServiceServer) LoadPackage(context.Context, *LoadPackageArgs) (*LoadPackageResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPackage not implemented")
}
func (*UnimplementedKclServiceServer) ListOptions(context.Context, *ParseProgramArgs) (*ListOptionsResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOptions not implemented")
}
func (*UnimplementedKclServiceServer) ListVariables(context.Context, *ListVariablesArgs) (*ListVariablesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVariables not implemented")
}
func (*UnimplementedKclServiceServer) ExecProgram(context.Context, *ExecProgramArgs) (*ExecProgramResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecProgram not implemented")
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (*UnimplementedKclServiceServer) BuildProgram(context.Context, *BuildProgramArgs) (*BuildProgramResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildProgram not implemented")
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (*UnimplementedKclServiceServer) ExecArtifact(context.Context, *ExecArtifactArgs) (*ExecProgramResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecArtifact not implemented")
}
func (*UnimplementedKclServiceServer) OverrideFile(context.Context, *OverrideFileArgs) (*OverrideFileResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OverrideFile not implemented")
}
func (*UnimplementedKclServiceServer) GetSchemaTypeMapping(context.Context, *GetSchemaTypeMappingArgs) (*GetSchemaTypeMappingResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchemaTypeMapping not implemented")
}
func (*UnimplementedKclServiceServer) FormatCode(context.Context, *FormatCodeArgs) (*FormatCodeResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FormatCode not implemented")
}
func (*UnimplementedKclServiceServer) FormatPath(context.Context, *FormatPathArgs) (*FormatPathResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FormatPath not implemented")
}
func (*UnimplementedKclServiceServer) LintPath(context.Context, *LintPathArgs) (*LintPathResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LintPath not implemented")
}
func (*UnimplementedKclServiceServer) ValidateCode(context.Context, *ValidateCodeArgs) (*ValidateCodeResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateCode not implemented")
}
func (*UnimplementedKclServiceServer) ListDepFiles(context.Context, *ListDepFilesArgs) (*ListDepFilesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDepFiles not implemented")
}
func (*UnimplementedKclServiceServer) LoadSettingsFiles(context.Context, *LoadSettingsFilesArgs) (*LoadSettingsFilesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadSettingsFiles not implemented")
}
func (*UnimplementedKclServiceServer) Rename(context.Context, *RenameArgs) (*RenameResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rename not implemented")
}
func (*UnimplementedKclServiceServer) RenameCode(context.Context, *RenameCodeArgs) (*RenameCodeResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenameCode not implemented")
}
func (*UnimplementedKclServiceServer) Test(context.Context, *TestArgs) (*TestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (*UnimplementedKclServiceServer) UpdateDependencies(context.Context, *UpdateDependenciesArgs) (*UpdateDependenciesResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDependencies not implemented")
}

func RegisterKclServiceServer(s *grpc.Server, srv KclServiceServer) {
	s.RegisterService(&_KclService_serviceDesc, srv)
}

func _KclService_Ping_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(PingArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/Ping",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).Ping(ctx, req.(*PingArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_GetVersion_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(GetVersionArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/GetVersion",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).GetVersion(ctx, req.(*GetVersionArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ParseProgram_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ParseProgramArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ParseProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ParseProgram",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ParseProgram(ctx, req.(*ParseProgramArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ParseFile_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ParseFileArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ParseFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ParseFile",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ParseFile(ctx, req.(*ParseFileArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_LoadPackage_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(LoadPackageArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).LoadPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/LoadPackage",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).LoadPackage(ctx, req.(*LoadPackageArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ListOptions_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ParseProgramArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ListOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ListOptions",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ListOptions(ctx, req.(*ParseProgramArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ListVariables_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ListVariablesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ListVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ListVariables",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ListVariables(ctx, req.(*ListVariablesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ExecProgram_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ExecProgramArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ExecProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ExecProgram",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ExecProgram(ctx, req.(*ExecProgramArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func _KclService_BuildProgram_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(BuildProgramArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).BuildProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/BuildProgram",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).BuildProgram(ctx, req.(*BuildProgramArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func _KclService_ExecArtifact_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ExecArtifactArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ExecArtifact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ExecArtifact",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ExecArtifact(ctx, req.(*ExecArtifactArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_OverrideFile_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(OverrideFileArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).OverrideFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/OverrideFile",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).OverrideFile(ctx, req.(*OverrideFileArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_GetSchemaTypeMapping_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(GetSchemaTypeMappingArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).GetSchemaTypeMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/GetSchemaTypeMapping",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).GetSchemaTypeMapping(ctx, req.(*GetSchemaTypeMappingArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_FormatCode_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(FormatCodeArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).FormatCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/FormatCode",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).FormatCode(ctx, req.(*FormatCodeArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_FormatPath_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(FormatPathArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).FormatPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/FormatPath",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).FormatPath(ctx, req.(*FormatPathArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_LintPath_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(LintPathArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).LintPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/LintPath",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).LintPath(ctx, req.(*LintPathArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ValidateCode_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ValidateCodeArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ValidateCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ValidateCode",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ValidateCode(ctx, req.(*ValidateCodeArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_ListDepFiles_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ListDepFilesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).ListDepFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/ListDepFiles",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).ListDepFiles(ctx, req.(*ListDepFilesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_LoadSettingsFiles_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(LoadSettingsFilesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).LoadSettingsFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/LoadSettingsFiles",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).LoadSettingsFiles(ctx, req.(*LoadSettingsFilesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_Rename_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(RenameArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).Rename(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/Rename",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).Rename(ctx, req.(*RenameArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_RenameCode_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(RenameCodeArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).RenameCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/RenameCode",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).RenameCode(ctx, req.(*RenameCodeArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_Test_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(TestArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/Test",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).Test(ctx, req.(*TestArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclService_UpdateDependencies_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(UpdateDependenciesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclServiceServer).UpdateDependencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclService/UpdateDependencies",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KclServiceServer).UpdateDependencies(ctx, req.(*UpdateDependenciesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

var _KclService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gpyrpc.KclService",
	HandlerType: (*KclServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KclService_Ping_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _KclService_GetVersion_Handler,
		},
		{
			MethodName: "ParseProgram",
			Handler:    _KclService_ParseProgram_Handler,
		},
		{
			MethodName: "ParseFile",
			Handler:    _KclService_ParseFile_Handler,
		},
		{
			MethodName: "LoadPackage",
			Handler:    _KclService_LoadPackage_Handler,
		},
		{
			MethodName: "ListOptions",
			Handler:    _KclService_ListOptions_Handler,
		},
		{
			MethodName: "ListVariables",
			Handler:    _KclService_ListVariables_Handler,
		},
		{
			MethodName: "ExecProgram",
			Handler:    _KclService_ExecProgram_Handler,
		},
		{
			MethodName: "BuildProgram",
			Handler:    _KclService_BuildProgram_Handler,
		},
		{
			MethodName: "ExecArtifact",
			Handler:    _KclService_ExecArtifact_Handler,
		},
		{
			MethodName: "OverrideFile",
			Handler:    _KclService_OverrideFile_Handler,
		},
		{
			MethodName: "GetSchemaTypeMapping",
			Handler:    _KclService_GetSchemaTypeMapping_Handler,
		},
		{
			MethodName: "FormatCode",
			Handler:    _KclService_FormatCode_Handler,
		},
		{
			MethodName: "FormatPath",
			Handler:    _KclService_FormatPath_Handler,
		},
		{
			MethodName: "LintPath",
			Handler:    _KclService_LintPath_Handler,
		},
		{
			MethodName: "ValidateCode",
			Handler:    _KclService_ValidateCode_Handler,
		},
		{
			MethodName: "ListDepFiles",
			Handler:    _KclService_ListDepFiles_Handler,
		},
		{
			MethodName: "LoadSettingsFiles",
			Handler:    _KclService_LoadSettingsFiles_Handler,
		},
		{
			MethodName: "Rename",
			Handler:    _KclService_Rename_Handler,
		},
		{
			MethodName: "RenameCode",
			Handler:    _KclService_RenameCode_Handler,
		},
		{
			MethodName: "Test",
			Handler:    _KclService_Test_Handler,
		},
		{
			MethodName: "UpdateDependencies",
			Handler:    _KclService_UpdateDependencies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gpyrpc.proto",
}
