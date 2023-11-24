package testing

import "io"

type indentingWriter struct {
	w io.Writer
}

func NewIndentingWriter(w io.Writer) indentingWriter {
	return indentingWriter{
		w: w,
	}
}

func (w indentingWriter) Write(bs []byte) (int, error) {
	var written int
	indent := true
	for _, b := range bs {
		if indent {
			wrote, err := w.w.Write([]byte("  "))
			if err != nil {
				return written, err
			}
			written += wrote
		}
		wrote, err := w.w.Write([]byte{b})
		if err != nil {
			return written, err
		}
		written += wrote
		indent = b == '\n'
	}
	return written, nil
}
