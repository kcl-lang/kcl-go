// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func newTourCmd() *cli.Command {
	return &cli.Command{
		Hidden: true,
		Name:   "tour",
		Usage:  "kclvm command tour",
		Action: func(c *cli.Context) error {
			fmt.Println("TODO")
			return nil
		},
	}
}
