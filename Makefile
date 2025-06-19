# Project Variables
PROJECT_NAME := ttymer
VERSION := 1.0.2
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Tools
GO ?= go

# Flags
GFLAGS := -v
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Directories
BIN_DIR := bin

.PHONY: all build clean test fmt deps install build-linux build-darwin build-windows

all: build

build: fmt
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GFLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME) ./ttymer.go

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build $(GFLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-linux-amd64 ./ttymer.go

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GFLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 ./ttymer.go

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GFLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-windows-amd64.exe ./ttymer.go

clean:
	rm -rf $(BIN_DIR) $(PROJECT_NAME)
	$(GO) clean

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

deps:
	$(GO) mod download

install: build
	install -Dm755 $(BIN_DIR)/$(PROJECT_NAME) /usr/bin/$(PROJECT_NAME)
