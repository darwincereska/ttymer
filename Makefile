BINARY_NAME=ttymer

GFLAGS=-v

.PHONY: all
all: build


.PHONY: build
build:
	go build $(GFLAGS) -o $(BINARY_NAME) ttymer.go

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: test
test:
	go test -v ./ ...

.PHONY: nix-build
nix-build:
	nix-build -E 'with import <nixpkgs> {}; callPackage ./default.nix {}'
	

.PHONY: deps
deps:
	go mod download

.PHONY: fmt
fmt:
	go fmt ./ ...

.PHONY: install
install: build
	go install

