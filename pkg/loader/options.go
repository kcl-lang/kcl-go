package loader

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type OptionHelps = []*gpyrpc.OptionHelp

// ListFileOptions provides users with the ability to parse kcl program and get all option
// calling information.
func ListFileOptions(filename string) (OptionHelps, error) {
	svc := kcl.Service()
	resp, err := svc.ListOptions(&gpyrpc.ParseProgram_Args{
		Paths: []string{filename},
	})
	if err != nil {
		return nil, err
	}
	return resp.Options, nil
}
