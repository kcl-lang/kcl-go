package kpm

import (
	"encoding/json"
	"github.com/orangebees/go-oneutils/Convert"
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/orangebees/go-oneutils/Semver"
	"github.com/orangebees/go-oneutils/Set"
	"os"
	"sort"
	"strings"
)

type Metadata struct {
	//包名
	Name string `json:"name"`
	//版本
	Version Semver.VersionString `json:"version"`
	//包完整性校验信息
	Integrity GlobalStore.Integrity `json:"integrity"`
	//包大小
	PackageSize int64 `json:"package_size"`
	//子包名称
	SubPkgName []string `json:"sub_pkg_name,omitempty"`
	//文件hash
	Files GlobalStore.FileInfoMap `json:"files,omitempty"`
}

// NewMetadata 生成新的包的元数据
func NewMetadata(pkgName, pkgPath string, pkgVersion string, gs *GlobalStore.FileStore) (*Metadata, error) {
	fim, err := gs.AddDir(pkgPath)
	if err != nil {
		return nil, err
	}
	ver, err := Semver.NewFromString(pkgVersion)
	if err != nil {
		return nil, err
	}
	m := Metadata{
		Name:       pkgName,
		Version:    Semver.VersionString(ver.TagString()),
		Integrity:  fim.GetIntegrity(),
		SubPkgName: nil,
		Files:      fim,
	}
	set := Set.AcquireSet()
	defer Set.ReleaseSet(set)
	for k, info := range fim {
		//计算包大小
		m.PackageSize += info.Size
		//添加子包名字
		//
		if strings.HasSuffix(k, ".k") {
			tmp := make([]byte, len(k))
			tmp = tmp[:0]
			index := 0
			for i := 0; i < len(k); i++ {
				if k[i] == '/' {
					index = i
					tmp = append(tmp, '.')
				} else {
					tmp = append(tmp, k[i])
				}
			}
			set.SAdd(Convert.B2S(tmp[:index]))
		}
		m.SubPkgName = set.SMembers()
		sort.Strings(m.SubPkgName)
	}
	return &m, nil
}

// LoadLocalMetadata 加载本地元数据
func LoadLocalMetadata(pkgName string, pkgVersion Semver.VersionString, gs *GlobalStore.FileStore) (*Metadata, error) {
	path, err := gs.GetMetadataPath(pkgName + "@" + string(pkgVersion))
	if err != nil {
		return nil, err
	}
	path += ".json"
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
	err := gs.BuildDir(md.Files, md.Name+"@"+string(md.Version))
	if err != nil {
		return err
	}
	return nil
}

// Save 保存元数据
func (md *Metadata) Save(gs *GlobalStore.FileStore) error {
	path, err := gs.GetMetadataPath(md.Name + "@" + string(md.Version))
	if err != nil {
		return err
	}
	path = path + ".json"
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
