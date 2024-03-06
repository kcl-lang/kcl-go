// Copyright The KCL Authors. All rights reserved.

package validate

import (
	"errors"
	"os"

	"kcl-lang.io/kcl-go/pkg/service"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

// ValidateOptions represents the options for the Validate function.
type ValidateOptions struct {
	Schema        string // The schema to validate against.
	AttributeName string // The attribute name to validate.
	Format        string // The format of the data.
}

// Validate validates the given data file against the specified
// schema file with the provided options.
func Validate(dataFile, schemaFile string, opts *ValidateOptions) (ok bool, err error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return false, err
	}
	if opts == nil {
		opts = &ValidateOptions{}
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.ValidateCode(&gpyrpc.ValidateCode_Args{
		File:          schemaFile,
		Data:          string(data),
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}

func ValidateCode(data, code string, opts *ValidateOptions) (ok bool, err error) {
	if opts == nil {
		opts = &ValidateOptions{}
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.ValidateCode(&gpyrpc.ValidateCode_Args{
		Data:          data,
		Code:          code,
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}

func ValidateCodeFile(dataFile, data, code string, opts *ValidateOptions) (ok bool, err error) {
	if opts == nil {
		opts = &ValidateOptions{}
	}
	client := service.NewKclvmServiceClient()
	resp, err := client.ValidateCode(&gpyrpc.ValidateCode_Args{
		Datafile:      dataFile,
		Data:          data,
		Code:          code,
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}
