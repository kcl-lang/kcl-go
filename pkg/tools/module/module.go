package module

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type UpdateDependenciesArgs = gpyrpc.UpdateDependencies_Args
type UpdateDependenciesResult = gpyrpc.UpdateDependencies_Result

// Download and update dependencies defined in the kcl.mod file and return the external package name and location list.
func UpdateDependencies(args *UpdateDependenciesArgs) (*UpdateDependenciesResult, error) {
	svc := kcl.Service()
	return svc.UpdateDependencies(args)
}
