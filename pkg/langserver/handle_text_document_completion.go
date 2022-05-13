package langserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"

	"kusionstack.io/kclvm-go/pkg/service"
	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

func (h *langHandler) handleTextDocumentCompletion(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.completion(params.TextDocument.URI, &params)
}

func (h *langHandler) completion(uri DocumentURI, params *CompletionParams) ([]CompletionItem, error) {
	filename, f, err := h.openOrLoadFile(uri)
	if err != nil {
		return nil, err
	}
	word := f.WordAt(params.Position)
	if params.CompletionContext.TriggerKind == TriggerCharacter || params.CompletionContext.TriggerKind == TriggerForIncompleteCompletions {
		return nil, fmt.Errorf("completion trigger kind unsupported: %v", params.CompletionContext.TriggerKind)
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.Complete(&gpyrpc.Complete_Args{
		Pos: &gpyrpc.Position{
			Line:     int64(params.Position.Line),
			Column:   int64(params.Position.Character),
			Filename: filename,
		},
		Name: word,
		Code: f.Text,
	})
	if err != nil {
		return nil, err
	}
	completeItemsStr := resp.CompleteItems
	var result []CompletionItem
	err = json.Unmarshal([]byte(completeItemsStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
