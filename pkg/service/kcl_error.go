// Copyright The KCL Authors. All rights reserved.

package service

import (
	"encoding/json"
	"fmt"
)

const KCLVM_SERVER_ERROR_CODE int64 = 0x4B434C // the ASCII code of "KCL"

type ServerError struct {
	Code int64  `json:"code,omitempty"`
	Msg  string `json:"message,omitempty"`
}

func wrapKclvmServerError(err error) error {
	error_with_code := ServerError{}
	serde_error := json.Unmarshal([]byte(err.Error()), &error_with_code)
	if serde_error == nil && error_with_code.Code == KCLVM_SERVER_ERROR_CODE {
		err = fmt.Errorf("%s", error_with_code.Msg)
	}
	return err
}
