package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/ExecCmd"
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/orangebees/go-oneutils/PathHandle"
	"os"
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

func (c CliClient) Get(rb *RequireBase) error {
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

	//找不到，开始查找元文件
	metadata, err := LoadLocalMetadata(rb.Name, rb.Version, store)
	if err != nil {
		//找不到元文件,下载
		err = c.PkgDownload(rb)
		if err != nil {
			return err
		}
		return err
	}
	if exist {
		//找到包
		println("found", rb.GetPkgString())
		if rb.Integrity == "" {
			rb.Integrity = metadata.Integrity
		}
		return nil
	}
	println("not found pkg", rb.GetPkgString())
	//找到元文件
	err = metadata.Build(store)
	if err != nil {
		return err
	}
	//下载成功则得到元数据，开始检查hash文件是否缺失

	return nil
}

// PkgDownload 下载包
//
//	先检测包类型
func (c CliClient) PkgDownload(rb *RequireBase) error {
	println("downloading pkg", rb.GetPkgString())
	if rb.Type == "git" {
		//git版本
		err := PathHandle.RunInTempDir(func(tmppath string) error {
			err := ExecCmd.Run(tmppath, "git", "clone", "--branch", string(rb.Version), "https://"+rb.Name)
			if err != nil {
				return err
			}
			t2 := tmppath + PathHandle.Separator + rb.GetShortName()
			metadata, err := NewMetadata(rb.Name, t2, string(rb.Version), c.GitStore)
			if err != nil {
				return err
			}
			err = metadata.Save(c.GitStore)
			if err != nil {
				return err
			}
			err = metadata.Build(c.GitStore)
			if err != nil {
				return err
			}
			rb.Integrity = metadata.Integrity
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		//仓库版本

	}
	return nil
}

func (c CliClient) LoadKpmFileStruct(rb *RequireBase) (*KpmFile, error) {
	var store *GlobalStore.FileStore
	if rb.Type == "git" {
		store = kpmC.GitStore
	} else {
		store = kpmC.RegistryStore
	}
	path, err := store.GetDirPath(rb.GetPkgString())
	if err != nil {
		return nil, err
	}
	filebytes, err := os.ReadFile(path + PathHandle.Separator + "kpm.json")
	if err != nil {
		return nil, err
	}
	kf := KpmFile{}
	err = json.Unmarshal(filebytes, &kf)
	if err != nil {
		return nil, err
	}
	return &kf, nil
}

func (c CliClient) LoadKpmFileStructInWorkdir() (*KpmFile, error) {
	filebytes, err := os.ReadFile(c.WorkDir + PathHandle.Separator + "kpm.json")
	if err != nil {
		return nil, err
	}
	kf := KpmFile{}
	err = json.Unmarshal(filebytes, &kf)
	if err != nil {
		return nil, err
	}
	return &kf, nil
}

func (c CliClient) SaveKpmFileInWorkdir(kf *KpmFile) error {
	marshal, err := json.Marshal(&kf)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.WorkDir+PathHandle.Separator+"kpm.json", marshal, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
