## What is `protoc-gen-ent`

The `protoc-gen-ent` can generate [ent](https://github.com/ent/ent)'s schemas.

## Why `protoc-gen-ent`

A project may have services that require or do not require ORM. 

For services that do not require ORM, can define them in [protobuf](https://github.com/protocolbuffers/protobuf).

This will cause handwritten [protobuf](https://github.com/protocolbuffers/protobuf) `.proto` file is inconsistent with the handwritten [ent](https://github.com/ent/ent) schema `.go` file.

If can generate [ent](https://github.com/ent/ent)'s schema `.go` from [protobuf](https://github.com/protocolbuffers/protobuf) `.proto`, In this way, all data structs and services will be as consistent as possible.

The `protoc-gen-ent` is such a tool.

## How to use

Can see `Makefile`.

Can directly:
``` 
	go install github.com/go-woo/protoc-gen-ent@latest
	go get entgo.io/ent/cmd/ent
	go run entgo.io/ent/cmd/ent init Todo
	rm ./ent/schema/todo.go

	protoc --proto_path=. \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:. \
	       ./gent/gent.proto
```
## Todo

* Enum more detailed handling.

* Support Value, JSON, Custom ID, Hook.