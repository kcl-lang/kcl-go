package scripts

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func UnTarGz(tarGzFile, outputDir string) error {
	os.MkdirAll(outputDir, 0666)

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
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return fmt.Errorf("UnTarGz: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
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
