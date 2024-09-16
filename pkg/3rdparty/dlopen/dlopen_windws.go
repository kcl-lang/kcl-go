//go:build windows
// +build windows

// Package dlopen provides some convenience functions to dlopen a library and
// get its symbols.
package dlopen

import (
	"errors"
	"unsafe"

	"syscall"
)

var ErrSoNotFound = errors.New("unable to open a handle to the library")

// LibHandle represents an open handle to a library (.so)
type LibHandle struct {
	Handle  syscall.Handle
	LibName string
}

func GetHandle(libs []string) (*LibHandle, error) {
	if len(libs) == 0 {
		return nil, ErrSoNotFound
	}
	name := libs[0]
	dll, err := syscall.LoadLibrary(name)
	if err != nil {
		return nil, err
	}
	return &LibHandle{
		Handle:  dll,
		LibName: name,
	}, nil
}

// GetSymbolPointer takes a symbol name and returns a pointer to the symbol.
func (l *LibHandle) GetSymbolPointer(symbol string) (unsafe.Pointer, error) {
	proc, err := syscall.GetProcAddress(l.Handle, symbol)
	if err != nil {
		return nil, err
	}
	return unsafe.Pointer(proc), nil
}

// Close closes a LibHandle.
func (l *LibHandle) Close() error {
	_ = syscall.FreeLibrary(l.Handle)
	return nil
}
