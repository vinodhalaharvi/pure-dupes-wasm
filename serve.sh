#!/bin/bash

# serve.sh - Simple HTTP server for Phase 1
# Handles all the file protocol issues

PORT=${1:-8080}

echo "🚀 Starting Phase 1 Server"
echo "=========================="
echo ""

# Check if files exist
if [ ! -f "index_phase1.html" ]; then
    echo "❌ index_phase1.html not found"
    echo "💡 Run ./build_phase1.sh first"
    exit 1
fi

if [ ! -f "main.wasm" ]; then
    echo "❌ main.wasm not found"
    echo "💡 Run ./build_phase1.sh first"
    exit 1
fi

# Find available port
check_port() {
    lsof -i :$1 >/dev/null 2>&1
}

while check_port $PORT; do
    echo "⚠️  Port $PORT in use, trying $((PORT+1))"
    PORT=$((PORT+1))
done

echo "✅ All files present"
echo "📡 Starting server on port $PORT"
echo ""
echo "🌐 Open in browser:"
echo "   http://localhost:$PORT/index_phase1.html"
echo ""
echo "📚 Features available:"
echo "   ✅ Web Workers"
echo "   ✅ IndexedDB Caching"
echo "   ✅ Progress Reporting"
echo "   ✅ Smart Groups"
echo ""
echo "Press Ctrl+C to stop"
echo "---"
echo ""

# Start server based on what's available
if command -v python3 &> /dev/null; then
    python3 -m http.server $PORT
elif command -v python &> /dev/null; then
    python -m SimpleHTTPServer $PORT
elif command -v php &> /dev/null; then
    php -S localhost:$PORT
else
    echo "❌ No HTTP server found"
    echo "💡 Install Python or PHP:"
    echo "   brew install python3  # Mac"
    echo "   apt-get install python3  # Linux"
    exit 1
fi
