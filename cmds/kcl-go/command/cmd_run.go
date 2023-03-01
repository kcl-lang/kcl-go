// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go"
	"kusionstack.io/kclvm-go/pkg/kcl"
)

// keep same as kcl command
var runRunFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "argument",
		Aliases: []string{"D"},
		Usage:   "Specify the top-level argument",
	},
	&cli.StringSliceFlag{
		Name:    "overrides",
		Aliases: []string{"O"},
		Usage:   "Specify the configuration override path and value",
	},

	&cli.StringFlag{
		Name:    "setting",
		Aliases: []string{"Y"},
		Usage:   "Specify the command line setting file",
	},
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "Specify the output file",
	},

	&cli.BoolFlag{
		Name:    "disable-none",
		Aliases: []string{"n"},
		Usage:   "Disable dumping None values",
	},
	&cli.BoolFlag{
		Name:  "sort-keys",
		Usage: "Sort result keys",
	},
	&cli.BoolFlag{
		Name:    "strict-range-check",
		Aliases: []string{"r"},
		Usage:   "Do perform strict numeric range check",
	},
	&cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "Run in debug mode (for developers only)",
	},
	&cli.StringFlag{
		Name:  "output-type",
		Value: "yaml",
		Usage: "set output type (json|yaml)",
	},
	&cli.IntFlag{
		Name:  "max-proc",
		Usage: "set max kclvm process",
		Value: 1,
	},
}

func NewRunCmd() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "compile and run KCL program",
		ArgsUsage: "kfiles...",
		Flags:     runRunFlags,
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			kclvm.InitKclvmRuntime(c.Int("max-proc"))

			if c.Bool("debug") {
				opt, err := kcl.ParseArgs(
					c.Args().Slice(),
					kcl.WithOptions(c.StringSlice("argument")...),
					kcl.WithOverrides(c.StringSlice("overrides")...),
					kcl.WithSettings(c.String("setting")),
				)
				if err != nil {
					fmt.Print(err)
					os.Exit(1)
				}

				fmt.Println("======== args begin ========")
				fmt.Println(opt.JSONString())
				fmt.Println("======== args end ========")
			}

			start := time.Now()
			result, err := kcl.RunFiles(
				c.Args().Slice(),
				kcl.WithOptions(c.StringSlice("argument")...),
				kcl.WithOverrides(c.StringSlice("overrides")...),
				kcl.WithSettings(c.String("setting")),
				kcl.WithSortKeys(c.Bool("sort-keys")),
			)

			if c.Bool("debug") {
				fmt.Println("======== EscapedTime begin ========")
				fmt.Println("Python:", result.GetPyEscapedTime())
				fmt.Println("Golang:", time.Since(start).Seconds())
				fmt.Println("======== EscapedTime end ========")
			}

			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}

			switch strings.ToLower(c.String("output-type")) {
			case "json":
				fmt.Println(result.GetRawJsonResult())

			case "yaml":
				fmt.Println(result.GetRawYamlResult())

			default:
				fmt.Println(result.GetRawYamlResult())

			}
			return nil
		},
	}
}
