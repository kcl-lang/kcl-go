// Copyright The KCL Authors. All rights reserved.

package override

import (
	"testing"
)

func TestOverrideFile(t *testing.T) {
	_, err := OverrideFile("./testdata/test.k", []string{"config.image=\"image/image:v1\""}, []string{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = OverrideFile("./testdata/test.k", []string{"config.replicas:1"}, []string{})
	if err != nil {
		t.Fatal(err)
	}
}
