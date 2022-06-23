package list

import "strings"

var standardSystemModules = []string{
	"collection",
	"net",
	"math",
	"datetime",
	"regex",
	"yaml",
	"json",
	"crypto",
	"base64",
	"testing",
	"units",
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
	return strings.HasPrefix(pkgpath, "kcl_plugin/")
}
