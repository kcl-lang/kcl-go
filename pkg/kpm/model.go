package kpm

import (
	"github.com/orangebees/go-oneutils/GlobalStore"
)

type CliClient struct {
	GitStore         *GlobalStore.FileStore
	RegistryStore    *GlobalStore.FileStore
	WorkDir          string
	Root             string
	RegistryAddr     string
	RegistryAddrPath string
	KclVmVersion     string
}

func (c CliClient) Get(rb RequireBase) error {
	var store *GlobalStore.FileStore
	if rb.Type == "git" {
		store = c.GitStore
	} else {
		store = c.RegistryStore
	}
	exist, err := store.DirIsExist(rb.Name + "@" + string(rb.Version))
	if err != nil {
		return err
	}
	if exist {
		//找到包
		println("found", rb.GetPkgString())
		return nil
	}
	println("not found pkg", rb.GetPkgString())
	//找不到，开始查找元文件
	metadata, err := LoadLocalMetadata(rb.Name, string(rb.Version), store)
	if err != nil {
		//找不到元文件,下载
		return err
	}
	err = metadata.Build(store)
	if err != nil {
		return err
	}
	//如果元文件找不到，则下载
	//下载成功则得到元数据，开始检查hash文件是否缺失

	return nil
}

//func (r *Require) Get(kpmroot, kpmserver string) error {
//	kpmserverurl, err := url.Parse(kpmserver)
//	if err != nil {
//		return err
//	}
//	kpmserverpath := kpmserverurl.Host
//	//检测包目录是否存在，如果不存在则使用本地元文件构建，如果没有元文件，则下载
//	if r.IsInLocal(kpmroot, kpmserverpath) != nil {
//		println("not found pkg", r.ToString())
//		if r.PkgInfoIsInLocal(kpmroot, kpmserverpath) != nil {
//			println("not found pkginfo", r.ToString())
//			err = r.PkgDownload(kpmroot, kpmserver)
//			if err != nil {
//				return err
//			}
//			println("downloading", r.ToString())
//		}
//		println("building", r.ToString())
//		err = r.Build(kpmroot, kpmserverpath)
//		if err != nil {
//			return err
//		}
//	} else {
//		if r.PkgInfoIsInLocal(kpmroot, kpmserverpath) != nil {
//
//		}
//		println("found", r.ToString())
//	}
//
//	return nil
//}

func (c CliClient) PkgDownload(rb RequireBase) error {

	return nil
}
func (c CliClient) Build(rb RequireBase) error {
	if rb.Type == "git" {

	}
	return nil
}
func (c CliClient) PkgInfoIsInLocal(rb RequireBase) error {
	c.GitStore.GetMetadataPath(rb.Name)
	return nil
}
