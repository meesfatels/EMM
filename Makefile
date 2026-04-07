.PHONY: build install test clean srcinfo

VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -s -w -X main.version=$(VERSION)

build:
	go build -trimpath -ldflags="$(LDFLAGS)" -o emm ./cmd/emm

install: build
	install -Dm755 emm $(DESTDIR)/usr/bin/emm

test:
	go test ./...

clean:
	rm -f emm

srcinfo:
	cd packaging/aur && makepkg --printsrcinfo > .SRCINFO
