package loader

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type LoadPackageArgs = gpyrpc.LoadPackageArgs
type LoadPackageResult = gpyrpc.LoadPackageResult

// LoadPackage provides users with the ability to parse KCL program and semantic model
// information including symbols, types, definitions, etc.
func LoadPackage(args *LoadPackageArgs) (*LoadPackageResult, error) {
	svc := kcl.Service()
	return svc.LoadPackage(args)
}
