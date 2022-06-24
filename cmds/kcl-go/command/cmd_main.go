// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/logger"
)

func Main() {
	app := cli.NewApp()
	app.Name = "kcl-go"
	app.Usage = "K Configuration Language Virtual Machine"
	app.Version = "0.0.1"

	// kclvm -m kclvm
	// kclvm -m kclvm.tools.plugin
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
		newSetpupKclvmCmd(),

		newBuildInfoCmd(),

		newKclCmd(), // go run main.go kcl -D aa=11 -D bb=22 main.k other.k

		newRunCmd(),
		newValidateCmd(),
		newTestCmd(),
		newPluginCmd(),
		newCleanCmd(),

		newLintCmd(),
		newFmtCmd(),
		newDocCmd(),

		newGrpcServerCmd(),
		newRestServerCmd(),

		newAstCmd(),

		newListCmd(),
		newLispAppCmd(),
		newTourCmd(),
		newLspCmd(),
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
