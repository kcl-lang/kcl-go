package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"regexp"
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
	name string
	ty   types.Type
	tag  string
}

type genKclTypeContext struct {
	context
	// Go package path.
	pkgPath string
	// Go structs in all package path
	goStructs map[*types.TypeName]goStruct
	// All pkg path -> package mapping
	packages map[string]*packages.Package
	// Semantic type -> AST struct type mapping
	tyMapping map[types.Type]*ast.StructType
	// Semantic type -> AST struct type mapping
	tySpecMapping map[string]string
	// Generate all go structs into one KCL file.
	oneFile bool
}

func (k *kclGenerator) genSchemaFromGoStruct(w io.Writer, filename string, _ interface{}) error {
	ctx := genKclTypeContext{
		pkgPath: filename,
		context: context{
			resultMap: make(map[string]convertResult),
			imports:   make(map[string]struct{}),
			paths:     []string{},
		},
		goStructs:     map[*types.TypeName]goStruct{},
		packages:      map[string]*packages.Package{},
		tyMapping:     map[types.Type]*ast.StructType{},
		tySpecMapping: map[string]string{},
		oneFile:       true,
	}
	results, err := ctx.convertSchemaFromGoPackage()
	if err != nil {
		return err
	}
	kclSch := kclFile{
		Schemas: []schema{},
	}
	for _, result := range results {
		if result.IsSchema {
			kclSch.Schemas = append(kclSch.Schemas, result.schema)
		}
	}
	// generate kcl schema code
	return k.genKcl(w, kclSch)
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
						} else {
							// Struct from current package
							ty := typeCustom{
								Name: pkgPath + "." + obj.Name(),
							}
							return ty
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
				tagName, tagTy, err := parserGoStructFieldTag(field.tag)
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
	var results []convertResult
	for _, s := range ctx.goStructs {
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
				tagName, tagTy, err := parserGoStructFieldTag(field.tag)
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
	return results, nil
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

func (ctx *genKclTypeContext) fetchStructs(pkgPath string) error {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedTypes |
			packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return err
	}
	// Check Go module loader errors
	var errs []string
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, e := range pkg.Errors {
				errs = append(errs, fmt.Sprintf("\t%s: %v", pkg.PkgPath, e))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("could not load Go packages:\n%s", strings.Join(errs, "\n"))
	}
	for _, p := range pkgs {
		ctx.addPackage(p)
	}
	for _, pkg := range pkgs {
		ctx.fetchStructsFromPkg(pkg)
	}
	return nil
}

func (ctx *genKclTypeContext) fetchStructsFromPkg(pkg *packages.Package) error {
	ctx.recordTypeInfo(pkg)
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
		if structType, ok := ctx.tyMapping[structType]; ok {
			for _, field := range structType.Fields.List {
				for _, fieldName := range field.Names {
					if fieldName.Name == f.Name() {
						if field.Doc != nil {
							fieldDocs[fieldName.Name] = field.Doc.Text()
						}
						if field.Tag != nil {
							tag = field.Tag.Value
						}
					}
				}
			}
		}
		if f.Embedded() {
			embeddedFields, embeddedFieldDocs := ctx.getEmbeddedFieldsAndDocs(f.Type())
			fields = append(fields, embeddedFields...)
			for k, v := range embeddedFieldDocs {
				fieldDocs[k] = v
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

func parserGoStructFieldTag(tag string) (string, typeInterface, error) {
	tagMap := make(map[string]string, 0)
	sp := strings.Split(tag, "`")
	if len(sp) == 1 {
		return "", nil, errors.New("this field not found tag string like ``")
	}
	value, ok := lookupTag(sp[1], "kcl")
	if !ok {
		value, ok = lookupTag(sp[1], "json")
		if !ok {
			value, ok = lookupTag(sp[1], "yaml")
			if !ok {
				return "", nil, errors.New("not found tag key named json, yaml or kcl")
			}
		}
		// Deal json or yaml tags
		tagInfos := strings.Split(value, ",")
		if len(tagInfos) > 0 {
			return tagInfos[0], nil, nil
		} else {
			return "", nil, errors.New("invalid tag key named json")
		}
	}
	// Deal kcl tags
	reg := "name=.*,type=.*"
	match, err := regexp.Match(reg, []byte(value))
	if err != nil {
		return "", nil, err
	}
	if !match {
		return "", nil, errors.New("don't match the kcl tag info, the tag info style is name=NAME,type=TYPE")
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
	}, nil
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
