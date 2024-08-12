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
	pkgPath   string
	goStructs map[*types.TypeName]goStruct
	oneFile   bool
}

func (k *kclGenerator) genSchemaFromGoStruct(w io.Writer, filename string, _ interface{}) error {
	ctx := genKclTypeContext{
		pkgPath: filename,
		context: context{
			resultMap: make(map[string]convertResult),
			imports:   make(map[string]struct{}),
			paths:     []string{},
		},
		oneFile: true,
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

func (ctx *genKclTypeContext) typeName(defName string, fieldName string, ty types.Type) typeInterface {
	switch ty := ty.(type) {
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
		return ctx.typeName(defName, fieldName, ty.Elem())
	case *types.Named:
		obj := ty.Obj()
		switch {
		case obj.Pkg().Path() == "time" && obj.Name() == "Time":
			return typePrimitive(typStr)
		case obj.Pkg().Path() == "time" && obj.Name() == "Duration":
			return typePrimitive(typInt)
		case obj.Pkg().Path() == "math/big" && obj.Name() == "Int":
			return typePrimitive(typInt)
		default:
			if _, ok := ctx.goStructs[obj]; !ok {
				return ctx.typeName(defName, fieldName, ty.Underlying())
			} else {
				return typeCustom{
					Name: obj.Name(),
				}
			}
		}
	case *types.Array:
		return typeArray{
			Items: ctx.typeName(defName, fieldName, ty.Elem()),
		}
	case *types.Slice:
		return typeArray{
			Items: ctx.typeName(defName, fieldName, ty.Elem()),
		}
	case *types.Map:
		return typeDict{
			Key:   ctx.typeName(defName, fieldName, ty.Key()),
			Value: ctx.typeName(defName, fieldName, ty.Elem()),
		}
	case *types.Struct:
		schemaName := fmt.Sprintf("%s%s", defName, strcase.ToCamel(fieldName))
		if _, ok := ctx.resultMap[schemaName]; !ok {
			result := convertResult{IsSchema: true}
			ctx.resultMap[schemaName] = result
			for i := 0; i < ty.NumFields(); i++ {
				sf := ty.Field(i)
				typeName := ctx.typeName(schemaName, sf.Name(), sf.Type())
				result.schema.Name = schemaName
				result.schema.Properties = append(result.Properties, property{
					Name: formatName(sf.Name()),
					Type: typeName,
				})
				ctx.resultMap[schemaName] = result
			}
		}
		return typeCustom{
			Name: schemaName,
		}
	case *types.Union:
		var types []typeInterface
		for i := 0; i < ty.Len(); i++ {
			types = append(types, ctx.typeName(defName, fieldName, ty.Term(i).Type()))
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
			types = append(types, ctx.typeName(defName, fieldName, ty.EmbeddedType(i)))
		}
		return typeUnion{
			Items: types,
		}
	default:
		return typePrimitive(typAny)
	}
}

func (ctx *genKclTypeContext) convertSchemaFromGoPackage() ([]convertResult, error) {
	structs, error := fetchStructs(ctx.pkgPath)
	ctx.goStructs = structs
	if error != nil {
		return nil, error
	}
	var results []convertResult
	for _, s := range structs {
		name := s.name
		if _, ok := ctx.resultMap[name]; !ok {
			result := convertResult{IsSchema: true}
			result.schema.Name = name
			result.schema.Description = s.doc
			for _, field := range s.fields {
				typeName := ctx.typeName(name, field.name, field.ty)
				fieldName := formatName(field.name)
				tagName, tagTy, err := parserGoStructFieldTag(field.tag)
				if err == nil && tagName != "" && tagTy != nil {
					fieldName = tagName
					typeName = tagTy
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

func fetchStructs(pkgPath string) (map[*types.TypeName]goStruct, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, err
	}
	structs := make(map[*types.TypeName]goStruct)
	for _, pkg := range pkgs {
		astFiles := pkg.Syntax
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if obj, ok := obj.(*types.TypeName); ok {
				named, _ := obj.Type().(*types.Named)
				if structType, ok := named.Underlying().(*types.Struct); ok {
					structDoc := getStructDoc(name, astFiles)
					fields, fieldDocs := getStructFieldsAndDocs(structType, name, astFiles)
					pkgPath := named.Obj().Pkg().Path()
					pkgName := named.Obj().Pkg().Name()
					structs[named.Obj()] = goStruct{
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
	return structs, nil
}

func getStructDoc(structName string, astFiles []*ast.File) string {
	for _, file := range astFiles {
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
						if genDecl.Doc != nil {
							return genDecl.Doc.Text()
						}
					}
				}
			}
		}
	}
	return ""
}

func getStructFieldsAndDocs(structType *types.Struct, structName string, astFiles []*ast.File) ([]field, map[string]string) {
	fieldDocs := make(map[string]string)
	var fields []field
	for i := 0; i < structType.NumFields(); i++ {
		f := structType.Field(i)
		var tag string
		for _, file := range astFiles {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == structName {
							if structType, ok := typeSpec.Type.(*ast.StructType); ok {
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
						}
					}
				}
			}
		}
		if f.Embedded() {
			embeddedFields, embeddedFieldDocs := getEmbeddedFieldsAndDocs(f.Type(), astFiles, structName)
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

func getEmbeddedFieldsAndDocs(t types.Type, astFiles []*ast.File, structName string) ([]field, map[string]string) {
	fieldDocs := make(map[string]string)
	var fields []field
	switch t := t.(type) {
	case *types.Pointer:
		fields, fieldDocs = getEmbeddedFieldsAndDocs(t.Elem(), astFiles, structName)
	case *types.Named:
		if structType, ok := t.Underlying().(*types.Struct); ok {
			fields, fieldDocs = getStructFieldsAndDocs(structType, structName, astFiles)
		}
	case *types.Struct:
		fields, fieldDocs = getStructFieldsAndDocs(t, structName, astFiles)
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
		return "", nil, errors.New("not found tag key named kcl")
	}
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
