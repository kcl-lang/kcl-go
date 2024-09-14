//go:build rpc || !cgo
// +build rpc !cgo

// Copyright The KCL Authors. All rights reserved.

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"

	"kcl-lang.io/kcl-go/pkg/3rdparty/grpc_gateway_util"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

var _ = fmt.Sprint

// Client represents an restful method result.
type RestfulResult struct {
	Error  string        `json:"error"`
	Result proto.Message `json:"result"`
}

type restServer struct {
	address string
	router  *httprouter.Router
	builtin *BuiltinServiceClient
	c       *KclvmServiceClient
}

func RunRestServer(address string) error {
	s := newRestServer(address)
	return s.Run()
}

func newRestServer(address string) *restServer {
	if strings.HasPrefix(address, ":") {
		address = "127.0.0.1" + address
	}
	p := &restServer{
		address: address,
		router:  httprouter.New(),
		builtin: NewBuiltinServiceClient(),
		c:       newKclvmServiceClient(),
	}
	p.initHttpRrouter()
	return p
}

func (p *restServer) Run() error {
	fmt.Printf("listen on http://%s ...\n", p.address)
	return http.ListenAndServe(p.address, p.router)
}

func (p *restServer) initHttpRrouter() {
	p.router.GET("/api:protorpc/BuiltinService.Ping", p.handle_Ping)
	p.router.GET("/api:protorpc/BuiltinService.ListMethod", p.handle_ListMethod)

	p.router.GET("/api:protorpc/KclvmService.ExecProgram", p.handle_ExecProgram)
	p.router.GET("/api:protorpc/KclvmService.BuildProgram", p.handle_BuildProgram)
	p.router.GET("/api:protorpc/KclvmService.ExecArtifact", p.handle_ExecArtifact)
	p.router.GET("/api:protorpc/KclvmService.ParseFile", p.handle_ParseFile)
	p.router.GET("/api:protorpc/KclvmService.ParseProgram", p.handle_ParseProgram)
	p.router.GET("/api:protorpc/KclvmService.ListOptions", p.handle_ListOptions)
	p.router.GET("/api:protorpc/KclvmService.ListVariables", p.handle_ListVariables)
	p.router.GET("/api:protorpc/KclvmService.LoadPackage", p.handle_LoadPackage)
	p.router.GET("/api:protorpc/KclvmService.FormatCode", p.handle_FormatCode)
	p.router.GET("/api:protorpc/KclvmService.FormatPath", p.handle_FormatPath)
	p.router.GET("/api:protorpc/KclvmService.LintPath", p.handle_LintPath)
	p.router.GET("/api:protorpc/KclvmService.OverrideFile", p.handle_OverrideFile)
	p.router.GET("/api:protorpc/KclvmService.GetSchemaTypeMapping", p.handle_GetSchemaTypeMapping)
	p.router.GET("/api:protorpc/KclvmService.ValidateCode", p.handle_ValidateCode)
	p.router.GET("/api:protorpc/KclvmService.Rename", p.handle_Rename)
	p.router.GET("/api:protorpc/KclvmService.RenameCode", p.handle_RenameCode)
	p.router.GET("/api:protorpc/KclvmService.Test", p.handle_Test)
	p.router.GET("/api:protorpc/KclvmService.UpdateDependencies", p.handle_UpdateDependencies)
	p.router.GET("/api:protorpc/KclvmService.GetVersion", p.handle_GetVersion)

	p.router.POST("/api:protorpc/BuiltinService.Ping", p.handle_Ping)
	p.router.POST("/api:protorpc/BuiltinService.ListMethod", p.handle_ListMethod)

	p.router.POST("/api:protorpc/KclvmService.ExecProgram", p.handle_ExecProgram)
	p.router.POST("/api:protorpc/KclvmService.BuildProgram", p.handle_BuildProgram)
	p.router.POST("/api:protorpc/KclvmService.ExecArtifact", p.handle_ExecArtifact)
	p.router.POST("/api:protorpc/KclvmService.ParseFile", p.handle_ParseFile)
	p.router.POST("/api:protorpc/KclvmService.ParseProgram", p.handle_ParseProgram)
	p.router.POST("/api:protorpc/KclvmService.ListOptions", p.handle_ListOptions)
	p.router.POST("/api:protorpc/KclvmService.ListVariables", p.handle_ListVariables)
	p.router.POST("/api:protorpc/KclvmService.LoadPackage", p.handle_LoadPackage)
	p.router.POST("/api:protorpc/KclvmService.FormatCode", p.handle_FormatCode)
	p.router.POST("/api:protorpc/KclvmService.FormatPath", p.handle_FormatPath)
	p.router.POST("/api:protorpc/KclvmService.LintPath", p.handle_LintPath)
	p.router.POST("/api:protorpc/KclvmService.OverrideFile", p.handle_OverrideFile)
	p.router.POST("/api:protorpc/KclvmService.GetSchemaTypeMapping", p.handle_GetSchemaTypeMapping)
	p.router.POST("/api:protorpc/KclvmService.ValidateCode", p.handle_ValidateCode)
	p.router.POST("/api:protorpc/KclvmService.Rename", p.handle_Rename)
	p.router.POST("/api:protorpc/KclvmService.RenameCode", p.handle_RenameCode)
	p.router.POST("/api:protorpc/KclvmService.Test", p.handle_Test)
	p.router.POST("/api:protorpc/KclvmService.UpdateDependencies", p.handle_UpdateDependencies)
	p.router.POST("/api:protorpc/KclvmService.GetVersion", p.handle_GetVersion)
}

func (p *restServer) handle(
	w http.ResponseWriter, r *http.Request,
	args proto.Message, fn func() (proto.Message, error),
) {
	switch r.Method {
	case "GET":
		if err := grpc_gateway_util.PopulateQueryParameters(args, r.URL.Query()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		if err := json.NewDecoder(r.Body).Decode(args); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var result RestfulResult
	if x, err := fn(); err != nil {
		result.Error = err.Error()
	} else {
		result.Result = x // OK
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *restServer) handle_Ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.Ping_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.builtin.Ping(args)
	})
}

func (p *restServer) handle_ListMethod(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ListMethod_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.builtin.ListMethod(args)
	})
}

func (p *restServer) handle_ExecProgram(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ExecProgram_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ExecProgram(args)
	})
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecutProgram method and will be removed in v0.11.0.
func (p *restServer) handle_BuildProgram(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.BuildProgram_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.BuildProgram(args)
	})
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecutProgram method and will be removed in v0.11.0.
func (p *restServer) handle_ExecArtifact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ExecArtifact_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ExecArtifact(args)
	})
}

func (p *restServer) handle_ParseFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ParseFile_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ParseFile(args)
	})
}

func (p *restServer) handle_ParseProgram(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ParseProgram_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ParseProgram(args)
	})
}

func (p *restServer) handle_ListOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ParseProgram_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListOptions(args)
	})
}

func (p *restServer) handle_ListVariables(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ListVariables_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListVariables(args)
	})
}

func (p *restServer) handle_LoadPackage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.LoadPackage_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.LoadPackage(args)
	})
}

func (p *restServer) handle_FormatCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.FormatCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.FormatCode(args)
	})
}

func (p *restServer) handle_FormatPath(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.FormatPath_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.FormatPath(args)
	})
}

func (p *restServer) handle_LintPath(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.LintPath_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.LintPath(args)
	})
}

func (p *restServer) handle_OverrideFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.OverrideFile_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.OverrideFile(args)
	})
}

func (p *restServer) handle_GetSchemaTypeMapping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.GetSchemaTypeMapping_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.GetSchemaTypeMapping(args)
	})
}

func (p *restServer) handle_ValidateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ValidateCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ValidateCode(args)
	})
}

func (p *restServer) handle_ListDepFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ListDepFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListDepFiles(args)
	})
}

func (p *restServer) handle_LoadSettingsFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.LoadSettingsFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.LoadSettingsFiles(args)
	})
}

func (p *restServer) handle_Rename(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.Rename_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.Rename(args)
	})
}

func (p *restServer) handle_RenameCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.RenameCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.RenameCode(args)
	})
}

func (p *restServer) handle_Test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.Test_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.Test(args)
	})
}

func (p *restServer) handle_UpdateDependencies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.UpdateDependencies_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.UpdateDependencies(args)
	})
}

func (p *restServer) handle_GetVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.GetVersion_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.GetVersion(args)
	})
}
