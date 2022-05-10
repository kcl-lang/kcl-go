// Copyright 2021 The KCL Authors. All rights reserved.

package scripts

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	kcl_plugin "kusionstack.io/kcl-plugin"
)

const KclvmDownloadUrlBase = "https://github.com/KusionStack/KCLVM/releases/download/"

var DefaultKclvmTriple = getKclvmTriple()

var KclvmDownloadUrlBase_mirrors = []string{
	// test: python3 -m http.server => http://127.0.0.1:8000/

	// "http://127.0.0.1:8000/downloads",
}

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

var KclvmTripleList = []string{
	"kclvm-centos",
	"kclvm-Darwin",
	"kclvm-Darwin-arm64",
	"kclvm-ubuntu",
}

func SetupKclvmAll(outdir string) error {
	defaultBackup := DefaultKclvmTriple
	defer func() {
		DefaultKclvmTriple = defaultBackup
	}()

	for _, triple := range KclvmTripleList {
		DefaultKclvmTriple = triple
		root := filepath.Join(outdir, triple)

		err := SetupKclvm(root)
		if err != nil {
			return err
		}

		fmt.Println(root, "ok")
	}

	return nil
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

func InstallKclvm(kclvmRoot string) (err error) {
	md5sumFile := JoinedPath(kclvmRoot, "md5sum.txt")
	if FileExists(md5sumFile) {
		return nil
	}

	var triple = DefaultKclvmTriple
	var localFilename = "zz_download-" + GetKclvmFilename(triple)
	defer func() {
		if err == nil {
			os.Remove(localFilename)
		}
	}()

	if err := DownloadKclvm(triple, localFilename); err != nil {
		return err
	}

	if strings.HasSuffix(localFilename, ".zip") {
		if err := Unzip(localFilename, kclvmRoot); err != nil {
			return err
		}
	} else {
		if err := UnTarGz(localFilename, "kclvm", kclvmRoot); err != nil {
			return err
		}
	}

	// write md5sum
	if s := filepath.Join(kclvmRoot, "md5sum.txt"); !FileExists(s) {
		txt := fmt.Sprintf("%s *%s\n", GetKclvmMd5um(triple), GetKclvmFilename(triple))
		if err := ioutil.WriteFile(s, []byte(txt), 0666); err != nil {
			return err
		}
	}

	// write VERSION
	if s := filepath.Join(kclvmRoot, "VERSION"); !FileExists(s) {
		if err := ioutil.WriteFile(s, []byte(KclvmVersion), 0666); err != nil {
			return err
		}
	}

	return nil
}

func getKclvmTriple() string {
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

func GetKclvmFilename(triple string) string {
	ext := "tar.gz"
	if strings.Contains(strings.ToLower(triple), "windows") {
		ext = "zip"
	}
	return fmt.Sprintf("%s-%s.%s", triple, KclvmVersion, ext)
}

func GetKclvmMd5um(triple string) string {
	kclvmFilename := GetKclvmFilename(triple)
	return KclvmMd5sum[kclvmFilename]
}

func DownloadKclvm(triple, localFilename string) error {
	if triple == "" {
		triple = DefaultKclvmTriple
	}
	if triple == "" {
		return fmt.Errorf("triple missing")
	}

	kclvmFilename := GetKclvmFilename(triple)
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

	for _, s := range KclvmDownloadUrlBase_mirrors {
		mirrorBase := strings.TrimSpace(s)
		if mirrorBase != "" {
			if !strings.HasSuffix(mirrorBase, "/") {
				mirrorBase += "/"
			}
			urls = append(urls, mirrorBase+kclvmFilename)
		}
	}

	var errs = make(chan error, len(urls))
	var okfiles = make(chan string, len(urls))
	var wg sync.WaitGroup

	wg.Add(len(urls))
	for i, s := range urls {
		go func(id int, url, localFilename string) {
			defer wg.Done()
			tmpname := fmt.Sprintf("%s.%d", localFilename, id)
			err := HttpGetFile(ctx, url, tmpname)
			if err != nil {
				errs <- err
				return
			}
			if got := MD5File(tmpname); got != md5sum {
				errs <- fmt.Errorf("md5 mismatch: expect=%v, got=%v, local=%s", md5sum, got, localFilename)
				return
			}

			// OK
			okfiles <- tmpname
			cancel()
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
			return fmt.Errorf("md5 mismatch: expect=%v, got=%v, local=%s", md5sum, got, localFilename)
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
