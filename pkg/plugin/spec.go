// Copyright The KCL Authors. All rights reserved.
package plugin

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Plugin represents a KCL Plugin with metadata and methods.
// It contains the plugin name, version, a reset function, and a map of methods.
type Plugin struct {
	Name      string                // Name of the plugin
	Version   string                // Version of the plugin
	ResetFunc func()                // Reset function for the plugin
	MethodMap map[string]MethodSpec // Map of method names to their specifications
}

// MethodSpec defines the specification for a KCL Plugin method.
// It includes the method type and the body function which executes the method logic.
type MethodSpec struct {
	Type *MethodType                                   // Specification of the method's type
	Body func(args *MethodArgs) (*MethodResult, error) // Function to execute the method's logic
}

// MethodType describes the type of a KCL Plugin method's arguments, keyword arguments, and result.
// It specifies the types of positional arguments, keyword arguments, and the result type.
type MethodType struct {
	ArgsType   []string          // List of types for positional arguments
	KwArgsType map[string]string // Map of keyword argument names to their types
	ResultType string            // Type of the result
}

// MethodArgs represents the arguments passed to a KCL Plugin method.
// It includes a list of positional arguments and a map of keyword arguments.
type MethodArgs struct {
	Args   []interface{}          // List of positional arguments
	KwArgs map[string]interface{} // Map of keyword arguments
}

// MethodResult represents the result returned from a KCL Plugin method.
// It holds the value of the result.
type MethodResult struct {
	V interface{} // Result value
}

// ParseMethodArgs parses JSON strings for positional and keyword arguments
// and returns a MethodArgs object.
// args_json: JSON string of positional arguments
// kwargs_json: JSON string of keyword arguments
func ParseMethodArgs(args_json, kwargs_json string) (*MethodArgs, error) {
	p := &MethodArgs{
		KwArgs: make(map[string]interface{}),
	}
	if args_json != "" {
		if err := json.Unmarshal([]byte(args_json), &p.Args); err != nil {
			return nil, err
		}
	}
	if kwargs_json != "" {
		if err := json.Unmarshal([]byte(kwargs_json), &p.KwArgs); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// GetCallArg retrieves an argument by index or key.
// If the key exists in KwArgs, it returns the corresponding value.
// Otherwise, it returns the positional argument at the given index.
func (p *MethodArgs) GetCallArg(index int, key string) any {
	if val, ok := p.KwArgs[key]; ok {
		return val
	}
	if index < len(p.Args) {
		return p.Args[index]
	}
	return nil
}

// Arg returns the positional argument at the specified index.
func (p *MethodArgs) Arg(i int) interface{} {
	return p.Args[i]
}

// KwArg returns the keyword argument with the given name.
func (p *MethodArgs) KwArg(name string) interface{} {
	return p.KwArgs[name]
}

// IntArg returns the positional argument at the specified index
// as an int64. It panics if the conversion fails.
func (p *MethodArgs) IntArg(i int) int64 {
	s := fmt.Sprint(p.Args[i])
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// FloatArg returns the positional argument at the specified index
// as a float64. It panics if the conversion fails.
func (p *MethodArgs) FloatArg(i int) float64 {
	s := fmt.Sprint(p.Args[i])
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// BoolArg returns the positional argument at the specified index
// as a bool. It panics if the conversion fails.
func (p *MethodArgs) BoolArg(i int) bool {
	s := fmt.Sprint(p.Args[i])
	v, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return v
}

// StrArg returns the positional argument at the specified index
// as a string.
func (p *MethodArgs) StrArg(i int) string {
	s := fmt.Sprint(p.Args[i])
	return s
}

// ListArg returns the positional argument at the specified index
// as a list of any type.
func (p *MethodArgs) ListArg(i int) []any {
	return p.Args[i].([]any)
}

// MapArg returns the positional argument at the specified index
// as a map with string keys and any type values.
func (p *MethodArgs) MapArg(i int) map[string]any {
	return p.Args[i].(map[string]any)
}

// IntKwArg returns the keyword argument with the given name
// as an int64. It panics if the conversion fails.
func (p *MethodArgs) IntKwArg(name string) int64 {
	s := fmt.Sprint(p.KwArgs[name])
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// FloatKwArg returns the keyword argument with the given name
// as a float64. It panics if the conversion fails.
func (p *MethodArgs) FloatKwArg(name string) float64 {
	s := fmt.Sprint(p.KwArgs[name])
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// BoolKwArg returns the keyword argument with the given name
// as a bool. It panics if the conversion fails.
func (p *MethodArgs) BoolKwArg(name string) bool {
	s := fmt.Sprint(p.KwArgs[name])
	v, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return v
}

// StrKwArg returns the keyword argument with the given name
// as a string.
func (p *MethodArgs) StrKwArg(name string) string {
	s := fmt.Sprint(p.KwArgs[name])
	return s
}

// ListKwArg returns the keyword argument with the given name
// as a list of any type.
func (p *MethodArgs) ListKwArg(name string) []any {
	return p.KwArgs[name].([]any)
}

// MapKwArg returns the keyword argument with the given name
// as a map with string keys and any type values.
func (p *MethodArgs) MapKwArg(name string) map[string]any {
	return p.KwArgs[name].(map[string]any)
}
