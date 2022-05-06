// Copyright 2021 The KCL Authors. All rights reserved.

package ktest

import (
	"regexp"
)

type Options struct {
	RunRegexp string
	QuietMode bool
	Verbose   bool
	Debug     bool
}

func (p *Options) shouldRun(name string) bool {
	if p.RunRegexp != "" {
		matched, _ := regexp.MatchString(p.RunRegexp, name)
		return matched
	}
	return true
}
