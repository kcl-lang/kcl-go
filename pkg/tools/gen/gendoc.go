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
	Name    string
	Version string `toml:"version,omitempty"` // kcl package version
	Schemas []*KclOpenAPIType
}

func (g *GenContext) render(spec *SwaggerV2Spec) error {
	// make directory
	err := os.MkdirAll(g.Target, 0755)
	if err != nil {
		return fmt.Errorf("failed to create docs/ directory under the target directory: %s", err)
	}

	// collect all the packages and schema list that they contain
	pkgs := make(map[string]*KclPackage)

	for _, schema := range spec.Definitions {
		pkgName := schema.KclExtensions.XKclModelType.Import.Package
		if _, ok := pkgs[pkgName]; ok {
			pkgs[pkgName].Schemas = append(pkgs[pkgName].Schemas, schema)
		} else {
			pkgs[pkgName] = &KclPackage{
				Name: pkgName,
			}
			pkgs[pkgName].Schemas = []*KclOpenAPIType{schema}
		}
	}

	err = g.renderPackage(pkgs)
	if err != nil {
		return err
	}

	for _, schema := range spec.Definitions {
		// create package directory if not exist
		pkgDir := schema.GetSchemaPkgDir(g.Target)
		err := os.MkdirAll(pkgDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create docs/%s directory under the target directory: %s", pkgDir, err)
		}
		// get doc file name
		fileName := fmt.Sprintf("%s.%s", schema.KclExtensions.XKclModelType.Type, g.Format)
		// render doc content
		content, err := g.renderSchemaDocContent(schema)
		if err != nil {
			return err
		}
		// write content to file
		err = os.WriteFile(filepath.Join(pkgDir, fileName), content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s in %s: %v", fileName, pkgDir, err)
		}
	}
	return nil
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
		"sortSchemas": func(schemas []*KclOpenAPIType) []*KclOpenAPIType {
			sort.Slice(schemas, func(i, j int) bool {
				return schemas[i].KclExtensions.XKclModelType.Type < schemas[j].KclExtensions.XKclModelType.Type
			})
			return schemas
		},
	}
}

func (g *GenContext) renderPackage(pkgs map[string]*KclPackage) error {
	for name, pkg := range pkgs {
		// create the package directory
		pkgDir := GetPkgDir(g.Target, name)
		err := os.MkdirAll(pkgDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create docs/%s directory under the target directory: %s", pkgDir, err)
		}
		indexFileName := fmt.Sprintf("%s.%s", "index", g.Format)
		// render index doc content
		var contentBuf bytes.Buffer
		err = tmpl.ExecuteTemplate(&contentBuf, "packageDoc", pkg)
		if err != nil {
			return fmt.Errorf("failed to render package %s with template, err: %s", name, err)
		}
		if err != nil {
			return err
		}
		// write content to file
		err = os.WriteFile(filepath.Join(pkgDir, indexFileName), contentBuf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s in %s: %v", indexFileName, pkgDir, err)
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
	//todo: deal err
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
