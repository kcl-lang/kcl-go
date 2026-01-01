// Copyright The KCL Authors. All rights reserved.

package utils

import (
	"fmt"
	"testing"
)

func Assert(condition bool, args ...any) {
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			panic("Assert failed, " + msg)
		} else {
			panic("Assert failed")
		}
	}
}

func TAssert(tb testing.TB, condition bool, args ...any) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("Assert failed, %s", msg)
		} else {
			tb.Fatalf("Assert failed")
		}
	}
}

func TAssertf(tb testing.TB, condition bool, format string, a ...any) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatalf("tAssert failed, %s", msg)
		} else {
			tb.Fatalf("tAssert failed")
		}
	}
}
