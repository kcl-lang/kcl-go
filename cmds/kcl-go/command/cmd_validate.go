// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	"kusionstack.io/kclvm-go/pkg/tools/validate"
)

var cmdValidateFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "data",
		Usage: "A JSON or YAML data string.",
	},
	&cli.StringFlag{
		Name:  "code",
		Usage: "A KCL code string.",
	},
	&cli.StringFlag{
		Name:  "schema",
		Usage: "The schema name required for verification.",
	},
	&cli.StringFlag{
		Name:  "attribute_name",
		Usage: "The validation attribute name, default is `value`.",
	},
	&cli.StringFlag{
		Name:  "format",
		Usage: "The data format, suppored json, JSON, yaml and YAML.",
	},
	&cli.IntFlag{
		Name:    "max-proc",
		Aliases: []string{"n"},
		Usage:   "set max kclvm process",
		Value:   1,
	},
}

func newValidateCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "vet",
		Usage:  "validate data match code.",
		Flags:  cmdValidateFlags,
		Action: func(c *cli.Context) error {
			data := c.String("data")
			code := c.String("code")

			kclvm_runtime.InitRuntime(c.Int("max-proc"))

			ok, err := validate.ValidateCode(data, code, &validate.ValidateOptions{
				Schema:        c.String("schema"),
				AttributeName: c.String("attribute_name"),
				Format:        c.String("format"),
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(ok)
			return nil
		},
	}
}
