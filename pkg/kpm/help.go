package kpm

import "kusionstack.io/kclvm-go"

const (
	CliHelp = `kpm  <command> [arguments]...`
)
const DefaultKclModContent = string(`[expected]
kclvm_version="` + kclvm.KclvmAbiVersion + `"`)
