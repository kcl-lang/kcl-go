// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"github.com/urfave/cli/v2"

	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
	"kcl-lang.io/kcl-go/pkg/service"
)

var cmdGrpcServerFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "http",
		Usage: "set listen address",
		Value: ":2021",
	},
	&cli.IntFlag{
		Name:    "max-proc",
		Aliases: []string{"n"},
		Usage:   "set max kclvm process",
		Value:   1,
	},
}

func NewGrpcServerCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "grpc-server",
		Usage:  "run grpc server",
		Flags:  cmdGrpcServerFlags,
		Action: func(c *cli.Context) error {
			kclvm_runtime.InitRuntime(c.Int("max-proc"))
			return service.RunGrpcServer(c.String("http"))
		},
	}
}
