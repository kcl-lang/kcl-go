// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/tools/ktest"
)

func NewTestCmd() *cli.Command {
	return &cli.Command{
		Name:      "test",
		Usage:     "test packages",
		ArgsUsage: "[packages]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "run",
				Usage: "Run only those tests matching the regular expression.",
			},
			&cli.IntFlag{
				Name:    "max-proc",
				Aliases: []string{"n"},
				Usage:   "set max kclvm process",
				Value:   1,
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Set quiet mode",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Log all tests as they are run",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Run in debug mode (for developers only)",
			},
		},
		Action: func(c *cli.Context) error {
			kclvm_runtime.InitRuntime(c.Int("max-proc"))
			return ktest.RunTest(c.Args().First(), ktest.Options{
				RunRegexp: c.String("run"),
				QuietMode: c.Bool("quiet"),
				Verbose:   c.Bool("verbose"),
				Debug:     c.Bool("debug"),
			})
		},
	}
}
