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
				//Add packages to the global store
				ps := c.Args().Slice()[c.Args().Len()-1]
				if c.Bool("git") {
					ps = "git:" + ps
				} else {
					ps = "registry:" + ps
				}
				println(ps)
				pkgStruct, err := GetRequirePkgStruct(ps)
				if err != nil {
					return err
				}
				rb := RequireBase{
					RequirePkgStruct: *pkgStruct,
				}
				err = kpmC.Get(&rb)
				if err != nil {
					return err
				}
			case "addfile":
				//Add the current working directory to the global store
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
