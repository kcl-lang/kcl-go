// Copyright 2021 The KCL Authors. All rights reserved.

package utils

import (
	"fmt"
	"testing"
)

func Assert(condition bool, args ...interface{}) {
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			panic("Assert failed, " + msg)
		} else {
			panic("Assert failed")
		}
	}
}

func TAssert(tb testing.TB, condition bool, args ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("Assert failed, %s", msg)
		} else {
			tb.Fatalf("Assert failed")
		}
	}
}

func TAssertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatalf("tAssert failed, %s", msg)
		} else {
			tb.Fatalf("tAssert failed")
		}
	}
}
