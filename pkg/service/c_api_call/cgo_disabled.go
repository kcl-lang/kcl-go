//go:build !cgo || !kclvm_service_capi
// +build !cgo !kclvm_service_capi

package capicall

import "kusionstack.io/kclvm-go/pkg/spec/gpyrpc"

type PROTOCAPI_KclvmServiceClient struct {
}

func PROTOCAPI_NewKclvmServiceClient() *PROTOCAPI_KclvmServiceClient {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) Ping(in *gpyrpc.Ping_Args) (out *gpyrpc.Ping_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ParseFile_LarkTree(in *gpyrpc.ParseFile_LarkTree_Args) (out *gpyrpc.ParseFile_LarkTree_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ParseFile_AST(in *gpyrpc.ParseFile_AST_Args) (out *gpyrpc.ParseFile_AST_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ParseProgram_AST(in *gpyrpc.ParseProgram_AST_Args) (out *gpyrpc.ParseProgram_AST_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ExecProgram(in *gpyrpc.ExecProgram_Args) (out *gpyrpc.ExecProgram_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ResetPlugin(in *gpyrpc.ResetPlugin_Args) (out *gpyrpc.ResetPlugin_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) FormatCode(in *gpyrpc.FormatCode_Args) (out *gpyrpc.FormatCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) FormatPath(in *gpyrpc.FormatPath_Args) (out *gpyrpc.FormatPath_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) LintPath(in *gpyrpc.LintPath_Args) (out *gpyrpc.LintPath_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) OverrideFile(in *gpyrpc.OverrideFile_Args) (out *gpyrpc.OverrideFile_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) EvalCode(in *gpyrpc.EvalCode_Args) (out *gpyrpc.EvalCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ResolveCode(in *gpyrpc.ResolveCode_Args) (out *gpyrpc.ResolveCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) GetSchemaType(in *gpyrpc.GetSchemaType_Args) (out *gpyrpc.GetSchemaType_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ValidateCode(in *gpyrpc.ValidateCode_Args) (out *gpyrpc.ValidateCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) SpliceCode(in *gpyrpc.SpliceCode_Args) (out *gpyrpc.SpliceCode_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) Complete(in *gpyrpc.Complete_Args) (out *gpyrpc.Complete_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) GoToDef(in *gpyrpc.GoToDef_Args) (out *gpyrpc.GoToDef_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) DocumentSymbol(in *gpyrpc.DocumentSymbol_Args) (out *gpyrpc.DocumentSymbol_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) Hover(in *gpyrpc.Hover_Args) (out *gpyrpc.Hover_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ListDepFiles(in *gpyrpc.ListDepFiles_Args) (out *gpyrpc.ListDepFiles_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ListUpStreamFiles(in *gpyrpc.ListUpStreamFiles_Args) (out *gpyrpc.ListUpStreamFiles_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) ListDownStreamFiles(in *gpyrpc.ListDownStreamFiles_Args) (out *gpyrpc.ListDownStreamFiles_Result, err error) {
	panic("unsupport cgo")
}

func (c *PROTOCAPI_KclvmServiceClient) LoadSettingsFiles(in *gpyrpc.LoadSettingsFiles_Args) (out *gpyrpc.LoadSettingsFiles_Result, err error) {
	panic("unsupport cgo")
}
