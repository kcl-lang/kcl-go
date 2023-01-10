package kpm

import (
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/orangebees/go-oneutils/Semver"
)

type Require struct {
	RequireBase
	//别名
	Alias string `json:"alias,omitempty"`
}

type RequireBase struct {
	////包类型 git，registry

	RequirePkgStruct
	//校验和 sha512
	Integrity GlobalStore.Integrity `json:"integrity"`

	//git:github.com/a/b@v0.0.1
	//git:github.com/a/b@v0.0.0#asdfghjkl
	//registry:github.com/a/b@v0.0.1
}

type RequirePkgStruct struct {
	//包类型 git，registry
	Type string `json:"type"`
	//包名，确定包的命名空间
	Name string `json:"name"`
	//确定此包的版本
	Version Semver.VersionString `json:"version"`

	//git:github.com/a/b@v0.0.1
	//git:github.com/a/b@v0.0.0#asdfghjkl
	//reg:github.com/a/b@v0.0.1
}

func (rps *RequirePkgStruct) GetPkgString() PkgString {
	return rps.Type + ":" + rps.Name + "@" + string(rps.Version)
}
func (rps *RequirePkgStruct) GetShortName() string {
	for i := len(rps.Name) - 1; i >= 0; i-- {
		if rps.Name[i] == '/' {
			return rps.Name[i+1:]
		}
	}
	return ""
}
