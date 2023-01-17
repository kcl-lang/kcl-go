package kpm

import "kusionstack.io/kclvm-go"

const (
	CliHelp = `kpm  <command> [arguments]...`
)
const DefaultKclModContent = `[expected]
kclvm_version="` + kclvm.KclvmAbiVersion + `"`
