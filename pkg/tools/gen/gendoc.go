package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	htmlTmpl "html/template"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed templates/doc/schemaDoc.gotmpl
var schemaDocTmpl string

//go:embed templates/doc/packageDoc.gotmpl
var packageDocTmpl string

const (
	schemaDocTmplFile  = "schemaDoc.gotmpl"
	packageDocTmplFile = "packageDoc.gotmpl"
)

// GenContext defines the context during the generation
type GenContext struct {
	// PackagePath is the package path to the package or module to generate docs for
	PackagePath string
	// Format is the doc format to output
	Format Format
	// Target is the target directory to output the docs
	Target string
	// IgnoreDeprecated defines whether to generate documentation for deprecated schemas
	IgnoreDeprecated bool
	// EscapeHtml defines whether to escape html symbols when the output format is markdown
	EscapeHtml bool
	// SchemaDocTmpl defines the content of the schemaDoc template
	SchemaDocTmpl string
	// PackageDocTmpl defines the content of the packageDoc template
	PackageDocTmpl string
	// Template is the doc render template
	Template *template.Template
}

// GenOpts is the user interface defines the doc generate options
type GenOpts struct {
	// Path is the path to the directory or file to generate docs for
	Path string
	// Format is the doc format to output
	Format string
	// Target is the target directory to output the docs
	Target string
	// IgnoreDeprecated defines whether to generate documentation for deprecated schemas
	IgnoreDeprecated bool
	// EscapeHtml defines whether to escape html symbols when the output format is markdown
	EscapeHtml bool
	// TemplateDir defines the relative path from the package root to the template directory
	TemplateDir string
}

type Format string

const (
	Html     Format = "html"
	Markdown Format = "md"
)

// KclPackage contains package information of package metadata(such as name, version, description, ...) and exported models(such as schemas)
type KclPackage struct {
	Name              string `json:"name,omitempty"`        // kcl package name
	Version           string `json:"version,omitempty"`     // kcl package version
	Description       string `json:"description,omitempty"` // summary of the kcl package
	schemaMapping     map[string]*KclOpenAPIType
	subPackageMapping map[string]*KclPackage
	SchemaList        []*KclOpenAPIType `json:"schemaList,omitempty"`     // the schema list sorted by name in the KCL package
	SubPackageList    []*KclPackage     `json:"subPackageList,omitempty"` // the sub package list sorted by name in the KCL package
}

func (g *GenContext) render(spec *SwaggerV2Spec) error {
	// make directory
	err := os.MkdirAll(g.Target, 0755)
	if err != nil {
		return fmt.Errorf("failed to create docs/ directory under the target directory: %s", err)
	}
	// extract kcl package from swaggerV2 spec
	rootPkg := spec.toKclPackage()
	// sort schemas and subpackages by their names
	rootPkg.sortSchemasAndPkgs()
	// render the package
	err = g.renderPackage(rootPkg, g.Target)
	if err != nil {
		return err
	}
	return nil
}

// toKclPackage extracts a kcl package and sub packages, schemas from a SwaggerV2 spec
func (spec SwaggerV2Spec) toKclPackage() *KclPackage {
	rootPkg := &KclPackage{
		Name:        spec.Info.Title,
		Version:     spec.Info.Version,
		Description: spec.Info.Description,
	}

	for schemaName, schema := range spec.Definitions {
		pkgName := schema.KclExtensions.XKclModelType.Import.Package
		if pkgName == "" {
			addOrCreateSchema(rootPkg, schemaName, schema)
			continue
		}
		parentPkg := rootPkg
		subs := strings.Split(pkgName, ".")
		for _, sub := range subs {
			if parentPkg.subPackageMapping == nil {
				parentPkg.subPackageMapping = map[string]*KclPackage{}
			}
			if _, ok := parentPkg.subPackageMapping[sub]; !ok {
				parentPkg.subPackageMapping[sub] = &KclPackage{
					Name: sub,
				}
			}
			parentPkg = parentPkg.subPackageMapping[sub]
		}

		addOrCreateSchema(parentPkg, schemaName, schema)
	}
	return rootPkg
}

func (pkg *KclPackage) sortSchemasAndPkgs() {
	pkg.SubPackageList = sortMapToSlice(pkg.subPackageMapping)
	pkg.SchemaList = sortMapToSlice(pkg.schemaMapping)
	for _, sub := range pkg.SubPackageList {
		sub.sortSchemasAndPkgs()
	}
}

func sortMapToSlice[T any](mapping map[string]T) []T {
	keys := make([]string, 0, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sorted := make([]T, 0, len(mapping))
	for _, k := range keys {
		sorted = append(sorted, mapping[k])
	}
	return sorted
}

func addOrCreateSchema(pkg *KclPackage, schemaName string, schema *KclOpenAPIType) {
	if pkg.schemaMapping == nil {
		pkg.schemaMapping = map[string]*KclOpenAPIType{schemaName: schema}
	} else {
		pkg.schemaMapping[schemaName] = schema
	}
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"containsString": func(list []string, elem string) bool {
			for _, s := range list {
				if s == elem {
					return true
				}
			}
			return false
		},
		"kclType": func(tpe KclOpenAPIType) string {
			return tpe.GetKclTypeName(false)
		},
		"fullTypeName": func(tpe KclOpenAPIType) string {
			if tpe.KclExtensions.XKclModelType.Import.Package != "" {
				return fmt.Sprintf("%s.%s", tpe.KclExtensions.XKclModelType.Import.Package, tpe.KclExtensions.XKclModelType.Type)
			}
			return tpe.KclExtensions.XKclModelType.Type
		},
		"escapeHtml": func(original string, escapeHtml bool) string {
			// escape html symbols if needed
			if escapeHtml {
				return htmlTmpl.HTMLEscapeString(original)
			}
			return original
		},
		"arr": func(els ...any) []any {
			return els
		},
		"sourcePath": func(tpe KclOpenAPIType) string {
			// todo: let users specify the source code base path
			return filepath.Join(tpe.GetSchemaPkgDir(""), tpe.KclExtensions.XKclModelType.Import.Alias)
		},
		"indexContent": func(pkg *KclPackage) string {
			return pkg.getIndexContent(0, "  ", "", false)
		},
		"indexContentIgnoreDirPath": func(pkg *KclPackage) string {
			return pkg.getIndexContent(0, "  ", "", true)
		},
	}
}

func (pkg *KclPackage) getPackageIndexContent(level int, indentation string, pkgPath string, ignoreDir bool) string {
	currentPkgPath := filepath.Join(pkgPath, pkg.Name)
	currentDocPath := pkg.Name
	if !ignoreDir {
		// get the full directory path
		currentDocPath = filepath.Join(currentPkgPath, fmt.Sprintf("%s.md", pkg.Name))
	}
	return fmt.Sprintf(`%s- [%s](%s)
%s`, strings.Repeat(indentation, level), pkg.Name, currentDocPath, pkg.getIndexContent(level+1, indentation, currentPkgPath, ignoreDir))
}

func (tpe *KclOpenAPIType) getSchemaIndexContent(level int, indentation string, pkgPath string, pkgName string, ignoreDir bool) string {
	docPath := pkgName
	if !ignoreDir {
		// get the full directory path
		docPath = filepath.Join(pkgPath, fmt.Sprintf("%s.md", pkgName))
	}
	if level == 0 {
		// the schema is defined in current package
		docPath = ""
	}

	return fmt.Sprintf(`%s- [%s](%s#%s)
`, strings.Repeat(indentation, level), tpe.KclExtensions.XKclModelType.Type, docPath, strings.ToLower(tpe.KclExtensions.XKclModelType.Type))
}

func (pkg *KclPackage) getIndexContent(level int, indentation string, pkgPath string, ignoreDir bool) string {
	var content string
	if len(pkg.SchemaList) > 0 {
		for _, sch := range pkg.SchemaList {
			content += sch.getSchemaIndexContent(level, indentation, pkgPath, pkg.Name, ignoreDir)
		}
	}
	if len(pkg.SubPackageList) > 0 {
		for _, pkg := range pkg.SubPackageList {
			content += pkg.getPackageIndexContent(level, indentation, pkgPath, ignoreDir)
		}
	}
	return content
}

func (g *GenContext) renderPackage(pkg *KclPackage, parentDir string) error {
	pkgName := pkg.Name
	if pkg.Name == "" {
		pkgName = "main"
	}
	fmt.Println(fmt.Sprintf("generating doc for package %s", pkgName))
	indexFileName := fmt.Sprintf("%s.%s", pkgName, g.Format)
	var contentBuf bytes.Buffer
	err := g.Template.ExecuteTemplate(&contentBuf, "packageDoc", struct {
		EscapeHtml bool
		Data       *KclPackage
	}{
		EscapeHtml: g.EscapeHtml,
		Data:       pkg,
	})
	if err != nil {
		return fmt.Errorf("failed to render package %s with template, err: %s", pkg.Name, err)
	}
	// write content to file
	err = os.WriteFile(filepath.Join(parentDir, indexFileName), contentBuf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s in %s: %v", indexFileName, parentDir, err)
	}

	for _, sub := range pkg.SubPackageList {
		pkgDir := GetPkgDir(parentDir, sub.Name)
		//fmt.Println(fmt.Sprintf("creating directory: %s", pkgDir))
		err := os.MkdirAll(pkgDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create docs/%s directory under the target directory: %s", pkgDir, err)
		}
		err = g.renderPackage(sub, pkgDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (opts *GenOpts) ValidateComplete() (*GenContext, error) {
	g := &GenContext{}
	// --- format ---
	switch strings.ToLower(opts.Format) {
	case string(Markdown):
		g.Format = Markdown
		break
	case string(Html):
		g.Format = Html
		break
	default:
		return nil, fmt.Errorf("invalid generate format. Allow values: %s", []Format{Markdown, Html})
	}

	// --- package path ---
	absPath, err := filepath.Abs(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("invalid file path(%s) to generate document from, can not get the absolute path: %s", opts.Path, err)
	}
	_, err = os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path(%s) to generate document from, path not exists: %s", opts.Path, err)
	}
	g.PackagePath = absPath

	// --- template directory ---
	g.SchemaDocTmpl = schemaDocTmpl
	g.PackageDocTmpl = packageDocTmpl
	if opts.TemplateDir != "" {
		tmplAbsPath := filepath.Join(g.PackagePath, opts.TemplateDir)
		templatesDirInfo, err := os.Stat(tmplAbsPath)
		if err != nil {
			return nil, fmt.Errorf("invalid template directory path: %s. error: %s", opts.TemplateDir, err.Error())
		}
		if !templatesDirInfo.IsDir() {
			return nil, fmt.Errorf("template path is not a directory: %s", opts.TemplateDir)
		}
		err = filepath.Walk(tmplAbsPath, func(path string, info os.FileInfo, _ error) error {
			if info.IsDir() {
				// skip directories
				return nil
			}
			rel, err := filepath.Rel(tmplAbsPath, path)
			if err != nil {
				return err
			}
			switch rel {
			case schemaDocTmplFile:
				// use custom schema Doc Template file
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				g.SchemaDocTmpl = string(content)
				return nil
			case packageDocTmplFile:
				// use custom package Doc Template file
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				g.PackageDocTmpl = string(content)
				return nil
			default:
				return fmt.Errorf("unexpected template file: %s", path)
			}
		})
		if err != nil {
			return nil, err
		}
	}
	// parse template
	g.Template = template.New("").Funcs(funcMap())
	_, err = g.Template.Parse(g.SchemaDocTmpl)
	if err != nil {
		return nil, err
	}
	_, err = g.Template.Parse(g.PackageDocTmpl)
	if err != nil {
		return nil, err
	}

	// --- target ---
	if opts.Target == "" {
		// complete target output directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get default target directory: %s", err)
		}
		g.Target = cwd
	} else {
		// check if the target output directory is a valid directory path
		file, err := os.Stat(opts.Target)
		if err != nil {
			return nil, fmt.Errorf("invalid target directory(%s) to output the doc files, path not exists: %s", opts.Target, err)
		}
		if !file.IsDir() {
			return nil, fmt.Errorf("invalid target directory(%s) to output the doc files: not a directory", opts.Target)
		}
		g.Target = opts.Target
	}
	g.Target = path.Join(g.Target, "docs")
	if _, err := os.Stat(g.Target); err == nil {
		// check and warn if the docs directory already exists
		fmt.Println(fmt.Sprintf("[Warn] path %s exists, all the content will be overwritten", g.Target))
		if err := os.RemoveAll(g.Target); err != nil {
			return nil, fmt.Errorf("failed to remove existing content in %s:%s", g.Target, err)
		}
	}
	g.EscapeHtml = opts.EscapeHtml
	return g, nil
}

// GenDoc generate document files from KCL source files
func (g *GenContext) GenDoc() error {
	spec, err := KclPackageToSwaggerV2Spec(g.PackagePath)
	if err != nil {
		return err
	}
	err = g.render(spec)
	if err != nil {
		return fmt.Errorf("render doc failed: %s", err)
	}
	return nil
}
