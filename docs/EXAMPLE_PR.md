# 📝 Example Pull Request & Issue Guide

This guide illustrates a real-world example of opening an **Issue** and submitting a **Pull Request (PR)** to Foltty.

---

## 1. Example GitHub Issue

### **Title:** `[FEATURE] Add support for Warp terminal emulator`
**Labels:** `enhancement`, `good first issue`

#### **Description:**
```markdown
### 💡 Problem Description
Currently, Foltty launches Ghostty, iTerm2, and Terminal, but several users on our team use Warp (https://www.warp.dev/).
When opening a project from Raycast or the TUI, there is no option to launch directly inside Warp.

### 🚀 Proposed Solution
1. Add `Warp` to the terminal emulator registry in both `raycast/src/core/terminal.ts` and `tui/internal/terminal/terminal.go`.
2. Launch command on macOS:
   `open -na Warp "<project_path>"`
3. Expose Warp in the Raycast preferences dropdown and the `foltty terminal set warp` CLI command.

### 🖥️ Affected Packages
- [x] Raycast Extension
- [x] Go CLI / TUI
- [x] Documentation
```

---

## 2. Example Pull Request (PR)

### **Title:** `feat(terminal): add native support for Warp terminal emulator`
**Branch:** `feat/warp-terminal-support` ➔ `main`

#### **PR Body:**
```markdown
## 📌 Description
Resolves #12.

Adds first-class support for launching projects inside the **Warp** terminal emulator (`/Applications/Warp.app`) across both the Raycast extension and the Go TUI/CLI.

## ⚡ Type of Change
- [x] ✨ New feature (non-breaking change adding functionality)
- [x] 📝 Documentation update

## 🖥️ Affected Packages
- [x] `raycast/` (Raycast Extension)
- [x] `tui/` (Go CLI / TUI)
- [x] Root (Makefile, CI/CD, Documentation)

## 🔍 Key Changes Made
1. **Raycast Extension:**
   - Added `warp` to `TERMINAL_REGISTRY` in `raycast/src/core/terminal.ts`.
   - Added `Warp` option to `preferences` dropdown in `raycast/package.json`.
2. **Go CLI & TUI:**
   - Added `Warp` to `terminal.GetAvailableTerminals()` in `tui/internal/terminal/terminal.go`.
   - Supported `foltty terminal set warp` and `foltty open <name> warp`.
3. **Documentation:**
   - Updated `README.md` to list Warp among the supported terminal emulators.

## 🧪 How Has This Been Tested?
- [x] `make lint-raycast` passed with 0 errors and 0 warnings.
- [x] `make build` succeeded for both Go binary and Raycast extension.
- [x] Tested opening `~/Documents` inside Warp via Raycast (`⌘ + T` ➔ Warp).
- [x] Tested opening via CLI: `foltty open docs warp`.

## 📸 Screenshots
<img width="600" alt="Foltty Raycast Terminal Submenu" src="https://user-images.githubusercontent.com/.../warp-demo.png">

## ✅ Checklist
- [x] My code adheres to the style guidelines of this project
- [x] Verified zero deprecations with ESLint 9 Flat Config
- [x] Updated README and CLI help output
- [x] All GitHub Actions CI checks are green
```
