// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"

	"kusionstack.io/kclvm-go/pkg/3rdparty/grpc_gateway_util"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ = fmt.Sprint

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
		c:       NewKclvmServiceClient(),
	}
	p.initHttpRrouter()
	return p
}

func (p *restServer) Run() error {
	fmt.Printf("listen on http://%s ...\n", p.address)
	return http.ListenAndServe(p.address, p.router)
}

func (p *restServer) initHttpRrouter() {
	methodMap := map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		// BuiltinService
		"BuiltinService.Ping":       p.handle_Ping,
		"BuiltinService.ListMethod": p.handle_ListMethod,

		// KclvmService
		"KclvmService.ParseFile_LarkTree":  p.handle_ParseFile_LarkTree,
		"KclvmService.ParseFile_AST":       p.handle_ParseFile_AST,
		"KclvmService.ParseProgram_AST":    p.handle_ParseProgram_AST,
		"KclvmService.ExecProgram":         p.handle_ExecProgram,
		"KclvmService.ResetPlugin":         p.handle_ResetPlugin,
		"KclvmService.FormatCode":          p.handle_FormatCode,
		"KclvmService.FormatPath":          p.handle_FormatPath,
		"KclvmService.LintPath":            p.handle_LintPath,
		"KclvmService.OverrideFile":        p.handle_OverrideFile,
		"KclvmService.EvalCode":            p.handle_EvalCode,
		"KclvmService.ResolveCode":         p.handle_ResolveCode,
		"KclvmService.GetSchemaType":       p.handle_GetSchemaType,
		"KclvmService.ValidateCode":        p.handle_ValidateCode,
		"KclvmService.SpliceCode":          p.handle_SpliceCode,
		"KclvmService.ListDepFiles":        p.handle_ListDepFiles,
		"KclvmService.ListUpStreamFiles":   p.handle_ListUpStreamFiles,
		"KclvmService.ListDownStreamFiles": p.handle_ListDownStreamFiles,
		"KclvmService.LoadSettingsFiles":   p.handle_LoadSettingsFiles,
	}

	for methodName, methodFunc := range methodMap {
		p.router.GET("/api:protorpc/"+methodName, methodFunc)
		p.router.POST("/api:protorpc/"+methodName, methodFunc)
	}
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

func (p *restServer) handle_ParseFile_LarkTree(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ParseFile_LarkTree_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ParseFile_LarkTree(args)
	})
}

func (p *restServer) handle_ParseFile_AST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ParseFile_AST_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ParseFile_AST(args)
	})
}

func (p *restServer) handle_ParseProgram_AST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ParseProgram_AST_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ParseProgram_AST(args)
	})
}

func (p *restServer) handle_ExecProgram(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ExecProgram_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ExecProgram(args)
	})
}

func (p *restServer) handle_ResetPlugin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var args = new(gpyrpc.ResetPlugin_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ResetPlugin(args)
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

func (p *restServer) handle_EvalCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.EvalCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.EvalCode(args)
	})
}

func (p *restServer) handle_ResolveCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ResolveCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ResolveCode(args)
	})
}

func (p *restServer) handle_GetSchemaType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.GetSchemaType_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.GetSchemaType(args)
	})
}

func (p *restServer) handle_ValidateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ValidateCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ValidateCode(args)
	})
}

func (p *restServer) handle_SpliceCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.SpliceCode_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.SpliceCode(args)
	})
}

func (p *restServer) handle_ListDepFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ListDepFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListDepFiles(args)
	})
}

func (p *restServer) handle_ListUpStreamFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ListUpStreamFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListUpStreamFiles(args)
	})
}

func (p *restServer) handle_ListDownStreamFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.ListDownStreamFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.ListDownStreamFiles(args)
	})
}

func (p *restServer) handle_LoadSettingsFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	args := new(gpyrpc.LoadSettingsFiles_Args)
	p.handle(w, r, args, func() (proto.Message, error) {
		return p.c.LoadSettingsFiles(args)
	})
}
