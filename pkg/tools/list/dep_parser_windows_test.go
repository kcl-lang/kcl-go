package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFilePath(t *testing.T) {
	_, err := newImportDepParser(".\\testdata\\complicate\\", DepOptions{Files: []string{"appops\\projectA\\invalid.k"}, UpStreams: []string{}})
	assert.EqualError(t, err, "invalid file path: stat appops\\projectA\\invalid.k: invalid argument", "err not match")
}
