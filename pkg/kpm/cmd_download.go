package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/PathHandle"
	"github.com/urfave/cli/v2"
	"os"
)

func NewDownloadCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "download",
		Usage:  "download dependencies pkg to local cache and link to workspace",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				//此处不应该有参数
				cli.ShowAppHelpAndExit(c, 0)
			}
			println("download...")
			filebytes, err := os.ReadFile(kpmC.WorkDir + PathHandle.Separator + "kpm.json")
			if err != nil {
				return err
			}
			kf := KpmFile{}
			err = json.Unmarshal(filebytes, &kf)
			if err != nil {
				return err
			}
			for _, rb := range kf.Direct {
				err = kpmC.Get(&rb)
				if err != nil {
					return err
				}
			}
			for ps, integrity := range kf.Indirect {
				pkgStruct, err := GetRequirePkgStruct(ps)
				if err != nil {
					return err
				}
				rb := RequireBase{
					RequirePkgStruct: *pkgStruct,
					Integrity:        integrity,
				}
				err = kpmC.Get(&rb)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
