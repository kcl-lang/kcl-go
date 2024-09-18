package gpyrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KclvmServiceServer is the server API for KclvmService service.
type KclvmServiceServer interface {
	// / Ping KclvmService, return the same value as the parameter
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
	Ping(context.Context, *Ping_Args) (*Ping_Result, error)
	// / GetVersion KclvmService, return the kclvm service version information
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
	GetVersion(context.Context, *GetVersion_Args) (*GetVersion_Result, error)
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
	ParseProgram(context.Context, *ParseProgram_Args) (*ParseProgram_Result, error)
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
	ParseFile(context.Context, *ParseFile_Args) (*ParseFile_Result, error)
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
	LoadPackage(context.Context, *LoadPackage_Args) (*LoadPackage_Result, error)
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
	ListOptions(context.Context, *ParseProgram_Args) (*ListOptions_Result, error)
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
	ListVariables(context.Context, *ListVariables_Args) (*ListVariables_Result, error)
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
	ExecProgram(context.Context, *ExecProgram_Args) (*ExecProgram_Result, error)
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
	BuildProgram(context.Context, *BuildProgram_Args) (*BuildProgram_Result, error)
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
	ExecArtifact(context.Context, *ExecArtifact_Args) (*ExecProgram_Result, error)
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
	OverrideFile(context.Context, *OverrideFile_Args) (*OverrideFile_Result, error)
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
	GetSchemaTypeMapping(context.Context, *GetSchemaTypeMapping_Args) (*GetSchemaTypeMapping_Result, error)
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
	FormatCode(context.Context, *FormatCode_Args) (*FormatCode_Result, error)
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
	FormatPath(context.Context, *FormatPath_Args) (*FormatPath_Result, error)
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
	LintPath(context.Context, *LintPath_Args) (*LintPath_Result, error)
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
	ValidateCode(context.Context, *ValidateCode_Args) (*ValidateCode_Result, error)
	ListDepFiles(context.Context, *ListDepFiles_Args) (*ListDepFiles_Result, error)
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
	LoadSettingsFiles(context.Context, *LoadSettingsFiles_Args) (*LoadSettingsFiles_Result, error)
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
	Rename(context.Context, *Rename_Args) (*Rename_Result, error)
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
	RenameCode(context.Context, *RenameCode_Args) (*RenameCode_Result, error)
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
	Test(context.Context, *Test_Args) (*Test_Result, error)
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
	UpdateDependencies(context.Context, *UpdateDependencies_Args) (*UpdateDependencies_Result, error)
}

// UnimplementedKclvmServiceServer can be embedded to have forward compatible implementations.
type UnimplementedKclvmServiceServer struct {
}

func (*UnimplementedKclvmServiceServer) Ping(context.Context, *Ping_Args) (*Ping_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedKclvmServiceServer) GetVersion(context.Context, *GetVersion_Args) (*GetVersion_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (*UnimplementedKclvmServiceServer) ParseProgram(context.Context, *ParseProgram_Args) (*ParseProgram_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseProgram not implemented")
}
func (*UnimplementedKclvmServiceServer) ParseFile(context.Context, *ParseFile_Args) (*ParseFile_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ParseFile not implemented")
}
func (*UnimplementedKclvmServiceServer) LoadPackage(context.Context, *LoadPackage_Args) (*LoadPackage_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPackage not implemented")
}
func (*UnimplementedKclvmServiceServer) ListOptions(context.Context, *ParseProgram_Args) (*ListOptions_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOptions not implemented")
}
func (*UnimplementedKclvmServiceServer) ListVariables(context.Context, *ListVariables_Args) (*ListVariables_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVariables not implemented")
}
func (*UnimplementedKclvmServiceServer) ExecProgram(context.Context, *ExecProgram_Args) (*ExecProgram_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecProgram not implemented")
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (*UnimplementedKclvmServiceServer) BuildProgram(context.Context, *BuildProgram_Args) (*BuildProgram_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuildProgram not implemented")
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (*UnimplementedKclvmServiceServer) ExecArtifact(context.Context, *ExecArtifact_Args) (*ExecProgram_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecArtifact not implemented")
}
func (*UnimplementedKclvmServiceServer) OverrideFile(context.Context, *OverrideFile_Args) (*OverrideFile_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OverrideFile not implemented")
}
func (*UnimplementedKclvmServiceServer) GetSchemaTypeMapping(context.Context, *GetSchemaTypeMapping_Args) (*GetSchemaTypeMapping_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchemaTypeMapping not implemented")
}
func (*UnimplementedKclvmServiceServer) FormatCode(context.Context, *FormatCode_Args) (*FormatCode_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FormatCode not implemented")
}
func (*UnimplementedKclvmServiceServer) FormatPath(context.Context, *FormatPath_Args) (*FormatPath_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FormatPath not implemented")
}
func (*UnimplementedKclvmServiceServer) LintPath(context.Context, *LintPath_Args) (*LintPath_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LintPath not implemented")
}
func (*UnimplementedKclvmServiceServer) ValidateCode(context.Context, *ValidateCode_Args) (*ValidateCode_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateCode not implemented")
}
func (*UnimplementedKclvmServiceServer) ListDepFiles(context.Context, *ListDepFiles_Args) (*ListDepFiles_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDepFiles not implemented")
}
func (*UnimplementedKclvmServiceServer) LoadSettingsFiles(context.Context, *LoadSettingsFiles_Args) (*LoadSettingsFiles_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadSettingsFiles not implemented")
}
func (*UnimplementedKclvmServiceServer) Rename(context.Context, *Rename_Args) (*Rename_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rename not implemented")
}
func (*UnimplementedKclvmServiceServer) RenameCode(context.Context, *RenameCode_Args) (*RenameCode_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenameCode not implemented")
}
func (*UnimplementedKclvmServiceServer) Test(context.Context, *Test_Args) (*Test_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (*UnimplementedKclvmServiceServer) UpdateDependencies(context.Context, *UpdateDependencies_Args) (*UpdateDependencies_Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDependencies not implemented")
}

func RegisterKclvmServiceServer(s *grpc.Server, srv KclvmServiceServer) {
	s.RegisterService(&_KclvmService_serviceDesc, srv)
}

func _KclvmService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).Ping(ctx, req.(*Ping_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVersion_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).GetVersion(ctx, req.(*GetVersion_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ParseProgram_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseProgram_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ParseProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ParseProgram",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ParseProgram(ctx, req.(*ParseProgram_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ParseFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseFile_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ParseFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ParseFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ParseFile(ctx, req.(*ParseFile_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_LoadPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadPackage_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).LoadPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/LoadPackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).LoadPackage(ctx, req.(*LoadPackage_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ListOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParseProgram_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ListOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ListOptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ListOptions(ctx, req.(*ParseProgram_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ListVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVariables_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ListVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ListVariables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ListVariables(ctx, req.(*ListVariables_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ExecProgram_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecProgram_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ExecProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ExecProgram",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ExecProgram(ctx, req.(*ExecProgram_Args))
	}
	return interceptor(ctx, in, info, handler)
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func _KclvmService_BuildProgram_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuildProgram_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).BuildProgram(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/BuildProgram",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).BuildProgram(ctx, req.(*BuildProgram_Args))
	}
	return interceptor(ctx, in, info, handler)
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func _KclvmService_ExecArtifact_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecArtifact_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ExecArtifact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ExecArtifact",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ExecArtifact(ctx, req.(*ExecArtifact_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_OverrideFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OverrideFile_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).OverrideFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/OverrideFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).OverrideFile(ctx, req.(*OverrideFile_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_GetSchemaTypeMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaTypeMapping_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).GetSchemaTypeMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/GetSchemaTypeMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).GetSchemaTypeMapping(ctx, req.(*GetSchemaTypeMapping_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_FormatCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FormatCode_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).FormatCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/FormatCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).FormatCode(ctx, req.(*FormatCode_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_FormatPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FormatPath_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).FormatPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/FormatPath",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).FormatPath(ctx, req.(*FormatPath_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_LintPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LintPath_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).LintPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/LintPath",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).LintPath(ctx, req.(*LintPath_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ValidateCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateCode_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ValidateCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ValidateCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ValidateCode(ctx, req.(*ValidateCode_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_ListDepFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDepFiles_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).ListDepFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/ListDepFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).ListDepFiles(ctx, req.(*ListDepFiles_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_LoadSettingsFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadSettingsFiles_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).LoadSettingsFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/LoadSettingsFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).LoadSettingsFiles(ctx, req.(*LoadSettingsFiles_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_Rename_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Rename_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).Rename(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/Rename",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).Rename(ctx, req.(*Rename_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_RenameCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenameCode_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).RenameCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/RenameCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).RenameCode(ctx, req.(*RenameCode_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Test_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/Test",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).Test(ctx, req.(*Test_Args))
	}
	return interceptor(ctx, in, info, handler)
}

func _KclvmService_UpdateDependencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDependencies_Args)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KclvmServiceServer).UpdateDependencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gpyrpc.KclvmService/UpdateDependencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KclvmServiceServer).UpdateDependencies(ctx, req.(*UpdateDependencies_Args))
	}
	return interceptor(ctx, in, info, handler)
}

var _KclvmService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gpyrpc.KclvmService",
	HandlerType: (*KclvmServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KclvmService_Ping_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _KclvmService_GetVersion_Handler,
		},
		{
			MethodName: "ParseProgram",
			Handler:    _KclvmService_ParseProgram_Handler,
		},
		{
			MethodName: "ParseFile",
			Handler:    _KclvmService_ParseFile_Handler,
		},
		{
			MethodName: "LoadPackage",
			Handler:    _KclvmService_LoadPackage_Handler,
		},
		{
			MethodName: "ListOptions",
			Handler:    _KclvmService_ListOptions_Handler,
		},
		{
			MethodName: "ListVariables",
			Handler:    _KclvmService_ListVariables_Handler,
		},
		{
			MethodName: "ExecProgram",
			Handler:    _KclvmService_ExecProgram_Handler,
		},
		{
			MethodName: "BuildProgram",
			Handler:    _KclvmService_BuildProgram_Handler,
		},
		{
			MethodName: "ExecArtifact",
			Handler:    _KclvmService_ExecArtifact_Handler,
		},
		{
			MethodName: "OverrideFile",
			Handler:    _KclvmService_OverrideFile_Handler,
		},
		{
			MethodName: "GetSchemaTypeMapping",
			Handler:    _KclvmService_GetSchemaTypeMapping_Handler,
		},
		{
			MethodName: "FormatCode",
			Handler:    _KclvmService_FormatCode_Handler,
		},
		{
			MethodName: "FormatPath",
			Handler:    _KclvmService_FormatPath_Handler,
		},
		{
			MethodName: "LintPath",
			Handler:    _KclvmService_LintPath_Handler,
		},
		{
			MethodName: "ValidateCode",
			Handler:    _KclvmService_ValidateCode_Handler,
		},
		{
			MethodName: "ListDepFiles",
			Handler:    _KclvmService_ListDepFiles_Handler,
		},
		{
			MethodName: "LoadSettingsFiles",
			Handler:    _KclvmService_LoadSettingsFiles_Handler,
		},
		{
			MethodName: "Rename",
			Handler:    _KclvmService_Rename_Handler,
		},
		{
			MethodName: "RenameCode",
			Handler:    _KclvmService_RenameCode_Handler,
		},
		{
			MethodName: "Test",
			Handler:    _KclvmService_Test_Handler,
		},
		{
			MethodName: "UpdateDependencies",
			Handler:    _KclvmService_UpdateDependencies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gpyrpc.proto",
}
