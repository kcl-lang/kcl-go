package langserver

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func (h *langHandler) handleTextDocumentHover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.hover(params.TextDocument.URI, &params)
}

func (h *langHandler) hover(uri DocumentURI, params *HoverParams) (*Hover, error) {
	filename, f, err := h.openOrLoadFile(uri)
	if err != nil {
		return nil, err
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.Hover(&gpyrpc.Hover_Args{
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
	hoverResult := resp.HoverResult
	var result Hover
	err = json.Unmarshal([]byte(hoverResult), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
