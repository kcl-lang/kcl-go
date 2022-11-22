// Package dlopen provides some convenience functions to dlopen a library and
// get its symbols.
package dlopen

var ErrSoNotFound = errors.New("unable to open a handle to the library")

// LibHandle represents an open handle to a library (.so)
type LibHandle struct {
	Handle  unsafe.Pointer
	Libname string
}

func GetHandle(libs []string) (*LibHandle, error) {
	panic("TODO: support windows")
}

// GetSymbolPointer takes a symbol name and returns a pointer to the symbol.
func (l *LibHandle) GetSymbolPointer(symbol string) (unsafe.Pointer, error) {
	panic("TODO: support windows")
}

// Close closes a LibHandle.
func (l *LibHandle) Close() error {
	panic("TODO: support windows")
}
