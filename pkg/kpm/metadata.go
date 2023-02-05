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
	Name        string                  `json:"name"`
	Version     Semver.VersionString    `json:"version"`
	Integrity   GlobalStore.Integrity   `json:"integrity"`
	PackageSize int64                   `json:"package_size"`
	SubPkgName  []string                `json:"sub_pkg_name,omitempty"`
	Files       GlobalStore.FileInfoMap `json:"files,omitempty"`
}

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
func (md *Metadata) Build(gs *GlobalStore.FileStore) error {
	err := gs.BuildDir(md.Files, md.Name+"@"+string(md.Version))
	if err != nil {
		return err
	}
	return nil
}
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
