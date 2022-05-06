// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ = fmt.Sprint

func RunGrpcServer(address string) error {
	grpcServer := grpc.NewServer()
	gpyrpc.RegisterKclvmServiceServer(grpcServer, newKclvmServiceImpl())

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer.Serve(lis)
	return nil
}

type _KclvmServiceImpl struct {
	c *KclvmServiceClient
}

func newKclvmServiceImpl() *_KclvmServiceImpl {
	return &_KclvmServiceImpl{
		c: NewKclvmServiceClient(),
	}
}

func (p *_KclvmServiceImpl) Ping(ctx context.Context, args *gpyrpc.Ping_Args) (*gpyrpc.Ping_Result, error) {
	return p.c.Ping(args)
}
func (p *_KclvmServiceImpl) ParseFile_LarkTree(ctx context.Context, args *gpyrpc.ParseFile_LarkTree_Args) (*gpyrpc.ParseFile_LarkTree_Result, error) {
	return p.c.ParseFile_LarkTree(args)
}
func (p *_KclvmServiceImpl) ParseFile_AST(ctx context.Context, args *gpyrpc.ParseFile_AST_Args) (*gpyrpc.ParseFile_AST_Result, error) {
	return p.c.ParseFile_AST(args)
}
func (p *_KclvmServiceImpl) ParseProgram_AST(ctx context.Context, args *gpyrpc.ParseProgram_AST_Args) (*gpyrpc.ParseProgram_AST_Result, error) {
	return p.c.ParseProgram_AST(args)
}
func (p *_KclvmServiceImpl) ExecProgram(ctx context.Context, args *gpyrpc.ExecProgram_Args) (*gpyrpc.ExecProgram_Result, error) {
	return p.c.ExecProgram(args)
}
func (p *_KclvmServiceImpl) ResetPlugin(ctx context.Context, args *gpyrpc.ResetPlugin_Args) (*gpyrpc.ResetPlugin_Result, error) {
	return p.c.ResetPlugin(args)
}
func (p *_KclvmServiceImpl) FormatCode(ctx context.Context, args *gpyrpc.FormatCode_Args) (*gpyrpc.FormatCode_Result, error) {
	return p.c.FormatCode(args)
}
func (p *_KclvmServiceImpl) FormatPath(ctx context.Context, args *gpyrpc.FormatPath_Args) (*gpyrpc.FormatPath_Result, error) {
	return p.c.FormatPath(args)
}
func (p *_KclvmServiceImpl) LintPath(ctx context.Context, args *gpyrpc.LintPath_Args) (*gpyrpc.LintPath_Result, error) {
	return p.c.LintPath(args)
}
func (p *_KclvmServiceImpl) OverrideFile(ctx context.Context, args *gpyrpc.OverrideFile_Args) (*gpyrpc.OverrideFile_Result, error) {
	return p.c.OverrideFile(args)
}
func (p *_KclvmServiceImpl) EvalCode(ctx context.Context, args *gpyrpc.EvalCode_Args) (*gpyrpc.EvalCode_Result, error) {
	return p.c.EvalCode(args)
}
func (p *_KclvmServiceImpl) ResolveCode(ctx context.Context, args *gpyrpc.ResolveCode_Args) (*gpyrpc.ResolveCode_Result, error) {
	return p.c.ResolveCode(args)
}
func (p *_KclvmServiceImpl) GetSchemaType(ctx context.Context, args *gpyrpc.GetSchemaType_Args) (*gpyrpc.GetSchemaType_Result, error) {
	return p.c.GetSchemaType(args)
}
func (p *_KclvmServiceImpl) ValidateCode(ctx context.Context, args *gpyrpc.ValidateCode_Args) (*gpyrpc.ValidateCode_Result, error) {
	return p.c.ValidateCode(args)
}
func (p *_KclvmServiceImpl) SpliceCode(ctx context.Context, args *gpyrpc.SpliceCode_Args) (*gpyrpc.SpliceCode_Result, error) {
	return p.c.SpliceCode(args)
}
func (p *_KclvmServiceImpl) GoToDef(ctx context.Context, args *gpyrpc.GoToDef_Args) (*gpyrpc.GoToDef_Result, error) {
	return p.c.GoToDef(args)
}
func (p *_KclvmServiceImpl) Complete(ctx context.Context, args *gpyrpc.Complete_Args) (*gpyrpc.Complete_Result, error) {
	return p.c.Complete(args)
}
func (p *_KclvmServiceImpl) DocumentSymbol(ctx context.Context, args *gpyrpc.DocumentSymbol_Args) (*gpyrpc.DocumentSymbol_Result, error) {
	return p.c.DocumentSymbol(args)
}
func (p *_KclvmServiceImpl) Hover(ctx context.Context, args *gpyrpc.Hover_Args) (*gpyrpc.Hover_Result, error) {
	return p.c.Hover(args)
}
func (p *_KclvmServiceImpl) ListDepFiles(ctx context.Context, args *gpyrpc.ListDepFiles_Args) (*gpyrpc.ListDepFiles_Result, error) {
	return p.c.ListDepFiles(args)
}
func (p *_KclvmServiceImpl) LoadSettingsFiles(ctx context.Context, args *gpyrpc.LoadSettingsFiles_Args) (*gpyrpc.LoadSettingsFiles_Result, error) {
	return p.c.LoadSettingsFiles(args)
}
