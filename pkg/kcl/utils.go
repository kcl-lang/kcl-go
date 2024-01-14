// Copyright The KCL Authors. All rights reserved.

package kcl

import (
	"os"
)

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
