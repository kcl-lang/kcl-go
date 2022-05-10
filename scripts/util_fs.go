package scripts

import (
	"os"
	"path/filepath"
)

func JoinedPath(elem ...string) string {
	return filepath.Join(elem...)
}

func FileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil || fi.IsDir() {
		return false
	}
	return true
}

func DirExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil || !fi.IsDir() {
		return false
	}
	return true
}
