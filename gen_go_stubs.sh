#!/bin/bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
PATH=$(go env GOPATH)/bin:$PATH
PACKAGE_PATH="echo-grpc-triton"
protoc -I proto --go-grpc_out=${PACKAGE_PATH} --go_out=${PACKAGE_PATH} proto/*.proto
