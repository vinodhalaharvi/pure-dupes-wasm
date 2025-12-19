#!/bin/bash

# check.sh - Verify build is complete before running

echo "🔍 Pre-flight Check"
echo "==================="
echo ""

ERRORS=0

# Check for required files
echo "Checking required files..."

if [ ! -f "main.wasm" ]; then
    echo "❌ main.wasm not found"
    echo "   Run: ./build.sh"
    ERRORS=$((ERRORS+1))
else
    SIZE=$(du -h main.wasm | cut -f1)
    echo "✅ main.wasm ($SIZE)"
fi

if [ ! -f "wasm_exec.js" ]; then
    echo "❌ wasm_exec.js not found"
    echo "   Run: ./build.sh"
    ERRORS=$((ERRORS+1))
else
    echo "✅ wasm_exec.js"
fi

if [ ! -f "wasm-worker.js" ]; then
    echo "❌ wasm-worker.js not found"
    ERRORS=$((ERRORS+1))
else
    echo "✅ wasm-worker.js"
fi

if [ ! -f "index.html" ]; then
    echo "❌ index.html not found"
    echo "   Run: ./build.sh"
    ERRORS=$((ERRORS+1))
else
    echo "✅ index.html"
fi

echo ""

# Check for test files
if [ -d "test-files" ]; then
    COUNT=$(find test-files -type f | wc -l)
    echo "✅ test-files/ ($COUNT files)"
else
    echo "⚠️  test-files/ not found (not critical)"
fi

echo ""

# Summary
if [ $ERRORS -eq 0 ]; then
    echo "✅ All checks passed!"
    echo ""
    echo "Ready to run:"
    echo "  ./serve.sh"
    echo ""
    exit 0
else
    echo "❌ $ERRORS error(s) found"
    echo ""
    echo "Please run:"
    echo "  ./build.sh"
    echo ""
    echo "Then try again:"
    echo "  ./check.sh"
    echo ""
    exit 1
fi
