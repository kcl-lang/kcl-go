package langserver

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/sourcegraph/jsonrpc2"
)

var logFatalf = log.Fatalf

func Run(config *Config) {
	var connOpt []jsonrpc2.ConnOpt
	if config.LogFile != "" {
		logDir := filepath.Dir(config.LogFile)
		if _, err := os.Stat(logDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(logDir, os.ModePerm)
			if err != nil {
				logFatalf("KCL language server: create log file dir failed. %v", err)
			}
		}
		f, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			logFatalf("KCL language server: create log file failed. %v", err)
		}
		defer f.Close()
		config.Logger = log.New(f, "", log.LstdFlags)
		if config.LogLevel >= 5 {
			connOpt = append(connOpt, jsonrpc2.LogMessages(config.Logger))
		}
	}

	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	if config.Quiet && (config.LogFile == "" || config.LogLevel < 5) {
		connOpt = append(connOpt, jsonrpc2.LogMessages(log.New(ioutil.Discard, "", 0)))
	}

	handler := NewHandler(config)
	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(config.Channel, jsonrpc2.VSCodeObjectCodec{}), // 支持tcp代理重连，方便调试
		handler, connOpt...,
	).DisconnectNotify()
}
