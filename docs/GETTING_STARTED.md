# 🚀 Getting Started - Phase 1

## The Problem You're Seeing

You tried to open `index_phase1.html` by **double-clicking** it, which opens it as:
```
file:///Users/vinodhalaharvi/Downloads/index_phase1.html
```

This **doesn't work** because:
- ❌ Browser blocks loading external scripts
- ❌ Web Workers can't load
- ❌ WASM might not initialize

## ✅ The Solution (2 Steps)

### Step 1: Build Everything

```bash
cd /path/to/pure-dupes

# Make scripts executable
chmod +x build_phase1.sh serve.sh

# Build
./build_phase1.sh
```

### Step 2: Run the Server

```bash
# Option A: Use serve script (easiest)
./serve.sh

# Option B: Use Python directly
python3 -m http.server 8080

# Option C: Use different port
python3 -m http.server 9000
```

### Step 3: Open in Browser

```
http://localhost:8080/index_phase1.html
```

**NOT** by double-clicking the file!

---

## 🎯 Quick Start (Copy-Paste)

```bash
# All in one go:
cd /path/to/pure-dupes
chmod +x *.sh
./build_phase1.sh
./serve.sh
```

Then open: **http://localhost:8080/index_phase1.html**

---

## 🔍 What You'll See When It Works

### In Terminal:
```
🚀 Starting Phase 1 Server
==========================

✅ All files present
📡 Starting server on port 8080

🌐 Open in browser:
   http://localhost:8080/index_phase1.html

📚 Features available:
   ✅ Web Workers
   ✅ IndexedDB Caching
   ✅ Progress Reporting
   ✅ Smart Groups

Press Ctrl+C to stop
```

### In Browser Console:
```
✅ Cache initialized: {cachedFiles: 0}
✅ Web Worker ready
🔍 pure-dupes WASM initialized
```

### In UI:
- No yellow warning banner
- "Web Worker ready" status
- Can select files
- Progress bar works

---

## 🐛 Still Having Issues?

### Error: "Port 8080 in use"

**Solution:**
```bash
# Use different port
python3 -m http.server 8081

# Or use serve script (auto-finds free port)
./serve.sh
```

### Error: "Files not found"

**Solution:**
```bash
# Make sure you're in the right directory
ls -la *.wasm *.js *.html

# If files missing, rebuild
./build_phase1.sh
```

### Error: "Worker not ready"

**Check:**
1. Are you using http:// (not file://)?
2. Is wasm-worker.js present?
3. Is main.wasm present?
4. Check browser console for errors

**Fix:**
```bash
# Rebuild everything
./build_phase1.sh

# Start server
./serve.sh

# Open in browser (http:// not file://)
```

---

## 📖 What Each File Does

| File | Purpose |
|------|---------|
| `main.wasm` | WASM module with analysis engine |
| `wasm_exec.js` | Go WASM runtime |
| `wasm-worker.js` | Web Worker for background processing |
| `index_phase1.html` | Main UI (cache code inlined) |
| `mcp-server` | MCP server for Claude |

---

## 💡 Pro Tips

### Tip 1: Auto-open browser

```bash
# Mac
./serve.sh && open http://localhost:8080/index_phase1.html

# Linux
./serve.sh && xdg-open http://localhost:8080/index_phase1.html
```

### Tip 2: Different terminal

```bash
# Terminal 1 (server)
./serve.sh

# Terminal 2 (development)
# Edit files, rebuild, etc.
```

### Tip 3: Check it's working

```bash
# In browser console:
cacheDB.getStats()
// Should return: Promise {<pending>}
// Then: {cachedFiles: 0}
```

---

## 🎉 Success Checklist

You're ready when:

- [ ] Server running (terminal shows "Starting server")
- [ ] Browser at http://localhost:8080 (not file://)
- [ ] No yellow warning banner
- [ ] Console shows "✅ Web Worker ready"
- [ ] Can click "Choose Files"
- [ ] Upload shows progress bar
- [ ] Results display
- [ ] Cache stats update

**If all checked: You're all set! 🚀**

---

## 🆘 Emergency One-Liner

If nothing works:

```bash
cd /path/to/pure-dupes && chmod +x *.sh && ./build_phase1.sh && python3 -m http.server 8080
```

Then open: **http://localhost:8080/index_phase1.html**

---

## 📞 Still Stuck?

Check these docs:
1. **TROUBLESHOOTING.md** - Detailed error solutions
2. **QUICK_START_PHASE1.md** - Comprehensive guide
3. **PHASE1_IMPLEMENTATION.md** - Full documentation

Or run:
```bash
./test_phase1.sh  # Automated diagnostics
```

---

**Remember: ALWAYS use http://, NEVER file://! 🎯**
