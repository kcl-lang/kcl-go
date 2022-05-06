// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
)

func NewKclCmd() *cli.Command {
	return &cli.Command{
		Hidden:          false,
		SkipFlagParsing: true,
		Name:            "kcl",
		Usage:           "kcl command",
		ArgsUsage:       "[-flags] [kfiles...]",
		Action: func(c *cli.Context) error {
			args := append([]string{"-m", "kclvm"}, c.Args().Slice()...)
			cmd := exec.Command(kclvm_runtime.MustGetKclvmPath(), args...)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Print(string(stdoutStderr))
				if c.Args().Len() > 0 {
					fmt.Println("ERR:", err)
					os.Exit(1)
				} else {
					os.Exit(0)
				}
			}
			fmt.Print(string(stdoutStderr))
			return nil
		},
	}
}
