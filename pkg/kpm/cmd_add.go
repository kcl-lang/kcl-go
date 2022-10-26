package kpm

import "github.com/urfave/cli/v2"

func NewAddCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "add",
		Usage:  "add dependencies pkg",
		Flags: []cli.Flag{&cli.BoolFlag{
			Name:  "git",
			Usage: "add git pkg",
		}},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}
			println("add...")
			return nil
		},
	}
}
