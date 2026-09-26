package kcl

import (
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	got, err := Ping("hello")
	if err != nil {
		t.Fatal(err)
	}
	assert2.Equal(t, "hello", got)
}
