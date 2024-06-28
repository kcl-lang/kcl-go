// Copyright 2023 The KCL Authors. All rights reserved.

package plugin

import (
	"encoding/json"
)

type PanicInfo struct {
	Message string `json:"__kcl_PanicInfo__"`
}

type BacktraceFrame struct {
	File string `json:"file,omitempty"`
	Func string `json:"func,omitempty"`
	Line int    `json:"line,omitempty"`
	Col  int    `json:"col,omitempty"`
}

func JSONError(err error) string {
	if x, ok := err.(*PanicInfo); ok {
		return x.JSONError()
	}
	if err != nil {
		x := &PanicInfo{
			Message: err.Error(),
		}
		return x.JSONError()
	}
	return ""
}

func (p *PanicInfo) JSONError() string {
	d, _ := json.Marshal(p)
	return string(d)
}

func (p *PanicInfo) Error() string {
	return p.JSONError()
}
