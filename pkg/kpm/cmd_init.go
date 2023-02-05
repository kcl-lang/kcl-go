package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/Convert"
	"github.com/orangebees/go-oneutils/PathHandle"
	"github.com/urfave/cli/v2"
	"os"
)

func NewInitCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "init",
		Usage:  "initialize new module in current directory",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				println("not args...")
				cli.ShowAppHelpAndExit(c, 0)
			}
			println("init...")
			_, err := os.Stat(kpmC.WorkDir + PathHandle.Separator + "kpm.json")
			if err == nil {
				println("kpm.json is exist")
				return nil
			}

			marshal, err := json.Marshal(KpmFile{
				PackageName:     c.Args().First(),
				KclvmMinVersion: kpmC.KclVmVersion,
				Direct:          nil,
				Indirect:        nil,
			})
			if err != nil {
				return err
			}
			err = os.WriteFile(kpmC.WorkDir+PathHandle.Separator+"kpm.json", marshal, 0777)
			if err != nil {
				return err
			}
			println("Create kpm.json success!")
			_, err = os.Stat(kpmC.WorkDir + PathHandle.Separator + "kcl.mod")
			if err == nil {
				return nil
			}
			//The file does not exist, so this file is created
			err = os.WriteFile(kpmC.WorkDir+PathHandle.Separator+"kcl.mod", Convert.S2B(DefaultKclModContent), 0777)
			if err != nil {
				return err
			}
			println("Create kcl.mod success!")
			return nil
		},
	}
}
