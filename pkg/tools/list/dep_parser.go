// Copyright The KCL Authors. All rights reserved.

package list

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	pathpkg "path"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"kcl-lang.io/kcl-go/pkg/3rdparty/toml"
)

var _ = fmt.Sprint

const KCL_MOD_PATH_ENV = "${KCL_MOD}"

const (
	Default_KclMod      = "kcl.mod"
	Default_KclYaml     = "kcl.yaml"
	Default_ProjectYaml = "project.yaml"
)

type Option struct {
	KclMod                 string // default: Default_kclMod
	KclYaml                string // default: Default_KclYaml
	ProjectYaml            string // default: Default_ProjectYaml
	FlagAll                bool
	UseAbsPath             bool
	ExcludeExternalPackage bool
	ExcludeBuiltin         bool
	IgnoreImportError      bool
}

func (p *Option) merge(other *Option) {
	if other.KclMod != "" {
		p.KclMod = other.KclMod
	}
	if other.KclYaml != "" {
		p.KclYaml = other.KclYaml
	}
	if other.ProjectYaml != "" {
		p.ProjectYaml = other.ProjectYaml
	}
	if other.FlagAll {
		p.FlagAll = true
	}
	if other.UseAbsPath {
		p.UseAbsPath = true
	}
	if other.ExcludeExternalPackage {
		p.ExcludeExternalPackage = true
	}
	if other.ExcludeBuiltin {
		p.ExcludeBuiltin = true
	}
	if other.IgnoreImportError {
		p.IgnoreImportError = true
	}
}

func (p *Option) adjust() {
	if p.KclMod == "" {
		p.KclMod = Default_KclMod
	}
	if p.KclYaml == "" {
		p.KclYaml = Default_KclYaml
	}
	if p.ProjectYaml == "" {
		p.ProjectYaml = Default_ProjectYaml
	}
}

type DepParser struct {
	opt Option
	vfs fs.FS

	projectYamlDirList []string
	kclYamlList        []string
	mainKList          []string
	kList              []string

	importMap   map[string][]string
	pkgFilesMap map[string][]string

	touchedFiles    []string
	touchedFilesDag map[string]color
	touchedApps     []string
	untouchedApps   []string

	err error
}

type color uint32

const (
	white color = iota
	black
	grey // must be > white and black
)

func (c color) String() string {
	switch c {
	case white:
		return "white"
	case black:
		return "black"
	default:
		return "grey"
	}
}

func NewDepParser(root string, opt ...Option) *DepParser {
	root = pathpkg.Clean(root)
	return NewDepParserWithFS(os.DirFS(root), opt...)
}

func NewDepParserWithFS(vfs fs.FS, opts ...Option) *DepParser {
	p := &DepParser{
		vfs:             vfs,
		importMap:       make(map[string][]string),
		pkgFilesMap:     make(map[string][]string),
		touchedFilesDag: make(map[string]color),
	}

	for _, opt := range opts {
		p.opt.merge(&opt)
	}
	p.opt.adjust()

	p.kList = p.getKList()
	p.mainKList = p.getMainKList()
	p.kclYamlList = p.getKclYamlList()
	p.projectYamlDirList = p.getProjectYamlDirList()

	for _, main_k := range p.mainKList {
		if err := p.loadImportMap(pathpkg.Dir(main_k), p.importMap); err != nil {
			p.err = err
			break
		}
	}
	for _, kcl_yaml := range p.kclYamlList {
		if err := p.loadImportMap(pathpkg.Dir(kcl_yaml), p.importMap); err != nil {
			p.err = err
			break
		}
	}

	return p
}

// GetError return parser error.
func (p *DepParser) GetError() error { return p.err }

func (p *DepParser) GetAppFiles(pkgpath string, includeDependFiles bool) []string {
	if !includeDependFiles {
		return p.pkgFilesMap[pkgpath]
	}

	var files []string
	for k := range p.scanAppAllFiles(pkgpath, nil) {
		files = append(files, k)
	}
	sort.Strings(files)

	return files
}

func (p *DepParser) scanAppAllFiles(pkgpath string, files map[string]string) map[string]string {
	if files == nil {
		files = make(map[string]string)
	}

	for _, s := range p.pkgFilesMap[pkgpath] {
		files[s] = s
	}

	for _, importPkg := range p.importMap[pkgpath] {
		p.scanAppAllFiles(importPkg, files)
	}

	return files
}

func (p *DepParser) GetAppPkgs(pkgpath string, includeDependFiles bool) []string {
	if !includeDependFiles {
		return p.importMap[pkgpath]
	}

	var pkgs []string
	for k := range p.scanAppAllPkgs(pkgpath, nil) {
		pkgs = append(pkgs, k)
	}
	sort.Strings(pkgs)

	return pkgs
}

func (p *DepParser) scanAppAllPkgs(pkgpath string, pkgs map[string]string) map[string]string {
	if pkgs == nil {
		pkgs = make(map[string]string)
	}

	for _, s := range p.importMap[pkgpath] {
		pkgs[s] = s
	}

	for _, importPkg := range p.importMap[pkgpath] {
		p.scanAppAllPkgs(importPkg, pkgs)
	}

	return pkgs
}

func (p *DepParser) GetTouchedApps(touchedFiles ...string) (touchedApps, untouchedApps []string) {
	if len(touchedFiles) == 0 {
		return nil, nil
	}
	if reflect.DeepEqual(p.touchedFiles, touchedFiles) {
		return p.touchedApps, p.untouchedApps
	}

	p.touchedFiles = touchedFiles
	p.touchedFilesDag = make(map[string]color)
	p.touchedApps = []string{}
	p.untouchedApps = []string{}

	// init grey color
	for _, s := range p.touchedFiles {
		p.touchedFilesDag[pathpkg.Dir(s)] = grey
		p.touchedFilesDag[strings.TrimSuffix(s, ".k")] = grey
	}

	// if dir/project.yaml exists, set the grey color
	for _, s := range p.touchedFiles {
		if projYamlDir := p.getProjectYamlDir(s); projYamlDir != "" {
			for _, s := range p.kList {
				if s == projYamlDir || strings.HasPrefix(s, projYamlDir+"/") {
					p.touchedFilesDag[pathpkg.Dir(s)] = grey
					p.touchedFilesDag[strings.TrimSuffix(s, ".k")] = grey
				}
			}
		}
	}

	for _, main_k := range p.mainKList {
		if app := pathpkg.Dir(main_k); p.checkPkgColor(app) != black {
			p.touchedApps = append(p.touchedApps, app)
		} else {
			p.untouchedApps = append(p.untouchedApps, app)
		}
	}
	return p.touchedApps, p.untouchedApps
}

func (p *DepParser) checkPkgColor(pkgpath string) color {
	if !strings.ContainsAny(pkgpath, `\/`) {
		return black
	}

	if isBuiltinPkg(pkgpath) || isPluginPkg(pkgpath) {
		return black
	}
	if p.opt.ExcludeExternalPackage {
		if isExternalPkg(p.vfs, pkgpath) {
			return black
		}
	}

	if c := p.touchedFilesDag[pkgpath]; c != white {
		return c
	}

	for _, s := range p.importMap[pkgpath] {
		if p.checkPkgColor(s) != black {
			p.touchedFilesDag[pkgpath] = grey
			return grey
		}
	}

	p.touchedFilesDag[pkgpath] = black
	return black
}

func (p *DepParser) IsApp(pkgpath string) bool {
	if fi, _ := fs.Stat(p.vfs, pkgpath+"/main.k"); fi != nil && !fi.IsDir() {
		return true
	}
	if fi, _ := fs.Stat(p.vfs, pathpkg.Join(pkgpath, p.opt.KclYaml)); fi != nil && !fi.IsDir() {
		return true
	}
	return false
}

func (p *DepParser) getProjectYamlDir(pkgpath string) string {
	for _, s := range p.projectYamlDirList {
		if pkgpath == s || strings.HasPrefix(pkgpath, s+"/") {
			return s
		}
	}
	return ""
}

func isExternalPkg(vfs fs.FS, pkgpath string) bool {
	// Read the kcl.mod file
	modFileContent, err := fs.ReadFile(vfs, Default_KclMod)
	if err != nil {
		return false
	}
	// Parse the TOML content
	var modFileData map[string]any
	if err := toml.Unmarshal([]byte(modFileContent), &modFileData); err != nil {
		return false
	}
	// Extract dependency information
	if deps, ok := modFileData["dependencies"].(map[string]any); ok {
		for dep := range deps {
			if strings.HasPrefix(pkgpath, dep) || strings.HasPrefix(pkgpath, strings.Replace(dep, "-", "_", -1)) {
				return true
			}
		}
	}
	return false
}

func (p *DepParser) GetDepPkgList(pkgpath string) []string {
	return p.importMap[pkgpath]
}

func (p *DepParser) GetPkgFileList(pkgpath string) []string {
	files, _ := loadKFileList(p.vfs, pkgpath, p.opt)
	return files
}

func (p *DepParser) GetMainKList() []string {
	return p.mainKList
}

func (p *DepParser) GetPkgList() []string {
	var ss []string
	for k := range p.importMap {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return ss
}

func (p *DepParser) GetKList() []string {
	return p.kList
}

func (p *DepParser) GetImportMap() map[string][]string {
	return p.importMap
}

func (p *DepParser) GetImportMapString() string {
	json, err := json.MarshalIndent(p.importMap, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (p *DepParser) getProjectYamlDirList() []string {
	var list []string
	fs.WalkDir(p.vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(path, ".git") {
			return nil
		}
		if strings.HasSuffix(path, "/"+p.opt.ProjectYaml) {
			list = append(list, strings.TrimSuffix(path, "/"+p.opt.ProjectYaml))
		}
		return nil
	})

	return list
}

func (p *DepParser) getKclYamlList() []string {
	var list []string
	fs.WalkDir(p.vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(path, ".git") {
			return nil
		}
		if strings.HasSuffix(path, "/"+p.opt.KclYaml) {
			list = append(list, path)
			return nil
		}
		return nil
	})

	return list
}

func (p *DepParser) getMainKList() []string {
	var list []string
	fs.WalkDir(p.vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(path, ".git") {
			return nil
		}
		if strings.HasSuffix(path, "/main.k") {
			list = append(list, path)
		}
		return nil
	})

	return list
}

func (p *DepParser) getKList() []string {
	var list []string
	fs.WalkDir(p.vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(path, ".git") {
			return nil
		}

		if !strings.HasSuffix(path, ".k") {
			return nil
		}

		// skip xxx_test.k
		if strings.HasSuffix(path, "_test.k") {
			return nil
		}

		// skip _xxx.k
		if strings.HasPrefix(pathpkg.Base(path), "_") {
			return nil
		}

		// OK
		list = append(list, path)
		return nil
	})

	return list
}

func (p *DepParser) loadImportMap(path string, import_map map[string][]string) error {
	if p.err != nil {
		return p.err
	}

	pkgpath := path
	if strings.HasSuffix(path, ".k") {
		pkgpath = pathpkg.Dir(path)
	}

	if isBuiltinPkg(pkgpath) || isPluginPkg(pkgpath) {
		return nil
	}

	if p.opt.ExcludeExternalPackage {
		if isExternalPkg(p.vfs, pkgpath) {
			return nil
		}
	}

	if _, ok := import_map[pkgpath]; ok {
		return nil
	}

	var k_files []string
	if x, ok := p.pkgFilesMap[pkgpath]; !ok {
		var err error
		k_files, err = loadKFileList(p.vfs, pkgpath, p.opt)
		if err != nil {
			return fmt.Errorf("package %s: %w", pkgpath, err)
		}
		p.pkgFilesMap[pkgpath] = k_files
	} else {
		k_files = x
	}

	for _, file := range k_files {
		src, err := fs.ReadFile(p.vfs, file)
		if err != nil {
			return fmt.Errorf("package %s: %w", pkgpath, err)
		}

	Loop:
		for _, import_path := range parseImport(string(src)) {
			import_path := fixImportPath(file, import_path)

			for _, s := range import_map[pkgpath] {
				if s == import_path {
					continue Loop
				}
			}

			import_map[pkgpath] = append(import_map[pkgpath], import_path)

			if err := p.loadImportMap(import_path, import_map); err != nil {
				return err
			}
		}
	}

	sort.Strings(import_map[pkgpath])
	return nil
}

func loadKFileList(vfs fs.FS, path string, opt Option) ([]string, error) {
	if strings.HasSuffix(path, ".k") {
		return []string{path}, nil
	}

	if fi, _ := fs.Stat(vfs, path+".k"); fi != nil && !fi.IsDir() {
		return []string{path + ".k"}, nil
	}

	kclYamlPath := pathpkg.Join(path, opt.KclYaml)
	if fi, _ := fs.Stat(vfs, kclYamlPath); fi != nil && !fi.IsDir() {
		// kcl_cli_configs:
		//   file:
		//     - ../../../../base/pkg/kusion_models/app_configuration/sofa/sofa_app_configuration_render.k
		//     - ../../../../base/pkg/kusion_models/app_configuration/metadata_render.k
		//     - ../../../../base/pkg/kusion_models/app_configuration/deploy_topology_render.k
		//     - ../base/base.k
		//     - main.k
		//     - ../../../../base/pkg/kusion_models/app_configuration/sofa/sofa_app_configuration_backend.k
		//   disable_none: true

		var settings struct {
			Config struct {
				Files []string `yaml:"file"`
			} `yaml:"kcl_cli_configs"`
		}

		data, err := fs.ReadFile(vfs, kclYamlPath)
		if err != nil {
			panic(fmt.Errorf("%s: %v", kclYamlPath, err))
		}
		if err := yaml.Unmarshal([]byte(data), &settings); err != nil {
			panic(fmt.Errorf("%s: %v", kclYamlPath, err))
		}
		for i, s := range settings.Config.Files {
			switch {
			case strings.HasPrefix(s, KCL_MOD_PATH_ENV):
				goldenPath := strings.Replace(s, KCL_MOD_PATH_ENV+"/", "/", -1)
				goldenPath = strings.Trim(goldenPath, "/")
				goldenPath = pathpkg.Clean(goldenPath)

				if _, err := fs.Stat(vfs, goldenPath); err != nil {
					panic(fmt.Errorf("%s: %v", kclYamlPath, err))
				}

				settings.Config.Files[i] = goldenPath

			default:
				goldenPath := pathpkg.Join(path, s)
				goldenPath = strings.Trim(goldenPath, "/")
				goldenPath = pathpkg.Clean(goldenPath)

				if _, err := fs.Stat(vfs, goldenPath); err != nil {
					panic(fmt.Errorf("%s: %v", kclYamlPath, err))
				}

				settings.Config.Files[i] = goldenPath
			}
		}

		if len(settings.Config.Files) > 0 || opt.IgnoreImportError {
			return settings.Config.Files, nil
		} else {
			return nil, fmt.Errorf("no kcl file")
		}
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

	if len(k_files) > 0 || opt.IgnoreImportError {
		return k_files, nil
	} else {
		return nil, fmt.Errorf("no kcl file")
	}
}
