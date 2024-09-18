//go:build cgo
// +build cgo

package main

import (
	"fmt"

	"kcl-lang.io/kcl-go/pkg/kcl"                   // Import the native API
	_ "kcl-lang.io/kcl-go/pkg/plugin/hello_plugin" // Import the hello plugin
)

func main() {
	yaml := kcl.MustRun("main.k", kcl.WithCode(code)).GetRawYamlResult()
	fmt.Println(yaml)
}

const code = `
import kcl_plugin.hello

name = "kcl"
three = hello.add(1,2)
`
