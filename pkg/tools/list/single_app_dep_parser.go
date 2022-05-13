// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"io/fs"
	"os"
	pathpkg "path"
	"sort"
)

type SingleAppDepParser struct {
	opt Option
	vfs fs.FS

	appPkgpath  string
	importMap   map[string][]string
	pkgFilesMap map[string][]string

	allfiles []string
}

func NewSingleAppDepParser(root string, opt ...Option) *SingleAppDepParser {
	root = pathpkg.Clean(root)
	return NewSingleAppDepParserWithFS(os.DirFS(root), opt...)
}

func NewSingleAppDepParserWithFS(vfs fs.FS, opts ...Option) *SingleAppDepParser {
	p := &SingleAppDepParser{
		vfs:         vfs,
		importMap:   make(map[string][]string),
		pkgFilesMap: make(map[string][]string),
	}

	for _, opt := range opts {
		p.opt.merge(&opt)
	}
	p.opt.adjust()

	return p
}

func (p *SingleAppDepParser) GetAppFiles(appPkgpath string, includeDependFiles bool) []string {
	p.parseOnce(appPkgpath)

	if includeDependFiles {
		return p.allfiles
	}
	return p.pkgFilesMap[appPkgpath]
}

func (p *SingleAppDepParser) GetAppPkgs(appPkgpath string, includeDependFiles bool) []string {
	p.parseOnce(appPkgpath)

	if includeDependFiles {
		var pkgs []string
		for k := range p.importMap {
			pkgs = append(pkgs, k)
		}
		sort.Strings(pkgs)
		return pkgs
	}

	return p.importMap[appPkgpath]
}

func (p *SingleAppDepParser) parseOnce(appPkgpath string) {
	if p.appPkgpath == appPkgpath {
		return
	}

	p.appPkgpath = appPkgpath
	p.importMap = make(map[string][]string)
	p.pkgFilesMap = make(map[string][]string)
	p.allfiles = []string{}

	p.scanAppFiles(appPkgpath)

	var filesMap = make(map[string]struct{})
	for _, files := range p.pkgFilesMap {
		for _, s := range files {
			filesMap[s] = struct{}{}
		}
	}
	for s := range filesMap {
		p.allfiles = append(p.allfiles, s)
	}
	sort.Strings(p.allfiles)
}

func (p *SingleAppDepParser) scanAppFiles(pkgpath string) {
	if _, ok := p.pkgFilesMap[pkgpath]; ok {
		return
	}

	// 1. loadKFileList
	k_files := loadKFileList(p.vfs, pkgpath, p.opt)
	p.pkgFilesMap[pkgpath] = k_files

	// 2. parse import
	var importMap = make(map[string]string)
	for _, file := range k_files {
		src, err := fs.ReadFile(p.vfs, file)
		if err != nil {
			panic(err)
		}

		for _, import_path := range parseImport(string(src)) {
			import_path := fixImportPath(file, import_path)
			importMap[import_path] = import_path
		}
	}

	// 3. save import list
	var importList []string
	for import_path := range importMap {
		importList = append(importList, import_path)
	}
	sort.Strings(importList)
	p.importMap[pkgpath] = importList

	// 4. scan import
	for _, import_path := range importList {
		p.scanAppFiles(import_path)
	}
}
