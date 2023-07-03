// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
)

func NewDocCmd() *cli.Command {
	return &cli.Command{
		Hidden:          false,
		SkipFlagParsing: true,
		Name:            "doc",
		Usage:           "show documentation for package or symbol",
		Action: func(c *cli.Context) error {
			args := append([]string{"-m", "kclvm.tools.docs"}, c.Args().Slice()...)
			cmd := exec.Command(kclvm_runtime.MustGetKclvmPath(), args...)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Print(string(stdoutStderr))
				fmt.Println("ERR:", err)
				os.Exit(1)
			}
			fmt.Print(string(stdoutStderr))
			return nil
		},
	}
}
