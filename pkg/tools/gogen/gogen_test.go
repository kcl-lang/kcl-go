package gogen

import (
	"testing"

	"kusionstack.io/kclvm-go/pkg/kcl"
)

func TestExample(t *testing.T) {
	result, err := kcl.GetSchemaType("./testdata/main.k", "", "")
	if err != nil {
		t.Fatal(err)
	}
	goCode := GenGoCodeFromKclType(result)
	t.Logf("go code: \n%s", goCode)
}
