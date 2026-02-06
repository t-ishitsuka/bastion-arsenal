.PHONY: build install test clean lint dev build-all release

BINARY := bastion-arsenal
BUILD_DIR := ./cmd/arsenal

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(BUILD_DIR)

install:
	go install -ldflags "$(LDFLAGS)" $(BUILD_DIR)

test:
	go test ./... -v

clean:
	rm -f $(BINARY)
	rm -rf dist/

lint:
	golangci-lint run ./...

# Dev: build and install to ~/.local/bin
dev: build
	cp $(BINARY) ~/.local/bin/$(BINARY)

# Build for all platforms
build-all:
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64 $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64 $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64 $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64 $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-arm64.exe $(BUILD_DIR)

# Create release archives
release: build-all
	cd dist && tar czf $(BINARY)-$(VERSION)-linux-amd64.tar.gz $(BINARY)-linux-amd64
	cd dist && tar czf $(BINARY)-$(VERSION)-linux-arm64.tar.gz $(BINARY)-linux-arm64
	cd dist && tar czf $(BINARY)-$(VERSION)-darwin-amd64.tar.gz $(BINARY)-darwin-amd64
	cd dist && tar czf $(BINARY)-$(VERSION)-darwin-arm64.tar.gz $(BINARY)-darwin-arm64
	cd dist && zip $(BINARY)-$(VERSION)-windows-amd64.zip $(BINARY)-windows-amd64.exe
	cd dist && zip $(BINARY)-$(VERSION)-windows-arm64.zip $(BINARY)-windows-arm64.exe
	cd dist && shasum -a 256 *.tar.gz *.zip > SHA256SUMS
