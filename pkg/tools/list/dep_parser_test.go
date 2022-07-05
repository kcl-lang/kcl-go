// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestFindPkgInfo(t *testing.T) {
	pkgroot, pkgpath, err := FindPkgInfo("testdata/mymod/sub/app")
	if err != nil {
		t.Fatal(err)
	}

	wd, _ := os.Getwd()
	expectPkgRoot := filepath.Join(wd, "testdata/mymod")

	if pkgroot != expectPkgRoot {
		t.Fatalf("pkgroot: expect = %s, got = %s", expectPkgRoot, pkgroot)
	}
	if expect := "sub/app"; pkgpath != expect {
		t.Fatalf("pkgpath: expect = %s, got = %s", expect, pkgpath)
	}
}

func TestFindPkgInfo_failed(t *testing.T) {
	if _, _, err := FindPkgInfo("./testdata/no-kcl-mod"); err == nil {
		t.Fatal("expect error, got nil")
	}
}

func TestImportDepParser_ListUpstreamFiles(t *testing.T) {
	for _, tc := range importDepParserTestCases {
		t.Run(tc.name, func(t *testing.T) {
			depParser, err := newImportDepParser(tc.root, DepOptions{Files: tc.files})
			assert.Nil(t, err, "NewDepParser failed")
			deps := depParser.upstreamFiles()
			assert.ElementsMatch(t, tc.upStreams, deps)
		})
	}
}

func TestImportDepParser_ListDownstreamFiles(t *testing.T) {
	for _, tc := range importDepParserTestCases {
		t.Run(tc.name, func(t *testing.T) {
			depParser, err := newImportDepParser(tc.root, DepOptions{Files: tc.files, UpStreams: tc.changed})
			assert.Nil(t, err, "NewDepParser failed")
			affected := depParser.downStreamFiles()
			assert.ElementsMatch(t, tc.downStreams, affected)
		})
	}
}

func BenchmarkImportDepParser_walkDownStream(b *testing.B) {
	tc := importDepParserTestCases[0]
	depParser, err := newImportDepParser(tc.root, DepOptions{Files: tc.files, UpStreams: tc.changed})
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depParser.downStreamFiles()
	}
}

func BenchmarkImportDepParser_walkUpStreamFiles(b *testing.B) {
	tc := importDepParserTestCases[0]
	depParser, err := newImportDepParser(tc.root, DepOptions{Files: tc.files, UpStreams: tc.changed})
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depParser.upstreamFiles()
	}
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
			expect:     "frontend",
		},
		{
			name:       "invalid import: out of program bound",
			filePath:   "base/a.k",
			importPath: "...frontend",
			expect:     "frontend",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, fixImportPath(tc.filePath, tc.importPath))
		})
	}
}

func Test_fixPath(t *testing.T) {
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

func Test_listKFiles(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		expect   []string
	}{
		{
			name:     "path to a KCL file",
			filePath: "base/frontend/container/container.k",
			expect:   []string{"base/frontend/container/container.k"},
		},
		{
			name:     "path to a KCL file without suffix",
			filePath: "base/frontend/container/container",
			expect:   []string{"base/frontend/container/container.k"},
		},
		{
			name:     "path to a KCL package",
			filePath: "base/frontend/container",
			expect:   []string{"base/frontend/container/container.k", "base/frontend/container/container_port.k"},
		},
		{
			name:     "path to a KCL package containing test/internal/non-kcl files",
			filePath: "base/frontend/container/probe",
			expect: []string{
				"base/frontend/container/probe/probe.k",
				"base/frontend/container/probe/exec.k",
				"base/frontend/container/probe/http.k",
				"base/frontend/container/probe/tcp.k",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vfs := os.DirFS("./testdata/complicate")
			assert.ElementsMatch(t, tc.expect, listKFiles(vfs, tc.filePath))
		})
	}
}
func Test_invalidFilePath(t *testing.T) {
	_, err := newImportDepParser("./testdata/complicate/", DepOptions{Files: []string{"appops/projectA/invalid.k"}, UpStreams: []string{}})
	assert.EqualError(t, err, "invalid file path: stat testdata/complicate/appops/projectA/invalid.k: no such file or directory", "err not match")
}

var importDepParserTestCases = []struct {
	name        string
	root        string
	files       []string
	upStreams   []string
	changed     []string
	downStreams []string
}{
	{
		name:  "projectA",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectA/base/base.k", "appops/projectA/dev/main.k", "base/render/server/server_render.k"},
		upStreams: []string{
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{"base/frontend/container/container_port.k"},
		downStreams: []string{
			"base/frontend/container",
			"base/frontend/server/server.k",
			"base/frontend/server",
			"appops/projectA/base",
			"appops/projectA/base/base.k",
		},
	},
	{
		name:  "projectB",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectB/base/base.k", "appops/projectB/dev/main.k", "base/render/job/job_render.k"},
		upStreams: []string{
			"base/frontend/job",
			"base/frontend/job/job.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{"base/render/job/job_render.k"},
		downStreams: []string{
			"base/render/job",
		},
	},
	{
		name:  "projectAB",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectA/base/base.k", "appops/projectA/dev/main.k", "base/render/server/server_render.k", "appops/projectB/base/base.k", "appops/projectB/dev/main.k", "base/render/job/job_render.k"},
		upStreams: []string{
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
			"base/frontend/job",
			"base/frontend/job/job.k",
		},
	},
	{
		name:  "projectE_no_repeat_process_same_import",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectE/base/base.k"},
		upStreams: []string{
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{
			"base/frontend/container/container.k",
		},
		downStreams: []string{
			"appops/projectE/base/base.k",
			"appops/projectE/base",
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
		},
	},
	{
		name:  "projectF-relative-import",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectF/dev/main.k"},
		upStreams: []string{
			"appops/projectF/base/base.k",
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{"base/frontend/container/container.k"},
		downStreams: []string{
			"base/frontend/container",
			"base/frontend/server/server.k",
			"base/frontend/server",
			"appops/projectF/base/base.k",
			"appops/projectF/base",
			"appops/projectF/dev/main.k",
			"appops/projectF/dev",
		},
	},
	{
		name:  "projectG-absolute-import-module",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectG/dev/main.k"},
		upStreams: []string{
			"appops/projectG/base/base.k",
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{"base/frontend/container/container.k"},
		downStreams: []string{
			"base/frontend/container",
			"base/frontend/server/server.k",
			"base/frontend/server",
			"appops/projectG/base/base.k",
			"appops/projectG/base",
			"appops/projectG/dev/main.k",
			"appops/projectG/dev",
		},
	},
	{
		name:      "projectC-delete-unused-file",
		root:      "./testdata/complicate/",
		files:     []string{"appops/projectC/dev/main.k", "base/render/server/server_render.k"},
		upStreams: []string{},
		changed:   []string{"appops/projectC/base/base.k"},
		downStreams: []string{
			"appops/projectC/base",
		},
	},
	{
		name:  "projectD-delete-imported",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist/deleted_file.k"},
		downStreams: []string{
			"base/frontend/not_exist",
			"appops/projectD/base/base.k",
			"appops/projectD/base",
		},
	},
	{
		name:  "projectD-delete-test-file",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"base/frontend/not_exist",
		},
		changed:     []string{"base/frontend/not_exist/deleted_test.k"},
		downStreams: []string{},
	},
	{
		name:  "projectD-delete-imported-pkg",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist"},
		downStreams: []string{
			"appops/projectD/base/base.k",
			"appops/projectD/base",
		},
	},
	{
		name:  "projectD-delete-imported-file",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist.k"},
		downStreams: []string{
			"base/frontend",
			"appops/projectD/base/base.k",
			"appops/projectD/base",
		},
	},
	{
		name:  "projectH-file-directory-name-conflict",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectH/dev/main.k"},
		upStreams: []string{
			"appops/projectH/base/base/base.k",
			"appops/projectH/base/base",
			"base/frontend/server/server/server.k",
			"base/frontend/server/server",
		},
	},
}

func TestDepParser_kclModEnv(t *testing.T) {
	depParser := NewDepParser("./testdata/kcl_mod_env/")
	appFiles := depParser.GetPkgFileList(".")
	expect := []string{"main1.k", "main2.k"}
	for i, file := range appFiles {
		if file != expect[i] {
			t.Fatalf("pkgroot: expect = %s, got = %s", appFiles, expect)
		}
	}
}

func TestDepParser_listDepFiles(t *testing.T) {
	pkgroot := "../../../testdata"
	pkgpath := "app0"

	depParser := NewDepParser(pkgroot, Option{})

	files := depParser.GetAppFiles(pkgpath, true)

	expect := []string{
		"main.k",
		"app0/before/base.k",
		"app0/main.k",
		"app0/sub/sub.k",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}

func TestDepParser_listDepFiles_failed(t *testing.T) {
	pkgroot := "../../../testdata"

	depParser := NewDepParser(pkgroot, Option{})

	err := depParser.GetError()
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	expectErrMsg := "package app0-failed/sub_not_found: no kcl file"
	if !strings.Contains(err.Error(), expectErrMsg) {
		t.Fatalf("expect %q, got %q", expectErrMsg, err)
	}
}

func TestSingleAppDepParser_listDepFiles(t *testing.T) {
	pkgroot := "../../../testdata"
	pkgpath := "app0"

	depParser := NewSingleAppDepParser(pkgroot, Option{})

	files, err := depParser.GetAppFiles(pkgpath, true)
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"main.k",
		"app0/before/base.k",
		"app0/main.k",
		"app0/sub/sub.k",
	}

	sort.Strings(files)
	sort.Strings(expect)

	if !reflect.DeepEqual(files, expect) {
		t.Fatalf("\nexpect = %v\ngot    = %v", expect, files)
	}
}

func TestDepParser_vfs(t *testing.T) {
	for _, tt := range testsVFS {
		depParser := NewDepParserWithFS(tt.vfs)

		appFiles := depParser.GetAppFiles(tt.app, false)
		if !reflect.DeepEqual(appFiles, tt.appFiles) {
			t.Fatalf("appFiles: expect = %v, got = %v", tt.appFiles, appFiles)
		}

		appAllFiles := depParser.GetAppFiles(tt.app, true)
		if !reflect.DeepEqual(appAllFiles, tt.appAllFiles) {
			t.Fatalf("appAllFiles: expect = %v, got = %v", tt.appAllFiles, appAllFiles)
		}

		appAllPkgs := depParser.GetAppPkgs(tt.app, true)
		if !reflect.DeepEqual(appAllPkgs, tt.appAllPkgs) {
			t.Fatalf("appAllPkgs: expect = %v, got = %v", tt.appAllPkgs, appAllPkgs)
		}

	}
}

var testsVFS = []struct {
	app         string
	appFiles    []string
	appAllFiles []string
	appAllPkgs  []string
	importPkgs  []string
	vfs         *fstest.MapFS
}{
	{
		app: "myapp",
		appFiles: []string{
			"myapp/base.k",
			"myapp/main.k",
		},
		appAllFiles: []string{
			"myapp/base.k",
			"myapp/main.k",
			"mypkg/a.k",
			"mypkg/subpkg/b.k",
		},
		appAllPkgs: []string{
			"mypkg",
			"mypkg/subpkg",
		},
		vfs: &fstest.MapFS{
			"kcl.mod": {},

			// myapp/*.k
			"myapp/main.k": {
				Data: []byte("import mypkg"),
			},
			"myapp/base.k": {
				Data: []byte(""),
			},

			// mypkg
			"mypkg/a.k": {
				Data: []byte("import .subpkg"),
			},

			// mypkg/subpkg
			"mypkg/subpkg/b.k": {
				Data: []byte("a = 1"),
			},
		},
	},
}
