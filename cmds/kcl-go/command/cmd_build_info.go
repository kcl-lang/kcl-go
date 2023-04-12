// Copyright 2023 The KCL Authors. All rights reserved.

package command

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
)

func NewBuildInfoCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "build-info",
		Usage:  "print build info",
		Action: func(c *cli.Context) error {
			if info, ok := debug.ReadBuildInfo(); ok {
				if d, err := json.MarshalIndent(info, "", "\t"); err == nil {
					fmt.Println(string(d))
					return nil
				}
			}

			fmt.Println("ERROR: read build info failed")
			os.Exit(1)

			return nil
		},
	}
}
