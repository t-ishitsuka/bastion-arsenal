.PHONY: build install test clean

BINARY := arsenal
BUILD_DIR := ./cmd/arsenal

build:
	go build -o $(BINARY) $(BUILD_DIR)

install:
	go install $(BUILD_DIR)

test:
	go test ./... -v

clean:
	rm -f $(BINARY)

lint:
	golangci-lint run ./...

# Dev: build and install to ~/.local/bin
dev: build
	cp $(BINARY) ~/.local/bin/$(BINARY)
