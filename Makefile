BINARY_NAME=ttymer
VERSION=1.0.0
PACKAGE_NAME=ttymer_$(VERSION)_amd64
GFLAGS=-v

.PHONY: all
all: build


.PHONY: build
build:
	go build $(GFLAGS) -o $(BINARY_NAME) ttymer.go

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(PACKAGE_NAME)

.PHONY: test
test:
	go test -v ./ ...

.PHONY: nix-build
nix-build:
	nix-build -E 'with import <nixpkgs> {}; callPackage ./default.nix {}'
	
.PHONY: deb
deb: build
	mkdir -p $(PACKAGE_NAME)/DEBIAN
	mkdir -p $(PACKAGE_NAME)/usr/local/bin
	cp $(BINARY_NAME) $(PACKAGE_NAME)/usr/local/bin/
	echo "Package: ttymer" > $(PACKAGE_NAME)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(PACKAGE_NAME)/DEBIAN/control
	echo "Section: base" >> $(PACKAGE_NAME)/DEBIAN/control
	echo "Priority: optional" >> $(PACKAGE_NAME)/DEBIAN/control
	echo "Architecture: amd64" >> $(PACKAGE_NAME)/DEBIAN/control
	echo "Maintainer: Your Name <your.email@example.com>" >> $(PACKAGE_NAME)/DEBIAN/control
	echo "Description: A terminal timer application" >> $(PACKAGE_NAME)/DEBIAN/control
	dpkg-deb --build $(PACKAGE_NAME)
	rm -rf $(PACKAGE_NAME)

.PHONY: deps
deps:
	go mod download

.PHONY: fmt
fmt:
	go fmt ./ ...

.PHONY: install
install: build
	go install

.PHONY: arch
arch:
	BUILDDIR=/tmp/ttymer makepkg -si