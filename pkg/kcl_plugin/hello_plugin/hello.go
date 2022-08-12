// Copyright 2022 The KCL Authors. All rights reserved.

package hello_plugin

import (
	"fmt"
	"strings"
	"sync/atomic"

	"kusionstack.io/kclvm-go/pkg/kcl_plugin"
)

var global_int int64

func _reset_plugin() {
	atomic.SwapInt64(&global_int, 0)
}

func set_global_int(v int64) {
	if kcl_plugin.DebugMode {
		fmt.Printf("plugin-go://hello_plugin.set_global_int(%d)\n", v)
	}
	atomic.SwapInt64(&global_int, v)
}

func get_global_int() int64 {
	if kcl_plugin.DebugMode {
		fmt.Printf("plugin-go://hello_plugin.get_global_int(): %d\n", global_int)
	}
	return atomic.LoadInt64(&global_int)
}

func say_hello(msg string) {
	fmt.Println("hello.say_hello:", msg)
}

func add(a, b int64) int64 {
	if kcl_plugin.DebugMode {
		fmt.Printf("plugin-go://hello_plugin.add(%d, %d)\n", a, b)
	}
	return a + b
}

func tolower(s string) string {
	return strings.ToLower(s)
}

func update_dict(d map[string]interface{}, key, value string) map[string]interface{} {
	d[key] = value
	return d
}

func list_append(list []interface{}, values ...interface{}) []interface{} {
	return append(list, values...)
}

func foo(
	a, b interface{},
	__sep__ *struct{},
	x []interface{},
	values map[string]interface{},
) interface{} {
	m := make(map[string]interface{})

	m["a"] = a
	m["b"] = b
	m["x"] = x

	for k, v := range values {
		m[k] = v
	}

	return m
}
