# Makefile for pure-dupes Phase 1
.PHONY: help build wasm runtime mcp test-files serve clean check all install

# Colors
GREEN  := \033[0;32m
BLUE   := \033[0;34m
YELLOW := \033[1;33m
RED    := \033[0;31m
NC     := \033[0m

# Variables
WASM_FILE := main.wasm
WASM_SRC := main_wasm_enhanced.go
WASM_EXEC := wasm_exec.js
MCP_SERVER := mcp-server
MCP_SRC := mcp-server.go
INDEX := index.html
PORT := 8080

# Default target
help:
	@echo "$(BLUE)🔍 Pure Dupes - Phase 1 Build System$(NC)"
	@echo ""
	@echo "$(GREEN)Quick Start:$(NC)"
	@echo "  make all      - Build everything"
	@echo "  make serve    - Start HTTP server"
	@echo ""
	@echo "$(GREEN)Build Targets:$(NC)"
	@echo "  make build    - Build all components"
	@echo "  make wasm     - Build WASM module only"
	@echo "  make runtime  - Get wasm_exec.js"
	@echo "  make mcp      - Build MCP server"
	@echo ""
	@echo "$(GREEN)Utility Targets:$(NC)"
	@echo "  make check    - Verify all files exist"
	@echo "  make test-files - Create test data"
	@echo "  make clean    - Remove built files"
	@echo "  make install  - Install dependencies"
	@echo ""
	@echo "$(GREEN)Run Targets:$(NC)"
	@echo "  make serve    - Start server on port $(PORT)"
	@echo "  make test     - Run automated tests"
	@echo ""

# Build everything
all: build
	@echo ""
	@echo "$(GREEN)✅ Build complete!$(NC)"
	@echo ""
	@echo "$(BLUE)To run:$(NC)"
	@echo "  make serve"
	@echo ""

# Build all components
build: check-go wasm runtime mcp test-files $(INDEX)
	@echo "$(GREEN)✅ All components built$(NC)"

# Build WASM module
wasm: $(WASM_FILE)

$(WASM_FILE): $(WASM_SRC)
	@echo "$(BLUE)📦 Building WASM module...$(NC)"
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $(WASM_FILE) $(WASM_SRC)
	@echo "$(GREEN)✅ WASM built: $$(du -h $(WASM_FILE) | cut -f1)$(NC)"

# Get wasm_exec.js runtime
runtime: $(WASM_EXEC)

$(WASM_EXEC):
	@echo "$(BLUE)📥 Getting WASM runtime...$(NC)"
	@if [ -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" .; \
		echo "$(GREEN)✅ Copied from Go installation$(NC)"; \
	elif [ -f "./go-projects/pure-wasm/web/static/wasm_exec.js" ]; then \
		cp ./go-projects/pure-wasm/web/static/wasm_exec.js .; \
		echo "$(GREEN)✅ Copied from existing project$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  Downloading from Go repository...$(NC)"; \
		curl -f -s -o $(WASM_EXEC) https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js; \
		if [ $$? -eq 0 ]; then \
			echo "$(GREEN)✅ Downloaded wasm_exec.js$(NC)"; \
		else \
			echo "$(RED)❌ Failed to get wasm_exec.js$(NC)"; \
			exit 1; \
		fi \
	fi

# Build MCP server
mcp: $(MCP_SERVER)

$(MCP_SERVER): $(MCP_SRC)
	@echo "$(BLUE)🤖 Building MCP server...$(NC)"
	go build -o $(MCP_SERVER) $(MCP_SRC)
	@echo "$(GREEN)✅ MCP server built$(NC)"

# Create index.html
$(INDEX): index_phase1.html
	@echo "$(BLUE)📄 Preparing HTML...$(NC)"
	cp index_phase1.html $(INDEX)
	@echo "$(GREEN)✅ HTML ready$(NC)"

# Create test files
test-files:
	@echo "$(BLUE)📁 Creating test files...$(NC)"
	@mkdir -p test-files/subdir
	@echo "test content 1" > test-files/file1.txt
	@echo "test content 1" > test-files/file1_duplicate.txt
	@echo "test content 2" > test-files/file2.txt
	@echo "different content" > test-files/file3.txt
	@echo "test content 1" > test-files/subdir/file1.txt
	@echo "$(GREEN)✅ Test files created$(NC)"

# Check if all required files exist
check:
	@echo "$(BLUE)🔍 Checking files...$(NC)"
	@error=0; \
	if [ ! -f "$(WASM_FILE)" ]; then \
		echo "$(RED)❌ $(WASM_FILE) not found$(NC)"; \
		error=1; \
	else \
		echo "$(GREEN)✅ $(WASM_FILE) ($$(du -h $(WASM_FILE) | cut -f1))$(NC)"; \
	fi; \
	if [ ! -f "$(WASM_EXEC)" ]; then \
		echo "$(RED)❌ $(WASM_EXEC) not found$(NC)"; \
		error=1; \
	else \
		echo "$(GREEN)✅ $(WASM_EXEC)$(NC)"; \
	fi; \
	if [ ! -f "wasm-worker.js" ]; then \
		echo "$(RED)❌ wasm-worker.js not found$(NC)"; \
		error=1; \
	else \
		echo "$(GREEN)✅ wasm-worker.js$(NC)"; \
	fi; \
	if [ ! -f "$(INDEX)" ]; then \
		echo "$(RED)❌ $(INDEX) not found$(NC)"; \
		error=1; \
	else \
		echo "$(GREEN)✅ $(INDEX)$(NC)"; \
	fi; \
	if [ -d "test-files" ]; then \
		echo "$(GREEN)✅ test-files/ ($$(find test-files -type f | wc -l | tr -d ' ') files)$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  test-files/ not found (optional)$(NC)"; \
	fi; \
	echo ""; \
	if [ $$error -eq 0 ]; then \
		echo "$(GREEN)✅ All checks passed!$(NC)"; \
		echo ""; \
		echo "$(BLUE)Ready to run:$(NC)"; \
		echo "  make serve"; \
		echo ""; \
	else \
		echo "$(RED)❌ Some files missing$(NC)"; \
		echo ""; \
		echo "$(BLUE)Run:$(NC)"; \
		echo "  make build"; \
		echo ""; \
		exit 1; \
	fi

# Start HTTP server
serve: check
	@echo "$(BLUE)🚀 Starting server on port $(PORT)...$(NC)"
	@echo ""
	@echo "$(GREEN)Open in browser:$(NC)"
	@echo "  http://localhost:$(PORT)/$(INDEX)"
	@echo ""
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@echo ""
	@python3 -m http.server $(PORT) 2>&1 | grep -v "Traceback\|Error\|error" || true

# Run tests
test:
	@if [ -x "./test.sh" ]; then \
		./test.sh; \
	else \
		echo "$(YELLOW)⚠️  test.sh not found or not executable$(NC)"; \
	fi

# Clean built files
clean:
	@echo "$(BLUE)🧹 Cleaning...$(NC)"
	@rm -f $(WASM_FILE) $(WASM_EXEC) $(INDEX) $(MCP_SERVER)
	@rm -rf test-files/
	@echo "$(GREEN)✅ Cleaned$(NC)"

# Install dependencies (check Go version)
install: check-go
	@echo "$(GREEN)✅ Dependencies OK$(NC)"
	@echo ""
	@echo "$(BLUE)Go version:$(NC) $$(go version)"
	@echo ""

# Check if Go is installed
check-go:
	@which go > /dev/null || (echo "$(RED)❌ Go not found. Install from https://go.dev/dl/$(NC)" && exit 1)
	@echo "$(GREEN)✅ Go installed$(NC)"

# Development: watch and rebuild
watch:
	@echo "$(BLUE)👀 Watching for changes...$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@while true; do \
		inotifywait -e modify $(WASM_SRC) 2>/dev/null || fswatch -1 $(WASM_SRC) 2>/dev/null || sleep 5; \
		echo "$(BLUE)🔄 Rebuilding...$(NC)"; \
		make wasm; \
	done

# Show file sizes
sizes:
	@echo "$(BLUE)📊 File Sizes:$(NC)"
	@ls -lh $(WASM_FILE) $(WASM_EXEC) $(INDEX) 2>/dev/null || echo "Run 'make build' first"

# Quick serve (skip checks)
serve-quick:
	@echo "$(BLUE)🚀 Starting server (no checks)...$(NC)"
	@python3 -m http.server $(PORT)
