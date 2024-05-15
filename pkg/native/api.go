//go:build cgo
// +build cgo

package native

import (
	"kcl-lang.io/kcl-go/pkg/kcl"
)

func MustRun(path string, opts ...kcl.Option) *kcl.KCLResultList[kcl.KCLResultType] {
	v, err := Run[kcl.KCLResultType](path, opts...)
	if err != nil {
		panic(err)
	}

	return v
}

func Run[T kcl.KCLResultType](path string, opts ...kcl.Option) (*kcl.KCLResultList[T], error) {
	return run[T]([]string{path}, opts...)
}

func run[T kcl.KCLResultType](pathList []string, opts ...kcl.Option) (*kcl.KCLResultList[T], error) {
	args, err := kcl.ParseArgs(pathList, opts...)
	if err != nil {
		return nil, err
	}

	client := NewNativeServiceClient()
	resp, err := client.ExecProgram(args.ExecProgram_Args)
	if err != nil {
		return nil, err
	}
	return kcl.ExecResultToKCLResult[T](&args, resp, args.GetLogger(), kcl.DefaultHooks)
}
