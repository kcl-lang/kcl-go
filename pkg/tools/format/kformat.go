package format

import (
	"errors"
	"io"

	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

func FormatCode(code interface{}) ([]byte, error) {
	var codeStr string
	switch code := code.(type) {
	case []byte:
		codeStr = string(code)
	case string:
		codeStr = code
	case io.Reader:
		var p []byte
		_, err := code.Read(p)
		if err != nil {
			return nil, err
		}
		codeStr = string(p)
	default:
		return nil, errors.New("unsupported source code format. valid formats: []byte, string, io.Reader")
	}

	client := service.NewKclvmServiceClient()
	resp, err := client.FormatCode(&gpyrpc.FormatCode_Args{
		Source: codeStr,
	})
	if err != nil {
		return nil, err
	}
	return resp.Formatted, nil
}

func FormatPath(path string) (changedPaths []string, err error) {
	client := service.NewKclvmServiceClient()
	resp, err := client.FormatPath(&gpyrpc.FormatPath_Args{
		Path: path,
	})
	if err != nil {
		return nil, err
	}
	return resp.ChangedPaths, nil
}
