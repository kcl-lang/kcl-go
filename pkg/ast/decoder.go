// Copyright 2022 The KCL Authors. All rights reserved.

package ast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var DebugMode bool

func DecodeModule(filename string, src interface{}) (module *Module, err error) {
	data, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return new(Module), nil
	}

	m, err := json_decodeMap(data)
	if err != nil {
		return nil, err
	}

	n, err := new(astNodeBuilder).build(m)
	if err != nil {
		return nil, err
	}
	module, ok := n.(*Module)
	if !ok {
		return nil, fmt.Errorf("not module type: %v", n)
	}

	return module, nil
}

type astNodeBuilder struct {
	root map[string]interface{}
}

func (p *astNodeBuilder) build(m map[string]interface{}) (n Node, err error) {
	if !DebugMode {
		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case error:
					err = fmt.Errorf("%w", r)
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}()
	}

	p.root = m
	n = p.buildNode(m)
	return
}

func (p *astNodeBuilder) buildNode(m map[string]interface{}) Node {
	typ := p.getAstType(m)
	if typ == _No_TypeName {
		return nil
	}

	n := MustNewNode(typ)

	for k, v := range m {
		if k == _json_ast_type_key {
			continue
		}

		var st = reflect.ValueOf(n).Elem()
		var field reflect.Value

		for i := 0; i < st.Type().NumField(); i++ {
			if p.getFieldJsonName(st, i) == k {
				field = st.FieldByIndex([]int{i})
			}
		}
		if !field.IsValid() {
			stMeta := st.FieldByName("Meta").Elem()
			for i := 0; i < stMeta.Type().NumField(); i++ {
				if p.getFieldJsonName(stMeta, i) == k {

					field = stMeta.FieldByIndex([]int{i})
				}
			}
		}
		if !field.IsValid() {
			continue
		}

		p.setStructField(field, v)
	}

	return n
}

func (p *astNodeBuilder) setStructField(field reflect.Value, v interface{}) {
	if v == nil || !field.IsValid() {
		return
	}
	switch v := v.(type) {
	case map[string]interface{}:
		if typ := p.getAstType(v); typ != _No_TypeName {
			field.Set(reflect.ValueOf(p.buildNode(v)))
		} else {
			data, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(data, field.Addr().Interface()); err != nil {
				panic(err)
			}
		}
	case []interface{}:
		elems := reflect.MakeSlice(field.Type(), len(v), len(v))
		for i, x := range v {
			p.setStructField(elems.Index(i), x)
		}
		field.Set(elems)
	case bool:
		switch field.Kind() {
		case reflect.Bool:
			field.SetBool(v)
		case reflect.String:
			field.SetString(fmt.Sprint(v))
		case reflect.Interface:
			field.Set(reflect.ValueOf(v))
		default:
			panic(fmt.Sprintf("unreachable: %T, %v;%v:AAA\n", v, v, field.Kind()))
		}
	case float64:
		switch field.Kind() {
		case reflect.Int:
			field.SetInt(int64(v))
		case reflect.Float64:
			field.SetFloat(v)
		case reflect.Interface:
			field.Set(reflect.ValueOf(v))
		default:
			panic(fmt.Sprintf("unreachable: %v, %v", v, field.Type()))
		}
	case string:
		switch field.Kind() {
		case reflect.String:
			field.SetString(v)
		case reflect.Interface:
			field.Set(reflect.ValueOf(v))
		default:
			panic(fmt.Sprintf("unreachable: %T", v))
		}
	default:
		panic(fmt.Sprintf("unreachable: %T", v))
	}
}

func (p *astNodeBuilder) getAstType(m map[string]interface{}) AstType {
	if len(m) == 0 {
		return _No_TypeName
	}
	if x, ok := m[_json_ast_type_key]; ok {
		if v, ok := x.(string); ok {
			return AstType(v)
		}
	}
	return _No_TypeName
}

func (p *astNodeBuilder) getFieldJsonName(st reflect.Value, i int) string {
	f := st.Type().Field(i)
	tag := f.Tag.Get("json")
	tag = strings.TrimSuffix(tag, ",omitempty")
	if tag == "-" {
		return ""
	}
	return tag
}
