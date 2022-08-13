GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find ./example/ent/schema -name *.proto")
	ENT_PROTO_FILES=$(shell $(Git_Bash) -c "find ./example/ent/schema -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find ./example/ent/schema -name *.proto)
	ENT_PROTO_FILES=$(shell find ./example/ent/schema -name *.proto)
endif

.PHONY: init
# init tools-chain
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/go-woo/protoc-gen-ent@latest
	go get entgo.io/ent/cmd/ent
	go run entgo.io/ent/cmd/ent init Todo
	rm ./ent/schema/todo.go

.PHONY: ent
# generate ent proto
ent:
	protoc --proto_path=. \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:. \
	       --ent_out=paths=source_relative:. \
	       $(ENT_PROTO_FILES)

.PHONY: rule
# generate rule proto
rule:
	protoc --proto_path=. \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:. \
	       ./gent/gent.proto

.PHONY: example
# example
example:
	make init;
	make ent;

.PHONY: all
# generate all
all:
	go install .
	make ent;

.PHONY: clear
# clear
clear:
	rm ./example/ent/schema/*.pb.go
