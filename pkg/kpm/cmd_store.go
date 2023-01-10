package kpm

import (
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/urfave/cli/v2"
)

func NewStoreCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "store",
		Usage:  "Reads and performs actions on kpm store that is on the current filesystem",
		Flags: []cli.Flag{&cli.BoolFlag{
			Name:  "git",
			Usage: "add git pkg",
		}},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}
			switch c.Args().First() {
			case "add":
				//添加包到全局存储
			case "addfile":
				//添加当前工作目录到全局存储
				fim, err := kpmC.GitStore.AddDir(kpmC.WorkDir)
				if err != nil {
					return err
				}
				GlobalStore.ReleaseFileInfoMap(fim)
			default:
				cli.ShowAppHelpAndExit(c, 0)
				return nil
			}
			return nil
		},
	}
}
