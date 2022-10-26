package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/Convert"
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/orangebees/go-oneutils/PathHandle"
	"os"
	"strings"
)

type Metadata struct {
	Name        string
	Version     string
	Integrity   GlobalStore.Integrity
	PackageSize int64
	SubPkgName  []string
	Files       GlobalStore.FileInfoMap
}

// NewMetadata 生成新的包的元数据
func NewMetadata(pkgName, pkgVersion, pkgPath string, gs *GlobalStore.FileStore) (*Metadata, error) {
	fim, err := gs.AddDir(pkgPath)
	if err != nil {
		return nil, err
	}
	m := Metadata{
		Name:       pkgName,
		Version:    pkgVersion,
		Integrity:  fim.GetIntegrity(),
		SubPkgName: nil,
		Files:      fim,
	}
	for k, info := range fim {
		//计算包大小
		m.PackageSize += info.Size
		//添加子包名字
		if strings.HasSuffix(k, ".k") {
			tmps := strings.Split(k, ".")
			tmpslen := len(tmps)
			if tmpslen == 1 {
				continue
			}
			tmp := make([]byte, len(k))
			tmp = tmp[:0]
			for i := 0; i < tmpslen; i++ {
				if i != 0 {
					tmp = append(tmp, '.')
				}
				tmp = append(tmp, tmps[i]...)
			}
			m.SubPkgName = append(m.SubPkgName, Convert.B2S(tmp))
		}
	}
	return &m, nil
}

// LoadLocalMetadata 加载本地元数据
func LoadLocalMetadata(pkgName, pkgVersion string, gs *GlobalStore.FileStore) (*Metadata, error) {
	path, err := gs.GetMetadataPath(PathHandle.URLToLocalDirPath(pkgName + "@" + pkgVersion))
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	md := Metadata{}
	err = json.Unmarshal(file, &md)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

// Build 通过元数据构造包
func (md *Metadata) Build(gs *GlobalStore.FileStore) error {
	err := gs.BuildDir(md.Files, PathHandle.URLToLocalDirPath(md.Name)+"@"+md.Version)
	if err != nil {
		return err
	}
	return nil
}

// Save 保存元数据
func (md *Metadata) Save(gs *GlobalStore.FileStore) error {
	path, err := gs.GetMetadataPath(PathHandle.URLToLocalDirPath(md.Name + "@" + md.Version))
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(md)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, marshal, 0644)
	if err != nil {
		return err
	}
	return nil
}
