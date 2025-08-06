.PHONY: start start-with-password build clean help install

# Default target
start:
	@./start.sh

# Start with password
start-with-password:
	@./start.sh --password demo123

# Start on custom port
start-port:
	@./start.sh --port 9000

# Build only (no start)
build:
	@echo "🔨 Building GoShare..."
	@go build -o goshare .
	@cd frontend && npm install
	@echo "✅ Build complete!"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -f goshare goshare.exe
	@rm -rf frontend/node_modules
	@rm -rf frontend/build
	@echo "✅ Cleanup complete!"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@cd frontend && npm install
	@go mod download
	@echo "✅ Dependencies installed!"

# Show help
help:
	@echo "GoShare - Easy Commands"
	@echo "======================"
	@echo ""
	@echo "make start              # Start GoShare (no password)"
	@echo "make start-with-password # Start with demo password"
	@echo "make start-port         # Start on port 9000"
	@echo "make build              # Build without starting"
	@echo "make install            # Install dependencies"
	@echo "make clean              # Clean build files"
	@echo "make help               # Show this help"
	@echo ""
	@echo "Or use the script directly:"
	@echo "./start.sh --password yourpass --port 8080 ~/Documents"
