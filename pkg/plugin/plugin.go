// Copyright 2023 The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

package plugin

/*
#include <stdint.h>
#include <stdlib.h>

extern char* kcl_go_capi_InvokeJsonProxy(
    char* method,
    char* args_json,
    char* kwargs_json
);

static uint64_t kcl_go_capi_getInvokeJsonProxyPtr() {
	return (uint64_t)(kcl_go_capi_InvokeJsonProxy);
}
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
)

const CgoEnabled = true

//export kcl_go_capi_InvokeJsonProxy
func kcl_go_capi_InvokeJsonProxy(_method, _args_json, _kwargs_json *C.char) (result_json *C.char) {
	var method, args_json, kwargs_json string

	if _method != nil {
		method = C.GoString(_method)
	}
	if _args_json != nil {
		args_json = C.GoString(_args_json)
	}
	if _kwargs_json != nil {
		kwargs_json = C.GoString(_kwargs_json)
	}

	result := InvokeJson(method, args_json, kwargs_json)
	return c_String_new(result)
}

func GetInvokeJsonProxyPtr() uint64 {
	ptr := uint64(C.kcl_go_capi_getInvokeJsonProxyPtr())
	return ptr
}

func Invoke(method string, args []interface{}, kwargs map[string]interface{}) (result_json string) {
	var args_json, kwargs_json string

	if len(args) > 0 {
		d, err := json.Marshal(args)
		if err != nil {
			return JSONError(err)
		}
		args_json = string(d)
	}

	if kwargs != nil {
		d, err := json.Marshal(kwargs)
		if err != nil {
			return JSONError(err)
		}
		kwargs_json = string(d)
	}

	return _Invoke(method, args_json, kwargs_json)
}

func InvokeJson(method, args_json, kwargs_json string) (result_json string) {
	return _Invoke(method, args_json, kwargs_json)
}

func _Invoke(method, args_json, kwargs_json string) (result_json string) {
	defer func() {
		if r := recover(); r != nil {
			result_json = JSONError(errors.New(fmt.Sprint(r)))
		}
	}()

	// check method
	if method == "" {
		return JSONError(fmt.Errorf("empty method"))
	}

	// parse args, kwargs
	args, err := ParseMethodArgs(args_json, kwargs_json)
	if err != nil {
		return JSONError(err)
	}

	// todo: check args type
	// todo: check kwargs type

	// get method
	methodSpec, found := GetMethodSpec(method)
	if !found {
		return JSONError(fmt.Errorf("invalid method: %s, not found", method))
	}

	// call plugin method
	result, err := methodSpec.Body(args)
	if err != nil {
		return JSONError(err)
	}
	if result == nil {
		result = new(MethodResult)
	}

	// encode result
	data, err := json.Marshal(result.V)
	if err != nil {
		return JSONError(err)
	}

	return string(data)
}
