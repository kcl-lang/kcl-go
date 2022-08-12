// Copyright 2022 The KCL Authors. All rights reserved.

package kcl_plugin

import (
	"os"
	"strings"
	"sync"
)

// Debug flag for kcl-plugin
var DebugMode = false

func init() {
	if s := strings.ToLower(os.Getenv("KCLVM_PLUGIN_DEBUG")); s != "" && s != "false" && s != "0" {
		DebugMode = true
	}
}

var pluginManager struct {
	allPlugin     map[string]Plugin
	allMethodSpec map[string]MethodSpec
	sync.Mutex
}

func init() {
	pluginManager.allPlugin = make(map[string]Plugin)
	pluginManager.allMethodSpec = make(map[string]MethodSpec)
}

// Register register a new kcl plugin.
func RegisterPlugin(plugin Plugin) {
	pluginManager.Lock()
	defer pluginManager.Unlock()

	if plugin.Name == "" {
		panic("invalid plugin: empty name")
	}

	pluginManager.allPlugin[plugin.Name] = plugin
	for methodName, methodSpec := range plugin.MethodMap {
		methodAbsName := "kcl_plugin." + plugin.Name + "." + methodName
		pluginManager.allMethodSpec[methodAbsName] = methodSpec
	}
}

// GetPlugin get plugin object by name.
func GetPlugin(name string) (plugin Plugin, ok bool) {
	pluginManager.Lock()
	defer pluginManager.Unlock()

	x, ok := pluginManager.allPlugin[name]
	return x, ok
}

// GetMethodSpec get plugin method by name.
func GetMethodSpec(methodName string) (method MethodSpec, ok bool) {
	pluginManager.Lock()
	defer pluginManager.Unlock()

	idx := strings.LastIndex(methodName, ".")
	if idx <= 0 || idx >= len(methodName)-1 {
		return MethodSpec{}, false
	}

	x, ok := pluginManager.allMethodSpec[methodName]
	return x, ok
}

// ResetPlugin reset all kcl plugin state.
func ResetPlugin() {
	pluginManager.Lock()
	defer pluginManager.Unlock()

	for _, p := range pluginManager.allPlugin {
		if p.ResetFunc != nil {
			p.ResetFunc()
		}
	}
}
