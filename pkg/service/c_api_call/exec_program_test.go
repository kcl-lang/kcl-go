//go:build cgo
// +build cgo

package capicall

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

const EXEC_DATA_PATH = "./exec_data/"

func getFiles(root string, suffix string, sorted bool) []string {
	var files = []string{}
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, suffix) {
			files = append(files, path)
		}
		return nil
	})
	if sorted {
		sort.Strings(files)
	}
	return files
}

func TestExecSingleFile(t *testing.T) {
	client := PROTOCAPI_NewKclvmServiceClient()
	files := getFiles(EXEC_DATA_PATH, ".k", true)
	for _, file := range files {
		exec(t, file, client)
	}
}

func exec(t *testing.T, fileName string, client *PROTOCAPI_KclvmServiceClient) {
	args := &gpyrpc.ExecProgram_Args{
		KFilenameList: []string{fileName},
	}
	_, err := client.ExecProgram(args)
	assert.Nil(t, err)
}
