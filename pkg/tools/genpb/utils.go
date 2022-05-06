// Copyright 2022 The KCL Authors. All rights reserved.

package genpb

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
