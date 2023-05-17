//go:build linux || darwin
// +build linux darwin

package list

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFilePath(t *testing.T) {
	_, err := newImportDepParser("./testdata/complicate/", DepOptions{Files: []string{"appops/projectA/invalid.k"}, UpStreams: []string{}})
	assert.Equal(t, strings.Contains(err.Error(), "appops/projectA/invalid.k: no such file or directory"), true)
}
