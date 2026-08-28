// Package types is one of two sibling Go packages sharing the same
// final path component (`types`) referenced by `cross_pkg_collision`'s
// `Main` struct. The generator must give it a distinct alias so the
// generated KCL doesn't silently shadow one of the two references.
package types

// Tagged is a marker struct from package a/types.
type Tagged struct {
	Value string `json:"value"`
}
