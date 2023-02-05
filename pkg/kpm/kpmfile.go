package kpm

import "github.com/orangebees/go-oneutils/GlobalStore"

type KpmFile struct {
	//PackageName,this is usually the code source repository address
	PackageName string `json:"package_name"`
	//KclvmMinVersion,used to identify version conflict states
	KclvmMinVersion string `json:"kclvm_min_version"`
	//Dependencies that are directly needed, aliases are not duplicated
	Direct DirectRequire `json:"direct,omitempty"`
	//Indirect dependencies, do not look at the alias, the package name and version are unique
	Indirect IndirectRequire `json:"indirect,omitempty"`
}
type DirectRequire map[string]RequireBase
type IndirectRequire map[PkgString]GlobalStore.Integrity
