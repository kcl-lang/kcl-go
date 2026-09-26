// Copyright The KCL Authors. All rights reserved.

package rename

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// Rename renames all the occurrences of the target symbol in the files and
// returns the file paths that got changed.
func Rename(packageRoot, symbolPath string, filePaths []string, newName string) ([]string, error) {
	svc := kcl.Service()
	resp, err := svc.Rename(&gpyrpc.RenameArgs{
		PackageRoot: packageRoot,
		SymbolPath:  symbolPath,
		FilePaths:   filePaths,
		NewName:     newName,
	})
	if err != nil {
		return nil, err
	}
	return resp.ChangedFiles, nil
}

// RenameCode renames all the occurrences of the target symbol in the source
// codes and returns the changed codes. It does not rewrite files.
func RenameCode(packageRoot, symbolPath string, sourceCodes map[string]string, newName string) (map[string]string, error) {
	svc := kcl.Service()
	resp, err := svc.RenameCode(&gpyrpc.RenameCodeArgs{
		PackageRoot: packageRoot,
		SymbolPath:  symbolPath,
		SourceCodes: sourceCodes,
		NewName:     newName,
	})
	if err != nil {
		return nil, err
	}
	return resp.ChangedCodes, nil
}
