package kpm

import (
	"testing"
)

func TestInit(t *testing.T) {
	err := Setup()
	if err != nil {
		println(err.Error())
		return
	}
	// kpm init kkk
	err = CLI([]string{"init", "kkk"}...)
	if err != nil {
		println(err.Error())
		return
	}
}
func TestAdd(t *testing.T) {
	err := Setup()
	if err != nil {
		println(err.Error())
		return
	}
	// kpm add -git github.com/orangebees/konfig@v0.0.1
	err = CLI([]string{"add", "-git", "github.com/orangebees/konfig@v0.0.1"}...)
	if err != nil {
		println(err.Error())
		return
	}
}

func TestDownload(t *testing.T) {
	err := Setup()
	if err != nil {
		println(err.Error())
		return
	}
	// kpm download
	err = CLI([]string{"download"}...)
	if err != nil {
		println(err.Error())
		return
	}
}
func TestDel(t *testing.T) {
	err := Setup()
	if err != nil {
		println(err.Error())
		return
	}
	// kpm del konfig
	err = CLI([]string{"del", "konfig"}...)
	if err != nil {
		println(err.Error())
		return
	}
}
func TestStore(t *testing.T) {
	err := Setup()
	if err != nil {
		println(err.Error())
		return
	}
	// kpm store add -git github.com/orangebees/konfig@v0.0.1
	err = CLI([]string{"store", "add", "-git", "github.com/orangebees/konfig@v0.0.1"}...)
	if err != nil {
		println(err.Error())
		return
	}

}
