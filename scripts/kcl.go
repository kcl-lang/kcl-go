// Copyright The KCL Authors. All rights reserved.

package scripts

type (
	KclTripleType  string
	KclVersionType string
)

const (
	KclTripleType_unknown      KclTripleType = ""
	KclTripleType_linux_amd64  KclTripleType = "linux-amd64"
	KclTripleType_linux_arm64  KclTripleType = "linux-arm64"
	KclTripleType_darwin_amd64 KclTripleType = "darwin-amd64"
	KclTripleType_darwin_arm64 KclTripleType = "darwin-arm64"
	KclTripleType_windows      KclTripleType = "windows"
)

const (
	KclAbiVersion         KclVersionType = KclVersionType_v0_12_3
	KclVersionType_latest                = KclVersionType_v0_12_3

	KclVersionType_v0_12_3 KclVersionType = "v0.12.3"
	KclVersionType_v0_12_2 KclVersionType = "v0.12.2"
	KclVersionType_v0_12_1 KclVersionType = "v0.12.1"
	KclVersionType_v0_12_0 KclVersionType = "v0.12.0"
	KclVersionType_v0_11_2 KclVersionType = "v0.11.2"
	KclVersionType_v0_11_1 KclVersionType = "v0.11.1"
	KclVersionType_v0_11_0 KclVersionType = "v0.11.0"
	KclVersionType_v0_10_0 KclVersionType = "v0.10.0"
	KclVersionType_v0_9_0  KclVersionType = "v0.9.0"
	KclVersionType_v0_8_0  KclVersionType = "v0.8.0"
	KclVersionType_v0_7_5  KclVersionType = "v0.7.5"
	KclVersionType_v0_7_4  KclVersionType = "v0.7.4"
	KclVersionType_v0_7_3  KclVersionType = "v0.7.3"
	KclVersionType_v0_7_2  KclVersionType = "v0.7.2"
	KclVersionType_v0_7_1  KclVersionType = "v0.7.1"
	KclVersionType_v0_7_0  KclVersionType = "v0.7.0"
	KclVersionType_v0_6_0  KclVersionType = "v0.6.0"
	KclVersionType_v0_5_6  KclVersionType = "v0.5.6"
	KclVersionType_v0_5_5  KclVersionType = "v0.5.5"
	KclVersionType_v0_5_4  KclVersionType = "v0.5.4"
	KclVersionType_v0_5_3  KclVersionType = "v0.5.3"
	KclVersionType_v0_5_2  KclVersionType = "v0.5.2"
	KclVersionType_v0_5_1  KclVersionType = "v0.5.1"
	KclVersionType_v0_5_0  KclVersionType = "v0.5.0"
	KclVersionType_v0_4_6  KclVersionType = "v0.4.6"
	KclVersionType_v0_4_5  KclVersionType = "v0.4.5"
	KclVersionType_v0_4_4  KclVersionType = "v0.4.4"
	KclVersionType_v0_4_3  KclVersionType = "v0.4.3"
)
