// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	kcl "kcl-lang.io/kcl-go"
)

func NewFmtCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "fmt",
		Usage:  "format KCL sources",
		Action: func(c *cli.Context) error {
			ss, err := kcl.FormatPath(c.Args().First())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, s := range ss {
				fmt.Println(s)
			}
			return nil
		},
	}
}
