package command

import (
	"github.com/urfave/cli/v2"
	"kusionstack.io/kclvm-go/pkg/kpm"
)

func NewKpmCmd() *cli.Command {
	return &cli.Command{
		Hidden: true,
		Name:   "kpm",
		Usage:  "kpm is a kcl package manager",
		Action: func(c *cli.Context) error {
			err := kpm.CLI(c.Args().Slice()...)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
