# Copyright 2022 The KCL Authors. All rights reserved.

KCLVM_URL_MIRRORS:=

default:
	which kclvm
	kclvm -m kclvm --version

	go run ./cmds/kcl-go
	go run ./cmds/kcl-go run hello.k

setup-kclvm:
	go run ./cmds/kcl-go/ setup-kclvm  -mirrors=${KCLVM_URL_MIRRORS}

setup-kclvm-all:
	go run ./cmds/kcl-go/ setup-kclvm -all -mirrors=${KCLVM_URL_MIRRORS}

clean:
	-rm -rf ./_build
