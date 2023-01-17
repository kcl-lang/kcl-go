package kpm

import (
	"github.com/urfave/cli/v2"
)

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
			kf, err := kpmC.LoadKpmFileStructInWorkdir()
			if err != nil {
				return err
			}
			delete(kf.Direct, c.Args().Slice()[c.Args().Len()-1])
			err = kpmC.SaveKpmFileInWorkdir(kf)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
