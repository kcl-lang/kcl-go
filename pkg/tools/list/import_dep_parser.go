// Copyright The KCL Authors. All rights reserved.

package list

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strings"
)

var _ = fmt.Sprint

// DepOptions provides the files to inspect and list dependencies on
type DepOptions struct {
	// Files defines the scope to inspect the import dependency on.
	// Each value in the list should be a file or package path relative to the workdir.
	Files []string
	// UpStreams defines a list of UpStream file/package paths to ListDownStreamFiles on.
	// Each value in the list should be a file or package path relative to the workdir.
	// To list UpStream files/packages, this field will not be used and can be set nil or empty
	UpStreams []string
}

// importDepParser builds an import graph based on parsing the import statements in KCL files within a KCL work directory
type importDepParser struct {
	opt         DepOptions
	vfs         fs.FS        // vfs is the file system of the KCL work directory root
	importGraph *importGraph // importGraph an import dependency graph of files/packages based on files in opt and vfs
}

// newImportDepParser creates an importDepParser and then builds an import graph by calling the inspect function on each file.
// The DepOptions.Files defines the scope to inspect the import dependency on.
// Thus, only files/packages that have directly or indirectly UpStream/DownStream relations(or in the same package with them) with those files will be inspected.
func newImportDepParser(root string, opt DepOptions) (p *importDepParser, err error) {
	root = pathpkg.Clean(root)
	vfs := os.DirFS(root)
	for _, file := range opt.Files {
		// each file should be a valid file or package path relative to the vfs root directory
		if _, err := fs.Stat(vfs, file); err != nil {
			return nil, fmt.Errorf("invalid file path: %s", err)
		}
	}
	p = &importDepParser{
		vfs:         vfs,
		importGraph: newImportGraph(),
		opt:         opt,
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	for _, file := range opt.Files {
		p.inspect(file)
	}
	return p, nil
}

// importGraph stores index and reverted index on import dependency and package/file map to indicate upstream/downstream.
// By walking the import index and file index recursively, you can then walk through all the upstream paths of the start path.
// And on the contrary, by walking the inverted import index and inverted file index recursively, you can then walk through all the downstream paths of the start path.
//
// # Example
//
// For instance a KCL program that comes with three files: main.k and base/a.k, base/b.k,
// and the file main.k contains an import statement that imports base/b.k to it, while the file base/b.k imports base/a.k:
//
//	 # main.k imports base/b.k
//	 # base/b.k imports base/a.k
//
//		demo (KCL program root)
//		├── base
//		│   ├── a.k
//		│   └── b.k         # import .a
//		└── main.k          # import base.b
//
// the importGraph of the program will be:
//
//	 importIndex: map[string]stringSet{
//	 	"main.k":   {values: map[string]bool{"base/b.k": true}},
//	 	"base/b.k": {values: map[string]bool{"base/a.k": true}},
//	 },
//	 importIndexInverted: map[string]stringSet{
//	 	"base/b.k": {values: map[string]bool{"main.k": true}},
//			"base/a.k": {values: map[string]bool{"base/b.k": true}},
//	 },
//	 fileIndex: map[string]stringSet{
//	 	".":    {values: map[string]bool{"main.k": true}},
//	 	"base": {values: map[string]bool{"base/a.k": true, "base/b.k": true}},
//	 },
//	 fileIndexInverted: map[string]string{
//	 	"base/a.k": "base",
//	 	"base/b.k": "base",
//	 	"main.k":   ".",
//	 }
type importGraph struct {
	// importIndex indicates list of KCL files and which paths appear in them as import path
	importIndex map[string]stringSet
	// importIndexInverted indicates the list of paths, and the KCL files in which they appear as import paths
	importIndexInverted map[string]stringSet
	// fileIndex indicates list of package paths and which KCL files contained in them
	fileIndex map[string]stringSet
	// fileIndexInverted indicates the list of KCL files, and the package paths in which they are located
	fileIndexInverted map[string]string
	// processed is a list of processed package/module paths to make sure each file is inspected only once
	processed stringSet
}

// newImportGraph creates an empty import graph
func newImportGraph() *importGraph {
	return &importGraph{
		importIndex:         make(map[string]stringSet),
		importIndexInverted: make(map[string]stringSet),
		fileIndexInverted:   make(map[string]string),
		fileIndex:           make(map[string]stringSet),
		processed:           emptyStringSet(),
	}
}

// stringSet is a simple string set implementation by map
type stringSet map[string]bool

// emptyStringSet creates a string set with an empty value list
func emptyStringSet() stringSet {
	return make(stringSet)
}

// add a string value to the stringSet s
func (s stringSet) add(value string) {
	s[value] = true
}

// contains checks if the stringSet s contains certain string value
func (s stringSet) contains(value string) bool {
	_, ok := s[value]
	return ok
}

// toSlice generates a string slice containing all the string values in the stringSet s
func (s stringSet) toSlice() []string {
	var result []string
	for value := range s {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

// inspect function parses import stmts in files under current package path and builds the import graph incrementally and recursively.
// If path is a file path with .k suffix, all files under the same package path with that file are parsed.
// If path is a package path without .k suffix, all files under that package path are parsed.
func (p *importDepParser) inspect(path string) {
	// 1. find all KCL files in current package path
	var kFiles []string
	pkgPath := path
	if strings.HasSuffix(path, ".k") {
		pkgPath = pathpkg.Dir(path)
	}

	// 2. check and update the "processed" flag before actually inspecting to make sure each path is inspected only once
	if p.importGraph.processed.contains(pkgPath) {
		return
	}
	p.importGraph.processed.add(pkgPath)

	// 3. for each file in current package, parse import statement and build import index recursively
	kFiles = listKFiles(p.vfs, pkgPath)
	for _, f := range kFiles {
		// 3.1 record index and reverted index of file/packages
		p.importGraph.fileIndexInverted[f] = pkgPath
		addValueOnIndex(p.importGraph.fileIndex, pkgPath, f)

		// 3.2 read file content and extract import paths from it
		src, err := fs.ReadFile(p.vfs, f)
		if err != nil {
			panic(err)
		}
		for _, importPath := range parseImport(string(src)) {
			importPath = fixPath(p.vfs, fixImportPath(f, importPath))

			// 3.3 for each import path, record index and inverted index of importing
			addValueOnIndex(p.importGraph.importIndex, f, importPath)
			addValueOnIndex(p.importGraph.importIndexInverted, importPath, f)

			// 3.4 inspect the import path recursively
			p.inspect(importPath)
		}
	}
}

// addValueOnIndex checks if some index exists and initialize it if not, then adds a string value to that index in the given map[string]stringSet.
func addValueOnIndex(m map[string]stringSet, index string, value string) {
	if _, ok := m[index]; !ok {
		m[index] = emptyStringSet()
	}
	m[index].add(value)
}

// upstreamFiles walks through the import graph of the importDepParser p, and lists all the upstream files of the given file path list.
// The walk starts from all the files defined in the DepOptions of importDepParser.
func (p *importDepParser) upstreamFiles() []string {
	upFiles := emptyStringSet()
	for _, f := range p.opt.Files {
		p.importGraph.walkUpstream(f, func(filepath string) (walked bool) {
			if walked = upFiles.contains(filepath); !walked {
				upFiles.add(filepath)
			}
			return
		})
	}
	return upFiles.toSlice()
}

// downStreamFiles walks through the import graph of the importDepParser p, and lists all the downstream files of the given upstreams path list.
// The walk starts from all the upstream files defined in the DepOptions of importDepParser.
// But since the upstream files can be non-existent(The file might have been deleted, and you might want to know the DownStreams of that deleted file),
// before walking downStreams, the file index/inverted index, and the import index/inverted index of those non-existent files should be added to the import graph
func (p *importDepParser) downStreamFiles() []string {
	downFiles := emptyStringSet()
	for _, f := range p.opt.UpStreams {
		if !shouldIgnore(filepath.Base(f)) {
			// if the KCL file does not exist, that means it might have been deleted.
			// so we just infer the package path and the module path of that file and make sure to take them into account.
			if _, err := fs.Stat(p.vfs, f); errors.Is(err, os.ErrNotExist) {
				// 1. add the package path to the inverted file index to record that the package is in the downstream of that delete file
				pkgPath := pathpkg.Dir(f)
				p.importGraph.fileIndexInverted[f] = pkgPath
				// 2. add the module path to the UpStreams Opt to guarantee that files which import the module path will be involved to the walk
				modulePath := strings.TrimSuffix(f, ".k")
				p.opt.UpStreams = append(p.opt.UpStreams, modulePath)
			}
		}
	}
	for _, f := range p.opt.UpStreams {
		p.importGraph.walkDownstream(f, func(filepath string) (walked bool) {
			if walked = downFiles.contains(filepath); !walked {
				downFiles.add(filepath)
			}
			return
		})
	}
	return downFiles.toSlice()
}

// walkUpstream walks through the importGraph starting from the start file and up to the files the start file imports recursively
func (g *importGraph) walkUpstream(start string, walkFunc func(filepath string) (walked bool)) {
	// 1. collect all the upstream files/packages of the start file
	nexts := g.importIndex[start]
	if nexts == nil {
		return
	}
	for next := range nexts {
		// 2. for each path in the upstream list:
		// 2.1 walk that path if it has been walked
		if walkFunc(next) {
			continue
		}
		// 2.2 when the path is a package path, all the files in that package will be walked recursively
		if files, ok := g.fileIndex[next]; ok {
			for file := range files {
				if walkFunc(file) {
					continue
				}
				g.walkUpstream(file, walkFunc)
			}
		} else {
			// 2.2 when the path is a file path, walk it recursively
			g.walkUpstream(next, walkFunc)
		}
	}
}

// walkDownstream walks the importGraph starting from a start path and walks down to its downstream files recursively
func (g *importGraph) walkDownstream(start string, walkFunc func(filepath string) (walked bool)) {
	// 1. collect all the downstream files/packages of the start file
	// 1.1 list one step down stream files by searching the start file in the import inverted index
	nexts := g.importIndexInverted[start]
	if nexts == nil {
		nexts = emptyStringSet()
	}
	// 1.2 get the containing pkg path by searching the start file in the file inverted index
	pkg, ok := g.fileIndexInverted[start]
	if ok {
		// the package containing the file is in the file's DownStreams, too
		nexts.add(pkg)
	}
	// 2. call the walkFunc on each one of the DownStream files/packages, and recursively walk on them
	for next := range nexts {
		if walkFunc(next) {
			continue
		}
		g.walkDownstream(next, walkFunc)
	}
}

// fixPath fixes a path (an import path or a file path) to a valid file path
//
// That's how the path will be fixed:
// 1. a file path with .k suffix will be kept intact;
// 2. for an import path without .k suffix, such as a/b/c: if a/b/c exists as a dir, the path will be kept intact. And if a/b/c.k exists as a file, the fixed path will be a/b/c.k;
// 3. otherwise the path will be kept intact
func fixPath(vfs fs.FS, path string) string {
	if strings.HasSuffix(path, ".k") {
		return path
	}
	if fi, _ := fs.Stat(vfs, path); fi != nil && fi.IsDir() {
		return path
	}
	if fi, _ := fs.Stat(vfs, path+".k"); fi != nil && !fi.IsDir() {
		return path + ".k"
	}
	return path
}

// listKFiles lists all the KCL file paths by the given path.
// Those files will not be included in the result: Non-KCL files, private KCL files with "_" prefix, KCL test files with "_test.k" suffix.
//
// The given path might be a directory or a file path.
// If the path is a directory, the function will list all the KCL files in that dir.
// If the path is a file path, the function will return a list only containing the file path itself.
// If no KCL files found, an empty list will be returned.
func listKFiles(vfs fs.FS, path string) []string {
	if strings.HasSuffix(path, ".k") {
		return []string{path}
	}

	var kFiles []string
	if fi, _ := fs.Stat(vfs, path); fi != nil && fi.IsDir() {
		entryList, _ := fs.ReadDir(vfs, path)
		for _, info := range entryList {
			// just list KCL files directly under the path
			if info.IsDir() {
				continue
			}
			// skip files that should be ignored
			if shouldIgnore(info.Name()) {
				continue
			}
			// OK
			kFiles = append(kFiles, pathpkg.Join(path, info.Name()))
		}
		return kFiles
	}

	if fi, _ := fs.Stat(vfs, path+".k"); fi != nil && !fi.IsDir() {
		return []string{path + ".k"}
	}

	return kFiles
}

// shouldIgnore checks a file name and returns if the file should be ignored when parsing KCL import
func shouldIgnore(name string) bool {
	// ignore non-KCL files, _xxx.k(private kcl files), xxx_test.k(test files)
	if !strings.HasSuffix(name, ".k") || strings.HasPrefix(name, "_") || strings.HasSuffix(name, "_test.k") {
		return true
	}
	return false
}

// parseImport parses the KCL code and extracts the import paths from import statements
// For instance, the import paths parsed in following code with content will be: []string{"base.frontend", "base.api.core.v1"}
//
//	import base.frontend
//	import base.api.core.v1 as core_v1
//
//	main = frontend.Server{}
func parseImport(code string) []string {
	var m = make(map[string]string)
	var longStrPrefix string
	for _, line := range strings.Split(code, "\n") {
		lineCode := strings.TrimSpace(line)
		// remove commented code
		if idx := strings.Index(lineCode, "#"); idx >= 0 {
			lineCode = strings.TrimSpace(lineCode[:idx])
		}
		if lineCode == "" {
			continue
		}

		// skip long string
		if longStrPrefix != "" {
			if strings.HasSuffix(lineCode, longStrPrefix) {
				longStrPrefix = "" // long string end
			}
			continue // skip
		} else {
			if strings.HasPrefix(lineCode, `"""`) {
				longStrPrefix = `"""`
				if strings.HasSuffix(lineCode[len(`"""`):], longStrPrefix) {
					longStrPrefix = "" // long string end
				}
				continue
			}
			if strings.HasPrefix(lineCode, `'''`) {
				longStrPrefix = `'''`
				if strings.HasSuffix(lineCode[len(`'''`):], longStrPrefix) {
					longStrPrefix = "" // long string end
				}
				continue
			}

			// skip short string
			if strings.HasPrefix(lineCode, `"`) {
				continue
			}
			if strings.HasPrefix(lineCode, `'`) {
				continue
			}
		}

		ss := strings.Fields(lineCode)

		if len(ss) > 0 {
			// 'import xx' must at the beginning
			if !strings.HasPrefix(ss[0], "import") {
				break
			}

			// get import path from line code: "abc" from "import abc" or "import abc as bcd"
			if len(ss) >= 2 {
				pkgPath := strings.Trim(ss[1], `'"`)
				m[pkgPath] = pkgPath
			}
		}
	}

	var import_list []string
	for pkgpath := range m {
		import_list = append(import_list, pkgpath)
	}
	sort.Strings(import_list)
	return import_list
}

// fixImportPath fixes an original importPath to a file path (or package path)
// the filepath is the file path that defines an import and the importPath is the path part of an import statement
// suppose the filepath is a.b.c.k, and the import path is:
// 1. an absolute import path d.e, the result will be: d/e
// 2. a relative import path ..d.e, the result will be: a/d/e
func fixImportPath(path, importPath string) string {
	if !strings.HasPrefix(importPath, ".") {
		return strings.Replace(importPath, ".", "/", -1)
	}
	filepath.Join()

	pkgpath := path
	if strings.HasSuffix(pkgpath, ".k") {
		pkgpath = pathpkg.Dir(pkgpath)
	}

	// count leading dot
	var dotCount = len(importPath)
	for i, c := range importPath {
		if c != '.' {
			dotCount = i
			break
		}
	}

	// get importPath without leading dots and use "/" as path seperator instead of "."
	importPath = importPath[dotCount:]
	importPath = strings.Replace(importPath, ".", "/", -1)

	// one leading dot means the importPath is in the same package path with current file
	if dotCount == 1 {
		importPath = pkgpath + "/" + importPath
		return importPath
	}

	var ss = strings.Split(pkgpath, "/")
	// if the relative path is invalid as it imports a path that's out of the program root, fix it to a path just under the program root
	if (dotCount - 1) < len(ss) {
		// for relative import path, fix it to an absolute path
		importParts := append(ss[:len(ss)-(dotCount-1)], importPath)
		return strings.Join(importParts, "/")
	} else {
		// Use the relative filepath with ".."
		for i := 0; i < dotCount-1; i++ {
			pkgpath = pkgpath + "/.."
		}
		return pkgpath + "/" + importPath
	}
}
