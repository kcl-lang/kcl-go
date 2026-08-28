// Package app is the entry point of the gostruct_module Go fixture.
// It references the sibling `types` package (so the generator must
// emit a cross-package `import`) and `time.Time` (which the
// generator short-circuits to `str`).
//
// The `encoding/json.RawMessage` reference is deliberately
// out-of-module relative to the kcl-go module; under GenKclModule's
// module-root handling it must be degraded to `any` and no `.k`
// file or `import` is emitted for it.
package app

import (
	"encoding/json"
	"time"

	"kcl-lang.io/kcl-go/pkg/tools/gen/testdata/gostruct_module/app/types"
)

// App references the sibling `types.Inner`, the stdlib
// `time.Time`, and an out-of-module stdlib `json.RawMessage` so
// the GenKclModule assertions can exercise every code path:
// cross-package import, stdlib short-circuit, and out-of-module
// → any.
type App struct {
	Inner *types.Inner    `json:"inner"`
	When  time.Time       `json:"when"`
	Meta  json.RawMessage `json:"meta"`
}
