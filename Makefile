.PHONY: install example

install:
	go install ./protoc-gen-jrpc

example:
	protoc --go_out=plugins=grpc,paths=source_relative:. --jrpc_out=paths=source_relative:. ./example/pb/*.proto
