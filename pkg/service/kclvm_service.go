package service

import "kcl-lang.io/kcl-go/pkg/spec/gpyrpc"

type KclvmService interface {
	Ping(in *gpyrpc.Ping_Args) (out *gpyrpc.Ping_Result, err error)
	ExecProgram(in *gpyrpc.ExecProgram_Args) (out *gpyrpc.ExecProgram_Result, err error)
	ParseFile(in *gpyrpc.ParseFile_Args) (out *gpyrpc.ParseFile_Result, err error)
	ParseProgram(in *gpyrpc.ParseProgram_Args) (out *gpyrpc.ParseProgram_Result, err error)
	LoadPackage(in *gpyrpc.LoadPackage_Args) (out *gpyrpc.LoadPackage_Result, err error)
	FormatCode(in *gpyrpc.FormatCode_Args) (out *gpyrpc.FormatCode_Result, err error)
	FormatPath(in *gpyrpc.FormatPath_Args) (out *gpyrpc.FormatPath_Result, err error)
	LintPath(in *gpyrpc.LintPath_Args) (out *gpyrpc.LintPath_Result, err error)
	OverrideFile(in *gpyrpc.OverrideFile_Args) (out *gpyrpc.OverrideFile_Result, err error)
	GetSchemaType(in *gpyrpc.GetSchemaType_Args) (out *gpyrpc.GetSchemaType_Result, err error)
	GetSchemaTypeMapping(in *gpyrpc.GetSchemaTypeMapping_Args) (out *gpyrpc.GetSchemaTypeMapping_Result, err error)
	GetFullSchemaType(in *gpyrpc.GetFullSchemaType_Args) (out *gpyrpc.GetSchemaType_Result, err error)
	ValidateCode(in *gpyrpc.ValidateCode_Args) (out *gpyrpc.ValidateCode_Result, err error)
	ListDepFiles(in *gpyrpc.ListDepFiles_Args) (out *gpyrpc.ListDepFiles_Result, err error)
	LoadSettingsFiles(in *gpyrpc.LoadSettingsFiles_Args) (out *gpyrpc.LoadSettingsFiles_Result, err error)
	Rename(in *gpyrpc.Rename_Args) (out *gpyrpc.Rename_Result, err error)
	RenameCode(in *gpyrpc.RenameCode_Args) (out *gpyrpc.RenameCode_Result, err error)
	Test(in *gpyrpc.Test_Args) (out *gpyrpc.Test_Result, err error)
}
