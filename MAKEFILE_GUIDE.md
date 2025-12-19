# 🎯 Makefile Quick Reference

## ⚡ Super Simple (2 Commands)

```bash
make all     # Build everything
make serve   # Start server
```

Open: **http://localhost:8080**

---

## 🚀 Most Common Commands

| Command | What It Does |
|---------|-------------|
| `make` or `make help` | Show all commands |
| `make all` | Build everything (WASM, runtime, MCP, test files) |
| `make serve` | Start HTTP server |
| `make check` | Verify all files exist |
| `make clean` | Delete all built files |

---

## 🔨 Build Commands

| Command | What It Does |
|---------|-------------|
| `make build` | Build all components |
| `make wasm` | Build WASM module only |
| `make runtime` | Get wasm_exec.js |
| `make mcp` | Build MCP server |
| `make test-files` | Create sample data |

---

## 📋 Your Workflow

### First Time Setup
```bash
make all     # Build everything
make check   # Verify it worked
make serve   # Start server
```

### After Editing Code
```bash
make wasm    # Rebuild WASM only
# Refresh browser
```

### Clean Start
```bash
make clean   # Remove all built files
make all     # Rebuild from scratch
make serve   # Run
```

---

## 🎯 What Gets Built

After `make all`:
```
✅ main.wasm          (~2-4 MB)   - Your WASM module
✅ wasm_exec.js       (~13 KB)    - Go runtime
✅ index.html                     - UI file
✅ mcp-server                     - MCP binary
✅ test-files/                    - Sample data (5 files)
```

---

## 🔍 Detailed Commands

### `make all`
Builds everything and shows what to do next:
```bash
make all
# Output:
# ✅ Build complete!
# To run:
#   make serve
```

### `make serve`
Checks files exist, then starts server:
```bash
make serve
# Output:
# 🔍 Checking files...
# ✅ main.wasm (2.1M)
# ✅ wasm_exec.js
# ✅ wasm-worker.js
# ✅ index.html
# ✅ All checks passed!
#
# 🚀 Starting server on port 8080...
# Open in browser:
#   http://localhost:8080/index.html
```

### `make check`
Verifies all required files:
```bash
make check
# Shows which files exist and which are missing
```

### `make clean`
Removes all built files:
```bash
make clean
# Deletes: main.wasm, wasm_exec.js, index.html, mcp-server, test-files/
```

---

## 💡 Pro Tips

### Build and Run
```bash
make all && make serve
```

### Check Before Running
```bash
make check && make serve
```

### Quick Rebuild
```bash
# After editing main_wasm_enhanced.go:
make wasm
# Just rebuilds WASM, much faster!
```

### Different Port
```bash
# Edit Makefile, change:
PORT := 8080
# To:
PORT := 9000
```

### File Sizes
```bash
make sizes
# Shows size of all built files
```

---

## 🐛 Troubleshooting

### "Go not found"
```bash
# Install Go first
brew install go  # Mac
```

### "Make not found"
```bash
# Mac: Already installed
# Linux:
sudo apt-get install make  # Ubuntu/Debian
sudo dnf install make      # Fedora/RHEL
```

### "Permission denied"
```bash
# Make scripts executable
chmod +x *.sh
```

### Build fails
```bash
# Clean and retry
make clean
make all
```

---

## 🎨 Makefile vs Shell Scripts

| Task | Makefile | Shell Script |
|------|----------|-------------|
| Build all | `make all` | `./build.sh` |
| Just WASM | `make wasm` | Manual command |
| Check files | `make check` | `./check.sh` |
| Serve | `make serve` | `./serve.sh` |
| Clean | `make clean` | Manual deletion |

**Makefile Benefits:**
- ✅ Shorter commands
- ✅ Smart rebuilding (only if changed)
- ✅ Tab completion
- ✅ Industry standard

**Use whichever you prefer!** Both work great. 🎉

---

## 🚀 Quick Start Examples

### Complete Setup
```bash
# Clone/extract
cd pure-dupes-phase1-flat

# Build and run (2 commands!)
make all
make serve
```

### Development Workflow
```bash
# Edit main_wasm_enhanced.go

# Rebuild
make wasm

# Refresh browser
# Done!
```

### Clean Install
```bash
make clean     # Remove old files
make all       # Build fresh
make serve     # Run
```

---

## 📖 More Help

- Type `make` or `make help` to see all commands
- All commands show what they're doing
- Colors help identify status:
  - 🔵 Blue = Info
  - 🟢 Green = Success
  - 🟡 Yellow = Warning
  - 🔴 Red = Error

---

**TL;DR:**
```bash
make all      # Build
make serve    # Run
```

Done! 🎉
