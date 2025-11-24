package kcl

import (
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type VersionResult = gpyrpc.GetVersionResult

// GetVersion returns the KCL service version information.
func GetVersion() (*VersionResult, error) {
	svc := Service()
	resp, err := svc.GetVersion(&gpyrpc.GetVersionArgs{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
