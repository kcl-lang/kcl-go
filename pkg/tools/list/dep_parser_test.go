// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
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

func TestDepParser_graph(t *testing.T) {
	for _, testdata := range testDepParser {
		t.Run(testdata.name, func(t *testing.T) {
			depParser, err := NewImportDepParser(testdata.root, DepOption{Files: testdata.files})
			assert.Nil(t, err, "NewDepParser failed")
			deps := depParser.ListUpstreamFiles()
			assert.ElementsMatch(t, testdata.upStreams, deps)
		})
	}
}

func TestDepParser_affected(t *testing.T) {
	for _, testdata := range testDepParser {
		t.Run(testdata.name, func(t *testing.T) {
			depParser, err := NewImportDepParser(testdata.root, DepOption{Files: testdata.files, ChangedPaths: testdata.changed})
			assert.Nil(t, err, "NewDepParser failed")
			affected := depParser.ListDownStreamFiles()
			assert.ElementsMatch(t, testdata.downStreams, affected)
		})
	}
}

var testDepParser = []struct {
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
			"appops/projectA/base/base.k",
			"appops/projectA/dev/main.k",
			"base/render/server/server_render.k",
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{"base/frontend/container/container_port.k"},
		downStreams: []string{
			"base/frontend/container/container_port.k",
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
			"appops/projectB/base/base.k",
			"appops/projectB/dev/main.k",
			"base/render/job/job_render.k",
			"base/frontend/job",
			"base/frontend/job/job.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
		},
		changed: []string{"base/render/job/job_render.k"},
		downStreams: []string{
			"base/render/job/job_render.k",
			"base/render/job",
		},
	},
	{
		name:  "projectAB",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectA/base/base.k", "appops/projectA/dev/main.k", "base/render/server/server_render.k", "appops/projectB/base/base.k", "appops/projectB/dev/main.k", "base/render/job/job_render.k"},
		upStreams: []string{
			"appops/projectA/base/base.k",
			"appops/projectB/base/base.k",
			"appops/projectA/dev/main.k",
			"appops/projectB/dev/main.k",
			"base/frontend/server",
			"base/frontend/server/server.k",
			"base/frontend/container",
			"base/frontend/container/container.k",
			"base/frontend/container/container_port.k",
			"base/frontend/job",
			"base/frontend/job/job.k",
			"base/render/server/server_render.k",
			"base/render/job/job_render.k",
		},
	},
	{
		name:  "projectC-delete-unused-file",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectC/dev/main.k", "base/render/server/server_render.k"},
		upStreams: []string{
			"appops/projectC/dev/main.k",
			"base/render/server/server_render.k",
		},
		changed: []string{"appops/projectC/base/base.k"},
		downStreams: []string{
			"appops/projectC/base/base.k",
		},
	},
	{
		name:  "projectD-delete-imported",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"appops/projectD/base/base.k",
			"appops/projectD/dev/main.k",
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist/deleted_file.k"},
		downStreams: []string{
			"base/frontend/not_exist/deleted_file.k",
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
			"appops/projectD/base/base.k",
			"appops/projectD/dev/main.k",
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist/deleted_test.k"},
		downStreams: []string{
			"base/frontend/not_exist/deleted_test.k",
		},
	},
	{
		name:  "projectD-delete-imported-pkg",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"appops/projectD/base/base.k",
			"appops/projectD/dev/main.k",
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist"},
		downStreams: []string{
			"base/frontend/not_exist",
			"appops/projectD/base/base.k",
			"appops/projectD/base",
		},
	},
	{
		name:  "projectD-delete-imported-file",
		root:  "./testdata/complicate/",
		files: []string{"appops/projectD/base/base.k", "appops/projectD/dev/main.k"},
		upStreams: []string{
			"appops/projectD/base/base.k",
			"appops/projectD/dev/main.k",
			"base/frontend/not_exist",
		},
		changed: []string{"base/frontend/not_exist.k"},
		downStreams: []string{
			"base/frontend/not_exist.k",
			"base/frontend/not_exist",
			"appops/projectD/base/base.k",
			"appops/projectD/base",
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

func TestSingleAppDepParser_listDepFiles(t *testing.T) {
	pkgroot := "../../../testdata"
	pkgpath := "app0"

	depParser := NewSingleAppDepParser(pkgroot, Option{})

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
