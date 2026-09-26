// Copyright The KCL Authors. All rights reserved.

//go:build cgo
// +build cgo

// Package testpluginagent provides a self-contained KCL plugin agent used
// to exercise the WithPluginAgent option end to end: KCL code calls a host
// function through the agent pointer passed to the native client.
package testpluginagent

/*
#include <stdint.h>

uint64_t test_plugin_agent_get_proxy_ptr();
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

const addMethod = "kcl_plugin.testplugin.add"

var (
	calls      int64
	methodLock sync.Mutex
	lastMethod string
)

//export testPluginAgentInvoke
func testPluginAgentInvoke(method, argsJSON, kwargsJSON *C.char) *C.char {
	atomic.AddInt64(&calls, 1)
	methodLock.Lock()
	lastMethod = C.GoString(method)
	methodLock.Unlock()

	var args []any
	if argsJSON != nil {
		if err := json.Unmarshal([]byte(C.GoString(argsJSON)), &args); err != nil {
			return C.CString(errorJSON(err))
		}
	}
	result, err := invoke(C.GoString(method), args)
	if err != nil {
		return C.CString(errorJSON(err))
	}
	data, err := json.Marshal(result)
	if err != nil {
		return C.CString(errorJSON(err))
	}
	return C.CString(string(data))
}

func invoke(method string, args []any) (any, error) {
	if method != addMethod {
		return nil, fmt.Errorf("unknown method: %s", method)
	}
	if len(args) != 2 {
		return nil, fmt.Errorf("add expects 2 args, got %d", len(args))
	}
	a, err := toInt64(args[0])
	if err != nil {
		return nil, err
	}
	b, err := toInt64(args[1])
	if err != nil {
		return nil, err
	}
	return a + b, nil
}

func toInt64(v any) (int64, error) {
	switch v := v.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unexpected arg type: %T", v)
	}
}

func errorJSON(err error) string {
	data, _ := json.Marshal(map[string]any{"__kcl_PanicInfo__": err.Error()})
	return string(data)
}

// Ptr returns the plugin agent function pointer.
func Ptr() uint64 {
	return uint64(C.test_plugin_agent_get_proxy_ptr())
}

// Calls returns how many times the agent has been invoked.
func Calls() int64 {
	return atomic.LoadInt64(&calls)
}

// LastMethod returns the method name the agent was last invoked with.
func LastMethod() string {
	methodLock.Lock()
	defer methodLock.Unlock()
	return lastMethod
}
