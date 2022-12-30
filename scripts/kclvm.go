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
	KclvmAbiVersion         KclvmVersionType = "v0.4.4"
	KclvmVersionType_latest                  = KclvmVersionType_v0_4_4_beta_1

	KclvmVersionType_v0_4_4_beta_1  KclvmVersionType = "v0.4.4-beta.1"
	KclvmVersionType_v0_4_4_alpha_1 KclvmVersionType = "v0.4.4-alpha.1"
	KclvmVersionType_v0_4_3         KclvmVersionType = "v0.4.3"
	KclvmVersionType_v0_4_3_alpha_1 KclvmVersionType = "v0.4.3-alpha.1"
	KclvmVersionType_v0_4_2_alpha_5 KclvmVersionType = "v0.4.2-alpha.5"
	KclvmVersionType_v0_4_2_alpha_4 KclvmVersionType = "v0.4.2-alpha.4"
	KclvmVersionType_v0_4_2_alpha_3 KclvmVersionType = "v0.4.2-alpha.3"
	KclvmVersionType_v0_4_2_alpha_2 KclvmVersionType = "v0.4.2-alpha.2"
	KclvmVersionType_v0_4_2_alpha_1 KclvmVersionType = "v0.4.2-alpha.1"
)

var (
	DefaultKclvmTriple  KclvmTripleType  = getKclvmTripleType(runtime.GOOS, runtime.GOARCH)
	DefaultKclvmVersion KclvmVersionType = KclvmVersionType_latest

	KclvmTripleList = []KclvmTripleType{
		KclvmTripleType_darwin,
		KclvmTripleType_ubuntu,
		KclvmTripleType_centos,
		KclvmTripleType_darwin_arm64,
		KclvmTripleType_windows,
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
		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.4-beta.1
		"kclvm-v0.4.4-beta.1-Darwin.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.4-beta.1-Darwin-arm64.tar.gz": "", // read from *.md5.txt
		"kclvm-v0.4.4-beta.1-centos.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.4-beta.1-ubuntu.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.4-beta.1-windows.zip":         "", // read from *.md5.txt

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.4-alpha.1
		"kclvm-v0.4.4-alpha.1-Darwin.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.4-alpha.1-Darwin-arm64.tar.gz": "", // read from *.md5.txt
		"kclvm-v0.4.4-alpha.1-centos.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.4-alpha.1-ubuntu.tar.gz":       "", // read from *.md5.txt

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.3
		"kclvm-v0.4.3-Darwin.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.3-Darwin-arm64.tar.gz": "", // read from *.md5.txt
		"kclvm-v0.4.3-centos.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.3-ubuntu.tar.gz":       "", // read from *.md5.txt

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.3-alpha.1
		"kclvm-v0.4.3-alpha.1-Darwin.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.3-alpha.1-Darwin-arm64.tar.gz": "", // read from *.md5.txt
		"kclvm-v0.4.3-alpha.1-centos.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.3-alpha.1-ubuntu.tar.gz":       "", // read from *.md5.txt

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.2-alpha.5
		"kclvm-v0.4.2-alpha.5-Darwin.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.2-alpha.5-Darwin-arm64.tar.gz": "", // read from *.md5.txt
		"kclvm-v0.4.2-alpha.5-centos.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.2-alpha.5-ubuntu.tar.gz":       "", // read from *.md5.txt
		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.2-alpha.4
		"kclvm-v0.4.2-alpha.4-Darwin.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.2-alpha.4-Darwin-arm64.tar.gz": "", // read from *.md5.txt
		"kclvm-v0.4.2-alpha.4-centos.tar.gz":       "", // read from *.md5.txt
		"kclvm-v0.4.2-alpha.4-ubuntu.tar.gz":       "", // read from *.md5.txt

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.2-alpha.3
		"kclvm-v0.4.2-alpha.3-Darwin.tar.gz":       "9727d804b49f225682af9b7383c0ab6a",
		"kclvm-v0.4.2-alpha.3-Darwin-arm64.tar.gz": "bba2e8c40491d4305770ef7e67822cd6",
		"kclvm-v0.4.2-alpha.3-centos.tar.gz":       "afe52170ccd3b01ffefa48b73ac655d1",
		"kclvm-v0.4.2-alpha.3-ubuntu.tar.gz":       "a4a7fd7ced93cfb081de3da1207ce641",

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.2-alpha.2
		"kclvm-v0.4.2-alpha.2-Darwin.tar.gz":       "b9c4466c4fdedaade751633d6143bf62",
		"kclvm-v0.4.2-alpha.2-Darwin-arm64.tar.gz": "bf023bc212f9951532d6c08934f5783b",
		"kclvm-v0.4.2-alpha.2-centos.tar.gz":       "8dc2aa00f87ef974921a6094988caa0d",
		"kclvm-v0.4.2-alpha.2-ubuntu.tar.gz":       "35c12b2605e15b93ad053327f18db702",

		// https://github.com/KusionStack/KCLVM/releases/tag/v0.4.2-alpha.1
		"kclvm-v0.4.2-alpha.1-Darwin.tar.gz":       "16015b02d6b490d9091b194e5829f1c4",
		"kclvm-v0.4.2-alpha.1-Darwin-arm64.tar.gz": "1e31cd3c1061e5e18da2fbc460f64472",
		"kclvm-v0.4.2-alpha.1-centos.tar.gz":       "c94f3adc1d4cd9c3aa4df3e55775d7d8",
		"kclvm-v0.4.2-alpha.1-ubuntu.tar.gz":       "2aa6fba3f4d3466b660ee8fc4ca65bff",
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
			return KclvmTripleType_darwin_arm64
		}
	case "linux":
		return KclvmTripleType_ubuntu
	case "windows":
		return KclvmTripleType_windows
	}
	return KclvmTripleType_unknown
}
