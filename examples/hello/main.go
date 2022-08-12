// Copyright 2022 The KCL Authors. All rights reserved.

// Run CGO Mode:
// KCLVM_SERVICE_CLIENT_HANDLER=native KCLVM_PLUGIN_DEBUG=1 go run -tags=kclvm_service_capi .

package main

import (
	"fmt"

	"kusionstack.io/kclvm-go"
	_ "kusionstack.io/kclvm-go/pkg/kcl_plugin/hello_plugin"
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

hello.say_hello('hello KusionStack')

low_kcl = hello.tolower('KCL')

s = hello.update_dict({'name': 123}, 'name', 'kcl')['name']

# todo: l = hello.list_append(['abc'], 'name', 123)

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person {}
x1 = Person {
    age = 101
}
`
