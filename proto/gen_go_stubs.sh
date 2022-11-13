#!/bin/bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
PATH=$(go env GOPATH)/bin:$PATH
protoc -I . --go-grpc_out=.. --go_out=.. *.proto
