package service

import (
	"testing"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

func TestCallRestMethod_noMethod(t *testing.T) {
	var args = gpyrpc.Ping_Args{Value: "ping"}
	var result gpyrpc.Ping_Result

	var err = CallRestMethod(
		"http://"+tRestServerAddr, "UnknownService.noMethod",
		&args, &result,
	)
	if err == nil {
		t.Fatalf("expect error, got %v", err)
	}
}
