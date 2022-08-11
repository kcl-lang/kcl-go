//go:build cgo && kclvm_service_capi
// +build cgo,kclvm_service_capi

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

const CORRECT_DATA_PATH = "./exec_data/correct_data"

const ERROR_DATA_PATH = "./exec_data/error_data"

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

func TestExecCorrectSingleFile(t *testing.T) {
	client := PROTOCAPI_NewKclvmServiceClient()
	files := getFiles(CORRECT_DATA_PATH, ".k", true)
	for _, file := range files {
		_, err := exec(file, client)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestExecErrorSingleFile(t *testing.T) {
	client := PROTOCAPI_NewKclvmServiceClient()
	files := getFiles(ERROR_DATA_PATH, ".k", true)
	for _, file := range files {
		_, err := exec(file, client)
		assert.NotNil(t, err)
	}
}

func exec(fileName string, client *PROTOCAPI_KclvmServiceClient) (out *gpyrpc.ExecProgram_Result, err error) {
	args := &gpyrpc.ExecProgram_Args{
		KFilenameList: []string{fileName},
	}
	return client.ExecProgram(args)
}
