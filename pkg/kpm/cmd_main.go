package kpm

import (
	"github.com/urfave/cli/v2"
)

func CLI(args ...string) error {
	app := cli.NewApp()
	app.Name = "kpm"
	app.Usage = "kpm is a kcl package manager"
	app.Version = "v0.0.1-alpha.1"
	app.UsageText = CliHelp
	app.Commands = []*cli.Command{
		NewInitCmd(),
		NewAddCmd(),
		NewDelCmd(),
		NewDownloadCmd(),
		NewStoreCmd(),
	}
	err := Setup()
	if err != nil {
		return err
	}
	//添加一个参数确保与os.Args数量一致
	nargs := make([]string, len(args))
	nargs = nargs[:1]
	nargs = append(nargs, args...)
	err = app.Run(nargs)
	if err != nil {
		return err
	}
	return nil
}