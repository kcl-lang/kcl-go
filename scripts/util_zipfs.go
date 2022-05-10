package scripts

import (
	"archive/zip"
	"io"
	"io/fs"
)

func ZipFs(w io.Writer, vfs fs.FS) (err error) {
	zw := zip.NewWriter(w)
	defer zw.Close()

	return fs.WalkDir(vfs, ".", func(path string, dir fs.DirEntry, errx error) error {
		if dir.IsDir() {
			return nil
		}
		dst, err := zw.Create(path)
		if err != nil {
			return errx
		}

		src, err := vfs.Open(path)
		if err != nil {
			return errx
		}
		defer src.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	})
}
