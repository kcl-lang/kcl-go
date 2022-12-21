package format

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCode(t *testing.T) {
	tcases := [...]struct {
		source    string
		expect    string
		expectErr string
	}{
		{
			source: "a=a+1",
			expect: "a = a + 1\n",
		},
		{
			source:    "a=a+",
			expectErr: "KCL Syntax Error[E1001] : Invalid syntax\nInvalid syntax",
		},
	}

	for _, testCase := range tcases {
		actual, err := FormatCode(testCase.source)
		if testCase.expectErr != "" {
			assert.NotNil(t, err, "format code expect err, get no error")
			assert.Equal(t, testCase.expectErr, err.Error(), fmt.Sprintf("format code get wrong error result. expect: %s got: %s", testCase.expectErr, err.Error()))
		} else {
			assert.Equal(t, testCase.expect, string(actual), fmt.Sprintf("format file get wrong result. expect: %s got: %s", actual, testCase.expect))
		}
	}
}

// TODO: fix Broken pipe in kclvm_py, see in :
// https://github.com/KusionStack/kclvm-go/issues/75
/*
func TestFormatPath(t *testing.T) {
	successDir := "testdata/success"
	expectedFileSuffix := ".formatted"
	expectedFiles := findFiles(t, successDir, func(info fs.FileInfo) bool {
		return strings.HasSuffix(info.Name(), expectedFileSuffix)
	})

	sourceFiles := findFiles(t, successDir, func(info fs.FileInfo) bool {
		return strings.HasSuffix(info.Name(), ".k")
	})
	var sourceFilesBackup []kclFile

	for _, sourceFile := range sourceFiles {
		content, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			t.Fatalf("read source file content failed: %s", sourceFile)
		}
		sourceFilesBackup = append(sourceFilesBackup, kclFile{
			name:    sourceFile,
			content: content,
		})
	}

	changedPaths, err := FormatPath(successDir)
	// write back un-formatted file content
	defer writeFile(t, sourceFilesBackup)

	if err != nil {
		t.Fatalf("format path exec failed. %v", err)
	}

	var changedPathsRelative []string
	for _, p := range changedPaths {
		changedPathsRelative = append(changedPathsRelative, strings.TrimSuffix(path.Join(successDir, p), ".k")+expectedFileSuffix)
	}
	assert.ElementsMatchf(t, expectedFiles, changedPathsRelative, "format path get wrong result. changedPath mismatch, expect: %s, get: %s", expectedFiles, changedPathsRelative)

	for _, expectedFile := range expectedFiles {
		expected, err := ioutil.ReadFile(expectedFile)
		if err != nil {
			t.Fatalf("read expected formatted file failed: %s", expectedFile)
		}
		actualFile := strings.TrimSuffix(expectedFile, expectedFileSuffix) + ".k"
		get, err := ioutil.ReadFile(actualFile)
		if err != nil {
			t.Fatalf("read actual formatted file failed: %s", actualFile)
		}
		assert.Equal(t, expected, get, fmt.Sprintf("format path get wrong result. formatted content mismatch, file: %s, expect: %s, get: %s", actualFile, expected, get))
	}
}
*/

type filterFile func(fs.FileInfo) bool

func findFiles(t testing.TB, testDir string, filter filterFile) (names []string) {
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	for _, f := range files {
		if !f.IsDir() {
			if filter(f) {
				names = append(names, path.Join(testDir, f.Name()))
			}
		}
	}
	return names
}

type kclFile struct {
	name    string
	content []byte
}

func writeFile(t *testing.T, kclfiles []kclFile) {
	for _, backUpFile := range kclfiles {
		err := ioutil.WriteFile(backUpFile.name, backUpFile.content, 0666)
		if err != nil {
			t.Logf("write back formatted source file failed: %v", err)
		}
	}
}
