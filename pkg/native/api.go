//go:build cgo
// +build cgo

package native

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
)

func MustRun(path string, opts ...kcl.Option) *kcl.KCLResultList {
	v, err := Run(path, opts...)
	if err != nil {
		panic(err)
	}

	return v
}

func Run(path string, opts ...kcl.Option) (*kcl.KCLResultList, error) {
	return run([]string{path}, opts...)
}

func MustRunPaths(paths []string, opts ...kcl.Option) *kcl.KCLResultList {
	v, err := RunPaths(paths, opts...)
	if err != nil {
		panic(err)
	}

	return v
}

func RunPaths(paths []string, opts ...kcl.Option) (*kcl.KCLResultList, error) {
	return run(paths, opts...)
}

func run(pathList []string, opts ...kcl.Option) (*kcl.KCLResultList, error) {
	args, err := kcl.ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := NewNativeServiceClient()
	resp, err := client.ExecProgram(args.ExecProgram_Args)
	if err != nil {
		return nil, err
	}
	return kcl.ExecResultToKCLResult(&args, resp, args.GetLogger(), kcl.DefaultHooks)
}
