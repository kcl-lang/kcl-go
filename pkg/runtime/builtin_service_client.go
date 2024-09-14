//go:build rpc || !cgo
// +build rpc !cgo

// Copyright The KCL Authors. All rights reserved.

package runtime

import (
	"fmt"
	"io"
	"net/rpc"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type BuiltinServiceClient struct {
	*Runtime
}

func (p *BuiltinServiceClient) getClient(c *rpc.Client) *gpyrpc.PROTORPC_BuiltinServiceClient {
	return &gpyrpc.PROTORPC_BuiltinServiceClient{Client: c}
}
func (p *BuiltinServiceClient) wrapErr(err error, stderr io.Reader) error {
	if err != nil {
		if data, _ := io.ReadAll(stderr); len(data) != 0 {
			return fmt.Errorf("%w: stderr = %s", err, string(data))
		}
	}
	return err
}

func (p *BuiltinServiceClient) Ping(args *gpyrpc.Ping_Args) (resp *gpyrpc.Ping_Result, err error) {
	p.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).Ping(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *BuiltinServiceClient) ListMethod(args *gpyrpc.ListMethod_Args) (resp *gpyrpc.ListMethod_Result, err error) {
	p.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ListMethod(args)
		err = p.wrapErr(err, stderr)
	})
	return
}
