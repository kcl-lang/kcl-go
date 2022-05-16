// Copyright 2021 The KCL Authors. All rights reserved.

package scripts

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type (
	KclvmTripleType  string
	KclvmVersionType string
)

const (
	KclvmTripleType_unknown      KclvmTripleType = ""
	KclvmTripleType_centos       KclvmTripleType = "centos"
	KclvmTripleType_ubuntu       KclvmTripleType = "ubuntu"
	KclvmTripleType_darwin       KclvmTripleType = "Darwin"
	KclvmTripleType_darwin_arm64 KclvmTripleType = "Darwin-arm64"
	KclvmTripleType_windows      KclvmTripleType = "windows"
)

const (
	KclvmAbiVersion KclvmVersionType = "v0.4.2"

	KclvmVersionType_v0_4_2_alpha_1 KclvmVersionType = "v0.4.2-alpha.1"
	KclvmVersionType_latest                          = KclvmVersionType_v0_4_2_alpha_1
)

var (
	DefaultKclvmTriple  KclvmTripleType  = getKclvmTripleType(runtime.GOOS, runtime.GOARCH)
	DefaultKclvmVersion KclvmVersionType = KclvmVersionType_latest

	KclvmTripleList = []KclvmTripleType{
		KclvmTripleType_darwin,
		KclvmTripleType_ubuntu,
		KclvmTripleType_centos,

		// KclvmTripleType_darwin_arm64,
		// KclvmTripleType_windows,
	}

	KclvmVersionList = []KclvmVersionType{
		DefaultKclvmVersion,
	}

	// triple: centos, Darwin, Darwin-arm64, windows
	// Linux:   {baseUrl}/{version}/kclvm-{version}-{triple}.tar.gz
	// macOS:   {baseUrl}/{version}/kclvm-{version}-{triple}.tar.gz
	// Windows: {baseUrl}/{version}/kclvm-{version}-{triple}.zip

	KclvmDownloadUrlBase         = "https://github.com/KusionStack/KCLVM/releases/download/"
	KclvmDownloadUrlBase_mirrors = []string{}

	KclvmMd5sum = map[string]string{
		"kclvm-v0.4.2-alpha.1-ubuntu.tar.gz": "2aa6fba3f4d3466b660ee8fc4ca65bff",
		"kclvm-v0.4.2-alpha.1-Darwin.tar.gz": "16015b02d6b490d9091b194e5829f1c4",
		"kclvm-v0.4.2-alpha.1-centos.tar.gz": "c94f3adc1d4cd9c3aa4df3e55775d7d8",
	}
)

func SetupKclvmAll(outdir string) error {
	defaultBackup := DefaultKclvmTriple
	defer func() {
		DefaultKclvmTriple = defaultBackup
	}()

	for _, triple := range KclvmTripleList {
		DefaultKclvmTriple = triple
		root := filepath.Join(outdir, string(triple))

		err := SetupKclvm(root)
		if err != nil {
			return err
		}

		fmt.Println(root, "ok")
	}

	return nil
}

func SetupKclvm(kclvmRoot string) error {
	kclvmAsset := NewKclvmAssetHelper(DefaultKclvmTriple, DefaultKclvmVersion)
	return kclvmAsset.Install(kclvmRoot)
}

func getKclvmTripleType(goos, goarch string) KclvmTripleType {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return KclvmTripleType_darwin
		case "arm64":
			return KclvmTripleType_darwin // todo: KclvmTripleType_darwin_arm64
		}
	case "linux":
		return KclvmTripleType_ubuntu
	case "windows":
		return KclvmTripleType_windows
	}
	return KclvmTripleType_unknown
}
