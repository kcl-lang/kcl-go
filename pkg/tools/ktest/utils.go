// Copyright 2021 The KCL Authors. All rights reserved.

package ktest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kusionstack.io/kclvm-go/pkg/ast"
	settings_pkg "kusionstack.io/kclvm-go/pkg/settings"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ = fmt.Sprint

func readKclFiles(path string) (kFiles, kSettingFiles, kPluginFiles, kTestFiles []string) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		if name == "settings.yaml" {
			kSettingFiles = append(kSettingFiles, filepath.Join(path, name))
			continue
		}
		if name == "plugin.py" {
			kPluginFiles = append(kPluginFiles, filepath.Join(path, name))
			continue
		}
		if strings.HasSuffix(name, ".k") {
			if strings.HasSuffix(name, "_test.k") {
				kTestFiles = append(kTestFiles, filepath.Join(path, name))
			} else {
				kFiles = append(kFiles, filepath.Join(path, name))
			}
			continue
		}
	}

	sort.Strings(kFiles)
	sort.Strings(kTestFiles)

	return
}

func getTestSchemaNameList(p *ast.File) []string {
	if p.Module == nil || len(p.Module.Body) == 0 {
		return nil
	}
	var names []string
	for _, stmt := range p.Module.Body {
		switch stmt := stmt.(type) {
		case *ast.SchemaStmt:
			if strings.HasPrefix(stmt.Name, "Test") {
				names = append(names, stmt.Name)
			}
		}
	}
	return names
}

func getTestSchemaInfo(workDir string, p *ast.File, name string) (*kTestSchemaInfo, error) {
	var testingName = "testing"

	var schema *ast.SchemaStmt
	for _, stmt := range p.Module.Body {
		switch stmt := stmt.(type) {
		case *ast.SchemaStmt:
			if stmt.Name == name {
				schema = stmt
			}
		case *ast.ImportStmt:
			if stmt.Path == "testing" && stmt.AsName != "" {
				testingName = stmt.AsName
			}
		}
	}
	if schema == nil {
		return nil, nil
	}

	var fnGetFuncName = func(call *ast.CallExpr) string {
		fnIdent, ok := call.Func.(*ast.Identifier)
		if !ok {
			return ""
		}
		return strings.Join(fnIdent.Names, ".")
	}
	var fnGetFuncArgs = func(call *ast.CallExpr) ([]string, error) {
		var ss []string
		for _, x := range call.Args {
			switch x := x.(type) {
			case *ast.NameConstantLit:
				ss = append(ss, x.Value)
			case *ast.NumberLit:
				ss = append(ss, fmt.Sprintf("%v", x.Value))
			case *ast.StringLit:
				ss = append(ss, x.Value)
			default:
				return nil, fmt.Errorf("invalid argument: %v", x.JSONString())
			}
		}
		return ss, nil
	}

	var info kTestSchemaInfo

	for _, stmt := range schema.Body {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		for _, expr := range exprStmt.Exprs {
			callExpr, ok := expr.(*ast.CallExpr)
			if !ok {
				continue
			}

			switch fnName := fnGetFuncName(callExpr); fnName {
			case fmt.Sprintf("%s.arguments", testingName):
				args, err := fnGetFuncArgs(callExpr)
				if err != nil {
					return nil, fmt.Errorf("testiong.arguments: argument only support basic literal type(bool,number,str).")
				}
				if len(args) == 2 {
					info.Args = append(info.Args, &gpyrpc.CmdArgSpec{
						Name:  args[0],
						Value: args[1],
					})
				}

			case fmt.Sprintf("%s.setting_file", testingName):
				args, err := fnGetFuncArgs(callExpr)
				if err != nil {
					return nil, fmt.Errorf("testiong.setting_file: argument only support string literal type.")
				}
				if len(args) == 1 {
					info.SettingFile = args[0]
					if !filepath.IsAbs(info.SettingFile) {
						info.SettingFile = filepath.Join(workDir, info.SettingFile)
					}

					if settings, err := settings_pkg.LoadFile(info.SettingFile, nil); err == nil {
						for _, t := range settings.Options {
							var key string = t.Key
							var val string

							switch v := t.Value.(type) {
							case map[string]interface{}:
								if s, err := json.Marshal(v); err == nil {
									val = string(s)
								} else {
									val = fmt.Sprint(v)
								}
							case []interface{}:
								if s, err := json.Marshal(v); err == nil {
									val = string(s)
								} else {
									val = fmt.Sprint(v)
								}
							default:
								val = fmt.Sprint(v)
							}

							info.SettingFile_Args = append(info.SettingFile_Args, &gpyrpc.CmdArgSpec{
								Name:  key,
								Value: val,
							})
						}
					}
				}
			}
		}
	}

	var fnCmdArgSpecHasName = func(args []*gpyrpc.CmdArgSpec, name string) bool {
		for _, x := range args {
			if x.Name == name {
				return true
			}
		}
		return false
	}

	// merge setting_file args
	for _, x := range info.SettingFile_Args {
		if !fnCmdArgSpecHasName(info.Args, x.Name) {
			info.Args = append(info.Args, x)
		}
	}

	return &info, nil
}

func withLinePrefix(s, prefix string) string {
	var ss []string
	for _, s := range strings.Split(strings.TrimSpace(s), "\n") {
		ss = append(ss, prefix+s)
	}
	return strings.Join(ss, "\n")
}
