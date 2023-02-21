package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFilePath(t *testing.T) {
	_, err := newImportDepParser("./testdata/complicate/", DepOptions{Files: []string{"appops/projectA/invalid.k"}, UpStreams: []string{}})
	assert.EqualError(t, err, "invalid file path: CreateFile testdata/complicate/appops/projectA/invalid.k: The system cannot find the file specified.", "err not match")
}
