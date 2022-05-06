# Copyright 2021 The KCL Authors. All rights reserved.

K_FILES=$(wildcard testdata/*.k)
K_AST_FILES=$(patsubst %.k,%.k.ast.json,$(K_FILES))

default: $(K_AST_FILES)
	@echo "ok"

hello:
	go run ./testdata/a_gen_ast.go -kcl-file=a.k

clean:
	-rm ./testdata/*.ast.json

%.k.ast.json: %.k
	go run ./testdata/a_gen_ast.go -kcl-file=$< > $@
