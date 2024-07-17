package loader

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type ListVariablesArgs = gpyrpc.ListVariables_Args
type ListVariablesResult = gpyrpc.ListVariables_Result

// ListVariables provides users with the ability to parse KCL program and get all variables by specs.
func ListVariables(args *ListVariablesArgs) (*ListVariablesResult, error) {
	svc := kcl.Service()
	return svc.ListVariables(args)
}
