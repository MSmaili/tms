.PHONY: build install test clean release

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X 'github.com/MSmaili/tms/cmd.Version=$(VERSION)' \
           -X 'github.com/MSmaili/tms/cmd.GitCommit=$(COMMIT)' \
           -X 'github.com/MSmaili/tms/cmd.BuildDate=$(DATE)'

build:
	go build -ldflags "$(LDFLAGS)" -o tms .

install: build
	sudo mv tms /usr/local/bin/tms

test:
	go test -v ./...

clean:
	rm -f tms
	rm -rf dist/

# Build for multiple platforms (for releases)
release:
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/tms-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/tms-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/tms-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/tms-linux-arm64 .
	@echo "Release binaries built in dist/"
