#!/bin/bash

# Quick fix: Download wasm_exec.js directly from Go repository

echo "📥 Downloading wasm_exec.js from Go repository..."

# Get Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
MAJOR_MINOR=$(echo $GO_VERSION | cut -d. -f1,2)

echo "Go version detected: $GO_VERSION"
echo "Using branch: go$MAJOR_MINOR"
echo ""

# Try the release branch first
URL="https://raw.githubusercontent.com/golang/go/release-branch.go${MAJOR_MINOR}/misc/wasm/wasm_exec.js"

echo "Downloading from: $URL"
curl -f -s "$URL" -o wasm_exec.js

if [ $? -eq 0 ] && [ -f wasm_exec.js ] && [ -s wasm_exec.js ]; then
    SIZE=$(du -h wasm_exec.js | cut -f1)
    echo "✅ Successfully downloaded wasm_exec.js ($SIZE)"
    echo ""
    echo "You can now run the build script again:"
    echo "  ./build_wasm.sh"
else
    echo "❌ Failed to download from release branch"
    echo ""
    echo "Trying master branch..."
    
    URL="https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js"
    curl -f -s "$URL" -o wasm_exec.js
    
    if [ $? -eq 0 ] && [ -f wasm_exec.js ] && [ -s wasm_exec.js ]; then
        SIZE=$(du -h wasm_exec.js | cut -f1)
        echo "✅ Successfully downloaded wasm_exec.js ($SIZE)"
        echo ""
        echo "You can now run the build script again:"
        echo "  ./build_wasm.sh"
    else
        echo "❌ Failed to download wasm_exec.js"
        echo ""
        echo "Please download manually from:"
        echo "  https://github.com/golang/go/blob/master/misc/wasm/wasm_exec.js"
        exit 1
    fi
fi
