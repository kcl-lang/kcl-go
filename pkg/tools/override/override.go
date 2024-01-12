// Copyright The KCL Authors. All rights reserved.

package override

import (
	"errors"
	"fmt"
	"strings"

	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

const (
	DeleteAction         = "Delete"
	CreateOrUpdateAction = "CreateOrUpdate"
)

func OverrideFile(file string, specs, importPaths []string) (result bool, err error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.OverrideFile(&gpyrpc.OverrideFile_Args{
		File:        file,
		Specs:       specs,
		ImportPaths: importPaths,
	})
	if err != nil {
		return false, err
	}
	return resp.Result, nil
}

func ParseOverrideSpec(spec string) (*gpyrpc.CmdOverrideSpec, error) {
	if strings.Contains(spec, "=") {
		// Create or update the override value.
		splitValues := strings.SplitN(spec, "=", 2)
		if len(splitValues) < 2 {
			return nil, invalidSpecError(spec)
		}
		path := splitValues[0]
		fieldValue := splitValues[1]
		pkgpath, fieldPath, err := splitFieldPath(path)
		if err != nil {
			return nil, err
		}
		return &gpyrpc.CmdOverrideSpec{
			Pkgpath:    pkgpath,
			FieldPath:  fieldPath,
			FieldValue: fieldValue,
			Action:     CreateOrUpdateAction,
		}, nil
	} else if strippedSpec := strings.TrimSuffix(spec, "-"); strippedSpec != spec {
		// Delete the override value.
		pkgpath, fieldPath, err := splitFieldPath(strippedSpec)
		if err != nil {
			return nil, err
		}
		return &gpyrpc.CmdOverrideSpec{
			Pkgpath:    pkgpath,
			FieldPath:  fieldPath,
			FieldValue: "",
			Action:     DeleteAction,
		}, nil
	} else {
		return nil, invalidSpecError(spec)
	}
}

// Get field package path and identifier name from the path.
//
// split_field_path("pkg.to.path:field") -> ("pkg.to.path", "field")
func splitFieldPath(path string) (string, string, error) {
	err := errors.New("Invalid field path " + path)
	paths := strings.SplitN(path, ":", 2)
	if len(paths) == 1 {
		pkgpath := ""
		fieldPath := paths[0]
		if fieldPath == "" {
			return "", "", err
		}
		return pkgpath, fieldPath, nil
	} else if len(paths) == 2 {
		pkgpath := paths[0]
		fieldPath := paths[1]
		if fieldPath == "" {
			return "", "", err
		}
		return pkgpath, fieldPath, nil
	} else {
		return "", "", err
	}
}

// / Get the invalid spec error message.
func invalidSpecError(spec string) error {
	return fmt.Errorf("invalid spec format '%s', expected <pkgpath>:<field_path>=<filed_value> or <pkgpath>:<field_path>-", spec)
}
