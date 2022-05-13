// Copyright 2021 The KCL Authors. All rights reserved.

package kclvm_runtime

import (
	"bytes"
	"errors"
	"io"
	"net/rpc"
	"os/exec"

	"github.com/chai2010/protorpc"
)

type _Process struct {
	busy bool

	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr *limitBuffer
	c      *rpc.Client

	done chan error
}

// 创建新的进程, 可能失败
func createProcess(exe string, arg ...string) (p *_Process, err error) {
	p = new(_Process)

	p.cmd = exec.Command(exe, arg...)

	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	p.stderr = newLimitBuffer(10 * 1024)
	p.cmd.Stderr = p.stderr

	// 启动进程
	if err := p.cmd.Start(); err != nil {
		return nil, err
	}

	// 等待退出结果(2个缓存, 对应 Wait 和 Kill 返回值)
	p.done = make(chan error, 2)
	go func() {
		p.done <- p.cmd.Wait()
	}()

	// NewXxxServiceClient 会独占 信道(只能选择1个), 多个客户端需要手工构建 client
	conn := &procReadWriteCloser{proc: p, r: p.stdout, w: p.stdin}
	p.c = rpc.NewClientWithCodec(protorpc.NewClientCodec(conn))

	return p, nil
}

func (p *_Process) IsExited() bool { return len(p.done) > 0 }

func (p *_Process) IsFree() bool { return !p.IsExited() && !p.busy }
func (p *_Process) SetFree()     { p.busy = false }
func (p *_Process) SetBusy()     { p.busy = true }

func (p *_Process) GetClient() *rpc.Client { return p.c }
func (p *_Process) GetStderr() io.Reader   { return io.LimitReader(p.stderr, int64(p.stderr.cap)) }

func (p *_Process) Kill() error {
	if p.IsExited() {
		return nil
	}
	err := p.cmd.Process.Kill()
	p.done <- err
	return err
}

type procReadWriteCloser struct {
	proc *_Process
	r    io.ReadCloser
	w    io.WriteCloser
}

func (p *procReadWriteCloser) Read(data []byte) (n int, err error) {
	return p.r.Read(data)
}

func (p *procReadWriteCloser) Write(data []byte) (n int, err error) {
	return p.w.Write(data)
}

func (p *procReadWriteCloser) Close() error {
	return p.proc.Kill()
}

type limitBuffer struct {
	buf bytes.Buffer
	cap int
}

func newLimitBuffer(cap int) *limitBuffer {
	return &limitBuffer{cap: cap}
}

func (b *limitBuffer) Write(p []byte) (n int, err error) {
	n = b.cap - b.buf.Len()
	if n > 0 {
		b.buf.Write(p[:n])
	}
	if n < len(p) {
		err = errors.New("limitBuffer: overflow")
	}
	return n, err
}

func (b *limitBuffer) Read(p []byte) (n int, err error) {
	return b.buf.Read(p)
}

func (b *limitBuffer) String() string {
	return b.buf.String()
}
