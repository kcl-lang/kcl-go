package service

import "kusionstack.io/kclvm-go/pkg/spec/gpyrpc"

type KclvmService interface {
	Ping(in *gpyrpc.Ping_Args) (out *gpyrpc.Ping_Result, err error)
	ParseFile_LarkTree(in *gpyrpc.ParseFile_LarkTree_Args) (out *gpyrpc.ParseFile_LarkTree_Result, err error)
	ParseFile_AST(in *gpyrpc.ParseFile_AST_Args) (out *gpyrpc.ParseFile_AST_Result, err error)
	ParseProgram_AST(in *gpyrpc.ParseProgram_AST_Args) (out *gpyrpc.ParseProgram_AST_Result, err error)
	ExecProgram(in *gpyrpc.ExecProgram_Args) (out *gpyrpc.ExecProgram_Result, err error)
	ResetPlugin(in *gpyrpc.ResetPlugin_Args) (out *gpyrpc.ResetPlugin_Result, err error)
	FormatCode(in *gpyrpc.FormatCode_Args) (out *gpyrpc.FormatCode_Result, err error)
	FormatPath(in *gpyrpc.FormatPath_Args) (out *gpyrpc.FormatPath_Result, err error)
	LintPath(in *gpyrpc.LintPath_Args) (out *gpyrpc.LintPath_Result, err error)
	OverrideFile(in *gpyrpc.OverrideFile_Args) (out *gpyrpc.OverrideFile_Result, err error)
	EvalCode(in *gpyrpc.EvalCode_Args) (out *gpyrpc.EvalCode_Result, err error)
	ResolveCode(in *gpyrpc.ResolveCode_Args) (out *gpyrpc.ResolveCode_Result, err error)
	GetSchemaType(in *gpyrpc.GetSchemaType_Args) (out *gpyrpc.GetSchemaType_Result, err error)
	ValidateCode(in *gpyrpc.ValidateCode_Args) (out *gpyrpc.ValidateCode_Result, err error)
	SpliceCode(in *gpyrpc.SpliceCode_Args) (out *gpyrpc.SpliceCode_Result, err error)
	Complete(in *gpyrpc.Complete_Args) (out *gpyrpc.Complete_Result, err error)
	GoToDef(in *gpyrpc.GoToDef_Args) (out *gpyrpc.GoToDef_Result, err error)
	DocumentSymbol(in *gpyrpc.DocumentSymbol_Args) (out *gpyrpc.DocumentSymbol_Result, err error)
	Hover(in *gpyrpc.Hover_Args) (out *gpyrpc.Hover_Result, err error)
	ListDepFiles(in *gpyrpc.ListDepFiles_Args) (out *gpyrpc.ListDepFiles_Result, err error)
	ListUpStreamFiles(in *gpyrpc.ListUpStreamFiles_Args) (out *gpyrpc.ListUpStreamFiles_Result, err error)
	ListDownStreamFiles(in *gpyrpc.ListDownStreamFiles_Args) (out *gpyrpc.ListDownStreamFiles_Result, err error)
	LoadSettingsFiles(in *gpyrpc.LoadSettingsFiles_Args) (out *gpyrpc.LoadSettingsFiles_Result, err error)
}
