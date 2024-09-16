//go:build rpc || !cgo
// +build rpc !cgo

// Copyright The KCL Authors. All rights reserved.

package service

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
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
		c: newKclvmServiceClient(),
	}
}

func (p *_KclvmServiceImpl) Ping(ctx context.Context, args *gpyrpc.Ping_Args) (*gpyrpc.Ping_Result, error) {
	return p.c.Ping(args)
}
func (p *_KclvmServiceImpl) ExecProgram(ctx context.Context, args *gpyrpc.ExecProgram_Args) (*gpyrpc.ExecProgram_Result, error) {
	return p.c.ExecProgram(args)
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (p *_KclvmServiceImpl) BuildProgram(ctx context.Context, args *gpyrpc.BuildProgram_Args) (*gpyrpc.BuildProgram_Result, error) {
	return p.c.BuildProgram(args)
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (p *_KclvmServiceImpl) ExecArtifact(ctx context.Context, args *gpyrpc.ExecArtifact_Args) (*gpyrpc.ExecProgram_Result, error) {
	return p.c.ExecArtifact(args)
}
func (p *_KclvmServiceImpl) ParseFile(ctx context.Context, args *gpyrpc.ParseFile_Args) (*gpyrpc.ParseFile_Result, error) {
	return p.c.ParseFile(args)
}
func (p *_KclvmServiceImpl) ParseProgram(ctx context.Context, args *gpyrpc.ParseProgram_Args) (*gpyrpc.ParseProgram_Result, error) {
	return p.c.ParseProgram(args)
}
func (p *_KclvmServiceImpl) ListOptions(ctx context.Context, args *gpyrpc.ParseProgram_Args) (*gpyrpc.ListOptions_Result, error) {
	return p.c.ListOptions(args)
}
func (p *_KclvmServiceImpl) ListVariables(ctx context.Context, args *gpyrpc.ListVariables_Args) (*gpyrpc.ListVariables_Result, error) {
	return p.c.ListVariables(args)
}
func (p *_KclvmServiceImpl) LoadPackage(ctx context.Context, args *gpyrpc.LoadPackage_Args) (*gpyrpc.LoadPackage_Result, error) {
	return p.c.LoadPackage(args)
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
func (p *_KclvmServiceImpl) GetSchemaTypeMapping(ctx context.Context, args *gpyrpc.GetSchemaTypeMapping_Args) (*gpyrpc.GetSchemaTypeMapping_Result, error) {
	return p.c.GetSchemaTypeMapping(args)
}
func (p *_KclvmServiceImpl) ValidateCode(ctx context.Context, args *gpyrpc.ValidateCode_Args) (*gpyrpc.ValidateCode_Result, error) {
	return p.c.ValidateCode(args)
}
func (p *_KclvmServiceImpl) ListDepFiles(ctx context.Context, args *gpyrpc.ListDepFiles_Args) (*gpyrpc.ListDepFiles_Result, error) {
	return p.c.ListDepFiles(args)
}
func (p *_KclvmServiceImpl) LoadSettingsFiles(ctx context.Context, args *gpyrpc.LoadSettingsFiles_Args) (*gpyrpc.LoadSettingsFiles_Result, error) {
	return p.c.LoadSettingsFiles(args)
}
func (p *_KclvmServiceImpl) Rename(ctx context.Context, args *gpyrpc.Rename_Args) (*gpyrpc.Rename_Result, error) {
	return p.c.Rename(args)
}
func (p *_KclvmServiceImpl) RenameCode(ctx context.Context, args *gpyrpc.RenameCode_Args) (*gpyrpc.RenameCode_Result, error) {
	return p.c.RenameCode(args)
}
func (p *_KclvmServiceImpl) Test(ctx context.Context, args *gpyrpc.Test_Args) (*gpyrpc.Test_Result, error) {
	return p.c.Test(args)
}
func (p *_KclvmServiceImpl) UpdateDependencies(ctx context.Context, args *gpyrpc.UpdateDependencies_Args) (*gpyrpc.UpdateDependencies_Result, error) {
	return p.c.UpdateDependencies(args)
}
func (p *_KclvmServiceImpl) GetVersion(ctx context.Context, args *gpyrpc.GetVersion_Args) (*gpyrpc.GetVersion_Result, error) {
	return p.c.GetVersion(args)
}
