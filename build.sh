#!/bin/bash

set -e

echo "🚀 Building Phase 1: Foundation Features"
echo "========================================"
echo ""

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check prerequisites
echo -e "${BLUE}Checking prerequisites...${NC}"
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi
echo "✅ Go: $(go version)"
echo ""

# Initialize go.mod if not present
if [ ! -f "go.mod" ]; then
    echo "Creating go.mod..."
    cat > go.mod << 'GOMOD'
module github.com/vinodhalaharvi/pure-dupes

go 1.21
GOMOD
fi

# Step 1: Build Enhanced WASM
echo -e "${BLUE}Step 1: Building Enhanced WASM Module${NC}"
echo "Features: Progress reporting, Smart groups, Caching support"
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm main_wasm_enhanced.go

if [ $? -eq 0 ]; then
    SIZE=$(du -h main.wasm | cut -f1)
    echo -e "${GREEN}✅ WASM built successfully ($SIZE)${NC}"
else
    echo "❌ WASM build failed"
    exit 1
fi
echo ""

# Step 2: Get wasm_exec.js
echo -e "${BLUE}Step 2: Getting WASM runtime${NC}"
GOROOT=$(go env GOROOT)
WASM_EXEC=""

for loc in "$GOROOT/misc/wasm/wasm_exec.js" "/usr/local/go/misc/wasm/wasm_exec.js"; do
    if [ -f "$loc" ]; then
        WASM_EXEC="$loc"
        break
    fi
done

if [ -z "$WASM_EXEC" ]; then
    echo -e "${YELLOW}⚠️  Downloading from Go repository...${NC}"
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    MAJOR_MINOR=$(echo $GO_VERSION | cut -d. -f1,2)
    
    curl -f -s "https://raw.githubusercontent.com/golang/go/release-branch.go${MAJOR_MINOR}/misc/wasm/wasm_exec.js" -o wasm_exec.js
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Downloaded wasm_exec.js${NC}"
    else
        echo "❌ Failed to download wasm_exec.js"
        exit 1
    fi
else
    cp "$WASM_EXEC" .
    echo -e "${GREEN}✅ Copied from: $WASM_EXEC${NC}"
fi
echo ""

# Step 3: Build MCP Server
echo -e "${BLUE}Step 3: Building MCP Server${NC}"
go build -o mcp-server mcp-server.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ MCP Server built${NC}"
else
    echo "❌ MCP Server build failed"
    exit 1
fi
echo ""

# Step 4: Prepare files
echo -e "${BLUE}Step 4: Preparing deployment files${NC}"
cp index_phase1.html index.html
echo -e "${GREEN}✅ HTML ready${NC}"
echo ""

# Step 5: Create test directory
echo -e "${BLUE}Step 5: Creating test files${NC}"
mkdir -p test-files
echo "test content 1" > test-files/file1.txt
echo "test content 1" > test-files/file1_duplicate.txt
echo "test content 2" > test-files/file2.txt
echo "different content" > test-files/file3.txt
mkdir -p test-files/subdir
cp test-files/file1.txt test-files/subdir/
echo -e "${GREEN}✅ Test files created${NC}"
echo ""

# Summary
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Phase 1 Build Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "📦 Files created:"
echo "  ├─ main.wasm ($(du -h main.wasm 2>/dev/null | cut -f1 || echo 'N/A'))"
echo "  ├─ wasm_exec.js ($(du -h wasm_exec.js 2>/dev/null | cut -f1 || echo 'N/A'))"
echo "  ├─ wasm-worker.js"
echo "  ├─ cache-db.js (inlined in HTML)"
echo "  ├─ index.html"
echo "  └─ mcp-server"
echo ""
echo "🚀 To test:"
echo "  ./serve.sh"
echo "  Then open: http://localhost:8080"
echo ""
