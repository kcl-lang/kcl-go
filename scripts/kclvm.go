// Copyright 2021 The KCL Authors. All rights reserved.

package scripts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	kcl_plugin "kusionstack.io/kcl-plugin"
)

const (
	KclvmDownloadUrlBase        = "https://github.com/KusionStack/KCLVM/releases/download/"
	KclvmDownloadUrlBase_mirror = ""
)

const KclvmVersion = "0.4.1-alpha.4"

var KclvmVersionList = []string{
	"0.4.1-alpha.4",
}

var KclvmMd5sum = map[string]string{
	// 0.4.1-alpha.4
	"kclvm-centos-0.4.1-alpha.4.tar.gz":       "5329374c2cb336f34cacc4e088b88496",
	"kclvm-Darwin-0.4.1-alpha.4.tar.gz":       "409da9310cbcf5a7ef38c1895112f3ae",
	"kclvm-Darwin-arm64-0.4.1-alpha.4.tar.gz": "7dc7f293ec45870a75d49e5f5d6fd2d5",
	"kclvm-ubuntu-0.4.1-alpha.4.tar.gz":       "809f8a2f5b7721bee773457a03abfe90",
}

func SetupKclvm(kclvmRoot string) error {
	if err := InstallKclvm(kclvmRoot); err != nil {
		return err
	}

	kclvmPluginsPath := getPluginPath(kclvmRoot)
	if err := kcl_plugin.InstallPlugins(kclvmPluginsPath); err != nil {
		return err
	}

	return nil
}

func InstallKclvm(kclvmRoot string) error {
	if runtime.GOOS == "windows" {
		kclvmExe := JoinedPath(kclvmRoot, "kclvm.exe")
		if FileExists(kclvmExe) {
			return nil
		}
	} else {
		kclvmExe := JoinedPath(kclvmRoot, "bin", "kclvm")
		if FileExists(kclvmExe) {
			return nil
		}
	}

	var triple = GetKclvmTriple()
	var localFilename = "zz_kclvm.download.dat"

	defer os.Remove(localFilename)
	if err := DownloadKclvm(triple, localFilename); err != nil {
		return err
	}

	if strings.Contains(triple, "windows") {
		return Unzip(localFilename, kclvmRoot)
	} else {
		return UnTarGz(localFilename, kclvmRoot)
	}
}

func GetKclvmTriple() string {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "kclvm-Darwin"
		case "arm64":
			return "kclvm-Darwin-arm64"
		}
	case "linux":
		return "kclvm-ubuntu"
	case "windows":
		return "kclvm-windows"
	}
	return ""
}

func DownloadKclvm(triple, localFilename string) error {
	if triple == "" {
		triple = GetKclvmTriple()
	}
	if triple == "" {
		return fmt.Errorf("triple missing")
	}

	ext := "tar.gz"
	if strings.Contains(strings.ToLower(triple), "windows") {
		ext = "zip"
	}

	kclvmFilename := fmt.Sprintf("%s-%s.%s", triple, KclvmVersion, ext)
	md5sum := KclvmMd5sum[kclvmFilename]

	if md5sum == "" {
		return fmt.Errorf("%s: not found", kclvmFilename)
	}
	if MD5File(localFilename) == md5sum {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var urls = []string{KclvmDownloadUrlBase + kclvmFilename}
	if KclvmDownloadUrlBase_mirror != "" {
		urls = append(urls, KclvmDownloadUrlBase_mirror+kclvmFilename)
	}

	var errs = make(chan error, len(urls))
	var okfiles = make(chan string, len(urls))
	var wg sync.WaitGroup

	wg.Add(len(urls))
	for i, s := range urls {
		go func(id int, url, localFilename string) {
			defer wg.Done()
			tmpname := fmt.Sprintf("%s.%d", localFilename, id)
			if err := HttpGetFile(ctx, url, tmpname); err != nil {
				errs <- err
			} else {
				okfiles <- tmpname
				cancel()
			}
		}(i, s, localFilename)
	}
	wg.Wait()

	if len(okfiles) > 0 {
		tmpname := <-okfiles
		os.Rename(tmpname, localFilename)

		for id := range urls {
			tmpname := fmt.Sprintf("%s.%d", localFilename, id)
			os.Remove(tmpname)
		}

		if got := MD5File(localFilename); got != md5sum {
			return fmt.Errorf("md4 mismatch: expect=%v, got=%v", md5sum, got)
		}

		return nil
	}

	return <-errs
}

func getPluginPath(kclvmRoot string) string {
	if runtime.GOOS == "windows" {
		kclvmPluginPath := filepath.Join(kclvmRoot, "bin", "plugins")
		return kclvmPluginPath
	}
	kclvmPluginPath := filepath.Join(kclvmRoot, "plugins")
	return kclvmPluginPath
}
