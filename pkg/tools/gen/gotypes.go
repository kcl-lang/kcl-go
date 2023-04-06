package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type GoStruct struct {
	Name          string
	Fields        []*GoStructField
	FieldNum      int
	StructComment string
}

type GoStructField struct {
	FieldName    string
	FieldType    string
	FieldTag     string
	FieldTagKind string
	FieldComment string
}

// ParseGoSourceCode parse go source code from .go file path or source code
func ParseGoSourceCode(filename string, src interface{}) ([]*GoStruct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	goStructList := make([]*GoStruct, 0)
	for _, v := range f.Decls {
		goStruct := &GoStruct{}
		if stc, ok := v.(*ast.GenDecl); ok && stc.Tok == token.TYPE {
			if stc.Doc != nil {
				goStruct.StructComment = strings.TrimRight(stc.Doc.Text(), "\n")
			}
			for _, spec := range stc.Specs {
				goStructFields := make([]*GoStructField, 0)
				if tp, ok := spec.(*ast.TypeSpec); ok {
					goStruct.Name = tp.Name.Name
					if stp, ok := tp.Type.(*ast.StructType); ok {
						if !stp.Struct.IsValid() {
							continue
						}
						goStruct.FieldNum = stp.Fields.NumFields()
						for _, field := range stp.Fields.List {
							goStructField := &GoStructField{}

							// get field name
							if len(field.Names) == 1 {
								goStructField.FieldName = field.Names[0].Name
							} else if len(field.Names) > 1 {
								for _, name := range field.Names {
									goStructField.FieldName = goStructField.FieldName + name.String() + ","
								}
							}

							// get tag
							if field.Tag != nil {
								goStructField.FieldTag = field.Tag.Value
								goStructField.FieldTagKind = field.Tag.Kind.String()
							}

							// get comment
							if field.Comment != nil {
								goStructField.FieldComment = strings.TrimRight(field.Comment.Text(), "")
							}

							// get field type
							goStructField.FieldType = getTypeName(field.Type)
							goStructFields = append(goStructFields, goStructField)
						}
					}
					goStruct.Fields = goStructFields
				}
			}
			goStructList = append(goStructList, goStruct)
		}
	}
	return goStructList, nil
}

func getTypeName(f ast.Expr) string {
	if ft, ok := f.(*ast.Ident); ok {
		return ft.Name
	}
	if ft, ok := f.(*ast.ArrayType); ok {
		item := getTypeName(ft.Elt)
		return fmt.Sprintf("[]%s", item)
	}
	if ft, ok := f.(*ast.MapType); ok {
		key := getTypeName(ft.Key)
		value := getTypeName(ft.Value)
		return fmt.Sprintf("map[%s]%s", key, value)
	}
	if ft, ok := f.(*ast.StarExpr); ok {
		value := getTypeName(ft.X)
		return fmt.Sprintf("*%s", value)
	}
	if _, ok := f.(*ast.InterfaceType); ok {
		return "interface{}"
	}
	return ""
}
