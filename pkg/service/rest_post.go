// Copyright 2022 The KCL Authors. All rights reserved.

package service

import (
	"context"

	"github.com/golang/protobuf/proto"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ gpyrpc.KclvmServiceServer = (*RestClient)(nil)

type RestClient struct {
	host string
}

func NewRestClient(host string) *RestClient {
	return &RestClient{host: host}
}

func (p *RestClient) httpPost(method string, input, output proto.Message) error {
	var result RestfulResult
	result.Result = output
	return httpPost(p.host+"/api:protorpc/"+method, input, &result)
}

func (p *RestClient) Ping(ctx context.Context, args *gpyrpc.Ping_Args) (reply *gpyrpc.Ping_Result, err error) {
	reply = new(gpyrpc.Ping_Result)
	err = p.httpPost("KclvmService.Ping", args, reply)
	return
}

func (p *RestClient) ListMethod(ctx context.Context, args *gpyrpc.Ping_Args) (reply *gpyrpc.Ping_Result, err error) {
	reply = new(gpyrpc.Ping_Result)
	err = p.httpPost("BuiltinService.Ping", args, reply)
	return
}

func (p *RestClient) ExecProgram(ctx context.Context, args *gpyrpc.ExecProgram_Args) (reply *gpyrpc.ExecProgram_Result, err error) {
	reply = new(gpyrpc.ExecProgram_Result)
	err = p.httpPost("KclvmService.ExecProgram", args, reply)
	return
}

func (p *RestClient) FormatCode(ctx context.Context, args *gpyrpc.FormatCode_Args) (reply *gpyrpc.FormatCode_Result, err error) {
	reply = new(gpyrpc.FormatCode_Result)
	err = p.httpPost("KclvmService.FormatCode", args, reply)
	return
}

func (p *RestClient) FormatPath(ctx context.Context, args *gpyrpc.FormatPath_Args) (reply *gpyrpc.FormatPath_Result, err error) {
	reply = new(gpyrpc.FormatPath_Result)
	err = p.httpPost("KclvmService.FormatPath", args, reply)
	return
}

func (p *RestClient) LintPath(ctx context.Context, args *gpyrpc.LintPath_Args) (reply *gpyrpc.LintPath_Result, err error) {
	reply = new(gpyrpc.LintPath_Result)
	err = p.httpPost("KclvmService.LintPath", args, reply)
	return
}

func (p *RestClient) OverrideFile(ctx context.Context, args *gpyrpc.OverrideFile_Args) (reply *gpyrpc.OverrideFile_Result, err error) {
	reply = new(gpyrpc.OverrideFile_Result)
	err = p.httpPost("KclvmService.OverrideFile", args, reply)
	return
}

func (p *RestClient) EvalCode(ctx context.Context, args *gpyrpc.EvalCode_Args) (reply *gpyrpc.EvalCode_Result, err error) {
	reply = new(gpyrpc.EvalCode_Result)
	err = p.httpPost("KclvmService.EvalCode", args, reply)
	return
}
func (p *RestClient) ResolveCode(ctx context.Context, args *gpyrpc.ResolveCode_Args) (reply *gpyrpc.ResolveCode_Result, err error) {
	reply = new(gpyrpc.ResolveCode_Result)
	err = p.httpPost("KclvmService.ResolveCode", args, reply)
	return
}
func (p *RestClient) GetSchemaType(ctx context.Context, args *gpyrpc.GetSchemaType_Args) (reply *gpyrpc.GetSchemaType_Result, err error) {
	reply = new(gpyrpc.GetSchemaType_Result)
	err = p.httpPost("KclvmService.GetSchemaType", args, reply)
	return
}
func (p *RestClient) ValidateCode(ctx context.Context, args *gpyrpc.ValidateCode_Args) (reply *gpyrpc.ValidateCode_Result, err error) {
	reply = new(gpyrpc.ValidateCode_Result)
	err = p.httpPost("KclvmService.ValidateCode", args, reply)
	return
}
func (p *RestClient) SpliceCode(ctx context.Context, args *gpyrpc.SpliceCode_Args) (reply *gpyrpc.SpliceCode_Result, err error) {
	reply = new(gpyrpc.SpliceCode_Result)
	err = p.httpPost("KclvmService.SpliceCode", args, reply)
	return
}

func (p *RestClient) Complete(ctx context.Context, args *gpyrpc.Complete_Args) (reply *gpyrpc.Complete_Result, err error) {
	reply = new(gpyrpc.Complete_Result)
	err = p.httpPost("KclvmService.Complete", args, reply)
	return
}

func (p *RestClient) GoToDef(ctx context.Context, args *gpyrpc.GoToDef_Args) (reply *gpyrpc.GoToDef_Result, err error) {
	reply = new(gpyrpc.GoToDef_Result)
	err = p.httpPost("KclvmService.GoToDef", args, reply)
	return
}

func (p *RestClient) DocumentSymbol(ctx context.Context, args *gpyrpc.DocumentSymbol_Args) (reply *gpyrpc.DocumentSymbol_Result, err error) {
	reply = new(gpyrpc.DocumentSymbol_Result)
	err = p.httpPost("KclvmService.DocumentSymbol", args, reply)
	return
}

func (p *RestClient) Hover(ctx context.Context, args *gpyrpc.Hover_Args) (reply *gpyrpc.Hover_Result, err error) {
	reply = new(gpyrpc.Hover_Result)
	err = p.httpPost("KclvmService.Hover", args, reply)
	return
}

func (p *RestClient) ListDepFiles(ctx context.Context, args *gpyrpc.ListDepFiles_Args) (reply *gpyrpc.ListDepFiles_Result, err error) {
	reply = new(gpyrpc.ListDepFiles_Result)
	err = p.httpPost("KclvmService.ListDepFiles", args, reply)
	return
}

func (p *RestClient) ListUpStreamFiles(ctx context.Context, args *gpyrpc.ListUpStreamFiles_Args) (reply *gpyrpc.ListUpStreamFiles_Result, err error) {
	reply = new(gpyrpc.ListUpStreamFiles_Result)
	err = p.httpPost("KclvmService.ListUpStreamFiles", args, reply)
	return
}
func (p *RestClient) ListDownStreamFiles(ctx context.Context, args *gpyrpc.ListDownStreamFiles_Args) (reply *gpyrpc.ListDownStreamFiles_Result, err error) {
	reply = new(gpyrpc.ListDownStreamFiles_Result)
	err = p.httpPost("KclvmService.ListDownStreamFiles", args, reply)
	return
}
func (p *RestClient) LoadSettingsFiles(ctx context.Context, args *gpyrpc.LoadSettingsFiles_Args) (reply *gpyrpc.LoadSettingsFiles_Result, err error) {
	reply = new(gpyrpc.LoadSettingsFiles_Result)
	err = p.httpPost("KclvmService.LoadSettingsFiles", args, reply)
	return
}

func (p *RestClient) ParseFile_LarkTree(ctx context.Context, args *gpyrpc.ParseFile_LarkTree_Args) (reply *gpyrpc.ParseFile_LarkTree_Result, err error) {
	reply = new(gpyrpc.ParseFile_LarkTree_Result)
	err = p.httpPost("KclvmService.ParseFile_LarkTree", args, reply)
	return
}
func (p *RestClient) ParseFile_AST(ctx context.Context, args *gpyrpc.ParseFile_AST_Args) (reply *gpyrpc.ParseFile_AST_Result, err error) {
	reply = new(gpyrpc.ParseFile_AST_Result)
	err = p.httpPost("KclvmService.ParseFile_AST", args, reply)
	return
}
func (p *RestClient) ParseProgram_AST(ctx context.Context, args *gpyrpc.ParseProgram_AST_Args) (reply *gpyrpc.ParseProgram_AST_Result, err error) {
	reply = new(gpyrpc.ParseProgram_AST_Result)
	err = p.httpPost("KclvmService.ParseProgram_AST", args, reply)
	return
}
func (p *RestClient) ResetPlugin(ctx context.Context, args *gpyrpc.ResetPlugin_Args) (reply *gpyrpc.ResetPlugin_Result, err error) {
	reply = new(gpyrpc.ResetPlugin_Result)
	err = p.httpPost("KclvmService.ResetPlugin", args, reply)
	return
}
