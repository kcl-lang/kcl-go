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
	"strings"
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

// writeKFile writes content into dir/name and returns the full path.
func writeKFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	assert2.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// personK is a schema fixture reused by the Generate* endpoint tests.
const personK = "schema Person:\n    name: str\n    age: int\n"

// TestRestServerListMethod checks that the endpoint surfaces the method list
// reported by the embedded libkcl runtime. libkcl v0.13.1 registers
// BuiltinService.ListMethod in its dispatcher, so the handler calls
// api.ServiceClient.ListMethod directly and gets the real names back.
// The assertions anchor on a few stable names rather than the whole list so a
// later lib bump adding RPCs does not break the test.
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
			names := result.Result.GetMethodNameList()
			assert2.NotEmpty(t, names, "%s ListMethod list", method)
			for _, want := range []string{
				"BuiltinService.ListMethod",
				"KclService.ExecProgram",
				"KclService.ParseFile",
			} {
				assert2.Contains(t, names, want, "%s ListMethod list", method)
			}
		}
	}
}

func TestRestServerGenerateKcl(t *testing.T) {
	ts := newTestRestServer(t)

	var resp struct {
		Error  string                    `json:"error"`
		Result *gpyrpc.GenerateKclResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.GenerateKcl",
		map[string]any{"source": `{"name":"alice","age":30}`, "filename": "data.json"}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		got := resp.Result.GetKcl()
		assert2.Contains(t, got, `name = "alice"`)
		assert2.Contains(t, got, "age = 30")
	}
}

func TestRestServerGenerateToml(t *testing.T) {
	ts := newTestRestServer(t)

	// A bare schema declaration evaluates to no top-level value, so the
	// fixture needs a real assignment for there being anything to serialize.
	dir := t.TempDir()
	file := writeKFile(t, dir, "main.k", personK+"\nalice = Person {name = \"alice\", age = 30}\n")

	var resp struct {
		Error  string                     `json:"error"`
		Result *gpyrpc.GenerateTomlResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.GenerateToml",
		map[string]any{
			"exec_args": map[string]any{"k_filename_list": []string{file}},
			"sort_keys": true,
		}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		// sort_keys=true orders `age` before `name`.
		got := resp.Result.GetToml()
		assert2.Contains(t, got, "[alice]")
		assert2.Less(t, strings.Index(got, "age"), strings.Index(got, "name"),
			"sort_keys=true must order keys alphabetically: %q", got)
	}
}

func TestRestServerGenerateOpenAPI(t *testing.T) {
	ts := newTestRestServer(t)

	dir := t.TempDir()
	file := writeKFile(t, dir, "person.k", personK)

	var resp struct {
		Error  string                        `json:"error"`
		Result *gpyrpc.GenerateOpenAPIResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.GenerateOpenAPI",
		map[string]any{"parse_args": map[string]any{"paths": []string{file}}}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		var spec struct {
			OpenAPI    string `json:"openapi"`
			Components struct {
				Schemas map[string]any `json:"schemas"`
			} `json:"components"`
		}
		if assert2.NoError(t, json.Unmarshal([]byte(resp.Result.GetSpec()), &spec)) {
			assert2.Equal(t, "3.0.0", spec.OpenAPI, "v3 is the default spec version")
			assert2.Contains(t, spec.Components.Schemas, "Person")
		}
	}
}

func TestRestServerGenerateProto(t *testing.T) {
	ts := newTestRestServer(t)

	dir := t.TempDir()
	file := writeKFile(t, dir, "person.k", personK)

	var resp struct {
		Error  string                      `json:"error"`
		Result *gpyrpc.GenerateProtoResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.GenerateProto",
		map[string]any{
			"parse_args": map[string]any{"paths": []string{file}},
			"package":    "example.v1",
		}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		got := resp.Result.GetProto()
		assert2.Contains(t, got, `syntax = "proto3";`)
		assert2.Contains(t, got, "package example.v1;")
		assert2.Contains(t, got, "message Person {")
	}
}

func TestRestServerGenerateDoc(t *testing.T) {
	ts := newTestRestServer(t)

	dir := t.TempDir()
	file := writeKFile(t, dir, "person.k", personK)

	var resp struct {
		Error  string                    `json:"error"`
		Result *gpyrpc.GenerateDocResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.GenerateDoc",
		map[string]any{"parse_args": map[string]any{"paths": []string{file}}}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		// Markdown is the default format.
		got := resp.Result.GetContent()
		assert2.Contains(t, got, "### Person")
		assert2.Contains(t, got, "| name | str |")
	}
}

func TestRestServerFormatTestReport(t *testing.T) {
	ts := newTestRestServer(t)

	var resp struct {
		Error  string                         `json:"error"`
		Result *gpyrpc.FormatTestReportResult `json:"result"`
	}
	postRPC(t, ts, "/api:protorpc/KclService.FormatTestReport",
		map[string]any{"result": map[string]any{
			"info": []map[string]any{{
				"name":        "TestPerson",
				"log_message": "hello from the report",
			}},
		}}, &resp)

	assert2.Empty(t, resp.Error)
	if assert2.NotNil(t, resp.Result) {
		got := resp.Result.GetReport()
		assert2.Contains(t, got, "TestPerson")
		assert2.Contains(t, got, "hello from the report")
	}
}

// TestRestServerRouteParity checks that the REST route table covers every
// method the embedded runtime advertises through ListMethod, on both GET and
// POST. Deriving the list from the runtime rather than hard-coding it means a
// future lib bump that adds an RPC fails here instead of silently shipping an
// endpoint that is missing from the route table.
func TestRestServerRouteParity(t *testing.T) {
	ts := newTestRestServer(t)

	names, err := newRestServer("127.0.0.1:0").service.ListMethod(new(gpyrpc.ListMethodArgs))
	assert2.NoError(t, err)
	assert2.NotEmpty(t, names.GetMethodNameList(), "runtime advertised no methods")

	for _, full := range names.GetMethodNameList() {
		service, method, ok := strings.Cut(full, ".")
		if !assert2.True(t, ok, "malformed method name %q", full) {
			continue
		}
		route := "/api:protorpc/" + service + "." + method
		for _, verb := range []string{http.MethodGet, http.MethodPost} {
			req, err := http.NewRequest(verb, ts.URL+route, nil)
			assert2.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			assert2.NoError(t, err)
			resp.Body.Close()
			// A handler-level failure still answers 200 with the error in the
			// body; a 404 means the route was never registered.
			assert2.NotEqual(t, http.StatusNotFound, resp.StatusCode,
				"%s %s is advertised by the runtime but not served", verb, route)
		}
	}
}
