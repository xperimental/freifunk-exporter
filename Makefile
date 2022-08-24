SHELL = bash
DOCKER ?= docker
GO ?= go
GO_CMD := CGO_ENABLED=0 $(GO)
GIT_VERSION := $(shell git describe --tags --dirty || git rev-parse --short HEAD)
VERSION := $(GIT_VERSION:v%=%)
GIT_COMMIT := $(shell git rev-parse HEAD)

.PHONY: all
all: test build-binary

.PHONY: test
test:
	$(GO_CMD) test -cover ./...

.PHONY: build-binary
build-binary:
	$(GO_CMD) build -tags netgo -ldflags "-w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -o freifunk-exporter .

.PHONY: build-image
build-image:
	$(DOCKER) build -t ghcr.io/xperimental/freifunk-exporter:$(VERSION) .

.PHONY: clean
clean:
	rm -f freifunk-exporter
