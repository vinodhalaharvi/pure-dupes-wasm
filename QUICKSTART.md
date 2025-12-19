# ⚡ Quick Start Guide - Fixed Version

## 🎯 This Version Fixes

✅ **Go modules issue** - Added go.mod  
✅ **Simpler structure** - All files in one directory  
✅ **Easier builds** - Just run ./build.sh  

---

## 🚀 Installation (3 Commands)

```bash
# Extract
tar -xzf pure-dupes-phase1.tar.gz
cd pure-dupes-phase1-flat

# Build
./build.sh

# Run
./serve.sh
```

Then open: **http://localhost:8080**

---

## 📁 Simple Structure

```
pure-dupes-phase1-flat/
├── go.mod                    ← Go module file (NEW!)
├── main_wasm_enhanced.go     ← Main WASM code
├── wasm-worker.js            ← Web Worker
├── cache-db.js               ← Caching layer
├── index_phase1.html         ← UI
├── mcp-server.go             ← MCP server
│
├── build.sh                  ← Build everything
├── serve.sh                  ← Start server
├── test.sh                   ← Run tests
│
├── docs/                     ← All documentation
└── reference/                ← Original versions
```

**Everything in one place - no nested directories!**

---

## 🔨 Build Process

```bash
./build.sh
```

This will:
1. ✅ Check Go is installed
2. ✅ Use go.mod (no more "module not found" error)
3. ✅ Build main.wasm (~2-4 MB)
4. ✅ Download wasm_exec.js
5. ✅ Build mcp-server
6. ✅ Create index.html
7. ✅ Generate test files

**Output:**
```
main.wasm        ← Your WASM module
wasm_exec.js     ← Go runtime
index.html       ← Final UI
mcp-server       ← MCP binary
test-files/      ← Sample data
```

---

## 🧪 Test It Works

```bash
# Start server
./serve.sh

# Should see:
🚀 Starting Phase 1 Server
...
📡 Starting server on port 8080

🌐 Open in browser:
   http://localhost:8080
```

Then:
1. Open **http://localhost:8080**
2. Click "Choose Files"
3. Select `test-files/` directory
4. ✅ Should find 2 exact duplicates!

---

## ✅ What You Fixed

### Before (Broken):
```bash
./build.sh
❌ no required module provides package main_wasm_enhanced.go
```

### After (Working):
```bash
./build.sh
✅ WASM built successfully (2.1M)
✅ Phase 1 Build Complete!
```

**The fix:** Added `go.mod` file!

---

## 🎯 Features Ready

- ⚡ **Web Workers** - Non-blocking UI
- 💾 **Caching** - 10-100x faster re-scans
- 📊 **Smart Groups** - Intelligent duplicates
- 📈 **Progress** - Real-time updates
- 🤖 **MCP Server** - Claude integration

---

## 🐛 Still Having Issues?

### Go not found
```bash
# Install Go
brew install go  # Mac
```

### Permission denied
```bash
chmod +x *.sh
./build.sh
```

### Port 8080 in use
```bash
# Use different port
python3 -m http.server 8081
```

### Worker not loading
**Make sure you're using http:// not file://**
```bash
./serve.sh  # ← Always use this
```

---

## 📖 Documentation

All docs are in the `docs/` folder:

- `GETTING_STARTED.md` - Detailed setup
- `TROUBLESHOOTING.md` - Common errors
- `PHASE1_IMPLEMENTATION.md` - Full features
- `QUICK_START_PHASE1.md` - Reference guide

---

## 🎉 Success Checklist

You're ready when:

- [ ] `./build.sh` completes without errors
- [ ] You see "Phase 1 Build Complete!"
- [ ] Files exist: main.wasm, wasm_exec.js, index.html
- [ ] `./serve.sh` starts server
- [ ] Can open http://localhost:8080
- [ ] Can upload test-files/
- [ ] Finds 2 duplicate files
- [ ] Progress bar appears
- [ ] Results display

**If all checked: Perfect! 🚀**

---

## 💡 Pro Tips

### Quick rebuild
```bash
# Only rebuild WASM
GOOS=js GOARCH=wasm go build -o main.wasm main_wasm_enhanced.go
```

### Check what was built
```bash
ls -lah main.wasm wasm_exec.js index.html mcp-server
```

### Test MCP server
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | ./mcp-server
```

---

## 🚀 Deploy When Ready

### GitHub Pages
```bash
git add main.wasm wasm_exec.js wasm-worker.js index.html
git commit -m "Phase 1 ready"
git push
```

### Netlify/Vercel
Just upload these 4 files:
- main.wasm
- wasm_exec.js
- wasm-worker.js
- index.html

---

## 📞 Need Help?

1. Check `docs/TROUBLESHOOTING.md`
2. Run `./test.sh` for diagnostics
3. Check console for errors
4. Verify files: `ls -la`

---

**The build issue is now fixed! Just extract and run ./build.sh! 🎉**
