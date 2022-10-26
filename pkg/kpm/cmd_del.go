package kpm

import "github.com/urfave/cli/v2"

func NewDelCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "del",
		Usage:  "del  dependencies pkg",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}
			println("del...")
			return nil
		},
	}
}
