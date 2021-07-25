#!/bin/bash

.PHONY:
gen: clean
	@ echo generating protobuf
	@ buf generate --path proto/book

.PHONY:
fmt:
	@ echo formatting go code
	@ go fmt ./...

.PHONY:
install:
	@ installing dependencies
	@ go get \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/bufbuild/buf/cmd/buf

.PHONY:
gen-certs:
	@ echo generating tls certs
	@ cd certs; chmod +x gen.sh; ./gen.sh; cd ..

.PHONY:
run:
	@ echo running gRPC server
	@ go run server/main.go -port 8080

.PHONY:
unit-tests:
	@ echo running unit-tests
	@ go test -cover -race ./...

.PHONY:
clean:
	@ echo cleaning-up
	@ rm -rf gen/go/*/*.go
