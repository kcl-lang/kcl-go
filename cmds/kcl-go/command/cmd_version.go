// Copyright 2023 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
	"kcl-lang.io/kcl-go"
)

func NewVersionCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "version",
		Usage:  "print version info",
		Action: func(c *cli.Context) error {
			fmt.Printf("%s-%s-%s\n", runtime.GOOS, runtime.GOARCH, kclvm.KclvmAbiVersion)
			return nil
		},
	}
}
