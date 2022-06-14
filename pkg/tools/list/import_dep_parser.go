// Copyright 2021 The KCL Authors. All rights reserved.

package list

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	pathpkg "path"
	"sort"
	"strings"
)

var _ = fmt.Sprint

// DepOption defines the option to parse dependency info
type DepOption struct {
	// Files is a list of relative file paths. The deps parser will parse the import information starting from these files
	Files []string
	// ChangedPaths is a list of relative file paths whose content are changed.
	// The deps parser can filter the Files by the ChangedPaths to distinguish the downstream Files of the ChangedPaths
	ChangedPaths []string
}

// ImportDepParser parses the import statements in KCL files within the given work directory
type ImportDepParser struct {
	// opt is the list dependency option
	opt DepOption
	// vfs is the file system with the KCL program root as root
	vfs fs.FS
	// importsGraph is a graph of the file dependent relation according to the import statement
	importsGraph *importGraph
}

// NewImportDepParser initialize an import dependency parser from the given pkg root and the deps option
func NewImportDepParser(root string, opt DepOption) (*ImportDepParser, error) {
	root = pathpkg.Clean(root)
	vfs := os.DirFS(root)
	for _, file := range opt.Files {
		if _, err := fs.Stat(vfs, file); err != nil {
			return nil, fmt.Errorf("invalid file path: %s", err)
		}
	}
	return NewImportDepParserWithFS(vfs, opt), nil
}

func NewImportDepParserWithFS(vfs fs.FS, opt DepOption) *ImportDepParser {
	p := &ImportDepParser{
		vfs:          vfs,
		importsGraph: newImportGraph(),
		opt:          opt,
	}
	for _, file := range opt.Files {
		p.Inspect(file)
	}
	return p
}

// importGraph implements an incremental
type importGraph struct {
	// the key is the file path and the value is the set of files/pkgs that the key file imports
	imports map[string]*stringSet
	// the key is the file/pkg path and the value is a set of files which import the key file/pkg
	importedBy map[string]*stringSet
	// pkgMap is the file to package map. the key is a KCL file path and the value is the package path of the file
	pkgMap map[string]string
	// pkgFiles is the pkg files map. the key is a KCL package path and the value is a set of file paths within the package
	pkgFiles map[string]*stringSet
}

// newImportGraph initiates an import graph
func newImportGraph() *importGraph {
	return &importGraph{
		imports:    make(map[string]*stringSet),
		importedBy: make(map[string]*stringSet),
		pkgMap:     make(map[string]string),
		pkgFiles:   make(map[string]*stringSet),
	}
}

// stringSet is a simple string set implementation by map
type stringSet struct {
	values map[string]bool
}

// emptyStringSet returns an empty stringSet
func emptyStringSet() *stringSet {
	return &stringSet{
		values: map[string]bool{},
	}
}

// add an item to the stringSet
func (s *stringSet) add(value string) {
	s.values[value] = true
}

// check if the stringSet contains some value
func (s *stringSet) contains(value string) bool {
	_, ok := s.values[value]
	return ok
}

// toStringSlice returns a string slice of the values in the stringSet
func (s *stringSet) toStringSlice() []string {
	var result []string
	for value := range s.values {
		result = append(result, value)
	}
	return result
}

// Inspect will inspect current path: read the file content and parse import stmts, record the deps relation between the imported and importing.
// if path is a file, each file in the pkg dir containing the file will be parsed
// if path is a pkg path, each file in the pkg path will be parsed
func (p *ImportDepParser) Inspect(path string) {
	var kFiles []string
	pkgpath := path
	if strings.HasSuffix(path, ".k") {
		pkgpath = pathpkg.Dir(path)
	}
	kFiles = listKFiles(p.vfs, pkgpath)
	for _, f := range kFiles {
		p.importsGraph.pkgMap[f] = pkgpath
		if _, ok := p.importsGraph.pkgFiles[pkgpath]; !ok {
			p.importsGraph.pkgFiles[pkgpath] = emptyStringSet()
		}
		p.importsGraph.pkgFiles[pkgpath].add(f)
		src, err := fs.ReadFile(p.vfs, f)
		if err != nil {
			panic(err)
		}
		for _, importPath := range parseImport(string(src)) {
			importPath = p.fixPath(fixImportPath(f, importPath))
			if _, ok := p.importsGraph.imports[f]; !ok {
				p.importsGraph.imports[f] = emptyStringSet()
			}
			if _, ok := p.importsGraph.importedBy[importPath]; !ok {
				p.importsGraph.importedBy[importPath] = emptyStringSet()
			}
			p.importsGraph.imports[f].add(importPath)
			p.importsGraph.importedBy[importPath].add(f)

			if _, ok := p.importsGraph.imports[importPath]; ok {
				continue
			}
			p.Inspect(importPath)
		}
	}
}

// ListDownStreamFiles return a list of downstream dependent files from the given changed path list.
func (p *ImportDepParser) ListDownStreamFiles() []string {
	for _, f := range p.opt.ChangedPaths {
		if strings.HasSuffix(f, ".k") && !strings.HasSuffix(f, "_test.k") {
			if _, err := fs.Stat(p.vfs, f); errors.Is(err, os.ErrNotExist) {
				// changed KCL file (not test file) not exists, might be deleted
				pkgpath := pathpkg.Dir(f)
				p.importsGraph.pkgMap[f] = pkgpath
				_, ok := p.importsGraph.pkgFiles[pkgpath]
				if !ok {
					p.importsGraph.pkgFiles[pkgpath] = emptyStringSet()
				}
				p.importsGraph.pkgFiles[pkgpath].add(f)
				modulePath := strings.TrimSuffix(f, ".k")
				p.opt.ChangedPaths = append(p.opt.ChangedPaths, modulePath)
			}
		}
	}
	downFiles := emptyStringSet()
	for _, f := range p.opt.ChangedPaths {
		downFiles.add(f)
	}
	for _, f := range p.opt.ChangedPaths {
		p.importsGraph.walkDownstream(f, func(filepath string) {
			downFiles.add(filepath)
		})
	}
	return downFiles.toStringSlice()
}

// ListUpstreamFiles return a list of upstream dependent files from the given path list.
func (p *ImportDepParser) ListUpstreamFiles() []string {
	upFiles := emptyStringSet()
	for _, f := range p.opt.Files {
		upFiles.add(f)
	}
	for _, f := range p.opt.Files {
		p.importsGraph.walkUpstream(f, func(filepath string) {
			upFiles.add(filepath)
		})
	}
	return upFiles.toStringSlice()
}

// walkUpstream walks the importGraph from the fromPath and up to the files that the fromPath imports
func (g *importGraph) walkUpstream(fromPath string, f func(filepath string)) {
	nexts := g.imports[fromPath]
	if nexts == nil {
		return
	}
	for next := range nexts.values {
		f(next)
		if fileSet, ok := g.pkgFiles[next]; ok {
			for file := range fileSet.values {
				f(file)
				g.walkUpstream(file, f)
			}
		} else {
			g.walkUpstream(next, f)
		}
	}
}

// walkDownstream walks the importGraph from the fromPath and down to the files which import the fromPath
func (g *importGraph) walkDownstream(fromPath string, f func(filepath string)) {
	nexts := g.importedBy[fromPath]
	if nexts == nil {
		nexts = emptyStringSet()
	}
	pkg, ok := g.pkgMap[fromPath]
	if ok {
		nexts.add(pkg)
	}
	for next := range nexts.values {
		f(next)
		g.walkDownstream(next, f)
	}
}

// fixPath fixes a path (import path or file path) to a file path
// a/b/c.k -> a/b/c.k
// if a/b/c.k exists and is a file: a/b/c -> a/b/c.k
// if a/b/c.k not exists or is a dir: a/b/c -> a/b/c
func (p *ImportDepParser) fixPath(path string) string {
	if strings.HasSuffix(path, ".k") {
		return path
	}
	if fi, _ := fs.Stat(p.vfs, path+".k"); fi != nil && !fi.IsDir() {
		return path + ".k"
	}
	return path
}

// listKFiles returns a list of KCL file paths under the given package path or by the given file path. It will return an empty list if no KCL files found
func listKFiles(vfs fs.FS, path string) []string {
	if strings.HasSuffix(path, ".k") {
		return []string{path}
	}
	if fi, _ := fs.Stat(vfs, path+".k"); fi != nil && !fi.IsDir() {
		return []string{path + ".k"}
	}
	var k_files []string
	entryList, _ := fs.ReadDir(vfs, path)
	for _, info := range entryList {
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), ".k") {
			continue
		}
		// skip _xxx.k
		if strings.HasPrefix(info.Name(), "_") {
			continue
		}
		// skip xxx_test.k
		if strings.HasSuffix(info.Name(), "_test.k") {
			continue
		}
		// OK
		k_files = append(k_files, pathpkg.Join(path, info.Name()))
	}
	return k_files
}

// parseImport parses the import statements within the code and returns the import paths in it
func parseImport(code string) []string {
	var m = make(map[string]string)
	var longStrPrefix string
	for _, line := range strings.Split(code, "\n") {
		lineCode := strings.TrimSpace(line)
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
		if len(ss) == 0 {
			continue
		}

		// 'import xx' must at the begin
		if !strings.HasPrefix(ss[0], "import") {
			break
		}

		// import abc
		// import abc as bcd
		if len(ss) >= 0 {
			pkgpath := strings.Trim(ss[1], `'"`)
			m[pkgpath] = pkgpath
		}
	}

	var import_list []string
	for pkgpath := range m {
		import_list = append(import_list, pkgpath)
	}
	sort.Strings(import_list)
	return import_list
}

// fixImportPath fixes the original import_path by the filepath that defines the import and returns a filepath (or package path) of the file (or package) to import
func fixImportPath(filepath, import_path string) string {
	if !strings.HasPrefix(import_path, ".") {
		return strings.Replace(import_path, ".", "/", -1)
	}

	pkgpath := filepath
	if strings.HasSuffix(pkgpath, ".k") {
		pkgpath = pathpkg.Dir(pkgpath)
	}

	// count leading dot
	var dotCount = len(import_path)
	for i, c := range import_path {
		if c != '.' {
			dotCount = i
			break
		}
	}
	import_path = import_path[dotCount:]
	import_path = strings.Replace(import_path, ".", "/", -1)

	// import .metadata
	if dotCount == 1 {
		import_path = pkgpath + "/" + import_path
		return strings.Replace(import_path, ".", "/", -1)
	}

	var ss = strings.Split(pkgpath, "/")
	if (dotCount - 1) > len(ss) {
		dotCount = len(ss) + 1
	}

	import_parts := append(ss[:len(ss)-(dotCount-1)], import_path)
	return strings.Join(import_parts, "/")
}
