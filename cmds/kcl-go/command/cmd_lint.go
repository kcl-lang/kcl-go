// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	kcl "kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
)

func NewLintCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "lint",
		Usage:  "lints the KCL source files named on its command line.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "max-proc",
				Aliases: []string{"n"},
				Usage:   "set max kclvm process",
				Value:   1,
			},
		},
		Action: func(c *cli.Context) error {
			kclvm_runtime.InitRuntime(c.Int("max-proc"))
			results, err := kcl.LintPath(c.Args().Slice())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			wd, _ := os.Getwd()
			for _, s := range results {
				ss := strings.Split(s, " ")
				if len(ss) > 0 && filepath.IsAbs(ss[0]) && wd != "" {
					if x, err := filepath.Rel(wd, ss[0]); err == nil {
						ss[0] = "./" + x
					}
				}
				fmt.Println(strings.Join(ss, " "))
			}
			return nil
		},
	}
}
