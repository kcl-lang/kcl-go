//go:build linux || darwin
// +build linux darwin

package list

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFilePath(t *testing.T) {
	_, err := newImportDepParser("./testdata/complicate/", DepOptions{Files: []string{"appops/projectA/invalid.k"}, UpStreams: []string{}})
	assert.Equal(t, strings.Contains(err.Error(), "appops/projectA/invalid.k: no such file or directory"), true)
}

func TestImportDepParser_fixImportPath(t *testing.T) {
	testCases := []struct {
		name       string
		filePath   string
		importPath string
		expect     string
	}{
		{
			name:       "absolute import",
			filePath:   "main.k",
			importPath: "base.b",
			expect:     "base/b",
		},
		{
			name:       "relative import1",
			filePath:   "base/b.k",
			importPath: ".a",
			expect:     "base/a",
		},
		{
			name:       "relative import2",
			filePath:   "base/a.k",
			importPath: "..frontend",
			expect:     "base/../frontend",
		},
		{
			name:       "invalid import: out of program bound",
			filePath:   "base/a.k",
			importPath: "...frontend",
			expect:     "base/../../frontend",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, fixImportPath(tc.filePath, tc.importPath))
		})
	}
}

func TestFixPath(t *testing.T) {
	testCases := []struct {
		name    string
		oriPath string
		expect  string
	}{
		{
			name:    "file path with .k suffix",
			oriPath: "base/frontend/container/container.k",
			expect:  "base/frontend/container/container.k",
		},
		{
			name:    "file path without .k suffix",
			oriPath: "base/frontend/container/container",
			expect:  "base/frontend/container/container.k",
		},
		{
			name:    "dir path",
			oriPath: "base/frontend/container",
			expect:  "base/frontend/container",
		},
		{
			name:    "dir path",
			oriPath: "base/frontend/container/invalid",
			expect:  "base/frontend/container/invalid",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vfs := os.DirFS("./testdata/complicate")
			assert.Equal(t, tc.expect, fixPath(vfs, tc.oriPath))
		})
	}
}

func TestListDepPackages(t *testing.T) {
	files, err := ListDepPackages("./testdata/module_with_external/", &Option{
		ExcludeExternalPackage: true,
		ExcludeBuiltin:         true,
		IgnoreImportError:      true,
	})
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"./../../../relative_pkg4",
		"./../../relative_pkg3",
		"./../relative_pkg2",
		"./relative_pkg1",
		"pkg1",
		"pkg2",
		"pkg3/internal/pkg",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}
