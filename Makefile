.PHONY: build install clean test run help

# Build the binary
build:
	go build -o bookmarked ./cmd/bookmarked

# Build for all platforms
build-all:
	GOOS=windows GOARCH=amd64 go build -o dist/bookmarked-windows-amd64.exe ./cmd/bookmarked
	GOOS=darwin GOARCH=amd64 go build -o dist/bookmarked-darwin-amd64 ./cmd/bookmarked
	GOOS=darwin GOARCH=arm64 go build -o dist/bookmarked-darwin-arm64 ./cmd/bookmarked
	GOOS=linux GOARCH=amd64 go build -o dist/bookmarked-linux-amd64 ./cmd/bookmarked
	GOOS=linux GOARCH=arm64 go build -o dist/bookmarked-linux-arm64 ./cmd/bookmarked

# Install locally
install: build
	mv bookmarked $(GOPATH)/bin/

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f bookmarked bookmarked.exe
	rm -rf dist/

# Run the application
run: build
	./bookmarked

# Download dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  build-all  - Build for all platforms"
	@echo "  install    - Build and install to GOPATH/bin"
	@echo "  test       - Run tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  run        - Build and run"
	@echo "  deps       - Download dependencies"
	@echo "  help       - Show this help message"
