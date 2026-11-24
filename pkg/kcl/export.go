package kcl

import (
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// ExportRust generates Rust DTOs (with a `validate()` method mirroring each
// schema's `check:` block) from the KCL packages/files at pathList.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//
// Returns:
//   - The generated Rust source.
//   - An error if there is any failure in the process.
func ExportRust(pathList []string) (string, error) {
	resp, err := Service().GenerateRust(&gpyrpc.GenerateRustArgs{
		ParseArgs: &gpyrpc.ParseProgramArgs{Paths: pathList},
	})
	if err != nil {
		return "", err
	}
	return resp.RustSource, nil
}

// ExportPython generates Python dataclasses (with a `validate()` method
// mirroring each schema's `check:` block) from the KCL packages/files at
// pathList.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//
// Returns:
//   - The generated Python source.
//   - An error if there is any failure in the process.
func ExportPython(pathList []string) (string, error) {
	resp, err := Service().GeneratePython(&gpyrpc.GeneratePythonArgs{
		ParseArgs: &gpyrpc.ParseProgramArgs{Paths: pathList},
	})
	if err != nil {
		return "", err
	}
	return resp.PythonSource, nil
}

// ExportGo generates Go structs (with a `Validate()` method mirroring each
// schema's `check:` block) from the KCL packages/files at pathList.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//   - pkg: The Go package name for the generated file's `package` clause.
//     Empty defaults to "generated".
//
// Returns:
//   - The generated Go source.
//   - An error if there is any failure in the process.
func ExportGo(pathList []string, pkg string) (string, error) {
	resp, err := Service().GenerateGo(&gpyrpc.GenerateGoArgs{
		ParseArgs: &gpyrpc.ParseProgramArgs{Paths: pathList},
		Package:   pkg,
	})
	if err != nil {
		return "", err
	}
	return resp.GoSource, nil
}

// ExportJava generates Java classes (with a `validate()` method mirroring
// each schema's `check:` block) from the KCL packages/files at pathList.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//   - pkg: The Java package name for the generated file's `package` clause.
//     Empty omits the package clause entirely (Java's default package).
//
// Returns:
//   - The generated Java source.
//   - An error if there is any failure in the process.
func ExportJava(pathList []string, pkg string) (string, error) {
	resp, err := Service().GenerateJava(&gpyrpc.GenerateJavaArgs{
		ParseArgs: &gpyrpc.ParseProgramArgs{Paths: pathList},
		Package:   pkg,
	})
	if err != nil {
		return "", err
	}
	return resp.JavaSource, nil
}

// ExportTypeScript generates TypeScript classes (with a `validate()` method
// mirroring each schema's `check:` block) from the KCL packages/files at
// pathList.
//
// Parameters:
//   - pathList: A list of KCL file paths.
//
// Returns:
//   - The generated TypeScript source.
//   - An error if there is any failure in the process.
func ExportTypeScript(pathList []string) (string, error) {
	resp, err := Service().GenerateTypeScript(&gpyrpc.GenerateTypeScriptArgs{
		ParseArgs: &gpyrpc.ParseProgramArgs{Paths: pathList},
	})
	if err != nil {
		return "", err
	}
	return resp.TypescriptSource, nil
}
