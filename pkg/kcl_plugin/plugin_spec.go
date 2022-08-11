// Copyright 2022 The KCL Authors. All rights reserved.

package kcl_plugin

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// KCL Plugin object
type Plugin struct {
	Name      string
	ResetFunc func()
	MethodMap map[string]MethodSpec
}

// KCL Plugin method spec
type MethodSpec struct {
	Type *MethodType
	Body func(args *MethodArgs) (*MethodResult, error)
}

// KCL Plugin method type
type MethodType struct {
	ArgsType   []string
	KwArgsType map[string]string
	ResultType string
}

// plugin method args
type MethodArgs struct {
	Args   []interface{}
	KwArgs map[string]interface{}
}

// plugin method result
type MethodResult struct {
	V interface{}
}

func ParseMethodArgs(args_json, kwargs_json string) (*MethodArgs, error) {
	p := &MethodArgs{
		KwArgs: make(map[string]interface{}),
	}
	if args_json != "" {
		if err := json.Unmarshal([]byte(args_json), &p.Args); err != nil {
			return nil, err
		}
	}
	if kwargs_json != "" {
		if err := json.Unmarshal([]byte(kwargs_json), &p.KwArgs); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *MethodArgs) Arg(i int) interface{} {
	return p.Args[i]
}
func (p *MethodArgs) KwArg(name string) interface{} {
	return p.KwArgs[name]
}

func (p *MethodArgs) IntArg(i int) int64 {
	s := fmt.Sprint(p.Args[i])
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func (p *MethodArgs) FloatArg(i int) float64 {
	s := fmt.Sprint(p.Args[i])
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func (p *MethodArgs) StrArg(i int) string {
	s := fmt.Sprint(p.Args[i])
	return s
}

func (p *MethodArgs) IntKwArg(name string) int64 {
	s := fmt.Sprint(p.KwArgs[name])
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func (p *MethodArgs) FloatKwArg(name string) float64 {
	s := fmt.Sprint(p.KwArgs[name])
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func (p *MethodArgs) StrKwArg(name string) string {
	s := fmt.Sprint(p.KwArgs[name])
	return s
}
