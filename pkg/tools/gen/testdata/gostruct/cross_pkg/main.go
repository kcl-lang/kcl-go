// Package cross_pkg exercises the cross-package import emission path of
// genkcl_gostruct. Outer references Inner from the sibling `types`
// package via a Go import, with no `kcl:` tag override, so the Go type
// analyser is responsible for producing the KCL type reference.
package cross_pkg

import (
	"kcl-lang.io/kcl-go/pkg/tools/gen/testdata/gostruct/cross_pkg/types"
)

// Outer references types.Inner to verify that genkcl_gostruct emits an
// `import` statement and uses the alias-qualified type name
// (e.g. `Types.Inner`) when OneFile is false.
type Outer struct {
	Inner types.Inner `json:"inner"`
	Name  string      `json:"name"`
}

// Plain is a schema with no cross-package references. It exists so the
// generated file contains more than just the imported reference and so
// that the test can assert imports are scoped (i.e. only emitted for
// packages actually referenced).
type Plain struct {
	Age int `json:"age"`
}
