package service

import (
	"testing"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func TestCallRestMethod_ping(t *testing.T) {
	var args = gpyrpc.Ping_Args{Value: "ping"}
	var result gpyrpc.Ping_Result

	var err = CallRestMethod(
		"http://"+tRestServerAddr, "BuiltinService.Ping",
		&args, &result,
	)
	if err != nil {
		t.Fatal(err)
	}

	if result.Value != args.Value {
		t.Fatalf("expect %q, got %q", args.Value, result.Value)
	}
}

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
