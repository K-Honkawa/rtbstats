NAME1=rtbstats
NAME2=statsMonitor

REVISION:=$(shell git rev-parse --short HEAD)
LDFLAGS := -X main.revision=${REVISION}

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=goimports

.PHONY: help fmt build-linux build clean

## Show help
help:
	@make2help $(MAKEFILE_LIST)

## Format source codes
fmt:
	find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w

test:
	$(GOTEST) -v $(shell go list ./... | grep -v /vendor/)

BINDIR=./bin
NATIVEDIR=$(BINDIR)/native
## build binary
build:
	$(GOBUILD) -o $(NATIVEDIR)/$(NAME1) cmd/$(NAME1)/main.go
	$(GOBUILD) -o $(NATIVEDIR)/$(NAME2) cmd/$(NAME2)/main.go

## remove binary
clean:
	rm -rf $(BINDIR)
