// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// list pkgs
// list pkg.Objects
// list pkg.Types
// list options
// list schema attributes

func NewListCmd() *cli.Command {
	return &cli.Command{
		Hidden: true,
		Name:   "list",
		Usage:  "list packages/names/options/attributes",
		Action: func(c *cli.Context) error {
			fmt.Println("TODO")
			return nil
		},
	}
}
