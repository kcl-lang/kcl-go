// Copyright The KCL Authors. All rights reserved.

package utils

import "testing"

func TestFinePkgPath(t *testing.T) {
	pkgPath, err := GoodPkgPath("./testdata/sub/main.k")
	TAssert(t, err == nil, err)
	TAssert(t, pkgPath == "sub", pkgPath)

	pkgPath, err = GoodPkgPath("./testdata/a/b/x.k")
	TAssert(t, err == nil, err)
	TAssert(t, pkgPath == "a/b", pkgPath)
}
