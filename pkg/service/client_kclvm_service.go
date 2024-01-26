// Copyright The KCL Authors. All rights reserved.

package service

import (
	"fmt"
	"io"
	"net/rpc"

	"kcl-lang.io/kcl-go/pkg/kclvm_runtime"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type KclvmServiceClient struct {
	Runtime   *kclvm_runtime.Runtime
	pyRuntime *kclvm_runtime.Runtime
}

func NewKclvmServiceClient() *KclvmServiceClient {
	c := &KclvmServiceClient{
		Runtime:   kclvm_runtime.GetRuntime(),
		pyRuntime: kclvm_runtime.GetPyRuntime(),
	}
	return c
}

func (p *KclvmServiceClient) getClient(c *rpc.Client) KclvmService {
	return &gpyrpc.PROTORPC_KclvmServiceClient{Client: c}
}
func (p *KclvmServiceClient) wrapErr(err error, stderr io.Reader) error {
	if err != nil {
		err = wrapKclvmServerError(err)
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

func (p *KclvmServiceClient) ExecProgram(args *gpyrpc.ExecProgram_Args) (resp *gpyrpc.ExecProgram_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ExecProgram(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ParseFile(args *gpyrpc.ParseFile_Args) (resp *gpyrpc.ParseFile_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ParseFile(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) ParseProgram(args *gpyrpc.ParseProgram_Args) (resp *gpyrpc.ParseProgram_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ParseProgram(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) LoadPackage(args *gpyrpc.LoadPackage_Args) (resp *gpyrpc.LoadPackage_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).LoadPackage(args)
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
	p.pyRuntime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).OverrideFile(args)
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

func (p *KclvmServiceClient) GetSchemaTypeMapping(args *gpyrpc.GetSchemaTypeMapping_Args) (resp *gpyrpc.GetSchemaTypeMapping_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).GetSchemaTypeMapping(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) GetFullSchemaType(args *gpyrpc.GetFullSchemaType_Args) (resp *gpyrpc.GetSchemaType_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).GetFullSchemaType(args)
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

func (p *KclvmServiceClient) ListDepFiles(args *gpyrpc.ListDepFiles_Args) (resp *gpyrpc.ListDepFiles_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).ListDepFiles(args)
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

func (p *KclvmServiceClient) Rename(args *gpyrpc.Rename_Args) (resp *gpyrpc.Rename_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).Rename(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) RenameCode(args *gpyrpc.RenameCode_Args) (resp *gpyrpc.RenameCode_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).RenameCode(args)
		err = p.wrapErr(err, stderr)
	})
	return
}

func (p *KclvmServiceClient) Test(args *gpyrpc.Test_Args) (resp *gpyrpc.Test_Result, err error) {
	p.Runtime.DoTask(func(c *rpc.Client, stderr io.Reader) {
		resp, err = p.getClient(c).Test(args)
		err = p.wrapErr(err, stderr)
	})
	return
}
