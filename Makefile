# Copyright 2023 The KCL Authors. All rights reserved.

PROJECT_NAME = kcl-go

PWD:=$(shell pwd)

BUILD_IMAGE:=kcllang/kcl

# export DOCKER_DEFAULT_PLATFORM=linux/amd64
# or
# --platform linux/amd64

RUN_IN_DOCKER:=docker run -it --rm
RUN_IN_DOCKER+=-v ~/.ssh:/root/.ssh
RUN_IN_DOCKER+=-v ~/.gitconfig:/root/.gitconfig
RUN_IN_DOCKER+=-v ${PWD}:/root/kcl
RUN_IN_DOCKER+=-w /root/kcl ${BUILD_IMAGE}

clean:
	-rm -rf ./_build

test:
	go test ./...

fmt:
	go fmt ./...

gen-doc:
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
	gomarkdoc --output api.md .
	mv api.md docs

# ----------------
# Docker
# ----------------

sh-in-docker:
	${RUN_IN_DOCKER} bash
