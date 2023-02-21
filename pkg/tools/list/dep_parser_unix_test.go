//go:build linux || darwin
// +build linux darwin

package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFilePath(t *testing.T) {
	_, err := newImportDepParser("./testdata/complicate/", DepOptions{Files: []string{"appops/projectA/invalid.k"}, UpStreams: []string{}})
	assert.EqualError(t, err, "invalid file path: stat testdata/complicate/appops/projectA/invalid.k: no such file or directory", "err not match")
}
