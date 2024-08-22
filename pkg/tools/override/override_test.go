// Copyright The KCL Authors. All rights reserved.

package override

import (
	"os"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestOverrideFile(t *testing.T) {
	file := "./testdata/test.k"
	_, err := OverrideFile(file, []string{"config.image=\"image/image:v1\""}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = OverrideFile(file, []string{"config.replicas:1"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = OverrideFile(file, []string{"config.s=pkg.Service {}"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, strings.ReplaceAll(string(got), "\r\n", "\n"), strings.ReplaceAll(`import pkg

schema Config:
    image: str
    replicas: int

if True:
    configOther = Config {image = "image/other:v1"}

config: Config {
    image = "image/image:v1"
    replicas: 1
    s = pkg.Service {}
}
`, "\r\n", "\n"))
}

func TestOverrideFileWithRelativeImport(t *testing.T) {
	_, err := OverrideFile("./testdata/test_with_relative_import.k", []string{"config.replicas=1"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}
