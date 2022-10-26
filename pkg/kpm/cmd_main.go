package kpm

import (
	"github.com/urfave/cli/v2"
)

func CLI(args ...string) error {
	app := cli.NewApp()
	app.Name = "kpm"
	app.Usage = "kpm is a kcl package manager"
	app.Version = "v0.0.1-alpha.1"
	app.UsageText = CliHelp
	app.Commands = []*cli.Command{
		NewInitCmd(),
		NewAddCmd(),
		NewDelCmd(),
		NewDownloadCmd(),
		NewStoreCmd(),
	}
	err := Setup()
	if err != nil {
		return err
	}
	//添加一个参数确保与os.Args数量一致
	nargs := make([]string, len(args))
	nargs = nargs[:1]
	nargs = append(nargs, args...)
	err = app.Run(nargs)
	if err != nil {
		return err
	}
	return nil
}

type KpmFile struct {
	//包名，确定包的命名空间
	PackageName string `json:"package_name"`
	//确定此包的kcl最低运行版本
	KclvmMinVersion string `json:"kclvm_min_version"`
	//直接依赖，别名不重复
	Direct []Require `json:"direct,omitempty"`
	//间接依赖，不看别名，包名版本唯一即可
	Indirect []Require `json:"indirect,omitempty"`
}

type Require struct {
	//别名
	Alias string `json:"alias,omitempty"`
	//包名，确定包的命名空间
	Name string `json:"name,omitempty"`
	//确定此包的版本
	Version string `json:"version,omitempty"`
	//校验和 sha512
	Integrity string `json:"integrity"`
	//包类型 git，registry
	Type string `json:"type"`
	//git包地址
	GitAddress string `json:"git_address,omitempty"`
	//git包commit id
	GitCommit string `json:"git_commit,omitempty"`
}
