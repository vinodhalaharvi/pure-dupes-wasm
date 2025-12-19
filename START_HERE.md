# 🎯 START HERE

## Your Error: Worker Not Loading

**Why?** You ran `./serve.sh` **before** running `./build.sh`!

The files `main.wasm` and `wasm_exec.js` don't exist yet.

---

## ✅ The Fix (Do This Now)

```bash
# Step 1: Build (MUST DO FIRST!)
./build.sh

# Step 2: Verify
./check.sh

# Step 3: Run
./serve.sh
```

That's it!

---

## 📋 What Each Does

### `./build.sh`
Creates:
- ✅ main.wasm (WASM module)
- ✅ wasm_exec.js (Go runtime)
- ✅ index.html (UI)
- ✅ mcp-server (Claude)
- ✅ test-files/ (samples)

**Takes:** ~10-30 seconds

### `./check.sh`
Verifies all files exist before running server

**Takes:** < 1 second

### `./serve.sh`
Starts HTTP server on port 8080

**Takes:** Instant (runs until Ctrl+C)

---

## 🎯 Current Status

You are here:
```
❌ Step 1: Build      ← YOU SKIPPED THIS!
⏭️  Step 2: Verify
⏭️  Step 3: Run
```

**Fix it:**
```bash
./build.sh    ← Do this now!
./check.sh    ← Verify it worked
./serve.sh    ← Then run server
```

---

## 🚨 Common Mistake

**WRONG:**
```bash
./serve.sh    ← NO! Files don't exist yet!
```

**RIGHT:**
```bash
./build.sh    ← Build first
./serve.sh    ← Then serve
```

---

## ✅ Success Looks Like

### After `./build.sh`
```
✅ WASM built successfully (2.1M)
✅ Downloaded wasm_exec.js
✅ MCP Server built
✅ HTML ready
✅ Test files created
✅ Phase 1 Build Complete!
```

### After `./check.sh`
```
✅ main.wasm (2.1M)
✅ wasm_exec.js
✅ wasm-worker.js
✅ index.html
✅ test-files/ (5 files)

✅ All checks passed!
Ready to run: ./serve.sh
```

### After `./serve.sh`
```
✅ All files present
📡 Starting server on port 8080

🌐 Open in browser:
   http://localhost:8080
```

---

## 🐛 Still Getting Errors?

### "go: command not found"
**Fix:** Install Go
```bash
brew install go  # Mac
```

### "Permission denied"
**Fix:** Make scripts executable
```bash
chmod +x *.sh
```

### "Worker error: NetworkError"
**Fix:** You didn't build first!
```bash
./build.sh  # ← This creates the files the worker needs
```

---

## 🎉 Quick Start (Copy-Paste)

```bash
# Do these 3 commands in order:
./build.sh
./check.sh
./serve.sh
```

Then open: **http://localhost:8080**

---

## 📖 More Help

- **QUICKSTART.md** - Detailed guide
- **INSTALLATION.md** - Full instructions
- **docs/TROUBLESHOOTING.md** - Error solutions

---

**TL;DR: Run ./build.sh first! Then ./serve.sh!** 🚀
