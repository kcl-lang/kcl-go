package langserver

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func (h *langHandler) handleTextDocumentDefinition(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DocumentDefinitionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.definition(params.TextDocument.URI, &params.TextDocumentPositionParams)
}

func (h *langHandler) definition(uri DocumentURI, params *TextDocumentPositionParams) ([]Location, error) {
	filename, f, err := h.openOrLoadFile(uri)
	if err != nil {
		return nil, err
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.GoToDef(&gpyrpc.GoToDef_Args{
		Pos: &gpyrpc.Position{
			Line:     int64(params.Position.Line),
			Column:   int64(params.Position.Character),
			Filename: filename,
		},
		Code: f.Text,
	})

	if err != nil {
		return nil, err
	}
	loc := resp.Locations
	var result []Location
	err = json.Unmarshal([]byte(loc), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
