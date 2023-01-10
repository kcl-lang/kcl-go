package kpm

import (
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/orangebees/go-oneutils/PathHandle"
	"kusionstack.io/kclvm-go/scripts"
	"net/url"
	"os"
	"os/user"
)

var kpmC CliClient
var systemPkg = []string{"base64", "crypto", "json", "math", "net", "regex", "time", "units", "yaml"}

func Setup() error {
	var err error
	kpmC.WorkDir, err = os.Getwd()
	if err != nil {
		return nil
	}
	//加载环境变量
	if tmp := os.Getenv("KPM_ROOT"); tmp == "" {
		home := ""
		u, err := user.Current()
		if err != nil {
			if tmphome := os.Getenv("HOME"); tmphome != "" {
				home = tmphome
			} else {
				return nil
			}
		}
		home = u.HomeDir
		kpmC.Root = home + PathHandle.Separator + "kpm"
	}
	if tmp := os.Getenv("KPM_SERVER_ADDR"); tmp != "" {
		kpmC.RegistryAddr = tmp
	}
	parse, err := url.Parse(kpmC.RegistryAddr)
	if err != nil {
		return err
	}
	kpmC.RegistryAddrPath = parse.Host
	kpmC.GitStore, err = GlobalStore.NewFileStore(GlobalStore.FileStoreConfig{
		Root:                   kpmC.Root,
		Metadata:               "git" + PathHandle.Separator + "metadata",
		Build:                  "git" + PathHandle.Separator + "kcl_modules",
		Store:                  "store" + PathHandle.Separator + "v1" + PathHandle.Separator + "files",
		BucketCountIndexNumber: 2,
		BucketAllocationMethod: "hashStrPrefix",
		BucketHashType:         "sha512",
	}, GlobalStore.IgnoreDotGitPath)
	kpmC.RegistryStore, err = GlobalStore.NewFileStore(GlobalStore.FileStoreConfig{
		Root:                   kpmC.Root,
		Metadata:               "registry" + PathHandle.Separator + kpmC.RegistryAddrPath + PathHandle.Separator + "metadata",
		Build:                  "registry" + PathHandle.Separator + kpmC.RegistryAddrPath + PathHandle.Separator + "kcl_modules",
		Store:                  "store" + PathHandle.Separator + "v1" + PathHandle.Separator + "files",
		BucketCountIndexNumber: 2,
		BucketAllocationMethod: "hashStrPrefix",
		BucketHashType:         "sha512",
	}, GlobalStore.IgnoreDotGitPath)

	if err != nil {
		return err
	}
	kpmC.KclVmVersion = string(scripts.KclvmAbiVersion)
	return nil
}
