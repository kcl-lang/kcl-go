// Package consts exercises the Go `const` → KCL global variable emission
// path of genkcl_gostruct (kcl-lang/kcl-go#332 sub-task 3).
//
// The generator is expected to produce a top-level KCL global variable
// for every const in the main package whose value can be statically
// evaluated, regardless of whether the const is typed or untyped.
package consts

// Pi is an untyped floating-point const. The float64 approximation is
// what gets emitted — Go stores it as the exact rational 157/50 but
// users expect to see `3.14`.
const Pi = 3.14

// Name is an untyped string const.
const Name = "kcl"

// MaxRetries is a typed int const.
const MaxRetries int = 5

// Enabled is a typed bool const.
const Enabled bool = true

// HexPort demonstrates a typed int with a non-decimal literal that has
// to be evaluated by the Go type checker before the generator sees it.
const HexPort int = 0x1F90 // 8080

// TwoNames demonstrates a single ValueSpec with multiple names — both
// are emitted, each bound to its respective constant expression.
const TwoA, TwoB = 10, 20

// Container is here to ensure the generator still produces schemas
// alongside globals, and that the global block is emitted in the right
// place relative to them.
type Container struct {
	Label string `json:"label"`
}
