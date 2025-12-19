# Troubleshooting Guide - Phase 1 Errors

## 🔴 Common Errors & Solutions

### Error 1: "cdn.tailwindcss.com should not be used in production"

**What it means:**
- This is a WARNING, not an error
- Tailwind CDN is fine for development/testing
- For production, should use compiled CSS

**Fix for Production:**
```bash
# Install Tailwind
npm install -D tailwindcss
npx tailwindcss init

# Build CSS
npx tailwindcss -i ./src/input.css -o ./dist/output.css
```

**For now (testing):** 
✅ **Ignore this warning** - it doesn't break functionality

---

### Error 2: "Loading failed for cache-db.js"

**What it means:**
- File opened via `file://` protocol (double-clicked HTML)
- Browser blocks loading external scripts from file system

**Why it happens:**
```
file:///Users/you/Downloads/index.html
                             ↑
                    file:// protocol blocks script loading
```

**✅ Solution 1: Use HTTP Server (RECOMMENDED)**

```bash
# In the directory with your files:
python3 -m http.server 8080

# Then open:
http://localhost:8080
```

**✅ Solution 2: Fixed Version**

The updated `index_phase1.html` now has cache-db.js **inlined** - no external file needed!

**Test the fix:**
1. Download new `index_phase1.html`
2. Open in browser
3. Should work (but Worker still needs HTTP server)

---

### Error 3: "cacheDB is not defined"

**What it means:**
- cache-db.js didn't load
- So `cacheDB` variable doesn't exist

**Solution:**
✅ Use the fixed `index_phase1.html` - cache code is now inlined

---

### Error 4: "You are using the in-browser Babel transformer"

**What it means:**
- Babel is transpiling JSX in browser
- This is SLOW for production
- Fine for development

**Fix for Production:**
```bash
# Pre-compile JSX
npm install --save-dev @babel/core @babel/cli @babel/preset-react

# Compile
npx babel src --out-dir dist --presets @babel/preset-react
```

**For now (testing):**
✅ **Ignore this warning** - works fine for testing

---

### Error 5: Web Worker Not Loading

**What it means:**
- Worker can't load from `file://` protocol
- Needs HTTP server

**Solution:**
```bash
# MUST use HTTP server for Workers
python3 -m http.server 8080

# Then open:
http://localhost:8080
```

**What happens without Worker:**
- ⚠️ UI will freeze during processing
- ✅ But analysis will still work
- ✅ Caching still works

---

## 🎯 Quick Fix Checklist

### If you see ANY errors:

**Step 1: Use HTTP Server**
```bash
cd /path/to/your/files
python3 -m http.server 8080
```

**Step 2: Open in browser**
```
http://localhost:8080/index_phase1.html
```

**Step 3: Check console**
- Should see: "✅ Cache initialized"
- Should see: "✅ Web Worker ready"

**Step 4: Test upload**
- Click "Choose Files"
- Select test-files/ directory
- Should work perfectly!

---

## 🔍 Diagnostic Steps

### Check 1: Protocol
```javascript
// In browser console:
console.log(window.location.protocol);

// Should see: "http:"
// NOT: "file:"
```

### Check 2: Cache Working
```javascript
// In browser console:
cacheDB.getStats().then(console.log);

// Should see: {cachedFiles: 0} (or number)
```

### Check 3: Worker Status
```javascript
// In browser console:
// Should see these messages:
// "✅ Cache initialized"
// "✅ Web Worker ready"
```

### Check 4: Files Present
```bash
# Check these files exist:
ls -la main.wasm wasm_exec.js wasm-worker.js index_phase1.html

# All should be present
```

---

## 🛠️ Build Issues

### Issue: Build script fails

**Check Go version:**
```bash
go version
# Should be 1.21+
```

**Rebuild:**
```bash
chmod +x build_phase1.sh
./build_phase1.sh
```

### Issue: WASM doesn't load

**Check file size:**
```bash
ls -lh main.wasm
# Should be 2-4 MB
```

**Rebuild WASM:**
```bash
GOOS=js GOARCH=wasm go build -o main.wasm main_wasm_enhanced.go
```

### Issue: wasm_exec.js missing

**Download manually:**
```bash
curl -o wasm_exec.js https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js
```

---

## 🌐 Browser Compatibility

### Supported Browsers

| Browser | Version | Status |
|---------|---------|--------|
| Chrome  | 87+     | ✅ Full support |
| Firefox | 78+     | ✅ Full support |
| Safari  | 14+     | ✅ Full support |
| Edge    | 88+     | ✅ Full support |

### Check Browser Features

```javascript
// Test WASM support
if (typeof WebAssembly === 'object') {
    console.log('✅ WASM supported');
}

// Test Worker support
if (typeof Worker !== 'undefined') {
    console.log('✅ Web Workers supported');
}

// Test IndexedDB support
if (window.indexedDB) {
    console.log('✅ IndexedDB supported');
}
```

---

## 📱 Mobile Issues

### iOS Safari

**Issue:** May have memory limits

**Solution:**
- Process smaller batches
- Reduce chunk size
- Clear cache frequently

### Android Chrome

**Issue:** Worker may be slower

**Solution:**
- Use smaller file sets
- Increase chunk size
- Monitor memory usage

---

## 💾 Cache Issues

### Cache not saving

**Check storage:**
```javascript
// In console
navigator.storage.estimate().then(console.log);
```

**Clear and retry:**
```javascript
cacheDB.clear().then(() => {
    console.log('Cache cleared');
});
```

### Cache corruption

**Reset everything:**
```javascript
// In console
indexedDB.deleteDatabase('pure-dupes-cache');
location.reload();
```

---

## 🚨 Emergency Fixes

### Nothing works!

**Full reset:**
```bash
# 1. Clear browser cache
# Chrome: Ctrl+Shift+Del → Clear cache

# 2. Delete IndexedDB
# Chrome DevTools → Application → IndexedDB → Delete

# 3. Rebuild everything
./build_phase1.sh

# 4. Use HTTP server
python3 -m http.server 8080

# 5. Open fresh
http://localhost:8080
```

### Still broken?

**Check basics:**
```bash
# Files exist?
ls -la *.wasm *.js *.html

# Correct directory?
pwd

# Server running?
lsof -i :8080

# Port available?
python3 -m http.server 8081  # Try different port
```

---

## 📖 Reference: Correct Setup

### Files needed (minimum):
```
✅ main.wasm              # WASM module
✅ wasm_exec.js           # Go runtime
✅ wasm-worker.js         # Web Worker
✅ index_phase1.html      # Main HTML (cache inlined!)
```

### Files needed (with Worker):
```
All above, PLUS:
✅ HTTP server running
✅ Open via http:// not file://
```

### Correct URL:
```
✅ http://localhost:8080/index_phase1.html
❌ file:///Users/you/Downloads/index_phase1.html
```

---

## 🎓 Understanding the Errors

### Why file:// doesn't work

**Security reasons:**
1. CORS restrictions
2. Module loading blocked
3. Worker creation blocked
4. Fetch API limited

**What works on file://**
- ✅ Basic HTML/CSS
- ✅ Inline JavaScript
- ❌ External scripts (blocked)
- ❌ Web Workers (blocked)
- ❌ Fetch requests (limited)

### Why HTTP server needed

**Enables:**
- ✅ Proper MIME types
- ✅ CORS headers
- ✅ Worker loading
- ✅ Module imports
- ✅ Full WASM support

---

## 💡 Pro Tips

### Faster Development

**Use live server:**
```bash
npm install -g live-server
live-server --port=8080
```

**Auto-reload on changes:**
```bash
# Install watchexec
brew install watchexec  # Mac

# Watch and rebuild
watchexec -e go ./build_phase1.sh
```

### Better Debugging

**Enable verbose console:**
```javascript
// In console
localStorage.debug = '*';
location.reload();
```

**Check Worker messages:**
```javascript
// Worker logs go to console automatically
// Look for "✅ WASM Worker ready"
```

---

## 📞 Still Need Help?

### Checklist before asking:

- [ ] Used HTTP server (not file://)
- [ ] All files present
- [ ] Built with `./build_phase1.sh`
- [ ] Browser console checked
- [ ] Correct browser version
- [ ] Tried different browser
- [ ] Cleared cache
- [ ] Read this entire document

### Provide when asking:

1. **Browser & version**
2. **Console errors** (screenshot)
3. **File list** (`ls -la`)
4. **URL used** (http:// or file://)
5. **Build output**
6. **OS & Go version**

---

## ✅ Success Criteria

Everything working when you see:

**In Console:**
```
✅ Cache initialized: {cachedFiles: 0}
✅ Web Worker ready
🔍 pure-dupes WASM initialized
```

**In UI:**
- No error messages
- Can select files
- Progress bar works
- Results display
- Cache stats update

**Perfect! You're ready to use Phase 1! 🎉**
