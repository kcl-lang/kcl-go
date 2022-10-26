package kpm

import "github.com/urfave/cli/v2"

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
			return nil
		},
	}
}
