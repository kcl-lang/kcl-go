// Copyright 2021 The KCL Authors. All rights reserved.

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
	KclvmAbiVersion         KclvmVersionType = KclvmVersionType_v0_5_0
	KclvmVersionType_latest                  = KclvmVersionType_v0_5_0

	KclvmVersionType_v0_5_0         KclvmVersionType = "v0.5.0"
	KclvmVersionType_v0_5_0_beta_1  KclvmVersionType = "v0.5.0-beta.1"
	KclvmVersionType_v0_5_0_alpha_4 KclvmVersionType = "v0.5.0-alpha.4"
	KclvmVersionType_v0_5_0_alpha_3 KclvmVersionType = "v0.5.0-alpha.3"
	KclvmVersionType_v0_5_0_alpha_2 KclvmVersionType = "v0.5.0-alpha.2"
	KclvmVersionType_v0_5_0_alpha_1 KclvmVersionType = "v0.5.0-alpha.1"
	KclvmVersionType_v0_4_6         KclvmVersionType = "v0.4.6"
	KclvmVersionType_v0_4_5         KclvmVersionType = "v0.4.5"
	KclvmVersionType_v0_4_5_alpha_2 KclvmVersionType = "v0.4.5-alpha.2"
	KclvmVersionType_v0_4_5_alpha_1 KclvmVersionType = "v0.4.5-alpha.1"
	KclvmVersionType_v0_4_4         KclvmVersionType = "v0.4.4"
	KclvmVersionType_v0_4_4_beta_2  KclvmVersionType = "v0.4.4-beta.2"
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
