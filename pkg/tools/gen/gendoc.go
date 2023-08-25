package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	kpm "kcl-lang.io/kpm/pkg/api"
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

var tmpl *template.Template

func init() {
	var err error
	tmpl = template.New("").Funcs(funcMap())
	_, err = tmpl.Parse(schemaDocTmpl)
	if err != nil {
		panic(err)
	}
	_, err = tmpl.Parse(packageDocTmpl)
	if err != nil {
		panic(err)
	}
}

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
		"sourcePath": func(tpe KclOpenAPIType) string {
			// todo: let users specify the source code base path
			return filepath.Join(tpe.GetSchemaPkgDir(""), tpe.KclExtensions.XKclModelType.Import.Alias)
		},
		"index": func(pkg *KclPackage) string {
			return pkg.getIndexContent(0, "  ", "")
		},
	}
}

func (pkg *KclPackage) getPackageIndexContent(level int, indentation string, pkgPath string) string {
	currentPkgPath := filepath.Join(pkgPath, pkg.Name)
	currentDocPath := filepath.Join(currentPkgPath, "index.md")
	return fmt.Sprintf(`%s- [%s](%s)
%s`, strings.Repeat(indentation, level), pkg.Name, currentDocPath, pkg.getIndexContent(level+1, indentation, currentPkgPath))
}

func (tpe *KclOpenAPIType) getSchemaIndexContent(level int, indentation string, pkgPath string) string {
	docPath := filepath.Join(pkgPath, "index.md")
	if level == 0 {
		docPath = ""
	}
	return fmt.Sprintf(`%s- [%s](%s#schema-%s)
`, strings.Repeat(indentation, level), tpe.KclExtensions.XKclModelType.Type, docPath, tpe.KclExtensions.XKclModelType.Type)
}

func (pkg *KclPackage) getIndexContent(level int, indentation string, pkgPath string) string {
	var content string
	if len(pkg.SchemaList) > 0 {
		for _, sch := range pkg.SchemaList {
			content += sch.getSchemaIndexContent(level, indentation, pkgPath)
		}
	}
	if len(pkg.SubPackageList) > 0 {
		for _, pkg := range pkg.SubPackageList {
			content += pkg.getPackageIndexContent(level, indentation, pkgPath)
		}
	}
	return content
}

func (g *GenContext) renderPackage(pkg *KclPackage, parentDir string) error {
	// render the package's index.md page
	//fmt.Println(fmt.Sprintf("creating %s/index.md", parentDir))
	indexFileName := fmt.Sprintf("%s.%s", "index", g.Format)
	var contentBuf bytes.Buffer
	err := tmpl.ExecuteTemplate(&contentBuf, "packageDoc", pkg)
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

func (g *GenContext) renderSchemaDocContent(schema *KclOpenAPIType) ([]byte, error) {
	var contentBuf bytes.Buffer
	err := tmpl.ExecuteTemplate(&contentBuf, "schemaDoc", schema)
	if err != nil {
		return nil, fmt.Errorf("failed to render schema type %s.%s with template, err: %s", schema.KclExtensions.XKclModelType.Import.Package, schema.KclExtensions.XKclModelType.Type, err)
	}
	return contentBuf.Bytes(), nil
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
	}
	g.Target = path.Join(g.Target, "docs")
	if _, err := os.Stat(g.Target); err == nil {
		// check and warn if the docs directory already exists
		fmt.Println(fmt.Sprintf("[Warn] path %s exists, all the content will be overwritten", g.Target))
		if err := os.RemoveAll(g.Target); err != nil {
			return nil, fmt.Errorf("failed to remove existing content in %s:%s", g.Target, err)
		}
	}
	return g, nil
}

// GenDoc generate document files from KCL source files
func (g *GenContext) GenDoc() error {
	pkg, err := kpm.GetKclPackage(g.PackagePath)
	if err != nil {
		return fmt.Errorf("filePath is not a KCL package: %s", err)
	}
	spec, err := g.getSwagger2Spec(pkg)
	err = g.render(spec)
	if err != nil {
		return fmt.Errorf("render doc failed: %s", err)
	}
	return nil
}

func (g *GenContext) getSwagger2Spec(pkg *kpm.KclPackage) (*SwaggerV2Spec, error) {
	spec := &SwaggerV2Spec{
		Swagger:     "2.0",
		Definitions: make(map[string]*KclOpenAPIType),
		Info: SpecInfo{
			Title:   pkg.GetPkgName(),
			Version: pkg.GetVersion(),
		},
	}
	pkgMapping, err := pkg.GetAllSchemaTypeMapping()
	if err != nil {
		return spec, err
	}
	// package path -> package
	for packagePath, p := range pkgMapping {
		// schema name -> schema type
		for _, t := range p {
			id := SchemaId(packagePath, t.KclType)
			spec.Definitions[id] = GetKclOpenAPIType(packagePath, t.KclType, false)
			fmt.Println(fmt.Sprintf("generate docs for schema %s", id))
		}
	}
	return spec, nil
}
