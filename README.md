<p align="center">
  <img src="assets/logo.svg" width="160" height="160" alt="Foltty Logo" />
</p>

# 📁 Foltty

> **`folder` + `tty` = `foltty`** ⚡  
> The lightning-fast project and folder launcher for macOS terminal emulators, architected as a monorepo with a **Native Raycast Extension** and a **Standalone Go TUI / CLI**.

[![CI](https://github.com/odiador/foltty/actions/workflows/ci.yml/badge.svg)](https://github.com/odiador/foltty/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-purple.svg)](LICENSE)

---

## ⚡ Selective Installation (Interactive or 1-Line)

Foltty lets you install **only what you need**. You are not forced to install both components if you only use one!

### 🎯 Interactive Installer (Prompts you for your choice):
```bash
curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash
```
```text
What would you like to install?
  1) Both: Go CLI/TUI + Raycast Extension (Recommended)
  2) Go CLI/TUI only (Standalone binary, zero Node.js needed)
  3) Raycast Extension only (Visual launcher for Raycast)
Enter choice [1-3]:
```

### 🚀 Direct Flags (Automated / Headless):
- **Go CLI / TUI Only** (No Node.js or npm required):
  ```bash
  curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash -s -- --cli
  ```
- **Raycast Extension Only** (No Go required):
  ```bash
  curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash -s -- --raycast
  ```
- **Both Components**:
  ```bash
  curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash -s -- --all
  ```

---

## 💡 The Concept: What is Foltty?

The name comes directly from:
- **`folder`** (your local codebases, workspaces, and project directories)
- **`tty`** (your terminal emulator, with native support for **Ghostty**, **iTerm2**, **WezTerm**, **Warp**, and more)

**Foltty** eliminates friction when switching projects:
1. **Outside the terminal (Browser, Slack, Figma):** Trigger Raycast (`foltty`), fuzzy-search your project, and hit `Enter` to launch a clean terminal window directly in that directory.
2. **Inside the terminal (Ghostty, tmux, zsh):** Keep your hands on the keyboard—run `foltty` (or `Ctrl+P`) to navigate through an interactive TUI with instant live preview.
3. **Single Source of Truth:** Both environments share the same persistent storage (`~/.foltty_projects` and `~/.foltty_presets.json`). Anything added or removed in Raycast syncs immediately to the CLI, and vice versa.

---

## 🖥️ Multi-Terminal Emulator Support

Foltty is not limited to a single terminal emulator. It auto-detects and supports the top terminal emulators on macOS:

| Terminal | Supported | Launch Method |
| :--- | :---: | :--- |
| **Ghostty** | ✅ | Native window with `--working-directory` |
| **iTerm2** | ✅ | Direct folder launch / AppleScript tab support |
| **WezTerm** | ✅ | `wezterm start --cwd` |
| **Warp** | ✅ | Native directory targeting |
| **Alacritty**| ✅ | `--working-directory` |
| **Kitty** | ✅ | `--directory` |
| **Apple Terminal** | ✅ | Default macOS Terminal fallback |

### Terminal Configuration:
- **In Raycast:** Open Raycast Settings ➔ Extensions ➔ **Foltty** ➔ Select your **Default Terminal Emulator** in the dropdown.
- **In Go CLI / TUI:**
  ```bash
  foltty terminal            # List detected terminal emulators
  foltty terminal set ghostty # Set default terminal
  foltty open my-project iterm # Open project in a specific terminal
  ```

---

## 🌟 Built-in & Custom Presets

Foltty recognizes standard system shortcuts and allows you to define custom aliases:

### Default Presets:
- **`~` / `home`** ➔ Home directory (`$HOME`)
- **`docs` / `documents`** ➔ `~/Documents`
- **`downloads`** ➔ `~/Downloads`
- **`desktop`** ➔ `~/Desktop`
- **`dev`** ➔ Auto-detected development folder (`~/Development`, `~/Developer`, `~/Projects`)
- **`config`** ➔ `~/.config`

### Managing Custom Presets:
```bash
# Add custom alias
foltty preset add work ~/Development/enterprise-project
foltty preset add blog ~/Sites/my-blog

# Use the preset
foltty add work
foltty open work

# Remove a preset
foltty preset rm work
```

---

## 🚀 Usage Guide

### 1. Raycast Extension (GUI)

Run in local development mode with live Hot Reload:
```bash
make dev
```
1. Open Raycast (`⌥ + Space` or `⌘ + Space`).
2. Search for **`Foltty`** or **`Search Projects`** and press `Enter`.
3. **Key Shortcuts in Raycast:**
   - `Enter`: Open selected project in your preferred terminal emulator.
   - `⌘ + T`: Open submenu to choose a specific terminal (Ghostty, iTerm, WezTerm, Warp, etc.).
   - `⌘ + Enter`: Reveal folder in Finder.
   - `⌘ + ⇧ + A`: Quick Presets submenu (Home, Docs, Dev, etc.).
   - `⌘ + N`: Add new project form with preset dropdown or Finder picker.
   - `⌘ + ⌫`: Remove project from list.
   - Action Menu (`⌘ + K`) ➔ *"Save as Custom Preset..."*: Assign a reusable alias to any project.

---

### 2. Standalone Terminal TUI & CLI (Go)

Launch the interactive Spotlight-style TUI:
```bash
make tui
# or simply:
foltty
```

#### 🎨 Spotlight-like Centered TUI & Vim Mode:
- **Instant Search:** Starts directly in search mode with live filtering.
- **Vim Command Mode (`:`):** Type `:` to enter command mode:
  - `:h` or `:help`: Open full interactive help manual.
  - `:a <path> [name]`: Add any folder path with an optional custom display name.
  - `:r <new-name>`: Rename the selected project's display label.
  - `:d` or `:rm`: Delete selected project.
  - `:term <id>`: Switch default terminal emulator (`ghostty`, `iterm`, `wezterm`, `warp`, etc.).
  - `:q` or `:quit`: Exit Foltty.
- **Shortcuts:**
  - `↵ Enter`: Launch selected project in terminal.
  - `Ctrl + A`: Interactive prompt to add a project and display name.
  - `Ctrl + R`: Rename selected project.
  - `Ctrl + D`: Remove selected project.
  - `Esc`: Clear search / Exit.

Install `foltty` globally into your system (`~/.local/bin`):
```bash
make install-tui
```

**Common CLI Commands:**
```bash
foltty                     # Interactive Spotlight-style TUI with Vim mode
foltty add <path> [name]   # Register directory with optional display name
foltty rename <target> <n> # Rename a project's display name
foltty add docs            # Register ~/Documents via preset
foltty add ~               # Register Home directory
foltty list                # Print registered projects
foltty list --json         # Output in JSON format for automation
foltty rm docs             # Remove project or preset from list
foltty open docs           # Launch directly in your preferred terminal
```

---

## 🏗️ Monorepo Architecture

```tree
foltty/
├── raycast/                 # 🟣 Raycast Extension (React + TypeScript)
│   ├── assets/
│   │   └── icon.png         # Extension icon (512x512)
│   ├── src/
│   │   ├── core/
│   │   │   ├── terminal.ts  # Multi-terminal detection & execution engine
│   │   │   └── projects.ts  # Shared persistent storage & presets
│   │   └── index.tsx        # Raycast UI (ActionPanels, Form, List)
│   ├── package.json         # Manifest, Preferences schema, ESLint 9 Flat Config
│   ├── package-lock.json    # Store-compliant dependency lockfile
│   └── tsconfig.json
│
├── tui/                     # 🔵 Go CLI / TUI (Zero-dependency native binary)
│   ├── cmd/
│   │   └── foltty/
│   │       └── main.go      # CLI command router & argument parser
│   ├── internal/
│   │   ├── config/          # User preferences (~/.foltty_config.json)
│   │   ├── terminal/        # Terminal detection & execution
│   │   ├── project/         # Projects & presets storage engine
│   │   └── ui/              # FZF runner with preview & Native TUI
│   └── go.mod
│
├── .github/                 # 🚀 GitHub Workflows & Templates
│   ├── workflows/
│   │   ├── ci.yml           # CI validation (npm ci, ray lint, cross-platform Go build)
│   │   └── release.yml      # Automatic GitHub Releases on v* tags
│   ├── ISSUE_TEMPLATE/      # Feature Request & Bug Report templates
│   └── pull_request_template.md
│
├── docs/
│   └── EXAMPLE_PR.md        # Reference guide for issues and PR submissions
├── install.sh               # 📦 Selective installer (CLI, Raycast, or Both)
├── Makefile                 # 🛠️ Central monorepo orchestrator
└── README.md                # Project documentation
```

---

## 🛠️ Central Makefile Commands

| Command | Action |
| :--- | :--- |
| `make help` | Display available targets |
| `make install` | Install both Go binary and Raycast extension |
| `make install-tui` | Install only the Go `foltty` binary into `~/.local/bin` |
| `make install-raycast` | Install and build only the Raycast extension |
| `make build` | Build both Go binary (`bin/foltty`) and Raycast extension (`raycast/dist/`) |
| `make dev` | Start Raycast extension development server with Hot Reload |
| `make tui` | Build and launch the native interactive terminal UI |
| `make term` | Inspect all detected terminal emulators on this Mac |
| `make preset` | List all available built-in and custom presets |
| `make list` | List registered projects via CLI |
| `make add` | Register current directory in Foltty |
| `make lint-raycast`| Run official Raycast validator and linter |
| `make clean` | Clean build artifacts and dist folders |

---

## 🤝 Contributing & Community

Contributions are welcome! Please check our:
- [Pull Request Template](.github/pull_request_template.md)
- [Example PR & Issue Walkthrough](docs/EXAMPLE_PR.md)
