#!/bin/bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
PATH=$(go env GOPATH)/bin:$PATH
protoc -I . --go-grpc_out=.. --go_out=.. *.proto
