package kpm

import (
	"errors"
	"github.com/orangebees/go-oneutils/PathHandle"
)

type PkgString string

//konfig@v0.0.1
//konfig@v0.0.0#bbbbbb

func New() (ps PkgString) {
	return ""
}

func (ps PkgString) GetRequirePkgString() (*RequirePkgString, error) {
	b := string(ps)
	index1, index2 := 0, 0
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case ':':
			index1 = i
		case '@':
			index2 = i
		}
	}

	if index1 == 0 || index2 == 0 || index1 <= index2 {
		return nil, errors.New("parsing failed")
	}
	return &RequirePkgString{
		Type:    b[:index1],
		Name:    b[index1:index2],
		Version: PkgVersion(b[index2:]),
	}, nil

}
func (rps *RequirePkgString) Verify() error {
	if rps == nil {
		return errors.New("RequirePkgString is nil")
	}

	return nil
}

// VerifyExist VerifyLocal 验证这个包是否存在
func (c CliClient) VerifyExist(rps *RequirePkgString) error {

	PathHandle.URLToLocalDirPath(string(rps.GetPkgString()))

	return nil
}
