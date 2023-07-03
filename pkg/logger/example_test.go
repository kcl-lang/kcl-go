// Copyright 2021 The KCL Authors. All rights reserved.

package logger_test

import (
	"log"
	"os"

	"kcl-lang.io/kcl-go/pkg/logger"
)

func Example() {
	var logger = logger.GetLogger()

	logger.SetLevel("DEBUG")
	logger.Debug("1+1=2")
	logger.Info("hello")
}

func ExampleNewStdLogger() {
	var logger = logger.NewStdLogger(os.Stderr, "", "", log.Lshortfile)

	logger.Debug("1+1=2")
	logger.Info("hello")
}
