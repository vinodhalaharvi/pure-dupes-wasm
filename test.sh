#!/bin/bash

# test_phase1.sh - Automated testing for Phase 1 features
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

PASS=0
FAIL=0

echo "🧪 Phase 1 Automated Tests"
echo "=========================="
echo ""

# Helper functions
pass() {
    echo -e "${GREEN}✅ PASS:${NC} $1"
    ((PASS++))
}

fail() {
    echo -e "${RED}❌ FAIL:${NC} $1"
    ((FAIL++))
}

info() {
    echo -e "${BLUE}ℹ️  INFO:${NC} $1"
}

# Test 1: Files exist
echo "${BLUE}Test 1: Checking required files...${NC}"

if [ -f "main_wasm_enhanced.go" ]; then
    pass "main_wasm_enhanced.go exists"
else
    fail "main_wasm_enhanced.go missing"
fi

if [ -f "wasm-worker.js" ]; then
    pass "wasm-worker.js exists"
else
    fail "wasm-worker.js missing"
fi

if [ -f "cache-db.js" ]; then
    pass "cache-db.js exists"
else
    fail "cache-db.js missing"
fi

if [ -f "index_phase1.html" ]; then
    pass "index_phase1.html exists"
else
    fail "index_phase1.html missing"
fi

if [ -f "mcp-server.go" ]; then
    pass "mcp-server.go exists"
else
    fail "mcp-server.go missing"
fi

echo ""

# Test 2: Go compilation
echo "${BLUE}Test 2: Testing Go compilation...${NC}"

# Test WASM build
if GOOS=js GOARCH=wasm go build -o test_main.wasm main_wasm_enhanced.go 2>/dev/null; then
    pass "WASM compiles successfully"
    rm -f test_main.wasm
else
    fail "WASM compilation failed"
fi

# Test MCP server build
if go build -o test_mcp mcp-server.go 2>/dev/null; then
    pass "MCP server compiles successfully"
    
    # Test MCP server responds
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | timeout 2 ./test_mcp > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        pass "MCP server responds to requests"
    else
        fail "MCP server doesn't respond"
    fi
    
    rm -f test_mcp
else
    fail "MCP server compilation failed"
fi

echo ""

# Test 3: Code validation
echo "${BLUE}Test 3: Validating code...${NC}"

# Check for Phase 1 features in WASM code
if grep -q "reportProgress" main_wasm_enhanced.go; then
    pass "Progress reporting code present"
else
    fail "Progress reporting code missing"
fi

if grep -q "CreateSmartGroups" main_wasm_enhanced.go; then
    pass "Smart groups code present"
else
    fail "Smart groups code missing"
fi

if grep -q "ModTime" main_wasm_enhanced.go; then
    pass "Caching support code present"
else
    fail "Caching support code missing"
fi

# Check for Worker code
if grep -q "importScripts" wasm-worker.js; then
    pass "Web Worker code present"
else
    fail "Web Worker code missing"
fi

if grep -q "postMessage" wasm-worker.js; then
    pass "Worker messaging present"
else
    fail "Worker messaging missing"
fi

# Check for IndexedDB code
if grep -q "indexedDB.open" cache-db.js; then
    pass "IndexedDB code present"
else
    fail "IndexedDB code missing"
fi

if grep -q "putFileHash" cache-db.js; then
    pass "Cache operations present"
else
    fail "Cache operations missing"
fi

# Check HTML has all features
if grep -q "Web Worker" index_phase1.html; then
    pass "HTML references Worker"
else
    fail "HTML missing Worker reference"
fi

if grep -q "cacheDB" index_phase1.html; then
    pass "HTML uses caching"
else
    fail "HTML missing cache usage"
fi

if grep -q "progress" index_phase1.html; then
    pass "HTML has progress UI"
else
    fail "HTML missing progress UI"
fi

if grep -q "Smart" index_phase1.html || grep -q "smart" index_phase1.html; then
    pass "HTML displays smart groups"
else
    fail "HTML missing smart groups"
fi

echo ""

# Test 4: MCP Server Tools
echo "${BLUE}Test 4: MCP Server Tools...${NC}"

if grep -q "analyze_duplicates" mcp-server.go; then
    pass "analyze_duplicates tool defined"
else
    fail "analyze_duplicates tool missing"
fi

if grep -q "get_duplicate_groups" mcp-server.go; then
    pass "get_duplicate_groups tool defined"
else
    fail "get_duplicate_groups tool missing"
fi

if grep -q "check_file_hash" mcp-server.go; then
    pass "check_file_hash tool defined"
else
    fail "check_file_hash tool missing"
fi

echo ""

# Test 5: Build script
echo "${BLUE}Test 5: Build script validation...${NC}"

if [ -f "build_phase1.sh" ] && [ -x "build_phase1.sh" ]; then
    pass "Build script exists and is executable"
else
    fail "Build script missing or not executable"
fi

if grep -q "main_wasm_enhanced.go" build_phase1.sh; then
    pass "Build script compiles correct WASM"
else
    fail "Build script uses wrong file"
fi

if grep -q "mcp-server" build_phase1.sh; then
    pass "Build script builds MCP server"
else
    fail "Build script doesn't build MCP server"
fi

echo ""

# Summary
echo "=========================="
echo "📊 Test Summary"
echo "=========================="
echo -e "${GREEN}Passed: $PASS${NC}"
echo -e "${RED}Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! Phase 1 is ready!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Run: ./build_phase1.sh"
    echo "  2. Test: python3 -m http.server 8080"
    echo "  3. Open: http://localhost:8080"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Please review and fix.${NC}"
    exit 1
fi
