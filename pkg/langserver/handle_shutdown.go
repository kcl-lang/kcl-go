package langserver

import (
	"context"
	"log"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleShutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if h.lintTimer != nil {
		h.lintTimer.Stop()
	}
	h.mu.Lock()
	if h.shutdown {
		log.Printf("Warning: server received a shutdown request after it was already shut down.")
	}
	h.shutdown = true
	h.mu.Unlock()

	close(h.request)

	return nil, conn.Close()
}
