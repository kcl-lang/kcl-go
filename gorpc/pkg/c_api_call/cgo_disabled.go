//go:build !cgo
// +build !cgo

package capicall

import (
	"context"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

type PROTOCAPI_KclvmServiceClient struct {
}

func PROTOCAPI_NewKclvmServiceClient() *PROTOCAPI_KclvmServiceClient {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) Ping(ctx context.Context, in *gpyrpc.Ping_Args) (out *gpyrpc.Ping_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ParseFile_LarkTree(ctx context.Context, in *gpyrpc.ParseFile_LarkTree_Args) (out *gpyrpc.ParseFile_LarkTree_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ParseFile_AST(ctx context.Context, in *gpyrpc.ParseFile_AST_Args) (out *gpyrpc.ParseFile_AST_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ParseProgram_AST(ctx context.Context, in *gpyrpc.ParseProgram_AST_Args) (out *gpyrpc.ParseProgram_AST_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ExecProgram(ctx context.Context, in *gpyrpc.ExecProgram_Args) (out *gpyrpc.ExecProgram_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ResetPlugin(ctx context.Context, in *gpyrpc.ResetPlugin_Args) (out *gpyrpc.ResetPlugin_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) FormatCode(ctx context.Context, in *gpyrpc.FormatCode_Args) (out *gpyrpc.FormatCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) FormatPath(ctx context.Context, in *gpyrpc.FormatPath_Args) (out *gpyrpc.FormatPath_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) LintPath(ctx context.Context, in *gpyrpc.LintPath_Args) (out *gpyrpc.LintPath_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) OverrideFile(ctx context.Context, in *gpyrpc.OverrideFile_Args) (out *gpyrpc.OverrideFile_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) EvalCode(ctx context.Context, in *gpyrpc.EvalCode_Args) (out *gpyrpc.EvalCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ResolveCode(ctx context.Context, in *gpyrpc.ResolveCode_Args) (out *gpyrpc.ResolveCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) GetSchemaType(ctx context.Context, in *gpyrpc.GetSchemaType_Args) (out *gpyrpc.GetSchemaType_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ValidateCode(ctx context.Context, in *gpyrpc.ValidateCode_Args) (out *gpyrpc.ValidateCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) SpliceCode(ctx context.Context, in *gpyrpc.SpliceCode_Args) (out *gpyrpc.SpliceCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) Complete(ctx context.Context, in *gpyrpc.Complete_Args) (out *gpyrpc.Complete_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) GoToDef(ctx context.Context, in *gpyrpc.GoToDef_Args) (out *gpyrpc.GoToDef_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) DocumentSymbol(ctx context.Context, in *gpyrpc.DocumentSymbol_Args) (out *gpyrpc.DocumentSymbol_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) Hover(ctx context.Context, in *gpyrpc.Hover_Args) (out *gpyrpc.Hover_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ListDepFiles(ctx context.Context, in *gpyrpc.ListDepFiles_Args) (out *gpyrpc.ListDepFiles_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ListUpStreamFiles(ctx context.Context, in *gpyrpc.ListUpStreamFiles_Args) (out *gpyrpc.ListUpStreamFiles_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ListDownStreamFiles(ctx context.Context, in *gpyrpc.ListDownStreamFiles_Args) (out *gpyrpc.ListDownStreamFiles_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) LoadSettingsFiles(ctx context.Context, in *gpyrpc.LoadSettingsFiles_Args) (out *gpyrpc.LoadSettingsFiles_Result, err error) {
	panic("unsupport cgo")
}
