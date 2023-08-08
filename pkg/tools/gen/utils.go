// Copyright 2023 The KCL Authors. All rights reserved.

package gen

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

var _ = assert
var _ = assertf

func assert(ok bool, a ...interface{}) {
	if !ok {
		panic(fmt.Sprint(a...))
	}
}

func assertf(ok bool, format string, a ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(format, a...))
	}
}

func jsonString(p interface{}) string {
	x, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return ""
	}
	return string(x)
}

func readSource(filename string, src interface{}) (data []byte, err error) {
	if src == nil {
		src, err = os.ReadFile(filename)
		if err != nil {
			return
		}
	}

	switch src := src.(type) {
	case []byte:
		return src, nil
	case string:
		return []byte(src), nil
	case io.Reader:
		d, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported src type: %T", src)
	}
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
