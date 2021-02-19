SHELL := /bin/sh

REVISION := $(shell git rev-parse HEAD)
VERSION := $(shell git rev-parse --abbrev-ref HEAD)

export GO111MODULE=on
export GOPRIVATE=github.com/mariusmagureanu

all: build

check: test lint

.PHONY: test lint

lint:
	golint -min_confidence 0.3 -set_exit_status

build:
	@go generate
	@GOOS=linux GOARCH=amd64 go build -o bin/api-gw -ldflags "-s -w -X main.revision=${REVISION} -X main.version=${VERSION}"

osxbuild:
	@go generate
	@GOOS=darwin GOARCH=amd64 go build -o bin/api-gw -ldflags "-s -w -X main.revision=${REVISION} -X main.version=${VERSION}"

install:
	go install

clean:
	@rm -rf bin/
