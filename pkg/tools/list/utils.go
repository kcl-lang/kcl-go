package list

import "strings"

// TODO: read from kclvm rust.
var standardSystemModules = []string{
	"collection",
	"net",
	"manifests",
	"math",
	"datetime",
	"regex",
	"yaml",
	"json",
	"crypto",
	"base64",
	"units",
	"file",
}

func isBuiltinPkg(pkgpath string) bool {
	for _, s := range standardSystemModules {
		if s == pkgpath {
			return true
		}
	}
	return false
}

func isPluginPkg(pkgpath string) bool {
	return strings.HasPrefix(pkgpath, "kcl_plugin/") || strings.HasPrefix(pkgpath, "kcl_plugin.")
}
