# Copyright 2023 The KCL Authors. All rights reserved.

PROJECT_NAME = kcl-go

PWD:=$(shell pwd)

BUILD_IMAGE:=kusionstack/kclvm-builder

# export DOCKER_DEFAULT_PLATFORM=linux/amd64
# or
# --platform linux/amd64

RUN_IN_DOCKER:=docker run -it --rm
RUN_IN_DOCKER+=-v ~/.ssh:/root/.ssh
RUN_IN_DOCKER+=-v ~/.gitconfig:/root/.gitconfig
RUN_IN_DOCKER+=-v ${PWD}:/root/kcl
RUN_IN_DOCKER+=-w /root/kcl ${BUILD_IMAGE}

default:
	go run ./cmds/kcl-go run hello.k

clean:
	-rm -rf ./_build

test:
	go test ./...

# ----------------
# Docker
# ----------------

sh-in-docker:
	${RUN_IN_DOCKER} bash
