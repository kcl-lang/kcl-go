package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// GenKclModuleOptions configures GenKclModule. The embedded GenKclOptions
// is forced to ModeGoStruct/OneFile=false; any Mode/OneFile settings
// supplied by the caller are overridden silently so a single KCL module
// is always produced.
type GenKclModuleOptions struct {
	GenKclOptions
	// ModuleName populates the [package] name field in the emitted
	// kcl.mod. Defaults to the sanitised last segment of the Go module
	// path so the resulting file resolves as a fresh KCL package.
	ModuleName string
	// ModuleEdition populates the [package] edition field. Defaults to
	// "0.0.1", matching the canonical kcl.mod shape used elsewhere in
	// this repository.
	ModuleEdition string
	// ModuleVersion populates the [package] version field. Defaults to
	// "0.0.1" for the same reason.
	ModuleVersion string
}

// GenKclModule converts the Go package tree rooted at pkgDir into a
// resolvable multi-file KCL module written to outDir. Each non-main Go
// package that declares at least one exported struct (or any const, in
// OneFile=false mode) becomes a directory under outDir whose name
// matches the sanitised Go path, with a single `<name>.k` file inside
// holding the emitted schemas. Out-of-module Go references are
// degraded to `any` so the resulting module is self-contained.
//
// On failure, outDir may contain partial output (e.g. some `.k` files
// written before the failing package was reached); callers should
// remove the directory if a clean re-run is required.
func GenKclModule(outDir, pkgDir string, opts *GenKclModuleOptions) error {
	if opts == nil {
		opts = &GenKclModuleOptions{}
	}
	oneFile := false
	opts.GenKclOptions = GenKclOptions{
		Mode:                  opts.GenKclOptions.Mode,
		CastingOption:         opts.GenKclOptions.CastingOption,
		UseIntegersForNumbers: opts.GenKclOptions.UseIntegersForNumbers,
		OneFile:               &oneFile,
	}
	// Defensive normalisation: a caller who explicitly chose
	// ModeGoStruct keeps it; anything else is upgraded. We never want
	// a YAML/JSON/etc. mode to leak into a module-level emit.
	if opts.GenKclOptions.Mode != ModeGoStruct {
		opts.GenKclOptions.Mode = ModeGoStruct
	}
	pkgs, err := loadGoPackages(pkgDir + "/...")
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("no Go packages found under %s", pkgDir)
	}
	rootPkg := pkgs[0]
	if rootPkg.Module == nil {
		return fmt.Errorf("pkgDir %s is not inside a Go module", pkgDir)
	}
	moduleRoot := rootPkg.Module.Path
	moduleName := opts.ModuleName
	if moduleName == "" {
		moduleName = sanitizeKclPathSegment(lastSegment(moduleRoot, "/"))
	}
	moduleEdition := opts.ModuleEdition
	if moduleEdition == "" {
		moduleEdition = "0.0.1"
	}
	moduleVersion := opts.ModuleVersion
	if moduleVersion == "" {
		moduleVersion = "0.0.1"
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create outDir %s: %w", outDir, err)
	}

	// Generate one .k file per in-module, non-main package that has
	// something worth emitting (exported structs or consts).
	for _, p := range pkgs {
		if !strings.HasPrefix(p.PkgPath, moduleRoot) {
			// External (outside the caller's module): skip. typeName
			// will degrade any reference to it to `any` when
			// moduleRoot is set.
			continue
		}
		if p.PkgPath == "main" || strings.HasSuffix(p.PkgPath, "/main") || p.Name == "main" {
			continue
		}
		ctx := newModuleCtx(p, moduleRoot)
		ctx.populateFrom(p)
		if len(ctx.goStructs) == 0 && len(ctx.goConsts) == 0 {
			continue
		}
		results := ctx.convertFromGoStructs()
		kclSch, err := ctx.buildKclFile(results)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := genKclInto(&buf, kclSch); err != nil {
			return err
		}
		relPath := goPkgToKclImport(p.PkgPath)
		dir := filepath.Join(outDir, filepath.FromSlash(strings.ReplaceAll(relPath, ".", string(filepath.Separator))))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
		fileName := sanitizeKclPathSegment(p.Name)
		if fileName == "" || fileName[0] == '_' {
			fileName = "pkg"
		}
		outFile := filepath.Join(dir, fileName+".k")
		if err := os.WriteFile(outFile, buf.Bytes(), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", outFile, err)
		}
	}

	if err := writeKclMod(outDir, moduleName, moduleEdition, moduleVersion); err != nil {
		return err
	}
	return nil
}

// newModuleCtx builds a fresh genKclTypeContext with the configuration
// GenKclModule needs: ModeGoStruct, OneFile=false, and the supplied
// module root recorded so out-of-module Go types degrade to `any`.
func newModuleCtx(p *packages.Package, moduleRoot string) *genKclTypeContext {
	oneFile := false
	return &genKclTypeContext{
		pkgPath: p.PkgPath,
		context: context{
			resultMap: make(map[string]convertResult),
			imports:   make(map[string]struct{}),
			paths:     []string{},
		},
		goStructs:     map[*types.TypeName]goStruct{},
		goConsts:      map[string]global{},
		packages:      map[string]*packages.Package{},
		tyMapping:     map[types.Type]*ast.StructType{},
		tySpecMapping: map[string]string{},
		oneFile:       oneFile,
		moduleRoot:    moduleRoot,
		importAliases: map[string]string{},
	}
}

// genKclInto is a small indirection so we can route the rendered output
// through the existing (k *kclGenerator) genKcl renderer without having
// to instantiate a kclGenerator just to call the method.
func genKclInto(w *bytes.Buffer, s kclFile) error {
	g := newKclGenerator(&GenKclOptions{Mode: ModeGoStruct})
	return g.genKcl(w, s)
}

// writeKclMod writes a minimal `kcl.mod` file at the root of outDir.
// The format mirrors testdata/external/external_1/kcl.mod: a single
// [package] stanza with name/edition/version. No [dependencies] section
// is emitted because every Go reference that survives module-mode
// generation lives in a sibling directory under outDir, not in an
// upstream KCL registry.
func writeKclMod(outDir, name, edition, version string) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[package]\nname = %q\nedition = %q\nversion = %q\n", name, edition, version)
	return os.WriteFile(filepath.Join(outDir, "kcl.mod"), buf.Bytes(), 0o644)
}

// lastSegment returns the substring after the last occurrence of any
// rune in `seps`, or the whole input if no separator is present.
func lastSegment(s string, seps string) string {
	if i := strings.LastIndexAny(s, seps); i >= 0 {
		return s[i+1:]
	}
	return s
}
