# Copyright 2023 The KCL Authors. All rights reserved.

default:
	go install github.com/golang/protobuf/protoc-gen-go
	protoc --go_out=plugins=grpc:. gpyrpc.proto

	cp ./kusionstack.io/kclvm-go/pkg/spec/gpyrpc/gpyrpc.pb.go gpyrpc.pb.go
	-rm -rf ./kusionstack.io

clean: