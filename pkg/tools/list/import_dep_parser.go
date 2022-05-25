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
	opt DepOption
	vfs fs.FS

	depsGraph *DirectedAcyclicGraph
	parsed    map[string]bool
}

// FileVertex defines the file path in the dependency graph
type FileVertex struct {
	// the relative file path
	path string
}

// DepsEdge defines the dependency relation in the dependency graph. The target FileVertex depends on the source FileVertex,
// Put it another way, the content of the target FileVertex defines an import statement to the source FileVertex
type DepsEdge struct {
	source FileVertex
	target FileVertex
}

func (v *FileVertex) Id() string {
	return v.path
}

func (e *DepsEdge) Id() string {
	return fmt.Sprintf("%s->%s", e.source.Id(), e.target.Id())
}

func (e *DepsEdge) Source() Vertex {
	return &e.source
}

func (e *DepsEdge) Target() Vertex {
	return &e.target
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
		vfs:       vfs,
		depsGraph: NewDirectedAcyclicGraph(),
		opt:       opt,
		parsed:    make(map[string]bool),
	}
	for _, file := range opt.Files {
		p.Inspect(file)
	}
	return p
}

// Inspect will inspect current path: read the file content and parse import stmts, record the deps relation between the imported and importing.
// if path is a file, each file in the pkg dir containing the file will be parsed
// if path is a pkg path, each file in the pkg path will be parsed
func (p *ImportDepParser) Inspect(path string) {
	var kFiles []string
	pkgpath := path
	isKclFile := false
	if strings.HasSuffix(path, ".k") {
		pkgpath = pathpkg.Dir(path)
		isKclFile = true
	}
	pkgV := FileVertex{pkgpath}
	p.depsGraph.AddVertex(&pkgV)
	p.parsed[pkgpath] = true
	p.parsed[path] = true
	if isKclFile {
		fileV := FileVertex{path}
		p.depsGraph.AddVertex(&fileV)
		p.depsGraph.AddEdge(&DepsEdge{fileV, pkgV})
	}
	kFiles = listKFiles(p.vfs, pkgpath)

	for _, f := range kFiles {
		currentFileV := FileVertex{f}
		p.depsGraph.AddVertex(&currentFileV)
		p.depsGraph.AddEdge(&DepsEdge{currentFileV, pkgV})
		p.parsed[f] = true

		src, err := fs.ReadFile(p.vfs, f)
		if err != nil {
			panic(err)
		}
		for _, importPath := range parseImport(string(src)) {
			importPath = p.fixPath(fixImportPath(f, importPath))
			p.depsGraph.AddEdge(&DepsEdge{source: FileVertex{importPath}, target: currentFileV})
			if _, ok := p.parsed[importPath]; ok {
				continue
			}
			p.Inspect(importPath)
		}
	}
}

func (p *ImportDepParser) fixPath(path string) string {
	if strings.HasSuffix(path, ".k") {
		return path
	}
	if fi, _ := fs.Stat(p.vfs, path+".k"); fi != nil && !fi.IsDir() {
		return path + ".k"
	}
	return path
}

// ListDownStreamFiles return a list of downstream depend files from the given changed path list.
func (p *ImportDepParser) ListDownStreamFiles() []string {
	for _, path := range p.opt.ChangedPaths {
		if strings.HasSuffix(path, ".k") && !strings.HasSuffix(path, "_test.k") {
			if _, err := fs.Stat(p.vfs, path); errors.Is(err, os.ErrNotExist) {
				// changed KCL file (not test file) not exists, might be deleted
				pkgpath := pathpkg.Dir(path)
				importPath := strings.TrimSuffix(path, ".k")
				downStreamPaths := []string{pkgpath, importPath}
				for _, v := range downStreamPaths {
					if p.depsGraph.vertices.Contains(v) {
						currentFileV := FileVertex{path}
						p.depsGraph.AddVertex(&currentFileV)
						p.depsGraph.AddEdge(&DepsEdge{currentFileV, FileVertex{v}})
					}
				}
			}
		}
	}
	return p.GetRelatives(p.opt.ChangedPaths, true)
}

// ListUpstreamFiles return a list of upstream dependent files from the given path list.
func (p *ImportDepParser) ListUpstreamFiles() []string {
	return p.GetRelatives(p.opt.Files, false)
}

// GetRelatives returns a list of file paths that have import relation(directly or indirectly) with the focus paths.
// If the downStream is set to true, that means only the file paths that depend on the focus paths will be returned.
// Otherwise, only the file paths that the focus paths depend on will be returned.
func (p *ImportDepParser) GetRelatives(focusPaths []string, downStream bool) []string {
	visited := map[string]bool{}
	for _, focus := range focusPaths {
		visit(p.depsGraph, focus, visited, downStream)
	}
	relatives := make([]string, 0, len(visited))
	for affected := range visited {
		relatives = append(relatives, affected)
	}
	return relatives
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
