// Copyright The KCL Authors. All rights reserved.

package settings

import (
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestLoadSettingsFiles(t *testing.T) {
	result, err := LoadSettingsFiles("./test_data", []string{"./test_data/settings/kcl.yaml"})
	if err != nil {
		t.Fatal(err)
	}
	assert2.Empty(t, result.KclCliConfigs.Files)
	assert2.True(t, result.KclCliConfigs.StrictRangeCheck)
	assert2.Equal(t, 1, len(result.KclOptions))
	assert2.Equal(t, "key", result.KclOptions[0].Key)
	assert2.Equal(t, `"value"`, result.KclOptions[0].Value)
}
