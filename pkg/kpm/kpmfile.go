package kpm

type KpmFile struct {
	//包名，确定包的命名空间
	PackageName string `json:"package_name"`
	//确定此包的kcl最低运行版本
	KclvmMinVersion string `json:"kclvm_min_version"`
	//直接需要的依赖，别名不重复
	Direct DirectRequire `json:"direct,omitempty"`
	//间接需要的依赖，不看别名，包名版本唯一即可
	Indirect IndirectRequire `json:"indirect,omitempty"`
}
type DirectRequire map[string]RequireBase
type IndirectRequire map[PkgString]RequirePlus

// AddRequire 添加依赖到KpmFile
func (kf *KpmFile) AddRequire(r Require) {
	kf.Direct[r.Alias] = r.RequireBase
	kpmC.Get(r.RequireBase)
}
