// Copyright The KCL Authors. All rights reserved.

package scripts

type (
	KclvmTripleType  string
	KclvmVersionType string
)

const (
	KclvmTripleType_unknown      KclvmTripleType = ""
	KclvmTripleType_linux_amd64  KclvmTripleType = "linux-amd64"
	KclvmTripleType_linux_arm64  KclvmTripleType = "linux-arm64"
	KclvmTripleType_darwin_amd64 KclvmTripleType = "darwin-amd64"
	KclvmTripleType_darwin_arm64 KclvmTripleType = "darwin-arm64"
	KclvmTripleType_windows      KclvmTripleType = "windows"
)

const (
	KclvmAbiVersion         KclvmVersionType = KclvmVersionType_v0_10_0
	KclvmVersionType_latest                  = KclvmVersionType_v0_10_0

	KclvmVersionType_v0_10_0 KclvmVersionType = "v0.10.0"
	KclvmVersionType_v0_9_0  KclvmVersionType = "v0.9.0"
	KclvmVersionType_v0_8_0  KclvmVersionType = "v0.8.0"
	KclvmVersionType_v0_7_5  KclvmVersionType = "v0.7.5"
	KclvmVersionType_v0_7_4  KclvmVersionType = "v0.7.4"
	KclvmVersionType_v0_7_3  KclvmVersionType = "v0.7.3"
	KclvmVersionType_v0_7_2  KclvmVersionType = "v0.7.2"
	KclvmVersionType_v0_7_1  KclvmVersionType = "v0.7.1"
	KclvmVersionType_v0_7_0  KclvmVersionType = "v0.7.0"
	KclvmVersionType_v0_6_0  KclvmVersionType = "v0.6.0"
	KclvmVersionType_v0_5_6  KclvmVersionType = "v0.5.6"
	KclvmVersionType_v0_5_5  KclvmVersionType = "v0.5.5"
	KclvmVersionType_v0_5_4  KclvmVersionType = "v0.5.4"
	KclvmVersionType_v0_5_3  KclvmVersionType = "v0.5.3"
	KclvmVersionType_v0_5_2  KclvmVersionType = "v0.5.2"
	KclvmVersionType_v0_5_1  KclvmVersionType = "v0.5.1"
	KclvmVersionType_v0_5_0  KclvmVersionType = "v0.5.0"
	KclvmVersionType_v0_4_6  KclvmVersionType = "v0.4.6"
	KclvmVersionType_v0_4_5  KclvmVersionType = "v0.4.5"
	KclvmVersionType_v0_4_4  KclvmVersionType = "v0.4.4"
	KclvmVersionType_v0_4_3  KclvmVersionType = "v0.4.3"
)
