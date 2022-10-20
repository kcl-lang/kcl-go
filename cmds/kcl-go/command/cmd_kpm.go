package command

import (
	"github.com/urfave/cli/v2"
	"kusionstack.io/kclvm-go/pkg/kpm"
)

func NewKpmCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "kpm",
		Usage:  "kpm is a kcl package manager",
		Action: func(c *cli.Context) error {
			kpm.CLI(c.Args().Slice()...)
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}
			return nil
		},
	}
}
