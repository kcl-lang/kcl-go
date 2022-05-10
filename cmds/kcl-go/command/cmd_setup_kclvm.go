// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/scripts"
)

var cmdSetupKclvmFlags = []cli.Flag{}

func NewSetpupKclvmCmd() *cli.Command {
	return &cli.Command{
		Hidden:    false,
		Name:      "setup-kclvm",
		Usage:     "setup kclvm command",
		ArgsUsage: "kclvm-root",
		Flags:     cmdSetupKclvmFlags,
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			err := scripts.SetupKclvm(c.Args().First())
			if err != nil {
				return err
			}

			return nil
		},
	}
}
