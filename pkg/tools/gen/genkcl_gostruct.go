package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/packages"
)

type goStruct struct {
	pkgPath   string
	pkgName   string
	name      string
	fields    []field
	doc       string
	fieldDocs map[string]string
}

type field struct {
	name      string
	ty        types.Type
	tag       string
	anonymous bool // true for embedded fields
}

type genKclTypeContext struct {
	context
	// Go package path.
	pkgPath string
	// mainPkgPath is the Go import path of the package that contains the
	// structs we are generating KCL for. It is used to distinguish "main
	// package" structs from structs that are merely referenced (i.e. come
	// from imported Go packages).
	mainPkgPath string
	// Go structs in all package path
	goStructs map[*types.TypeName]goStruct
	// goConsts holds every statically-evaluable Go `const` declaration
	// from the main Go package, keyed by `<pkg_path>.<name>` so it
	// doesn't collide with struct spec docs which use the same key shape.
	goConsts map[string]global
	// All pkg path -> package mapping
	packages map[string]*packages.Package
	// Semantic type -> AST struct type mapping
	tyMapping map[types.Type]*ast.StructType
	// Semantic type -> AST struct type mapping
	tySpecMapping map[string]string
	// Generate all go structs into one KCL file.
	oneFile bool
	// moduleRoot, when non-empty, is the Go module path that bounds
	// `GenKclModule`. External Go packages whose path is outside
	// `moduleRoot` are degraded to `any` rather than emitted as `import`
	// statements. An empty value preserves the historical
	// `GenKcl(OneFile=false)` behaviour of emitting an `import` and
	// letting the caller wire the alias up to a real KCL module.
	moduleRoot string
	// importAliases maps the sanitised dotted Go path to the KCL
	// import alias we materialise for it. When two distinct Go paths
	// sanitise to the same preferred alias (e.g. `internal/a/types`
	// and `internal/b/types` both want `Types`), the second caller
	// receives `Types2`, the third `Types3`, etc. The map is populated
	// lazily by `aliasFor`; consult it whenever you need the alias
	// for a path so the alias-qualified type names and the `import ...
	// as` lines stay in sync.
	importAliases map[string]string
}

func (k *kclGenerator) genSchemaFromGoStruct(w io.Writer, filename string, _ any) error {
	// Default to OneFile=true to preserve the historical inlining behaviour
	// when the caller passes a nil/zero-valued GenKclOptions.
	oneFile := true
	if k.opts != nil && k.opts.OneFile != nil {
		oneFile = *k.opts.OneFile
	}
	ctx := genKclTypeContext{
		pkgPath: filename,
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
		importAliases: map[string]string{},
	}
	results, err := ctx.convertSchemaFromGoPackage()
	if err != nil {
		return err
	}
	kclSch, err := ctx.buildKclFile(results)
	if err != nil {
		return err
	}
	return k.genKcl(w, kclSch)
}

// buildKclFile assembles the in-memory representation of a single
// generated `.k` file from the populated context. It is shared between
// the single-package GenKcl path and the multi-file GenKclModule path
// so both emit identical imports/globals/schemas wiring.
func (ctx *genKclTypeContext) buildKclFile(results []convertResult) (kclFile, error) {
	kclSch := kclFile{
		Schemas: []schema{},
	}
	for _, result := range results {
		if result.IsSchema {
			kclSch.Schemas = append(kclSch.Schemas, result.schema)
		}
	}
	// Emit `import` statements for every external Go package referenced by
	// the generated schemas. Only meaningful when OneFile == false; in
	// OneFile mode `ctx.imports` stays empty because cross-package types
	// are inlined rather than aliased.
	if !ctx.oneFile {
		for _, imp := range getSortedKeys(ctx.imports) {
			kclSch.Imports = append(kclSch.Imports, kImport{
				PkgPath: imp,
				Alias:   ctx.aliasFor(imp),
			})
		}
	}
	// Emit top-level KCL global variables for every Go `const` we
	// recorded from the main package. Consts are emitted in both OneFile
	// modes because they are independent declarations, not cross-package
	// references.
	if len(ctx.goConsts) > 0 {
		for _, name := range getSortedKeys(ctx.goConsts) {
			kclSch.Globals = append(kclSch.Globals, ctx.goConsts[name])
		}
	}
	return kclSch, nil
}

func (ctx *genKclTypeContext) typeName(pkgPath, defName, fieldName string, typ types.Type) typeInterface {
	switch ty := typ.(type) {
	case *types.Basic:
		switch ty.Kind() {
		case types.Bool, types.UntypedBool:
			return typePrimitive(typBool)
		case types.Int,
			types.Int8,
			types.Int16,
			types.Int32,
			types.Int64,
			types.Uint,
			types.Uint8,
			types.Uint16,
			types.Uint32,
			types.Uint64,
			types.Uintptr,
			types.UnsafePointer,
			types.UntypedInt:
			return typePrimitive(typInt)
		case types.Float32,
			types.Float64,
			types.Complex64,
			types.Complex128,
			types.UntypedFloat,
			types.UntypedComplex:
			return typePrimitive(typFloat)
		case types.String, types.UntypedString, types.UntypedRune:
			return typePrimitive(typStr)
		default:
			return typePrimitive(typAny)
		}
	case *types.Pointer:
		return ctx.typeName(pkgPath, defName, fieldName, ty.Elem())
	case *types.Named:
		obj := ty.Obj()
		if obj != nil {
			pkg := obj.Pkg()
			switch {
			case pkg != nil && pkg.Path() == "time" && obj.Name() == "Time":
				return typePrimitive(typStr)
			case pkg != nil && pkg.Path() == "time" && obj.Name() == "Duration":
				return typePrimitive(typInt)
			case pkg != nil && pkg.Path() == "math/big" && obj.Name() == "Int":
				return typePrimitive(typInt)
			default:
				// Struct from external package in the Go module
				if _, ok := ctx.goStructs[obj]; !ok {
					if pkg != nil {
						// Record external package type information
						pkgPath := pkg.Path()
						if ctx.oneFile {
							ty := ctx.typeName(pkgPath, strcase.ToCamel(pkg.Name()), obj.Name(), ty.Underlying())
							return ty
						} else if ctx.moduleRoot != "" && !strings.HasPrefix(pkgPath, ctx.moduleRoot) {
							// Module mode: external (out-of-module) Go
							// packages are degraded to `any` and no `import`
							// is emitted, because the generated KCL module
							// cannot resolve them.
							return typePrimitive(typAny)
						} else {
							// Emit an `import` for this external Go package
							// and reference the type via its alias-qualified
							// KCL name, e.g. `Bar.Outer` for a struct
							// `Outer` declared in the Go package whose
							// last path segment is `bar`.
							imp := goPkgToKclImport(pkgPath)
							ctx.imports[imp] = struct{}{}
							return typeCustom{
								Name: ctx.aliasFor(imp) + "." + obj.Name(),
							}
						}
					} else {
						ty := ctx.typeName(pkgPath, defName, obj.Name(), ty.Underlying())
						return ty
					}
				} else {
					// Struct from current package
					return typeCustom{
						Name: obj.Name(),
					}
				}
			}
		} else {
			return typePrimitive(typAny)
		}
	case *types.Array:
		return typeArray{
			Items: ctx.typeName(pkgPath, defName, fieldName, ty.Elem()),
		}
	case *types.Slice:
		return typeArray{
			Items: ctx.typeName(pkgPath, defName, fieldName, ty.Elem()),
		}
	case *types.Map:
		return typeDict{
			Key:   ctx.typeName(pkgPath, defName, fieldName, ty.Key()),
			Value: ctx.typeName(pkgPath, defName, fieldName, ty.Elem()),
		}
	case *types.Struct:
		schemaName := fmt.Sprintf("%s%s", defName, strcase.ToCamel(fieldName))
		if _, ok := ctx.resultMap[schemaName]; !ok {
			result := convertResult{IsSchema: true}
			ctx.resultMap[schemaName] = result
			description := ""
			if doc, ok := ctx.tySpecMapping[pkgPath+"."+fieldName]; ok {
				description = doc
			}
			result.schema.Description = description
			result.schema.Name = schemaName
			fields, fieldDocs := ctx.getStructFieldsAndDocs(typ)
			for _, field := range fields {
				typeName := ctx.typeName(pkgPath, schemaName, field.name, field.ty)
				fieldName := formatName(field.name)
				fieldDoc := ""
				if doc, ok := fieldDocs[fieldName]; ok {
					fieldDoc = doc
				}
				// Use alias name and type defined in the `kcl` or `json`` tag
				tagName, tagTy, _, err := parserGoStructFieldTag(field.tag)
				if err == nil {
					if tagName != "" {
						fieldName = tagName
					}
					if tagTy != nil {
						typeName = tagTy
					}
				}
				result.schema.Properties = append(result.Properties, property{
					Name:        fieldName,
					Type:        typeName,
					Description: fieldDoc,
				})
			}
			ctx.resultMap[schemaName] = result
		}
		return typeCustom{
			Name: schemaName,
		}
	case *types.Union:
		var types []typeInterface
		for i := 0; i < ty.Len(); i++ {
			types = append(types, ctx.typeName(pkgPath, defName, fieldName, ty.Term(i).Type()))
		}
		return typeUnion{
			Items: types,
		}
	case *types.Interface:
		if !ty.IsComparable() {
			return typePrimitive(typAny)
		}
		var types []typeInterface
		for i := 0; i < ty.NumEmbeddeds(); i++ {
			types = append(types, ctx.typeName(pkgPath, defName, fieldName, ty.EmbeddedType(i)))
		}
		return typeUnion{
			Items: types,
		}
	default:
		return typePrimitive(typAny)
	}
}

func (ctx *genKclTypeContext) convertSchemaFromGoPackage() ([]convertResult, error) {
	err := ctx.fetchStructs(ctx.pkgPath)
	if err != nil {
		return nil, err
	}
	return ctx.convertFromGoStructs(), nil
}

// convertFromGoStructs walks the already-populated ctx.goStructs table
// (in deterministic pkgPath+name order) and translates each Go struct
// into a KCL schema in ctx.resultMap. It assumes ctx.goStructs and
// ctx.packages are populated; callers are responsible for that. This
// split lets GenKclModule reuse the conversion logic without the
// double-load that fetchStructs would otherwise introduce.
func (ctx *genKclTypeContext) convertFromGoStructs() []convertResult {
	// Sort the structs by (pkgPath, name) so iteration is deterministic.
	// This matters because ctx.aliasFor memoises alias assignment in
	// the order paths are encountered; a stable traversal ensures the
	// same input always produces the same output.
	sortedStructs := make([]goStruct, 0, len(ctx.goStructs))
	for _, s := range ctx.goStructs {
		sortedStructs = append(sortedStructs, s)
	}
	sort.Slice(sortedStructs, func(i, j int) bool {
		if sortedStructs[i].pkgPath != sortedStructs[j].pkgPath {
			return sortedStructs[i].pkgPath < sortedStructs[j].pkgPath
		}
		return sortedStructs[i].name < sortedStructs[j].name
	})
	var results []convertResult
	for _, s := range sortedStructs {
		name := s.name
		if _, ok := ctx.resultMap[name]; !ok {
			result := convertResult{IsSchema: true}
			result.schema.Name = name
			result.schema.Description = s.doc
			ctx.resultMap[name] = result
			for _, field := range s.fields {
				typeName := ctx.typeName(ctx.pkgPath, name, field.name, field.ty)
				fieldName := formatName(field.name)
				// Use alias name and type defined in the `kcl` or `json`` tag
				tagName, tagTy, _, err := parserGoStructFieldTag(field.tag)
				if err == nil {
					if tagName != "" {
						fieldName = tagName
					}
					if tagTy != nil {
						typeName = tagTy
					}
				}
				result.schema.Properties = append(result.Properties, property{
					Name:        fieldName,
					Type:        typeName,
					Description: s.fieldDocs[field.name],
				})
			}
			ctx.resultMap[name] = result
		}
	}
	// Append anonymous structs
	for _, key := range getSortedKeys(ctx.resultMap) {
		if ctx.resultMap[key].IsSchema {
			results = append(results, ctx.resultMap[key])
		}
	}
	return results
}

func (ctx *genKclTypeContext) recordTypeInfo(pkg *packages.Package) {
	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.StructType:
				ctx.tyMapping[pkg.TypesInfo.TypeOf(n)] = n
			case *ast.GenDecl:
				if n.Tok == token.TYPE {
					for _, spec := range n.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if n.Doc != nil && typeSpec.Name != nil {
								// <pkg_path>.<name>
								ctx.tySpecMapping[pkg.PkgPath+"."+typeSpec.Name.String()] = n.Doc.Text()
							}
						}
					}
				}
			}
			return true
		})
	}
}

func (ctx *genKclTypeContext) addPackage(p *packages.Package) {
	if pkg, ok := ctx.packages[p.PkgPath]; ok {
		if p != pkg {
			panic(fmt.Sprintf("duplicate package %s", p.PkgPath))
		}
		return
	}
	ctx.packages[p.PkgPath] = p
	ctx.recordTypeInfo(p)
	for _, pkg := range p.Imports {
		ctx.addPackage(pkg)
	}
}

// loadGoPackages wraps the standard packages.Load call shared by the
// single-package and module entry points. The returned slice is ordered
// with the caller-supplied package first when present, followed by its
// imports in BFS order (the order returned by packages.Load).
func loadGoPackages(pattern string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedTypes |
			packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, err
	}
	var errs []string
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, e := range pkg.Errors {
				errs = append(errs, fmt.Sprintf("\t%s: %v", pkg.PkgPath, e))
			}
		}
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("could not load Go packages:\n%s", strings.Join(errs, "\n"))
	}
	return pkgs, nil
}

// populateFrom records the supplied package as the ctx's `mainPkgPath`
// and pulls in every dependency that its reachable type graph needs.
// `addPackage` recurses `p.Imports` so type lookups resolve without us
// re-loading the dependency tree, `fetchStructsFromPkg` populates the
// struct table (filtered by oneFile/mainPkgPath), and `recordConstInfo`
// captures any Go `const`s declared in the main package.
func (ctx *genKclTypeContext) populateFrom(p *packages.Package) {
	if ctx.mainPkgPath == "" {
		ctx.mainPkgPath = p.PkgPath
	}
	ctx.addPackage(p)
	if err := ctx.fetchStructsFromPkg(p); err != nil {
		// fetchStructsFromPkg currently returns nil or panics; the
		// error parameter is kept for future expansion without breaking
		// the populateFrom signature.
		_ = err
	}
	if p.PkgPath == ctx.mainPkgPath {
		ctx.recordConstInfo(p)
	}
}

func (ctx *genKclTypeContext) fetchStructs(pkgPath string) error {
	pkgs, err := loadGoPackages(pkgPath)
	if err != nil {
		return err
	}
	// The first package is the one containing the input file; everything
	// else loaded by packages.Load is an imported dependency. We only want
	// to emit KCL schemas for structs from the main package, but we still
	// need the dependencies to be in `ctx.packages` so type lookups in
	// `typeName` succeed.
	if len(pkgs) > 0 {
		ctx.mainPkgPath = pkgs[0].PkgPath
	}
	for _, p := range pkgs {
		ctx.populateFrom(p)
	}
	return nil
}

func (ctx *genKclTypeContext) fetchStructsFromPkg(pkg *packages.Package) error {
	ctx.recordTypeInfo(pkg)
	// When OneFile is false we only want to emit KCL schemas for structs
	// from the main Go package; cross-package structs will be referenced
	// via `import` statements and their alias-qualified type names.
	// When OneFile is true we keep the historical behaviour of pulling
	// every reachable struct into a single self-contained KCL file.
	if !ctx.oneFile && pkg.PkgPath != ctx.mainPkgPath {
		return nil
	}
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj, ok := obj.(*types.TypeName); ok {
			if named, ok := obj.Type().(*types.Named); ok {
				if _, ok := named.Underlying().(*types.Struct); ok {
					pkgPath := named.Obj().Pkg().Path()
					pkgName := named.Obj().Pkg().Name()
					structDoc := ctx.getStructDoc(pkgPath, name)
					fields, fieldDocs := ctx.getStructFieldsAndDocs(obj.Type())
					ctx.goStructs[named.Obj()] = goStruct{
						pkgPath:   pkgPath,
						pkgName:   pkgName,
						name:      name,
						fields:    fields,
						doc:       structDoc,
						fieldDocs: fieldDocs,
					}
				}
			}
		}
	}

	return nil
}

func (ctx *genKclTypeContext) getStructDoc(pkgName, structName string) string {
	if spec, ok := ctx.tySpecMapping[pkgName+"."+structName]; ok {
		return spec
	}
	return ""
}

// recordConstInfo walks the AST of the supplied Go package and records
// every statically-evaluable `const` declaration as a KCL global
// variable in `ctx.goConsts`. Declarations that depend on runtime values
// or that the type checker could not constant-fold are silently skipped
// so the generator never produces invalid KCL output.
//
// Implicit repetitions inside an iota block (e.g. `B` and `C` in
//
//	const (
//	    A = iota
//	    B
//	    C
//	)
//
// ) are not represented in Go's AST, so they are skipped here. Callers
// who want each step emitted must spell out the iota expression
// explicitly, e.g. `B = iota`.
func (ctx *genKclTypeContext) recordConstInfo(pkg *packages.Package) {
	if pkg.TypesInfo == nil {
		return
	}
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				// When a single ValueSpec carries multiple Names
				// (e.g. `const A, B = 1, 2`) we require each Name to
				// have a matching value expression. Names without one
				// are skipped — see the doc comment above.
				for i, name := range vs.Names {
					if name == nil || name.Name == "_" {
						continue
					}
					if i >= len(vs.Values) {
						continue
					}
					valueExpr := vs.Values[i]
					if valueExpr == nil {
						continue
					}
					tv, ok := pkg.TypesInfo.Types[valueExpr]
					if !ok || tv.Value == nil {
						continue
					}
					doc := ""
					if vs.Doc != nil {
						doc = vs.Doc.Text()
					} else if gd.Doc != nil {
						doc = gd.Doc.Text()
					}
					ctx.goConsts[pkg.PkgPath+"."+name.Name] = global{
						Name:        name.Name,
						Type:        goTypeToKclTypeString(tv.Type),
						Value:       goConstToKclValue(tv.Value),
						Description: doc,
					}
				}
			}
		}
	}
}

// goTypeToKclTypeString maps a Go type to the KCL type annotation we
// want to emit for a `const`. Returns the empty string when the type is
// untyped (so the annotation can be omitted) or when the underlying
// kind is not one we know how to render in KCL.
func goTypeToKclTypeString(t types.Type) string {
	if t == nil {
		return ""
	}
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return ""
	}
	switch basic.Kind() {
	case types.UntypedBool, types.UntypedInt, types.UntypedRune,
		types.UntypedFloat, types.UntypedComplex, types.UntypedString,
		types.UntypedNil:
		return ""
	case types.Bool:
		return "bool"
	case types.String:
		return "str"
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.Uintptr:
		return "int"
	case types.Float32, types.Float64:
		return "float"
	}
	return ""
}

// goConstToKclValue evaluates a `go/constant` value into the native Go
// scalar that KCL's value formatter knows how to render. Returns nil
// when the kind is not representable (e.g. complex numbers); callers
// must check for nil before emitting.
//
// For floats the boolean returned by `Float64Val` is deliberately
// ignored: a Go literal such as `3.14` is stored internally as the
// exact rational `157/50` and `Float64Val` reports it as
// "not representable" because the nearest float64 is not bit-identical
// to that rational. We still want the familiar `3.14` rendering, so we
// always emit the float64 approximation the user expects.
func goConstToKclValue(c constant.Value) any {
	if c == nil {
		return nil
	}
	switch c.Kind() {
	case constant.Bool:
		return constant.BoolVal(c)
	case constant.String:
		return constant.StringVal(c)
	case constant.Int:
		if i, ok := constant.Int64Val(c); ok {
			return i
		}
		// Out of int64 range — fall back to the literal source so the
		// emitted KCL still preserves the original value.
		return c.ExactString()
	case constant.Float:
		if f, _ := constant.Float64Val(c); f != 0 || c.Kind() == constant.Float {
			return f
		}
		return c.ExactString()
	}
	return nil
}

func (ctx *genKclTypeContext) getStructFieldsAndDocs(typ types.Type) ([]field, map[string]string) {
	switch ty := typ.(type) {
	case *types.Pointer:
		return ctx.getStructFieldsAndDocs(ty.Elem())
	case *types.Named:
		if structType, ok := ty.Underlying().(*types.Struct); ok {
			return ctx.getStructTypeFieldsAndDocs(structType)
		}
	case *types.Struct:
		return ctx.getStructTypeFieldsAndDocs(ty)
	}
	return nil, nil
}

func (ctx *genKclTypeContext) getStructTypeFieldsAndDocs(structType *types.Struct) ([]field, map[string]string) {
	fieldDocs := make(map[string]string)
	var fields []field
	for i := 0; i < structType.NumFields(); i++ {
		f := structType.Field(i)
		var tag string
		if astStruct, ok := ctx.tyMapping[structType]; ok {
			// Match by field position to get the correct tag
			astFieldIndex := 0
			for _, field := range astStruct.Fields.List {
				if len(field.Names) == 0 {
					// This is an embedded field
					if astFieldIndex == i {
						if field.Tag != nil {
							tag = field.Tag.Value
						}
						break
					}
					astFieldIndex++
				} else {
					// Named fields - check if any match
					for _, fieldName := range field.Names {
						if fieldName.Name == f.Name() {
							if field.Doc != nil {
								fieldDocs[fieldName.Name] = field.Doc.Text()
							}
							if field.Tag != nil {
								tag = field.Tag.Value
							}
							break
						}
					}
					astFieldIndex += len(field.Names)
				}
			}
		}
		if f.Embedded() {
			// Parse tag to check if inline option is present
			_, _, tagOpts, _ := parserGoStructFieldTag(tag)
			if tagOpts.inline {
				// Only inline if the "inline" option is present in the tag
				embeddedFields, embeddedFieldDocs := ctx.getEmbeddedFieldsAndDocs(f.Type())
				fields = append(fields, embeddedFields...)
				for k, v := range embeddedFieldDocs {
					fieldDocs[k] = v
				}
			} else {
				// Don't inline - treat as a regular field
				// Use the name from the tag if available
				fieldName := f.Name()
				if tagName, _, _, err := parserGoStructFieldTag(tag); err == nil && tagName != "" {
					fieldName = tagName
				}
				if f.Exported() {
					fields = append(fields, field{
						name:      fieldName,
						ty:        f.Type(),
						tag:       tag,
						anonymous: true,
					})
				}
			}
		} else {
			if f.Exported() {
				fields = append(fields, field{
					name: f.Name(),
					ty:   f.Type(),
					tag:  tag,
				})
			}
		}
	}
	return fields, fieldDocs
}

func (ctx *genKclTypeContext) getEmbeddedFieldsAndDocs(typ types.Type) ([]field, map[string]string) {
	fieldDocs := make(map[string]string)
	var fields []field
	switch ty := typ.(type) {
	case *types.Pointer:
		fields, fieldDocs = ctx.getEmbeddedFieldsAndDocs(ty.Elem())
	case *types.Named:
		if _, ok := ty.Underlying().(*types.Struct); ok {
			fields, fieldDocs = ctx.getStructFieldsAndDocs(typ)
		}
	case *types.Struct:
		fields, fieldDocs = ctx.getStructFieldsAndDocs(typ)
	}
	return fields, fieldDocs
}

// tagOptions represents the parsed options from a struct tag
type tagOptions struct {
	inline    bool
	omitempty bool
}

func parserGoStructFieldTag(tag string) (string, typeInterface, tagOptions, error) {
	tagMap := make(map[string]string, 0)
	sp := strings.Split(tag, "`")
	if len(sp) == 1 {
		return "", nil, tagOptions{}, errors.New("this field not found tag string like ``")
	}
	value, ok := lookupTag(sp[1], "kcl")
	if !ok {
		value, ok = lookupTag(sp[1], "json")
		if !ok {
			value, ok = lookupTag(sp[1], "yaml")
			if !ok {
				return "", nil, tagOptions{}, errors.New("not found tag key named json, yaml or kcl")
			}
		}
		// Deal json or yaml tags
		tagInfos := strings.Split(value, ",")
		if len(tagInfos) > 0 {
			name := tagInfos[0]
			opts := parseTagOptions(tagInfos[1:])
			return name, nil, opts, nil
		} else {
			return "", nil, tagOptions{}, errors.New("invalid tag key named json")
		}
	}
	// Deal kcl tags
	reg := "name=.*,type=.*"
	match, err := regexp.Match(reg, []byte(value))
	if err != nil {
		return "", nil, tagOptions{}, err
	}
	if !match {
		return "", nil, tagOptions{}, errors.New("don't match the kcl tag info, the tag info style is name=NAME,type=TYPE")
	}
	tagInfo := strings.Split(value, ",")
	for _, s := range tagInfo {
		t := strings.Split(s, "=")
		tagMap[t[0]] = t[1]
	}
	fieldType := tagMap["type"]
	if strings.Contains(tagMap["type"], ")|") {
		typeUnionList := strings.Split(tagMap["type"], "|")
		var ss []string
		for _, u := range typeUnionList {
			_, _, litValue := isLitType(u)
			ss = append(ss, litValue)
		}
		fieldType = strings.Join(ss, "|")
	}
	return tagMap["name"], typeCustom{
		Name: fieldType,
	}, tagOptions{}, nil
}

// parseTagOptions parses tag options like "inline", "omitempty"
func parseTagOptions(options []string) tagOptions {
	var opts tagOptions
	for _, opt := range options {
		switch strings.TrimSpace(opt) {
		case "inline":
			opts.inline = true
		case "omitempty":
			opts.omitempty = true
		}
	}
	return opts
}

func isLitType(fieldType string) (ok bool, basicTyp, litValue string) {
	if !strings.HasSuffix(fieldType, ")") {
		return
	}

	i := strings.Index(fieldType, "(") + 1
	j := strings.LastIndex(fieldType, ")")

	switch {
	case strings.HasPrefix(fieldType, "bool("):
		return true, "bool", fieldType[i:j]
	case strings.HasPrefix(fieldType, "int("):
		return true, "int", fieldType[i:j]
	case strings.HasPrefix(fieldType, "float("):
		return true, "float", fieldType[i:j]
	case strings.HasPrefix(fieldType, "str("):
		return true, "str", strconv.Quote(fieldType[i:j])
	}
	return
}

// kclPathIdentRegexp matches a valid KCL identifier: an ASCII letter or
// underscore followed by zero or more ASCII letters, digits, or underscores.
// This matches the subset of XID_Start + XID_Continue that rustc_lexer
// accepts in KCL identifier positions (see
// crates/parser/src/parser/stmt.rs:436 → parse_identifier).
var kclPathIdentRegexp = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// sanitizeKclPathSegment rewrites a single Go package-path segment so it
// is a legal KCL identifier. Concretely:
//
//   - every rune outside [A-Za-z0-9_] is mapped to '_',
//   - a leading digit (Go segments such as "99designs") or KCL keyword
//     result is prefixed with "_" so the segment is a syntactically
//     valid KCL identifier.
//
// The function always returns a non-empty string; callers may use the
// returned string as both a KCL import-path segment and (after strcase
// conversion) as a candidate KCL import alias.
func sanitizeKclPathSegment(seg string) string {
	if seg == "" {
		return "_"
	}
	var b strings.Builder
	b.Grow(len(seg))
	for _, r := range seg {
		switch {
		case r >= 'A' && r <= 'Z',
			r >= 'a' && r <= 'z',
			r >= '0' && r <= '9',
			r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := b.String()
	if !kclPathIdentRegexp.MatchString(out) || isKclKeyword(out) {
		return "_" + out
	}
	return out
}

// goPkgToKclImport converts a Go package import path (e.g. `github.com/foo/bar`
// or a local module-relative path like `pkg/tools/gen/testdata/gostruct/external`)
// into the dotted form used in KCL `import` statements. Every segment is
// run through sanitizeKclPathSegment so hyphens, dots, leading digits, and
// other characters that are not legal in a KCL identifier are mapped to
// '_' (with a leading '_' prefix added when the sanitised segment would
// otherwise start with a digit or collide with a KCL keyword).
//
// We keep the full dotted path rather than just the last segment so that
// two Go packages that happen to share a final path component remain
// distinguishable in the generated KCL.
func goPkgToKclImport(goPkgPath string) string {
	parts := strings.Split(goPkgPath, "/")
	for i, p := range parts {
		parts[i] = sanitizeKclPathSegment(p)
	}
	return strings.Join(parts, ".")
}

// aliasForGoPkg derives a deterministic KCL import alias from a Go package
// path. We use the last path component (capitalised) as the alias; the
// surrounding code qualifies cross-package types with this alias, e.g.
// `Bar.Outer`. If the sanitised result is empty, starts with a digit, or
// is a KCL keyword we fall back to `Imported_<sanitised segment>` so the
// generated file still parses.
func aliasForGoPkg(goPkgPath string) string {
	// Accept either the original Go path (with `/`) or the dotted form
	// produced by goPkgToKclImport so callers can invoke this with
	// whichever representation they have handy.
	seg := goPkgPath
	if i := strings.LastIndexAny(seg, "/."); i >= 0 {
		seg = seg[i+1:]
	}
	if seg == "" {
		return "Imported"
	}
	seg = sanitizeKclPathSegment(seg)
	// Validate the POST-ToCamel result. Sanitize already added a leading
	// '_' for digit- and keyword-input segments, but strcase.ToCamel can
	// still produce an invalid identifier (e.g. `_99designs` → `99Designs`,
	// which starts with a digit) or collide with a KCL keyword (e.g.
	// `true` → `_true` → `True`). On either failure we fall back to a
	// stable, non-colliding alias so the generated file always parses.
	alias := strcase.ToCamel(seg)
	if alias == "" {
		return "Imported_" + seg
	}
	r := []rune(alias)
	if r[0] >= '0' && r[0] <= '9' {
		return "Imported_" + alias
	}
	if isKclKeyword(alias) {
		return "Imported_" + alias
	}
	return alias
}

// aliasFor returns the memoised KCL import alias for the supplied
// sanitised dotted Go path. The first caller for a given path receives
// the alias produced by `aliasForGoPkg`; subsequent callers with the
// same preferred alias (typically from a *different* Go path that
// shares a final segment such as `internal/a/types` vs
// `internal/b/types`) get `Types2`, `Types3`, … so the KCL `import`
// lines stay unambiguous and the alias-qualified type names can't
// silently shadow one another.
func (ctx *genKclTypeContext) aliasFor(dottedPath string) string {
	if existing, ok := ctx.importAliases[dottedPath]; ok {
		return existing
	}
	preferred := aliasForGoPkg(dottedPath)
	taken := make(map[string]struct{}, len(ctx.importAliases))
	for _, a := range ctx.importAliases {
		taken[a] = struct{}{}
	}
	candidate := preferred
	for n := 2; ; n++ {
		if _, takenAlready := taken[candidate]; !takenAlready {
			break
		}
		candidate = preferred + strconv.Itoa(n)
	}
	ctx.importAliases[dottedPath] = candidate
	return candidate
}

func lookupTag(tag, key string) (value string, ok bool) {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
