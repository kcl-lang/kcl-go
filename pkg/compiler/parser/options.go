// Copyright 2021 The KCL Authors. All rights reserved.

package parser

type Option interface {
	apply(*options)
}

type options struct{}

func (p *options) apply(do *options) {}
