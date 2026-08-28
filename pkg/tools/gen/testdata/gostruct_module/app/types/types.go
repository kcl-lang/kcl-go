// Package types holds the Inner struct referenced from the parent
// `app` package in the gostruct_module fixture. The generator is
// expected to emit it as a separate directory under the produced KCL
// module so the referencing schema can `import` it.
package types

// Inner is a simple struct from the sibling `types` package. The
// generated KCL module materialises a sibling `app/types/types.k`
// file that holds the Inner schema, so the parent `App` schema can
// reference it via an `import` statement.
type Inner struct {
	Value string `json:"value"`
}
