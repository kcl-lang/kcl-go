// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"testing/fstest"
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

	includeDependFiles := true
	files := depParser.GetAppFiles(pkgpath, includeDependFiles)

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

	includeDependFiles := true
	files := depParser.GetAppFiles(pkgpath, includeDependFiles)

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
