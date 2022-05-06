// Copyright 2022 The KCL Authors. All rights reserved.

package main

import (
	"fmt"

	"kusionstack.io/kclvm-go"
)

func main() {
	yaml := kclvm.MustRun("hello.k", kclvm.WithCode(k_code)).First().YAMLString()
	fmt.Println(yaml)
}

const k_code = `
import kcl_plugin.hello

name = "kcl"
age = 1
two = hello.add(1, 1)

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {}
x1 = Person {
    age = 101
}
`
