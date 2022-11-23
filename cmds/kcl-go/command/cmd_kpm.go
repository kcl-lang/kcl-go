package command

import (
	"github.com/urfave/cli/v2"
	"kusionstack.io/kclvm-go/pkg/kpm2"
)

func NewKpmCmd() *cli.Command {
	return &cli.Command{
		Hidden: true,
		Name:   "kpm",
		Usage:  "kpm is a kcl package manager",
		Action: func(c *cli.Context) error {
			kpm2.CLI(c.Args().Slice()...)
			//err := kpm2.CLI(c.Args().Slice()...)
			//if err != nil {
			//	return err
			//}

			return nil
		},
	}
}
