.PHONY: protos
protos:
	protoc --go_out=./api --go_opt=paths=source_relative \
	--go-grpc_out=./api --go-grpc_opt=paths=source_relative \
	--proto_path=api  logs.proto

