.PHONY: build install clean test run daemon

BINARY_NAME=cliptui
INSTALL_PATH=/usr/local/bin

build:
	@echo "Building clipTUI..."
	go build -o $(BINARY_NAME) ./cmd/cliptui

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "clipTUI installed successfully!"

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean

test:
	@echo "Running tests..."
	go test -v ./...

run: build
	./$(BINARY_NAME)

daemon: build
	./$(BINARY_NAME) daemon

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

release:
	@echo "Building release with goreleaser..."
	goreleaser release --snapshot --clean

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  install  - Build and install to $(INSTALL_PATH)"
	@echo "  clean    - Remove built files"
	@echo "  test     - Run tests"
	@echo "  run      - Build and run the TUI"
	@echo "  daemon   - Build and run the daemon"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  release  - Build release packages with goreleaser"
