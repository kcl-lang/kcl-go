// Copyright 2021 The KCL Authors. All rights reserved.

package parser

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFile(t *testing.T) {
	const (
		tfile_k    = "testdata/a.k"
		tfile_json = "testdata/a.k.ast.json"
	)

	f, err := ParseFile(tfile_k, nil)
	if err != nil {
		t.Fatal(err)
	}

	var (
		got  = f.Module.JSONMap()
		want map[string]interface{}
	)

	if x, err := os.ReadFile(tfile_json); err == nil {
		want = json_DecodeMap(t, x)
	} else {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, cmp_EquateEmpty()); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func json_DecodeMap(t *testing.T, data []byte) map[string]interface{} {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	return m
}
func cmp_EquateEmpty() cmp.Option {
	isEmpty := func(x interface{}) bool {
		if x == nil {
			return true
		}
		vx := reflect.ValueOf(x)
		switch {
		case x == nil:
			return true
		default:
			switch vx.Kind() {
			case reflect.Slice, reflect.Map:
				return vx.Len() == 0
			case reflect.String:
				return vx.Len() == 0
			}
		}
		return false
	}

	return cmp.FilterValues(
		func(x, y interface{}) bool {
			return isEmpty(x) && isEmpty(y)
		},
		cmp.Comparer(func(_, _ interface{}) bool {
			return true
		}),
	)
}
