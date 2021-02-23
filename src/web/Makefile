SHELL=bash
include variables.env
export

REVISION := $(shell git rev-parse HEAD)
VERSION := $(shell git rev-parse --abbrev-ref HEAD)

export GO111MODULE=on
export GOPRIVATE=github.com/mariusmagureanu

all: build prepare docker clean
start: build prepare docker

prepare:
	@mkdir -p ./burlan_tmp/vsm
	@cp nats.conf burlan_tmp
	@. variables.env

build:
	cd apigw && GOOS=linux GOARCH=amd64 go build -o bin/api-gw -ldflags "-s -w -X main.revision=${REVISION} -X main.version=${VERSION}" && cd ../ 

docker: prepare
	docker-compose pull
	docker-compose build
	docker-compose up -d

clean:
	docker-compose down
	rm -rf ./burlan_tmp
