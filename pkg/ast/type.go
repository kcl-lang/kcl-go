// Copyright 2022 The KCL Authors. All rights reserved.

package ast

import (
	"fmt"
	"sort"
)

type AstType string

const _No_TypeName AstType = ""

const _json_ast_type_key = "_ast_type"

var _ast_node_factory_map = map[AstType]func() Node{}

func GetTypeNameList() []AstType {
	l := make([]AstType, 0, len(_ast_node_factory_map))
	for k := range _ast_node_factory_map {
		l = append(l, k)
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	return l
}

func NewNode(typ AstType) (n Node, ok bool) {
	if fn, ok := _ast_node_factory_map[typ]; ok {
		return fn(), true
	}
	return nil, false
}

func MustNewNode(typ AstType) Node {
	if fn, ok := _ast_node_factory_map[typ]; ok {
		return fn()
	}
	panic(fmt.Sprintf("NewNode: unknown '%s' type", typ))
}
