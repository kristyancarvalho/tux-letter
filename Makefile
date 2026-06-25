VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

VERSION_PKG := github.com/kristyancarvalho/tux-letter/internal/version
LDFLAGS := -s -w \
	-X $(VERSION_PKG).Version=$(VERSION) \
	-X $(VERSION_PKG).Commit=$(COMMIT) \
	-X $(VERSION_PKG).Date=$(DATE)

BIN := bin/tux-letter
PKG := ./cmd/tux-letter
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all fmt fmt-check test lint build run coverage release-check \
	release-source-archive aur-srcinfo aur-verifysource aur-build clean tidy

all: build

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

fmt-check:
	@out="$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))"; \
	if [ -n "$$out" ]; then \
		echo "gofmt needs to be run on:"; echo "$$out"; exit 1; \
	fi

test:
	go test ./...

lint:
	go vet ./...

build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) $(PKG)

run:
	go run $(PKG) $(ARGS)

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func=coverage.txt | tail -1

release-check: fmt-check lint test build
	@echo "release-check passed"

release-source-archive:
	@scripts/release/source-archive.sh $(VERSION)

aur-srcinfo:
	cd packaging/aur && makepkg --printsrcinfo > .SRCINFO

aur-verifysource:
	cd packaging/aur && makepkg --verifysource

aur-build:
	cd packaging/aur && makepkg -f

tidy:
	go mod tidy

clean:
	rm -rf bin dist coverage.txt
