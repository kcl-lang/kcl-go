package kpm

import (
	"errors"
	"github.com/orangebees/go-oneutils/Semver"
	"github.com/urfave/cli/v2"
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
			kf, err := kpmC.LoadKpmFileStructInWorkdir()
			if err != nil {
				return err
			}
			//操作，先get包，再检测直接依赖中是否有同名包。检测包最小版本，检测
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
			shortname := rb.GetShortName()
			kf.Direct[shortname] = rb
			dkf, err := kpmC.LoadKpmFileStruct(&rb)
			if err == nil {
				//找到文件
				kfv, err := Semver.NewFromString(kf.KclvmMinVersion)
				if err != nil {
					return err
				}
				dkfv, err := Semver.NewFromString(dkf.KclvmMinVersion)
				if err != nil {
					return err
				}
				if kfv.Cmp(dkfv) == -1 {
					e := errors.New("the KclvmMinVersion of the added dependency " + shortname + " is greater than the KclvmMinVersion of the workspace")
					println(e.Error())
					return e
				}
				for k, v := range dkf.Indirect {
					kf.Indirect[k] = v
				}
				for _, v := range dkf.Direct {
					kf.Indirect[v.GetPkgString()] = v.Integrity
				}
			}
			//保存
			err = kpmC.SaveKpmFileInWorkdir(kf)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
