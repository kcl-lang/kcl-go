// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
)

func newPluginCmd() *cli.Command {
	return &cli.Command{
		SkipFlagParsing: true,
		Name:            "plugin",
		Usage:           "plugin tool",
		Action: func(c *cli.Context) error {
			args := append([]string{"-m", "kclvm.tools.plugin"}, c.Args().Slice()...)
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
