// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"fmt"
	"io"
	"net/rpc"
	"os"
	"strings"

	"kusionstack.io/kclvm-go/pkg/kclvm_runtime"
	capicall "kusionstack.io/kclvm-go/pkg/service/c_api_call"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ KclvmService = (*capicall.PROTOCAPI_KclvmServiceClient)(nil)

var Default_IsNative = false

type KclvmServiceClient struct {
	Runtime  *kclvm_runtime.Runtime
	IsNative bool //if true ,call service by C API
}

func NewKclvmServiceClient() *KclvmServiceClient {
	c := &KclvmServiceClient{
		Runtime: kclvm_runtime.GetRuntime(),
	}
	if Default_IsNative || strings.EqualFold(os.Getenv("KCLVM_SERVICE_CLIENT_HANDLER"), "native") {
		c.IsNative = true
	}
	return c
}

func (p *KclvmServiceClient) getClient(c *rpc.Client) KclvmService {
	if p.IsNative {
		return capicall.PROTOCAPI_NewKclvmServiceClient()
	}
	return &gpyrpc.PROTORPC_KclvmServiceClient{Client: c}
}
func (p *KclvmServiceClient) wrapErr(err error, stderr io.Reader) error {
	if err != nil {
		if data, _ := io.ReadAll(stderr); len(data) != 0 {
			return fmt.Errorf("%w: stderr = %s", err, string(data))
		}
	}
	return err
}

func (p *KclvmServiceClient) Ping(args *gpyrpc.Ping_Args) (resp *gpyrpc.Ping_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).Ping(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ParseFile_LarkTree(args *gpyrpc.ParseFile_LarkTree_Args) (resp *gpyrpc.ParseFile_LarkTree_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ParseFile_LarkTree(args)
		err = p.wrapErr(err, stderr)
	})
	return
}
func (p *KclvmServiceClient) ParseFile_AST(args *gpyrpc.ParseFile_AST_Args) (resp *gpyrpc.ParseFile_AST_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ParseFile_AST(args)
		err = p.wrapErr(err, stderr)
	})
	return
}
func (p *KclvmServiceClient) ParseProgram_AST(args *gpyrpc.ParseProgram_AST_Args) (resp *gpyrpc.ParseProgram_AST_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ParseProgram_AST(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ExecProgram(args *gpyrpc.ExecProgram_Args) (resp *gpyrpc.ExecProgram_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ExecProgram(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ResetPlugin(args *gpyrpc.ResetPlugin_Args) (resp *gpyrpc.ResetPlugin_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ResetPlugin(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) FormatCode(args *gpyrpc.FormatCode_Args) (resp *gpyrpc.FormatCode_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).FormatCode(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) FormatPath(args *gpyrpc.FormatPath_Args) (resp *gpyrpc.FormatPath_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).FormatPath(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) LintPath(args *gpyrpc.LintPath_Args) (resp *gpyrpc.LintPath_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).LintPath(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) OverrideFile(args *gpyrpc.OverrideFile_Args) (resp *gpyrpc.OverrideFile_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).OverrideFile(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) EvalCode(args *gpyrpc.EvalCode_Args) (resp *gpyrpc.EvalCode_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).EvalCode(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ResolveCode(args *gpyrpc.ResolveCode_Args) (resp *gpyrpc.ResolveCode_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ResolveCode(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) GetSchemaType(args *gpyrpc.GetSchemaType_Args) (resp *gpyrpc.GetSchemaType_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).GetSchemaType(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ValidateCode(args *gpyrpc.ValidateCode_Args) (resp *gpyrpc.ValidateCode_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ValidateCode(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) SpliceCode(args *gpyrpc.SpliceCode_Args) (resp *gpyrpc.SpliceCode_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).SpliceCode(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) Complete(args *gpyrpc.Complete_Args) (resp *gpyrpc.Complete_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).Complete(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) GoToDef(args *gpyrpc.GoToDef_Args) (resp *gpyrpc.GoToDef_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).GoToDef(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) DocumentSymbol(args *gpyrpc.DocumentSymbol_Args) (resp *gpyrpc.DocumentSymbol_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).DocumentSymbol(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) Hover(args *gpyrpc.Hover_Args) (resp *gpyrpc.Hover_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).Hover(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ListDepFiles(args *gpyrpc.ListDepFiles_Args) (resp *gpyrpc.ListDepFiles_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ListDepFiles(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ListUpStreamFiles(args *gpyrpc.ListUpStreamFiles_Args) (resp *gpyrpc.ListUpStreamFiles_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ListUpStreamFiles(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ListDownStreamFiles(args *gpyrpc.ListDownStreamFiles_Args) (resp *gpyrpc.ListDownStreamFiles_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ListDownStreamFiles(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) LoadSettingsFiles(args *gpyrpc.LoadSettingsFiles_Args) (resp *gpyrpc.LoadSettingsFiles_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).LoadSettingsFiles(args)
		err = p.wrapErr(err, stderr)
	})
	return
}
