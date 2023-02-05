package kpm

import (
	"testing"
)

var setupFlag bool

func TestInit(t *testing.T) {
	if !setupFlag {
		err := Setup()
		if err != nil {
			println(err.Error())
			return
		}
		setupFlag = true
	}

	// kpm init kkk
	err := CLI([]string{"init", "kkk"}...)
	if err != nil {
		println(err.Error())
		return
	}
}
func TestAdd(t *testing.T) {
	if !setupFlag {
		err := Setup()
		if err != nil {
			println(err.Error())
			return
		}
		setupFlag = true
	}
	// kpm add -git github.com/orangebees/konfig@v0.0.1
	err := CLI([]string{"add", "-git", "github.com/orangebees/konfig@v0.0.1"}...)
	if err != nil {
		println(err.Error())
		return
	}
}

func TestDownload(t *testing.T) {
	if !setupFlag {
		err := Setup()
		if err != nil {
			println(err.Error())
			return
		}
		setupFlag = true
	}
	// kpm download
	err := CLI([]string{"download"}...)
	if err != nil {
		println(err.Error())
		return
	}
}
func TestDel(t *testing.T) {
	if !setupFlag {
		err := Setup()
		if err != nil {
			println(err.Error())
			return
		}
		setupFlag = true
	}
	// kpm del konfig
	err := CLI([]string{"del", "konfig"}...)
	if err != nil {
		println(err.Error())
		return
	}
}
func TestStore(t *testing.T) {
	if !setupFlag {
		err := Setup()
		if err != nil {
			println(err.Error())
			return
		}
		setupFlag = true
	}
	// kpm store add -git github.com/orangebees/konfig@v0.0.1
	err := CLI([]string{"store", "add", "-git", "github.com/orangebees/konfig@v0.0.1"}...)
	if err != nil {
		println(err.Error())
		return
	}

}
