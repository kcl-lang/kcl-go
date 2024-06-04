package gen

import (
	"errors"
	"reflect"
	"strings"
)

func KclTypeOfGo(rv reflect.Value) (string, error) {
	if IsNil(rv) || !rv.IsValid() {
		return typAny, nil
	}

	if isMarshalTy(rv) {
		return typStr, nil
	}

	switch rv.Kind() {
	case reflect.Bool:
		return typBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return typInt, nil
	case reflect.Float32, reflect.Float64:
		return typFloat, nil
	case reflect.Array, reflect.Slice:
		return typList, nil
	case reflect.Ptr, reflect.Interface:
		return KclTypeOfGo(rv.Elem())
	case reflect.String:
		return typStr, nil
	case reflect.Map:
		return typDict, nil
	default:
		return "", errors.New("unsupported kcl type of go: " + rv.Kind().String())
	}
}

type TagOptions struct {
	skip      bool // "-"
	name      string
	omitempty bool
	omitzero  bool
}

func GetTagOptions(tag reflect.StructTag) TagOptions {
	t := tag.Get("kcl")
	if t == "-" {
		return TagOptions{skip: true}
	}
	var opts TagOptions
	parts := strings.Split(t, ",")
	opts.name = parts[0]
	for _, s := range parts[1:] {
		switch s {
		case "omitempty":
			opts.omitempty = true
		case "omitzero":
			opts.omitzero = true
		}
	}
	return opts
}

func IsZero(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0.0
	}
	return false
}

func IsEmpty(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return rv.Len() == 0
	case reflect.Struct:
		if rv.Type().Comparable() {
			return reflect.Zero(rv.Type()).Interface() == rv.Interface()
		}
		for i := 0; i < rv.NumField(); i++ {
			if !IsEmpty(rv.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Ptr:
		return rv.IsNil()
	}
	return false
}

func IsNil(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

func PtrTo(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return PtrTo(t.Elem())
	}
	return t
}

type Field struct {
	name  string
	ty    reflect.Type
	val   reflect.Value
	index []int
	opts  *TagOptions
}

// TypeFields returns a list of fields.
func TypeFields(t reflect.Type, v reflect.Value) []Field {
	// Anonymous fields to explore at the current level and the next.
	current := []Field{}
	next := []Field{{ty: t, val: v}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []Field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.ty] {
				continue
			}
			visited[f.ty] = true

			// Scan f.ty for fields to include.
			for i := 0; i < f.ty.NumField(); i++ {
				sf := f.ty.Field(i)
				val := f.val.Field(i)
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}

				opts := GetTagOptions(sf.Tag)
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					// Follow pointer.
					ft = ft.Elem()
				}
				// Record found field and index sequence.
				if !sf.Anonymous || ft.Kind() != reflect.Struct {
					field := Field{
						name:  sf.Name,
						index: index,
						ty:    ft,
						val:   val,
						opts:  &opts,
					}
					fields = append(fields, field)
					if count[f.ty] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 and 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, Field{name: ft.Name(), index: index, ty: ft, val: val})
				}
			}
		}
	}

	return fields
}
