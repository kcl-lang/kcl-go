package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/PathHandle"
	"github.com/urfave/cli/v2"
	"os"
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
			filebytes, err := os.ReadFile(kpmC.WorkDir + PathHandle.Separator + "kpm.json")
			if err != nil {
				return err
			}
			kf := KpmFile{}
			err = json.Unmarshal(filebytes, &kf)
			if err != nil {
				return err
			}
			delete(kf.Direct, c.Args().Slice()[c.Args().Len()-1])
			marshal, err := json.Marshal(&kf)
			if err != nil {
				return err
			}
			err = os.WriteFile(kpmC.WorkDir+PathHandle.Separator+"kpm.json", marshal, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
