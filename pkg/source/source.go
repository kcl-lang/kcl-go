package source

import (
	"fmt"
	"io"
	"os"
)

func ReadSource(filename string, src any) (data []byte, err error) {
	if src == nil {
		src, err = os.ReadFile(filename)
		if err != nil {
			return
		}
	}
	switch src := src.(type) {
	case []byte:
		return src, nil
	case string:
		return []byte(src), nil
	case io.Reader:
		d, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported src type: %T", src)
	}
}
