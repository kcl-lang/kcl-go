//go:build cgo
// +build cgo

package native

/*
#include <stdlib.h>
#include <stdint.h>
typedef struct kclvm_service kclvm_service;
*/
import "C"
import (
	"kcl-lang.io/lib/go/native"
)

type NativeServiceClient = native.NativeServiceClient

var (
	NewNativeServiceClient = native.NewNativeServiceClient
)
