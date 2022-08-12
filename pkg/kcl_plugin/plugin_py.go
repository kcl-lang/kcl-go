// Copyright 2022 The KCL Authors. All rights reserved.

package kcl_plugin

import "fmt"

func py_callPluginMethod(method, args_json, kwargs_json string) (result_json string) {
	return JSONError(fmt.Errorf("invalid method py: %s, not found", method))
}
