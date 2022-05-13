// Copyright 2021 The KCL Authors. All rights reserved.

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"kusionstack.io/kclvm-go/pkg/compiler/parser"
)

var (
	flagKFile = flag.String("kcl-file", "a.k", "set input kcl file")
)

func main() {
	flag.Parse()
	if *flagKFile == "" {
		fmt.Println("no kcl file")
		os.Exit(1)
	}

	f, err := parser.ParseFile(*flagKFile, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(JSONString(f.JSON))
}

func JSONString(v interface{}) string {
	if s, ok := v.(string); ok {
		v = []byte(s)
	}
	if x, ok := v.([]byte); ok {
		var m map[string]interface{}
		if err := json.Unmarshal(x, &m); err != nil {
			return string(x)
		}
		result, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			return string(x)
		}
		return string(result)
	}
	x, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}
	return string(x)
}
