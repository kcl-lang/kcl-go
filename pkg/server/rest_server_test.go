// Copyright The KCL Authors. All rights reserved.

//go:build cgo

package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// newTestRestServer serves the REST router through an httptest server so
// tests can exercise the endpoints end to end without binding a fixed port.
func newTestRestServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := newRestServer("127.0.0.1:0")
	ts := httptest.NewServer(s.router)
	t.Cleanup(ts.Close)
	return ts
}

// postRPC issues a POST against an /api:protorpc/ endpoint and decodes the
// RestfulResult-style response body into out.
func postRPC(t *testing.T, ts *httptest.Server, route string, body any, out any) {
	t.Helper()
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		assert2.NoError(t, err)
	}
	resp, err := http.Post(ts.URL+route, "application/json", bytes.NewReader(payload))
	assert2.NoError(t, err)
	defer resp.Body.Close()
	assert2.Equal(t, http.StatusOK, resp.StatusCode)
	if out != nil {
		assert2.NoError(t, json.NewDecoder(resp.Body).Decode(out))
	}
}

func TestRestServerLoadSettingsFiles(t *testing.T) {
	ts := newTestRestServer(t)

	dir := t.TempDir()
	settingsFile := filepath.Join(dir, "kcl.yaml")
	settings := "kcl_cli_configs:\n  strict_range_check: true\nkcl_options:\n  - key: key\n    value: value\n"
	assert2.NoError(t, os.WriteFile(settingsFile, []byte(settings), 0o644))

	var resp struct {
		Error  string                          `json:"error"`
		Result *gpyrpc.LoadSettingsFilesResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.LoadSettingsFiles",
		map[string]any{"work_dir": dir, "files": []string{settingsFile}}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		assert2.True(t, resp.Result.GetKclCliConfigs().GetStrictRangeCheck())
		if assert2.Len(t, resp.Result.GetKclOptions(), 1) {
			assert2.Equal(t, "key", resp.Result.GetKclOptions()[0].GetKey())
		}
	}
}

func TestRestServerGetSchemaTypeMappingUnderPath(t *testing.T) {
	ts := newTestRestServer(t)

	dir := t.TempDir()
	file := filepath.Join(dir, "person.k")
	assert2.NoError(t, os.WriteFile(file, []byte("schema Person:\n    name: str\n"), 0o644))

	var resp struct {
		Error  string                                      `json:"error"`
		Result *gpyrpc.GetSchemaTypeMappingUnderPathResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.GetSchemaTypeMappingUnderPath",
		map[string]any{"exec_args": map[string]any{"k_filename_list": []string{file}}}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		mapping := resp.Result.GetSchemaTypeMapping()
		if assert2.Contains(t, mapping, "__main__") {
			var names []string
			for _, st := range mapping["__main__"].GetSchemaType() {
				names = append(names, st.GetSchemaName())
			}
			assert2.Contains(t, names, "Person")
		}
	}
}

// TestRestServerListMethod documents the behavior measured against the
// embedded libkcl v0.13.0 runtime: its service dispatcher does not register
// BuiltinService.ListMethod (the Rust side answers with an empty buffer for
// unregistered methods), and the method is not part of api.ServiceClient.
// The REST endpoint must degrade gracefully instead of panicking: it returns
// HTTP 200 with an empty method list.
func TestRestServerListMethod(t *testing.T) {
	ts := newTestRestServer(t)

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		var req *http.Request
		var err error
		if method == http.MethodPost {
			req, err = http.NewRequest(method, ts.URL+"/api:protorpc/BuiltinService.ListMethod", bytes.NewReader(nil))
		} else {
			req, err = http.NewRequest(method, ts.URL+"/api:protorpc/BuiltinService.ListMethod", nil)
		}
		assert2.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert2.NoError(t, err)
		var result struct {
			Error  string                   `json:"error"`
			Result *gpyrpc.ListMethodResult `json:"result"`
		}
		assert2.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		resp.Body.Close()

		assert2.Equal(t, http.StatusOK, resp.StatusCode, "%s ListMethod status", method)
		assert2.Empty(t, result.Error, "%s ListMethod error", method)
		if assert2.NotNil(t, result.Result, "%s ListMethod result", method) {
			assert2.Empty(t, result.Result.GetMethodNameList(), "%s ListMethod list", method)
		}
	}
}
