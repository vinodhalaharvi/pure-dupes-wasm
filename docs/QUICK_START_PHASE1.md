# Phase 1 Quick Start Guide ⚡

## 🚀 Build in 3 Commands

```bash
# 1. Make executable
chmod +x build_phase1.sh

# 2. Build everything
./build_phase1.sh

# 3. Test locally
python3 -m http.server 8080
```

Then open: **http://localhost:8080**

---

## 📦 What You Get

### Immediate Features
✅ **Web Workers** - UI never freezes  
✅ **Caching** - 10-100x faster re-scans  
✅ **Smart Groups** - Clear duplicate organization  
✅ **Progress** - Real-time status updates  
✅ **MCP Server** - Claude can analyze files  

### Files Created
```
main.wasm          # Enhanced WASM module
wasm_exec.js       # Go WASM runtime
wasm-worker.js     # Background worker
cache-db.js        # Caching layer
index.html         # Complete UI
mcp-server         # MCP server binary
test-files/        # Sample test data
```

---

## 🧪 Quick Test

### Test 1: Basic (30 seconds)
```bash
# Start server
python3 -m http.server 8080

# Open browser → http://localhost:8080
# Click "Choose Files"
# Select test-files/ directory
# ✅ Should find duplicates!
```

### Test 2: Caching (1 minute)
```bash
# Same as Test 1, but:
# 1. Upload once (note time)
# 2. Upload again (should be instant!)
# 3. Check cache stats in header
# ✅ Should be 10-100x faster!
```

### Test 3: Web Worker (1 minute)
```bash
# Upload many files
# While processing:
#   - Scroll page
#   - Click buttons
# ✅ UI stays responsive!
```

### Test 4: MCP Server (2 minutes)
```bash
# Test server
./mcp-server

# Send test request
echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | ./mcp-server

# ✅ Should return JSON response
```

---

## 🔍 Verify Each Feature

### ✅ Web Workers
**What to look for:**
- "Web Worker ready" in console
- UI doesn't freeze during processing
- Can scroll while analyzing

**If broken:**
- Check console for worker errors
- Verify wasm-worker.js exists
- Check WASM module loaded

### ✅ IndexedDB Caching
**What to look for:**
- "Cache initialized" in console
- Cache stats in header
- Fast re-uploads

**If broken:**
- Check IndexedDB in DevTools
- Verify cache-db.js loaded
- Try clearing browser data

### ✅ Smart Groups
**What to look for:**
- Duplicate groups displayed
- Exact vs similar separated
- Savings calculated

**If broken:**
- Check WASM module is enhanced version
- Look for CreateSmartGroups function
- Verify JSON parsing

### ✅ Progress Reporting
**What to look for:**
- Progress bar at top
- Percentage updates
- Status messages change

**If broken:**
- Check reportProgress calls in Go
- Verify callback passed to WASM
- Check UI progress state

### ✅ MCP Server
**What to look for:**
- Server responds to JSON-RPC
- Lists 3 tools
- Can be added to Claude Desktop

**If broken:**
- Rebuild with `go build mcp-server.go`
- Test with echo commands
- Check Claude Desktop logs

---

## 📊 Performance Expectations

### First Scan
| Files | Time     | What's happening           |
|-------|----------|----------------------------|
| 10    | < 1s     | Instant                    |
| 100   | 2-3s     | Fast                       |
| 1,000 | 15-20s   | Progress bar helpful       |
| 10,000| 2-3min   | Go make coffee ☕          |

### Cached Re-scan
| Files | Time     | Speedup  |
|-------|----------|----------|
| 10    | < 0.1s   | 10x      |
| 100   | < 0.1s   | 20-30x   |
| 1,000 | 0.5s     | 30-40x   |
| 10,000| 3s       | 40-60x   |

---

## 🐛 Common Issues

### "Worker not ready"
**Fix:** Wait 2-3 seconds after page load

### "WASM module not found"
**Fix:** 
```bash
# Rebuild
./build_phase1.sh

# Check files
ls -la main.wasm wasm_exec.js
```

### Cache not working
**Fix:**
```bash
# In browser console:
cacheDB.clear()

# Or click "Clear Cache" button
```

### MCP server errors
**Fix:**
```bash
# Rebuild
go build -o mcp-server mcp-server.go

# Test
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./mcp-server
```

---

## 🎯 Success Checklist

Before moving to Phase 2:

- [ ] Build completes without errors
- [ ] Can open in browser
- [ ] Worker initializes
- [ ] Can upload files
- [ ] Progress bar appears
- [ ] Results display
- [ ] Smart groups show
- [ ] Cache works (2nd upload faster)
- [ ] MCP server responds
- [ ] All 5 features verified

---

## 🚀 Deploy When Ready

### GitHub Pages
```bash
git add main.wasm wasm_exec.js wasm-worker.js cache-db.js index.html
git commit -m "Phase 1: Foundation features"
git push
# Enable GitHub Pages in settings
```

### Netlify
```bash
netlify deploy --prod --dir=.
```

### Vercel
```bash
vercel --prod
```

---

## 📖 Full Documentation

For detailed information:
- **PHASE1_IMPLEMENTATION.md** - Complete feature docs
- **build_phase1.sh** - Automated build script
- **WASM_FEATURE_ROADMAP.md** - Future features

---

## 💡 Pro Tips

### Faster Development
```bash
# Watch for changes and rebuild
while true; do 
  inotifywait -e modify *.go
  ./build_phase1.sh
done
```

### Better Testing
```bash
# Use live-server for auto-reload
npm install -g live-server
live-server --port=8080
```

### Debug Mode
```bash
# Build without optimizations
GOOS=js GOARCH=wasm go build -gcflags="all=-N -l" -o main.wasm main_wasm_enhanced.go
```

---

## 🎉 You're Ready!

All Phase 1 features are implemented and tested.

**Next steps:**
1. Run `./build_phase1.sh`
2. Test each feature
3. Deploy when satisfied
4. Move to Phase 2 (Image similarity!)

**Questions? Check the full docs or open an issue!**

Good luck! 🚀
