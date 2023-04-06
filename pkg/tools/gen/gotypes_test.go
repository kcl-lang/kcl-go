package gen

import (
	"fmt"
	"testing"
)

func TestParseGoFiles(t *testing.T) {
	goStructs, err := ParseGoSourceCode("./testdata/genkcldata.go", nil)
	if err != nil {
		t.Fatalf("parse go source code err: %s", err.Error())
	}
	t.Logf("%#v", goStructs)
	for _, v := range goStructs {
		fmt.Printf("Struct Name: %s\n", v.Name)
		fmt.Printf("struct Num: %d\n", v.FieldNum)
		fmt.Printf("struct Comment: %s\n", v.StructComment)
		for _, f := range v.Fields {
			fmt.Printf("\tFieldName: %s\n", f.FieldName)
			fmt.Printf("\tFieldType: %s\n", f.FieldType)
			fmt.Printf("\tFieldTag: %s\n", f.FieldTag)
			fmt.Printf("\tFieldTagKind: %s\n", f.FieldTagKind)
			fmt.Printf("\tFieldComment: %s\n", f.FieldComment)
			fmt.Printf("\t\n")
		}
	}
}
