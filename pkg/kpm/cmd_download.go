package kpm

import (
	"github.com/urfave/cli/v2"
)

func NewDownloadCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "download",
		Usage:  "download dependencies pkg to local cache and link to workspace",
		Action: func(c *cli.Context) error {
			if c.NArg() != 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}
			println("download...")
			kf, err := kpmC.LoadKpmFileStructInWorkdir()
			if err != nil {
				return err
			}
			globalWriterFlag := false
			for rbn, rb := range kf.Direct {
				writerFlag := rb.Integrity == ""
				err = kpmC.Get(&rb)
				if err != nil {
					println(err.Error())
					return err
				}
				if writerFlag {
					globalWriterFlag = true
					kf.Direct[rbn] = rb
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
				writerFlag := rb.Integrity == ""
				err = kpmC.Get(&rb)
				if err != nil {
					println(err.Error())
					return err
				}
				if writerFlag {
					globalWriterFlag = true
					kf.Indirect[ps] = rb.Integrity
				}
			}
			if globalWriterFlag {
				err = kpmC.SaveKpmFileInWorkdir(kf)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
