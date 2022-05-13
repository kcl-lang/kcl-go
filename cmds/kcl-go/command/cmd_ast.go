// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var cmdAstFlags = []cli.Flag{}

func NewAstCmd() *cli.Command {
	return &cli.Command{
		Hidden:    false,
		Name:      "dev-ast",
		Usage:     "parse ast tree",
		ArgsUsage: "file.k",
		Flags:     cmdAstFlags,
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			client := service.NewKclvmServiceClient()
			resp, err := client.ParseFile_AST(&gpyrpc.ParseFile_AST_Args{
				Filename:   c.Args().First(),
				SourceCode: "",
			})
			if err != nil {
				return err
			}

			fmt.Println(resp.AstJson)
			return nil
		},
	}
}
