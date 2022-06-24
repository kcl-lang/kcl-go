// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/langserver"
)

const name = "kcl-language-server"
const version = "0.0.1"

var lspFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "log-file",
		Usage: "Specify the filepath to redirect logs to instead of writing to stdout/stderr",
	},
	&cli.IntFlag{
		Name:  "log-level",
		Value: 1,
		Usage: "Specify the log level. Default to 1. Valid log levels: 1 for info level; 6 for debug level.",
	},
	&cli.BoolFlag{
		Name:  "version",
		Value: false,
		Usage: "Print the kcl-language server version",
	},
	&cli.BoolFlag{
		Name:  "quiet",
		Value: false,
		Usage: "Run quieter",
	},
	&cli.BoolFlag{
		Name:  "tcp",
		Value: false,
		Usage: "Use TCP server instead of stdio",
	},
	&cli.StringFlag{
		Name:  "host",
		Value: "127.0.0.1",
		Usage: "Specify the host for TCP server",
	},
	&cli.StringFlag{
		Name:  "port",
		Value: "2088",
		Usage: "Specify the port for TCP server",
	},
}

func newLspCmd() *cli.Command {
	return &cli.Command{
		Hidden: true,
		Name:   "lsp",
		Usage:  "LSP server for KCL",
		Flags:  lspFlags,
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				fmt.Printf("%s %s\n", name, version)
				os.Exit(0)
			}

			if c.Bool("tcp") {
				fmt.Println("tcp not supported")
				os.Exit(1)
			}
			config := &langserver.Config{
				LogFile:  c.String("log-file"),
				LogLevel: c.Int("log-level"),
				Quiet:    c.Bool("quiet"),
				Channel:  stdrwc{},
			}
			langserver.Run(config)
			log.Println("kcl language server: connections closed")
			return nil
		},
	}
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (c stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (c stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
