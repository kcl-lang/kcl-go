# Copyright 2023 The KCL Authors. All rights reserved.

clean:
	-rm -rf ./_build

test:
	go test ./...

fmt:
	go fmt ./...

check:
	make build

build:
	go build ./...

gen-doc:
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
	gomarkdoc --output api.md .
	mv api.md docs
