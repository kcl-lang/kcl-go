// Copyright The KCL Authors. All rights reserved.

package kcl

import (
	"fmt"
	"testing"
)

func tAssert(tb testing.TB, condition bool, args ...interface{}) {
	if !condition {
		tb.Helper()
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("Assert failed, %s", msg)
		} else {
			tb.Fatalf("Assert failed")
		}
	}
}

func tAssertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	if !condition {
		tb.Helper()
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatalf("tAssert failed, %s", msg)
		} else {
			tb.Fatalf("tAssert failed")
		}
	}
}
