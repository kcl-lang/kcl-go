package langserver

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf16"

	"github.com/mattn/go-unicodeclass"
	"github.com/sourcegraph/jsonrpc2"
)

// Config is
type Config struct {
	LogFile  string
	LogLevel int
	Channel  io.ReadWriteCloser
	Quiet    bool
	Filename string      `yaml:"-"`
	Logger   *log.Logger `yaml:"-"`
}

// NewHandler create JSON-RPC handler for this language server.
func NewHandler(config *Config) jsonrpc2.Handler {
	if config.Logger == nil {
		config.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	var handler = &langHandler{
		loglevel:  config.LogLevel,
		logger:    config.Logger,
		files:     make(map[DocumentURI]*File),
		request:   make(chan DocumentURI),
		lintTimer: nil,
		conn:      nil,
		filename:  config.Filename,
	}
	return jsonrpc2.HandlerWithError(handler.handle)
}

type langHandler struct {
	mu           sync.Mutex // guards all fields
	loglevel     int
	logger       *log.Logger
	files        map[DocumentURI]*File
	request      chan DocumentURI
	lintTimer    *time.Timer
	lintDebounce time.Duration
	conn         *jsonrpc2.Conn
	rootPath     string
	filename     string
	folders      []string
	shutdown     bool
}

// File is
type File struct {
	LanguageID string
	Text       string
	Version    int
}

// WordAt is
func (f *File) WordAt(pos Position) string {
	isKclPunc := func(c rune) bool {
		kclPunc := []rune("=:;,?.()[]{}+-*/%&|^~<>!@")
		for _, p := range kclPunc {
			if c == p {
				return true
			}
		}
		return false
	}

	lines := strings.Split(f.Text, "\n")
	if pos.Line < 0 || pos.Line >= len(lines) {
		return ""
	}
	chars := utf16.Encode([]rune(lines[pos.Line]))
	if pos.Character < 0 || pos.Character > len(chars) {
		return ""
	}
	prevPos := 0
	currPos := -1
	prevCls := unicodeclass.Invalid
	for i, char := range chars {
		currCls := unicodeclass.Is(rune(char))
		if isKclPunc(rune(char)) {
			currCls = unicodeclass.Punctation
		}
		if currCls != prevCls {
			if i <= pos.Character {
				prevPos = i
			} else {
				if char == '_' {
					continue
				}
				currPos = i
				break
			}
		}
		prevCls = currCls
	}
	if currPos == -1 {
		currPos = len(chars)
	}
	return string(utf16.Decode(chars[prevPos:currPos]))
}

func isWindowsDrivePath(path string) bool {
	if len(path) < 4 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

func isWindowsDriveURI(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}

func fromURI(uri DocumentURI) (string, error) {
	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("only file URIs are supported, got %v", u.Scheme)
	}
	if isWindowsDriveURI(u.Path) {
		u.Path = u.Path[1:]
	}
	return u.Path, nil
}

func toURI(path string) DocumentURI {
	if isWindowsDrivePath(path) {
		path = "/" + path
	}
	return DocumentURI((&url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(path),
	}).String())
}

func (h *langHandler) lintRequest(uri DocumentURI) {
	//if h.lintTimer != nil {
	//	h.lintTimer.Reset(h.lintDebounce)
	//	return
	//}
	//h.lintTimer = time.AfterFunc(h.lintDebounce, func() {
	//	h.lintTimer = nil
	//	h.request <- uri
	//})
}

func (h *langHandler) logMessage(typ MessageType, message string) {
	_ = h.conn.Notify(
		context.Background(),
		"window/logMessage",
		&LogMessageParams{
			Type:    typ,
			Message: message,
		})
}

var rootMarkers = []string{"kcl.mod"}

func matchRootPath(fname string) string {
	dir := filepath.Dir(filepath.Clean(fname))
	var prev string
	for dir != prev {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			name := file.Name()
			isDir := file.IsDir()
			for _, marker := range rootMarkers {
				if strings.HasSuffix(marker, "/") {
					if !isDir {
						continue
					}
					marker = strings.TrimRight(marker, "/")
					if ok, _ := filepath.Match(marker, name); ok {
						return dir
					}
				} else {
					if isDir {
						continue
					}
					if ok, _ := filepath.Match(marker, name); ok {
						return dir
					}
				}
			}
		}
		prev = dir
		dir = filepath.Dir(dir)
	}
	return ""
}

func (h *langHandler) findRootPath(fname string) string {
	if dir := matchRootPath(fname); dir != "" {
		return dir
	}
	if dir := matchRootPath(fname); dir != "" {
		return dir
	}

	for _, folder := range h.folders {
		if len(fname) > len(folder) && strings.EqualFold(fname[:len(folder)], folder) {
			return folder
		}
	}

	return h.rootPath
}

func (h *langHandler) closeFile(uri DocumentURI) error {
	delete(h.files, uri)
	return nil
}

func (h *langHandler) saveFile(uri DocumentURI) error {
	h.lintRequest(uri)
	return nil
}

func (h *langHandler) openFile(uri DocumentURI, languageID string, version int, code string) {
	f := &File{
		Text:       code,
		LanguageID: languageID,
		Version:    version,
	}
	h.files[uri] = f
}

func (h *langHandler) openOrLoadFile(uri DocumentURI) (filename string, file *File, err error) {
	f, ok := h.files[uri]
	if !ok {
		h.logger.Printf("document not open: %v", uri)
	}
	filename, err = fromURI(uri)
	if err != nil {
		h.logger.Printf("invalid uri: %v", uri)
		return "", nil, fmt.Errorf("invalid uri: %v: %v", err, uri)
	}
	if f == nil {
		text, err := ioutil.ReadFile(filename)
		if err != nil {
			return filename, nil, fmt.Errorf("document not exist on file system: %v", err)
		}
		h.openFile(uri, "KCL", 1, string(text))
		f, _ = h.files[uri]
	}
	return filename, f, nil
}

func (h *langHandler) updateFile(uri DocumentURI, text string, version *int) error {
	f, ok := h.files[uri]
	if !ok {
		return fmt.Errorf("document not found: %v", uri)
	}
	f.Text = text
	if version != nil {
		f.Version = *version
	}

	h.lintRequest(uri)
	return nil
}

func (h *langHandler) addFolder(folder string) {
	folder = filepath.Clean(folder)
	found := false
	for _, cur := range h.folders {
		if cur == folder {
			found = true
			break
		}
	}
	if !found {
		h.folders = append(h.folders, folder)
	}
}

func (h *langHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Notif {
		switch req.Method {
		case "initialized":
			return
		case "textDocument/didOpen":
			return h.handleTextDocumentDidOpen(ctx, conn, req)
		case "textDocument/didChange":
			return h.handleTextDocumentDidChange(ctx, conn, req)
		case "textDocument/didSave":
			return h.handleTextDocumentDidSave(ctx, conn, req)
		case "textDocument/didClose":
			return h.handleTextDocumentDidClose(ctx, conn, req)
		case "exit":
			return h.handleExit(ctx, conn, req)
			//case "$/setTrace", "$/logTrace", "window/logMessage", "window/showMessage", "$/progress", "$/cancelRequest", "telemetry/event",
			//"window/workDoneProgress/cancel", "workspace/didChangeWorkspaceFolders", "workspace/didChangeConfiguration",
			//"workspace/didChangeWatchedFiles", "workspace/didCreateFiles", "workspace/didRenameFiles", "workspace/didDeleteFiles", "textDocument/willSave",
			//"textDocument/publishDiagnostics":
		}
		if h.loglevel > 5 {
			h.logger.Printf("unhandled notification: %s\n", req.Method)
		}
		return
	}
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, conn, req)
	case "shutdown":
		return h.handleShutdown(ctx, conn, req)
	case "textDocument/formatting":
		return h.handleTextDocumentFormatting(ctx, conn, req)
	case "textDocument/documentSymbol":
		return h.handleTextDocumentSymbol(ctx, conn, req)
	case "textDocument/completion":
		return h.handleTextDocumentCompletion(ctx, conn, req)
	case "textDocument/definition":
		return h.handleTextDocumentDefinition(ctx, conn, req)
	case "textDocument/references":
		return h.handleTextDocumentReference(ctx, conn, req)
	case "textDocument/hover":
		return h.handleTextDocumentHover(ctx, conn, req)
	case "textDocument/codeAction":
		return h.handleTextDocumentCodeAction(ctx, conn, req)
	case "workspace/executeCommand":
		return h.handleWorkspaceExecuteCommand(ctx, conn, req)
	case "workspace/didChangeConfiguration":
		return h.handleWorkspaceDidChangeConfiguration(ctx, conn, req)
	case "workspace/workspaceFolders":
		return h.handleWorkspaceWorkspaceFolders(ctx, conn, req)
	case "workspace/didChangeWorkspaceFolders":
		return h.handleDidChangeWorkspaceWorkspaceFolders(ctx, conn, req)
	}
	return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
}
