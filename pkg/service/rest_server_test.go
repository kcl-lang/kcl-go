// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"testing"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

func TestRestServer_ping(t *testing.T) {
	var args = gpyrpc.Ping_Args{Value: "abc"}
	var resp struct {
		Error  string             `json:"error"`
		Result gpyrpc.Ping_Result `json:"result"`
	}
	err := httpPost("http://"+tRestServerAddr+"/api:protorpc/BuiltinService.Ping", &args, &resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Result.Value != args.Value {
		t.Fatalf("expect = %v, got = %v", &args, &resp.Result)
	}
}
