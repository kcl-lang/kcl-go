package kpm

import (
	"github.com/orangebees/go-oneutils/GlobalStore"
	"github.com/orangebees/go-oneutils/Semver"
)

type Require struct {
	RequireBase
	Alias string `json:"alias,omitempty"`
}

type RequireBase struct {
	RequirePkgStruct
	Integrity GlobalStore.Integrity `json:"integrity"`
}

type RequirePkgStruct struct {
	Type    string               `json:"type"`
	Name    string               `json:"name"`
	Version Semver.VersionString `json:"version"`
}

func (rps *RequirePkgStruct) GetPkgString() PkgString {
	return rps.Type + ":" + rps.Name + "@" + string(rps.Version)
}
func (rps *RequirePkgStruct) GetShortName() string {
	for i := len(rps.Name) - 1; i >= 0; i-- {
		if rps.Name[i] == '/' {
			return rps.Name[i+1:]
		}
	}
	return ""
}
