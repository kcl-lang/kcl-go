// Copyright The KCL Authors. All rights reserved.

package rename

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestRename(t *testing.T) {
	bakContent, err := os.ReadFile("./testdata/rename/main.k.bak")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("./testdata/rename/main.k", bakContent, 0o644); err != nil {
		t.Fatal(err)
	}

	changedFiles, err := Rename("./testdata/rename", "a", []string{"./testdata/rename/main.k"}, "a2")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, 1, len(changedFiles))
	assert2.True(t, strings.Contains(filepath.ToSlash(changedFiles[0]), "testdata/rename/main.k"))
}

func TestRenameCode(t *testing.T) {
	changedCodes, err := RenameCode("/mock/path", "a", map[string]string{
		"/mock/path/main.k": "a = 1\nb = a",
	}, "a2")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, map[string]string{
		"/mock/path/main.k": "a2 = 1\nb = a2",
	}, changedCodes)
}
