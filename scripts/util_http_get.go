package scripts

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

func HttpGetData(ctx context.Context, url string, insecureSkipVerify bool) (data []byte, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %v", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %v", url, err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func HttpGetFile(ctx context.Context, url, localFilename string, insecureSkipVerify bool) error {
	if ctx == nil {
		ctx = context.Background()
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", url, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", url, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(localFilename)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", url, err)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		f.Close()
		os.Remove(localFilename)
		return fmt.Errorf("failed to download %s: %v", url, err)
	}
	f.Close()
	return nil
}
