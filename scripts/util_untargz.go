package scripts

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnTarGz(tarGzFile, trimPrefix, outputDir string) error {
	os.MkdirAll(outputDir, 0777)

	gzipStream, err := os.Open(tarGzFile)
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("UnTarGz: Next() failed: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			path := filepath.Join(outputDir, strings.TrimPrefix(header.Name, trimPrefix))
			if err := os.MkdirAll(path, 0777); err != nil {
				return fmt.Errorf("UnTarGz: MkdirAll() failed: %w", err)
			}
		case tar.TypeReg:
			path := filepath.Join(outputDir, strings.TrimPrefix(header.Name, trimPrefix))
			outFile, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("UnTarGz: Create() failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("UnTarGz: Copy() failed: %w", err)
			}
			outFile.Close()

		default:
			return fmt.Errorf(
				"UnTarGz: unknown type: %v in %v",
				header.Typeflag,
				header.Name,
			)
		}

	}
	return nil
}
