# Copyright 2021 The KCL Authors. All rights reserved.

K_FILES=$(wildcard ./*.k)
K_AST_FILES=$(patsubst %.k,%.k.ast.json,$(K_FILES))

default: $(K_AST_FILES)
	@echo "ok"

hello:
	go run ./a_gen_ast.go -kcl-file=a.k

clean:
	-rm *.ast.json

%.k.ast.json: %.k
	go run ./a_gen_ast.go -kcl-file=$< > $@
