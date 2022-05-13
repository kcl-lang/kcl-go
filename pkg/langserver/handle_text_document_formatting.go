package langserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	kfmt "kusionstack.io/kclvm-go/pkg/tools/format"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentFormatting(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DocumentFormattingParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return h.formatting(params.TextDocument.URI, params.Options)
}

func (h *langHandler) formatting(uri DocumentURI, options FormattingOptions) ([]TextEdit, error) {
	_, f, err := h.openOrLoadFile(uri)
	if err != nil {
		return nil, err
	}
	formatted, err := kfmt.FormatCode(f.Text)
	if err != nil {
		return nil, fmt.Errorf("format failed with exception: %s", err)
	}
	h.logMessage(LogInfo, "KCL Language Server: format succeed")
	if h.loglevel > 5 {
		h.logger.Println("formatted:", string(formatted))
	}
	text := strings.Replace(string(formatted), "\r", "", -1)
	return ComputeEdits(uri, f.Text, text), nil
}
