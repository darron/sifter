SIFTER_VERSION="0.8-dev"
GIT_COMMIT=$(shell git rev-parse HEAD)
COMPILE_DATE=$(shell date -u +%Y%m%d.%H%M%S)
BUILD_FLAGS=-X main.CompileDate=$(COMPILE_DATE) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(SIFTER_VERSION)
CONFIG_DIR=$(shell cd config && pwd)

all: build

deps:
	go get -u github.com/spf13/cobra
	go get -u github.com/hashicorp/consul/api
	go get -u github.com/PagerDuty/godspeed
	go get -u github.com/pmylund/sortutil
	go get -u github.com/darron/sifter

format:
	gofmt -w .

clean:
	rm -f bin/sifter || true

build: clean
	go build -ldflags "$(BUILD_FLAGS)" -o bin/sifter main.go

gziposx:
	gzip bin/sifter
	mv bin/sifter.gz bin/sifter-$(SIFTER_VERSION)-darwin.gz

linux: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "$(BUILD_FLAGS)" -o bin/sifter main.go

gziplinux:
	gzip bin/sifter
	mv bin/sifter.gz bin/sifter-$(SIFTER_VERSION)-linux-amd64.gz

consul:
	consul agent -data-dir `mktemp -d` -config-dir=$(CONFIG_DIR) -bootstrap -server -bind=127.0.0.1

consul_kill:
	pkill consul

release: clean build gziposx clean linux gziplinux clean
