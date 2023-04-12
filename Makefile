# Copyright 2023 The KCL Authors. All rights reserved.

KCLVM_URL_MIRRORS:=

default:
	go run ./cmds/kcl-go run hello.k

clean:
	-rm -rf ./_build
