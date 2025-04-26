# Makefile for godocs project

BINARY_NAME=godocs
INSTALL_PATH=/usr/local/bin

## Build the project binary
build:
	@echo "🚀 Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

## Install the binary to /usr/local/bin
install: build
	@echo "📦 Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)
	@echo "✅ Installed successfully! Now you can run '$(BINARY_NAME) --help'"

## Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME)
