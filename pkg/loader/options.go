package loader

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type OptionHelps = []*gpyrpc.OptionHelp
type ListOptionsArgs = gpyrpc.ParseProgramArgs
type ListOptionsResult = gpyrpc.ListOptionsResult

// ListFileOptions provides users with the ability to parse kcl program and get all option
// calling information.
func ListFileOptions(filename string) (OptionHelps, error) {
	svc := kcl.Service()
	resp, err := svc.ListOptions(&gpyrpc.ParseProgramArgs{
		Paths: []string{filename},
	})
	if err != nil {
		return nil, err
	}
	return resp.Options, nil
}

// ListOptions provides users with the ability to parse kcl program and get all option
// calling information.
func ListOptions(args *ListOptionsArgs) (*ListOptionsResult, error) {
	svc := kcl.Service()
	return svc.ListOptions(args)
}
