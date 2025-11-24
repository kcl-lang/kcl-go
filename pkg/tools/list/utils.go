package list

import "strings"

// TODO: read from kcl core.
var standardSystemModules = map[string]struct{}{
	"collection": {},
	"net":        {},
	"manifests":  {},
	"math":       {},
	"datetime":   {},
	"regex":      {},
	"yaml":       {},
	"json":       {},
	"crypto":     {},
	"base64":     {},
	"units":      {},
	"file":       {},
	"template":   {},
	"runtime":    {},
}

func isBuiltinPkg(pkgpath string) bool {
	if _, ok := standardSystemModules[pkgpath]; ok {
		return true
	}
	return false
}

func isPluginPkg(pkgpath string) bool {
	return strings.HasPrefix(pkgpath, "kcl_plugin/") || strings.HasPrefix(pkgpath, "kcl_plugin.")
}
