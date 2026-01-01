// Copyright 2023 The KCL Authors. All rights reserved.

package gen

import (
	"encoding/json"
	"fmt"
	"sort"
)

var _ = assert
var _ = assertf

func assert(ok bool, a ...any) {
	if !ok {
		panic(fmt.Sprint(a...))
	}
}

func assertf(ok bool, format string, a ...any) {
	if !ok {
		panic(fmt.Sprintf(format, a...))
	}
}

func jsonString(p any) string {
	x, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return ""
	}
	return string(x)
}

// getSortedKeys returns the keys sorted in alphabetical order of a map.
func getSortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
