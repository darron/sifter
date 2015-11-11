GIT_COMMIT=$(shell git rev-parse HEAD)
SIFTER_VERSION=$(shell ./version)
COMPILE_DATE=$(shell date -u +%Y%m%d.%H%M%S)
BUILD_FLAGS=-X main.CompileDate=$(COMPILE_DATE) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(SIFTER_VERSION)

all: build

deps:
	go get github.com/spf13/cobra

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

release: clean build gziposx clean linux gziplinux clean
