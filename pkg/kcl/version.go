package kcl

import (
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type VersionResult = gpyrpc.GetVersion_Result

// GetVersion returns the KCL service version information.
func GetVersion() (*VersionResult, error) {
	svc := Service()
	resp, err := svc.GetVersion(&gpyrpc.GetVersion_Args{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
