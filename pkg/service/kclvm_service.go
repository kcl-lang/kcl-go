package service

import "kcl-lang.io/lib/go/api"

// KCL service client interface.
// Deprecated: use `ServiceClient` at `kcl-lang.io/lib/go/api`
type KclvmService interface {
	api.ServiceClient
}
