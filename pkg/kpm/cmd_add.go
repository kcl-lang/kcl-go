package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/PathHandle"
	"github.com/urfave/cli/v2"
	"os"
)

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
			//读取
			filebytes, err := os.ReadFile(kpmC.WorkDir + PathHandle.Separator + "kpm.json")
			if err != nil {
				return err
			}
			kf := KpmFile{}
			err = json.Unmarshal(filebytes, &kf)
			if err != nil {
				return err
			}

			//操作，先get包，再检测直接依赖中是否有同名包。检测包版本，检测
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
			if kf.Direct == nil {
				kf.Direct = make(DirectRequire, 16)
			}
			if kf.Indirect == nil {
				kf.Indirect = make(IndirectRequire, 16)
			}
			kf.Direct[rb.GetShortName()] = rb
			//保存
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
