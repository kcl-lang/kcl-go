// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	"fmt"
	"io"
	"net/rpc"
	"sync"
	"sync/atomic"

	"github.com/chai2010/protorpc"
)

var _ = fmt.Sprint

func init() {
	protorpc.UseSnappy = false
	protorpc.UseCrc32ChecksumIEEE = false
}

type Runtime struct {
	maxProc int
	exe     string
	args    []string

	stoped int32
	procs  []*_Process
	limit  chan struct{}

	wg sync.WaitGroup
	mu sync.Mutex
}

func NewRuntime(maxProc int, exe string, args ...string) *Runtime {
	return &Runtime{
		maxProc: maxProc,
		exe:     exe,
		args:    args,

		procs: make([]*_Process, maxProc),
		limit: make(chan struct{}, maxProc),
	}
}

func (p *Runtime) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, proc := range p.procs {
		if proc == nil || proc.IsExited() {
			if proc, err := createProcess(p.exe, p.args...); err == nil {
				p.procs[i] = proc
			}
		}
	}
}

func (p *Runtime) enter() { p.limit <- struct{}{} }
func (p *Runtime) leave() { <-p.limit }

func (p *Runtime) isStoped() bool { return atomic.LoadInt32(&p.stoped) != 0 }
func (p *Runtime) setStop()       { atomic.StoreInt32(&p.stoped, 1) }

func (p *Runtime) Close() error {
	p.setStop()
	defer p.wg.Wait()

	var lastErr error
	for _, proc := range p.procs {
		if err := proc.Kill(); err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}

	return nil
}

func (p *Runtime) DoTask(task func(c *rpc.Client, stderr io.Reader)) {
	if p.isStoped() {
		return
	}

	p.enter()
	defer p.leave()

	proc := p.mustGetFreeProc()
	defer p.freeProc(proc)

	p.wg.Add(1)
	defer p.wg.Done()

	task(proc.GetClient(), proc.GetStderr())
}

func (p *Runtime) mustGetFreeProc() *_Process {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, proc := range p.procs {
		if proc != nil && proc.IsFree() {
			proc.SetBusy()
			return proc
		}
	}

	for i, proc := range p.procs {
		if proc == nil || proc.IsExited() {
			if proc, err := createProcess(p.exe, p.args...); err == nil {
				p.procs[i] = proc
				proc.SetBusy()
				return proc
			}
		}
	}

	if len(p.procs) < p.maxProc {
		proc, err := createProcess(p.exe, p.args...)
		if err != nil {
			return nil
		}
		p.procs = append(p.procs, proc)
		proc.SetBusy()
		return proc
	}

	fmt.Println("runtime.Runtime.mustGetFreeProc: unreachable")
	return nil
}

func (p *Runtime) freeProc(proc *_Process) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if proc != nil {
		proc.SetFree()
	}
}
