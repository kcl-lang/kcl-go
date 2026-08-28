// Package cross_pkg_collision exercises the cross-package import alias
// collision path of genkcl_gostruct. Main references structs from two
// sibling Go packages both named `types`; without alias disambiguation
// the generator would emit two `import … as Types` lines and silently
// shadow one of them.
package cross_pkg_collision

import (
	"kcl-lang.io/kcl-go/pkg/tools/gen/testdata/gostruct/cross_pkg_collision/a/types"
	btypes "kcl-lang.io/kcl-go/pkg/tools/gen/testdata/gostruct/cross_pkg_collision/b/types"
)

// Main holds references to both `types` packages. The generator is
// expected to emit two distinct `import … as <alias>` lines (the first
// takes the preferred `Types` alias, the second is renamed `Types2`)
// and reference each via its alias-qualified KCL type name.
type Main struct {
	A *types.Tagged   `json:"a"`
	B *btypes.Labeled `json:"b"`
}
