# ld-agent Go Makefile

.PHONY: all build-plugin build-example run clean test

# Default target
all: build-plugin build-example

# Build the calculator plugin as a shared library
build-plugin:
	@echo "🔨 Building calculator plugin..."
	@mkdir -p plugins
	go build -buildmode=plugin -o plugins/calculator.so ./plugins/calculator/

# Build the example program
build-example:
	@echo "🔨 Building example..."
	go build -o example/ld-agent-example ./example/

# Run the example (builds everything first)
run: build-plugin build-example
	@echo "🚀 Running ld-agent Go example..."
	@cd example && ./ld-agent-example

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f plugins/*.so
	rm -f example/ld-agent-example

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Initialize Go modules
init:
	@echo "📦 Initializing Go modules..."
	go mod tidy

# Show help
help:
	@echo "ld-agent Go Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all          - Build plugin and example (default)"
	@echo "  build-plugin - Build the calculator plugin"
	@echo "  build-example- Build the example program"
	@echo "  run          - Build and run the example"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  init         - Initialize Go modules"
	@echo "  help         - Show this help" 
