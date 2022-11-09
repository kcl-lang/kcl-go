package kpm

type Require struct {
	//包类型 git，registry
	Type string `json:"type"`
	//别名
	Alias string `json:"alias,omitempty"`
	//包名，确定包的命名空间
	Name string `json:"name"`
	//确定此包的版本
	Version PkgVersion `json:"version"`
	//校验和 sha512
	Integrity string `json:"integrity"`
	//引用计数
	Count int `json:"count,omitempty"`
}
