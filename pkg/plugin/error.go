// Copyright 2023 The KCL Authors. All rights reserved.

package plugin

import (
	"encoding/json"
)

type PanicInfo struct {
	IsPanic bool `json:"__kcl_PanicInfo__"`

	BackTrace []BacktraceFrame `json:"backtrace,omitempty"`
	RustFile  string           `json:"rust_file,omitempty"`
	RustLine  int              `json:"rust_line,omitempty"`
	RustCol   int              `json:"rust_col,omitempty"`

	KclPkgPath string `json:"kcl_pkgpath,omitempty"`
	KclFile    string `json:"kcl_file,omitempty"`
	KCLFunc    string `json:"kcl_func,omitempty"`
	KclLine    int    `json:"kcl_line,omitempty"`
	KclCol     int    `json:"kcl_col,omitempty"`
	KclArgMsg  string `json:"kcl_arg_msg,omitempty"`

	// only for schema check
	KclConfigMetaFile   string `json:"kcl_config_meta_file,omitempty"`
	KclConfigMetaLine   int    `json:"kcl_config_meta_line,omitempty"`
	KclConfigMetaCol    int    `json:"kcl_config_meta_col,omitempty"`
	KclConfigMetaArgMsg string `json:"kcl_config_meta_arg_msg,omitempty"`

	Message     string `json:"message"`
	ErrTypeCode string `json:"err_type_code,omitempty"`
	IsWarning   string `json:"is_warning,omitempty"`
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
			IsPanic: true,
			Message: err.Error(),
		}
		return x.JSONError()
	}
	return ""
}

func (p *PanicInfo) JSONError() string {
	p.IsPanic = true
	d, _ := json.Marshal(p)
	return string(d)
}

func (p *PanicInfo) Error() string {
	return p.JSONError()
}
