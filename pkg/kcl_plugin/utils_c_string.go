// Copyright 2022 The KCL Authors. All rights reserved.

package kcl_plugin

// #include <stdlib.h>
import "C"
import (
	"sync"
	"unsafe"
)

var c_String struct {
	sync.Mutex
	nextId int
	buf    []*C.char
}

func init() {
	c_String.buf = make([]*C.char, 100)
}

func c_String_new(s string) *C.char {
	c_String.Lock()
	defer c_String.Unlock()

	id := c_String.nextId % len(c_String.buf)
	c_String.nextId++

	if cs := c_String.buf[id]; cs != nil {
		C.free(unsafe.Pointer(cs))
	}

	cs := C.CString(s)
	c_String.buf[id] = cs
	return cs
}
