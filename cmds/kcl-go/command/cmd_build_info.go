// Copyright 2022 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

func NewBuildInfoCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "build-info",
		Usage:  "print build information",
		Action: func(c *cli.Context) error {
			info, ok := debug.ReadBuildInfo()
			if !ok {
				fmt.Println("ERR: ReadBuildInfo failed")
				os.Exit(1)
			}

			fmt.Println(info)
			return nil
		},
	}
}
