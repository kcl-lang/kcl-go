// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/logger"
)

func Main() {
	app := cli.NewApp()
	app.Name = "kcl-go"
	app.Usage = "K Configuration Language Virtual Machine"
	app.Version = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" {
				return info.Main.Version
			}
		}
		return "(devel)"
	}()

	// kclvm -m kclvm
	app.UsageText = `kcl-go
   kcl-go [global options] command [command options] [arguments...]

   kcl-go kcl -h
   kcl-go -h

    ___  __    ________  ___
   |\  \|\  \ |\   ____\|\  \
   \ \  \/  /|\ \  \___|\ \  \
    \ \   ___  \ \  \    \ \  \
     \ \  \\ \  \ \  \____\ \  \____
      \ \__\\ \__\ \_______\ \_______\
       \|__| \|__|\|_______|\|_______|`

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Run in debug mode (for developers only)",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logger.GetLogger().SetLevel("DEBUG")
			fmt.Println("kclvm path:", kclvm_runtime.MustGetKclvmPath())
		}
		return nil
	}

	app.Commands = []*cli.Command{
		NewBuildInfoCmd(),

		NewRunCmd(),
		NewValidateCmd(),
		NewCleanCmd(),

		NewLintCmd(),
		NewFmtCmd(),
		NewDocCmd(),

		NewRestServerCmd(),

		NewListCmd(),
		NewLispAppCmd(),
	}

	if len(os.Args) == 2 && os.Args[1] == "-gen-markdown" {
		md, err := app.ToMarkdown()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(md)
		os.Exit(0)
	}

	app.Run(os.Args)
}
