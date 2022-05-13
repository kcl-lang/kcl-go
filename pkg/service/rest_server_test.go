// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"testing"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
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

func TestRestServer_splice_pode(t *testing.T) {
	var codeSnippets []*gpyrpc.CodeSnippet = []*gpyrpc.CodeSnippet{
		{
			Schema: `schema Person:\n    age: int`,
			Rule:   `age > 0`,
		},
	}
	var args = gpyrpc.SpliceCode_Args{CodeSnippets: codeSnippets}
	var resp struct {
		Error  string                   `json:"error"`
		Result gpyrpc.SpliceCode_Result `json:"result"`
	}
	err := httpPost("http://"+tRestServerAddr+"/api:protorpc/KclvmService.SpliceCode", &args, &resp)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRestServer_resolve_pode(t *testing.T) {
	var args = gpyrpc.ResolveCode_Args{Code: "name = 'Alice'"}
	var resp struct {
		Error  string                    `json:"error"`
		Result gpyrpc.ResolveCode_Result `json:"result"`
	}
	err := httpPost("http://"+tRestServerAddr+"/api:protorpc/KclvmService.ResolveCode", &args, &resp)
	if err != nil {
		t.Fatal(err)
	}
}
