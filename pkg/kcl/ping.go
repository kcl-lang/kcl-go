package kcl

import (
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// Ping checks the availability of the KCL service and returns the echoed value.
func Ping(value string) (string, error) {
	svc := Service()
	resp, err := svc.Ping(&gpyrpc.PingArgs{Value: value})
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}
