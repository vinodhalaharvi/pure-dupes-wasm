# 📦 Installation Guide - Pure Dupes Phase 1

## 📥 What You Downloaded

You have **pure-dupes-phase1.tar.gz** (or .zip) - a complete package with all Phase 1 features!

**Package size:** ~54-70 KB  
**Contains:** 21 files organized in 5 directories  

---

## 🚀 Installation (3 Steps)

### Step 1: Extract

**Using .tar.gz:**
```bash
tar -xzf pure-dupes-phase1.tar.gz
cd pure-dupes-phase1
```

**Using .zip:**
```bash
unzip pure-dupes-phase1.zip
cd pure-dupes-phase1
```

### Step 2: Build

```bash
chmod +x *.sh scripts/*.sh
./build.sh
```

This will:
- Compile WASM module
- Download Go runtime
- Build MCP server
- Create test files
- Verify everything

### Step 3: Run

```bash
./serve.sh
```

Then open: **http://localhost:8080/index.html**

---

## 📁 What's Inside

```
pure-dupes-phase1/
├── README_FIRST.md          ← START HERE
├── build.sh                 ← Convenience: builds everything
├── serve.sh                 ← Convenience: starts server
├── test.sh                  ← Convenience: runs tests
│
├── core/                    ← Implementation files
│   ├── main_wasm_enhanced.go
│   ├── wasm-worker.js
│   ├── cache-db.js
│   ├── index_phase1.html
│   └── mcp-server.go
│
├── scripts/                 ← Build & test scripts
│   ├── build_phase1.sh
│   ├── serve.sh
│   ├── test_phase1.sh
│   ├── build_wasm.sh
│   └── download_wasm_exec.sh
│
├── docs/                    ← Documentation
│   ├── GETTING_STARTED.md
│   ├── QUICK_START_PHASE1.md
│   ├── PHASE1_IMPLEMENTATION.md
│   ├── PHASE1_COMPLETE_SUMMARY.md
│   ├── TROUBLESHOOTING.md
│   └── WASM_FEATURE_ROADMAP.md
│
└── reference/               ← Original versions
    ├── main.go
    ├── main_wasm.go
    ├── index_wasm.html
    └── README.md
```

---

## ⚡ Quick Start (Copy-Paste)

```bash
# Extract
tar -xzf pure-dupes-phase1.tar.gz
cd pure-dupes-phase1

# Make executable
chmod +x *.sh scripts/*.sh

# Build
./build.sh

# Run
./serve.sh
```

Open: **http://localhost:8080/index.html**

---

## 🎯 What to Read

**First time?** Read in this order:
1. `README_FIRST.md` - Overview
2. `docs/GETTING_STARTED.md` - Setup guide
3. `docs/QUICK_START_PHASE1.md` - Quick reference

**Having issues?**
- `docs/TROUBLESHOOTING.md` - Error solutions

**Want to learn more?**
- `docs/PHASE1_IMPLEMENTATION.md` - Full docs
- `docs/WASM_FEATURE_ROADMAP.md` - Future features

---

## 🧪 Testing

### Quick Test
```bash
# Run automated tests
./test.sh

# Expected output:
# ✅ All tests passed!
```

### Manual Test
```bash
# Start server
./serve.sh

# Open browser
http://localhost:8080/index.html

# Upload test files
# - Click "Choose Files"
# - Select test-files/ directory
# - Should find duplicates!
```

---

## 🎯 Features Included

### 1. ⚡ Web Workers
- Background processing
- Non-blocking UI
- Responsive interface

### 2. 💾 IndexedDB Caching
- 10-100x faster re-scans
- Persistent storage
- Incremental updates

### 3. 📊 Smart Duplicate Groups
- Intelligent grouping
- Exact vs similar
- Savings calculation

### 4. 📈 Progress Reporting
- Real-time updates
- Stage-by-stage progress
- Percentage display

### 5. 🤖 MCP Server
- Claude integration
- 3 tools available
- Directory analysis

---

## 📋 Prerequisites

**Required:**
- Go 1.21 or later
- Python 3 (for HTTP server)
- Modern browser (Chrome, Firefox, Safari, Edge)

**Check versions:**
```bash
go version        # Should be 1.21+
python3 --version # Should be 3.x
```

**Install if needed:**
```bash
# Mac
brew install go python3

# Linux (Ubuntu/Debian)
sudo apt-get install golang python3

# Linux (Fedora/RHEL)
sudo dnf install golang python3
```

---

## 🐛 Common Issues

### "go: command not found"
**Fix:** Install Go from https://go.dev/dl/

### "Permission denied"
**Fix:** Make scripts executable
```bash
chmod +x *.sh scripts/*.sh
```

### "Port 8080 in use"
**Fix:** Use different port
```bash
cd scripts
./serve.sh 8081  # Uses port 8081 instead
```

### "Worker not ready"
**Fix:** Use HTTP server (not file://)
```bash
./serve.sh
# Then open http://localhost:8080
```

See `docs/TROUBLESHOOTING.md` for more solutions.

---

## 📊 Build Output

After running `./build.sh`, you'll have:

```
main.wasm          # WASM module (~2-4 MB)
wasm_exec.js       # Go runtime (~13 KB)
index.html         # Final UI (copied from core/)
mcp-server         # MCP binary
test-files/        # Sample test data
```

**Total size:** ~2-5 MB

---

## 🎨 Project Structure After Build

```
pure-dupes-phase1/
├── main.wasm              ← Generated
├── wasm_exec.js           ← Generated
├── index.html             ← Generated
├── mcp-server             ← Generated
├── test-files/            ← Generated
│   ├── file1.txt
│   ├── file1_duplicate.txt
│   └── ...
│
├── (all original dirs remain)
```

---

## 🚀 Deployment

### Local Development
```bash
./serve.sh
# Access at http://localhost:8080
```

### GitHub Pages
```bash
git add main.wasm wasm_exec.js index.html
git commit -m "Deploy Phase 1"
git push
```

### Netlify/Vercel
Upload these files:
- main.wasm
- wasm_exec.js
- index.html
- wasm-worker.js (from core/)
- cache-db.js (from core/, if needed)

---

## 💡 Tips

### Faster Builds
```bash
# Skip tests
cd scripts && ./build_phase1.sh --skip-tests
```

### Different Port
```bash
# Use port 9000 instead
cd scripts && ./serve.sh 9000
```

### Clean Build
```bash
# Remove generated files
rm -f main.wasm wasm_exec.js index.html mcp-server
rm -rf test-files/

# Rebuild
./build.sh
```

---

## 🔄 Updates

### Get Latest Version
```bash
# Re-download package
# Extract to new directory
# Copy your changes if any
```

### Rebuild After Changes
```bash
# If you modified Go code
./build.sh

# If you only changed HTML/JS
cd core && cp index_phase1.html ../index.html
```

---

## 📖 Documentation

| File | Purpose |
|------|---------|
| README_FIRST.md | Package overview |
| GETTING_STARTED.md | Setup guide |
| QUICK_START_PHASE1.md | Quick reference |
| PHASE1_IMPLEMENTATION.md | Full feature docs |
| TROUBLESHOOTING.md | Error solutions |
| WASM_FEATURE_ROADMAP.md | Future features |

---

## 🎓 Next Steps

After Phase 1 is working:

1. **Test thoroughly** - Upload various file types
2. **Deploy** - Put on GitHub Pages or Netlify
3. **Phase 2** - Add image similarity (pHash)
4. **Phase 3** - Add audio/video deduplication
5. **Phase 4** - Add ML integration

See `docs/WASM_FEATURE_ROADMAP.md` for details.

---

## 📞 Support

**Having issues?**
1. Check `docs/TROUBLESHOOTING.md`
2. Run `./test.sh` for diagnostics
3. Check browser console for errors
4. Verify files: `ls -la`

**Still stuck?**
- Review all documentation in `docs/`
- Check file permissions: `chmod +x *.sh scripts/*.sh`
- Try clean rebuild: delete generated files, rebuild

---

## 🎉 You're Ready!

```bash
# Quick start:
tar -xzf pure-dupes-phase1.tar.gz
cd pure-dupes-phase1
chmod +x *.sh scripts/*.sh
./build.sh
./serve.sh
```

Then open: **http://localhost:8080/index.html**

**Happy deduplicating! 🔍✨**

---

**Package Version:** Phase 1 Complete  
**Release Date:** December 2024  
**License:** MIT  
**Made with:** Go + WASM + Functional Programming ❤️
