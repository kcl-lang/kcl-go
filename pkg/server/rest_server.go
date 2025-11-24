// Copyright The KCL Authors. All rights reserved.

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"

	"kcl-lang.io/kcl-go/pkg/3rdparty/grpc_gateway_util"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
	"kcl-lang.io/lib/go/api"
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
	service api.ServiceClient
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
		service: kcl.Service(),
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

	p.router.GET("/api:protorpc/KclService.ExecProgram", p.handle_ExecProgram)
	p.router.GET("/api:protorpc/KclService.BuildProgram", p.handle_BuildProgram)
	p.router.GET("/api:protorpc/KclService.ExecArtifact", p.handle_ExecArtifact)
	p.router.GET("/api:protorpc/KclService.ParseFile", p.handle_ParseFile)
	p.router.GET("/api:protorpc/KclService.ParseProgram", p.handle_ParseProgram)
	p.router.GET("/api:protorpc/KclService.ListOptions", p.handle_ListOptions)
	p.router.GET("/api:protorpc/KclService.ListVariables", p.handle_ListVariables)
	p.router.GET("/api:protorpc/KclService.LoadPackage", p.handle_LoadPackage)
	p.router.GET("/api:protorpc/KclService.FormatCode", p.handle_FormatCode)
	p.router.GET("/api:protorpc/KclService.FormatPath", p.handle_FormatPath)
	p.router.GET("/api:protorpc/KclService.LintPath", p.handle_LintPath)
	p.router.GET("/api:protorpc/KclService.OverrideFile", p.handle_OverrideFile)
	p.router.GET("/api:protorpc/KclService.GetSchemaTypeMapping", p.handle_GetSchemaTypeMapping)
	p.router.GET("/api:protorpc/KclService.ValidateCode", p.handle_ValidateCode)
	p.router.GET("/api:protorpc/KclService.Rename", p.handle_Rename)
	p.router.GET("/api:protorpc/KclService.RenameCode", p.handle_RenameCode)
	p.router.GET("/api:protorpc/KclService.Test", p.handle_Test)
	p.router.GET("/api:protorpc/KclService.UpdateDependencies", p.handle_UpdateDependencies)
	p.router.GET("/api:protorpc/KclService.GetVersion", p.handle_GetVersion)

	p.router.POST("/api:protorpc/BuiltinService.Ping", p.handle_Ping)

	p.router.POST("/api:protorpc/KclService.ExecProgram", p.handle_ExecProgram)
	p.router.POST("/api:protorpc/KclService.BuildProgram", p.handle_BuildProgram)
	p.router.POST("/api:protorpc/KclService.ExecArtifact", p.handle_ExecArtifact)
	p.router.POST("/api:protorpc/KclService.ParseFile", p.handle_ParseFile)
	p.router.POST("/api:protorpc/KclService.ParseProgram", p.handle_ParseProgram)
	p.router.POST("/api:protorpc/KclService.ListOptions", p.handle_ListOptions)
	p.router.POST("/api:protorpc/KclService.ListVariables", p.handle_ListVariables)
	p.router.POST("/api:protorpc/KclService.LoadPackage", p.handle_LoadPackage)
	p.router.POST("/api:protorpc/KclService.FormatCode", p.handle_FormatCode)
	p.router.POST("/api:protorpc/KclService.FormatPath", p.handle_FormatPath)
	p.router.POST("/api:protorpc/KclService.LintPath", p.handle_LintPath)
	p.router.POST("/api:protorpc/KclService.OverrideFile", p.handle_OverrideFile)
	p.router.POST("/api:protorpc/KclService.GetSchemaTypeMapping", p.handle_GetSchemaTypeMapping)
	p.router.POST("/api:protorpc/KclService.ValidateCode", p.handle_ValidateCode)
	p.router.POST("/api:protorpc/KclService.Rename", p.handle_Rename)
	p.router.POST("/api:protorpc/KclService.RenameCode", p.handle_RenameCode)
	p.router.POST("/api:protorpc/KclService.Test", p.handle_Test)
	p.router.POST("/api:protorpc/KclService.UpdateDependencies", p.handle_UpdateDependencies)
	p.router.POST("/api:protorpc/KclService.GetVersion", p.handle_GetVersion)
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

func (p *restServer) handle_Ping(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.PingArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.Ping(args)
	})
}

func (p *restServer) handle_ExecProgram(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.ExecProgramArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ExecProgram(args)
	})
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (p *restServer) handle_BuildProgram(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.BuildProgramArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.BuildProgram(args)
	})
}

// Depreciated: Please use the env.EnableFastEvalMode() and c.ExecuteProgram method and will be removed in v0.11.0.
func (p *restServer) handle_ExecArtifact(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.ExecArtifactArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ExecArtifact(args)
	})
}

func (p *restServer) handle_ParseFile(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.ParseFileArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ParseFile(args)
	})
}

func (p *restServer) handle_ParseProgram(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.ParseProgramArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ParseProgram(args)
	})
}

func (p *restServer) handle_ListOptions(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.ParseProgramArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ListOptions(args)
	})
}

func (p *restServer) handle_ListVariables(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.ListVariablesArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ListVariables(args)
	})
}

func (p *restServer) handle_LoadPackage(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.LoadPackageArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.LoadPackage(args)
	})
}

func (p *restServer) handle_FormatCode(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.FormatCodeArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.FormatCode(args)
	})
}

func (p *restServer) handle_FormatPath(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	var args = new(gpyrpc.FormatPathArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.FormatPath(args)
	})
}

func (p *restServer) handle_LintPath(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.LintPathArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.LintPath(args)
	})
}

func (p *restServer) handle_OverrideFile(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.OverrideFileArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.OverrideFile(args)
	})
}

func (p *restServer) handle_GetSchemaTypeMapping(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.GetSchemaTypeMappingArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.GetSchemaTypeMapping(args)
	})
}

func (p *restServer) handle_ValidateCode(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.ValidateCodeArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ValidateCode(args)
	})
}

func (p *restServer) handle_ListDepFiles(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.ListDepFilesArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.ListDepFiles(args)
	})
}

func (p *restServer) handle_LoadSettingsFiles(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.LoadSettingsFilesArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.LoadSettingsFiles(args)
	})
}

func (p *restServer) handle_Rename(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.RenameArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.Rename(args)
	})
}

func (p *restServer) handle_RenameCode(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.RenameCodeArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.RenameCode(args)
	})
}

func (p *restServer) handle_Test(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.TestArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.Test(args)
	})
}

func (p *restServer) handle_UpdateDependencies(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.UpdateDependenciesArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.UpdateDependencies(args)
	})
}

func (p *restServer) handle_GetVersion(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	args := new(gpyrpc.GetVersionArgs)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.service.GetVersion(args)
	})
}
