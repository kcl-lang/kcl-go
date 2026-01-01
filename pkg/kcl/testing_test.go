// Copyright The KCL Authors. All rights reserved.

package kcl

import (
	"fmt"
	"testing"
)

func tAssert(tb testing.TB, condition bool, args ...any) {
	if !condition {
		tb.Helper()
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("Assert failed, %s", msg)
		} else {
			tb.Fatalf("Assert failed")
		}
	}
}
