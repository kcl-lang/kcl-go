package kpm

import (
	"errors"
	"github.com/orangebees/go-oneutils/Semver"
)

type PkgString = string

func GetRequirePkgStruct(ps PkgString) (*RequirePkgStruct, error) {
	b := ps
	index1, index2 := 0, 0
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case ':':
			index1 = i
		case '@':
			index2 = i
			break
		}
	}
	if index1 == 0 || index2 == 0 || index1 >= index2 {

		return nil, errors.New("parsing failed")
	}
	if !(b[:index1] == "git" || b[:index1] == "registry") {
		return nil, errors.New("parsing failed")
	}
	return &RequirePkgStruct{
		Type:    b[:index1],
		Name:    b[index1+1 : index2],
		Version: Semver.VersionString(b[index2+1:]),
	}, nil

}
