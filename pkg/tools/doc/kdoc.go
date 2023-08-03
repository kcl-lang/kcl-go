package doc

import (
	"bytes"
	_ "embed"
	"fmt"
	kcl "kcl-lang.io/kcl-go"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/schema.gotmpl
var schemaDocTmpl string

var tmpl *template.Template

func init() {
	var err error
	// todo: change to nested template files
	tmpl, err = template.New("doc.md").Funcs(funcMap()).Parse(schemaDocTmpl)
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

func (g *GenContext) render(schemas map[string]*kcl.KclType) error {
	// make directory
	err := os.MkdirAll(g.Target, 0755)
	if err != nil {
		return fmt.Errorf("failed to create docs/ directory under the target directory: %s", err)
	}
	for _, schema := range schemas {
		// get doc file name
		fileName := fmt.Sprintf("%s.md", schema.SchemaName)
		// render doc content
		content, err := g.renderContent(schema)
		if err != nil {
			return err
		}
		// write content to file
		err = os.WriteFile(filepath.Join(g.Target, fileName), content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s in %s: %v", fileName, g.Target, err)
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
	}
}

func (g *GenContext) renderContent(schema *kcl.KclType) ([]byte, error) {
	var contentBuf bytes.Buffer
	err := tmpl.Execute(&contentBuf, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to render schema type %s with template", schema.SchemaName)
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

func (g *GenContext) GenDoc() error {
	typeMapping, err := kcl.GetSchemaTypeMapping(g.PackagePath, "", "")
	if err != nil {
		return fmt.Errorf("parse schema type from file failed: %s", err)
	}
	if len(typeMapping) == 0 {
		return fmt.Errorf("no schema found")
	}
	err = g.render(typeMapping)
	if err != nil {
		return fmt.Errorf("render doc failed: %s", err)
	}
	return nil
}
