// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/scripts"
)

var cmdSetupKclvmFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "triple",
		Usage: "set kclvm triple",
		Value: scripts.DefaultKclvmTriple,
	},
	&cli.StringFlag{
		Name:  "outdir",
		Usage: "set kclvm output dir",
		Value: "_" + scripts.DefaultKclvmTriple + "-root_",
	},
}

func NewSetpupKclvmCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "setup-kclvm",
		Usage:  "setup kclvm command",
		Flags:  cmdSetupKclvmFlags,
		Action: func(c *cli.Context) error {
			// go run ./cmds/kcl-go/ setup-kclvm --triple=kclvm-ubuntu

			triple := c.String("triple")
			outdir := c.String("outdir")

			if triple == "" || outdir == "" {
				cli.ShowCommandHelpAndExit(c, "setup-kclvm", 0)
			}

			if triple != scripts.DefaultKclvmTriple {
				if outdir == "" || outdir == "_"+scripts.DefaultKclvmTriple+"-root_" {
					outdir = "_" + triple + "-root_"
				}
			}

			scripts.DefaultKclvmTriple = triple

			err := scripts.SetupKclvm(outdir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			return nil
		},
	}
}
