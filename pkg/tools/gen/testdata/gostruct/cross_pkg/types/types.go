// Package types holds types that the main package references via a Go
// import. The generated KCL should pick them up and emit a matching
// `import` statement instead of inlining the structs.
package types

// Inner is a simple struct from an external Go package.
type Inner struct {
	Value string `json:"value"`
}
