# 👻 Foltty

> **`folder` + `tty` = `foltty`** ⚡  
> The lightning-fast project and folder launcher for macOS terminal emulators, architected as a monorepo with a **Native Raycast Extension** and a **Standalone Go TUI / CLI**.

[![CI](https://github.com/odiador/foltty/actions/workflows/ci.yml/badge.svg)](https://github.com/odiador/foltty/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-purple.svg)](LICENSE)

---

## ⚡ 1-Line Quick Install

Install the standalone Go CLI/TUI and build the Raycast extension with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash
```

**What the installer does:**
1. Detects your macOS architecture (Apple Silicon `arm64` or Intel `x86_64`).
2. Builds and installs the optimized `foltty` binary into `~/.local/bin/foltty`.
3. Ensures `~/.local/bin` is in your shell `$PATH` (`.zshrc` / `.bashrc`).
4. Bundles the Raycast extension ready to link into Raycast.

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

Launch the interactive TUI directly:
```bash
make tui
```

Install `foltty` globally into your system (`~/go/bin`):
```bash
make install-tui
```

**Common CLI Commands:**
```bash
foltty              # Interactive TUI with fuzzy search and preview
foltty add          # Register current directory in Foltty
foltty add docs     # Register ~/Documents via preset
foltty add ~        # Register Home directory
foltty list         # Print registered projects
foltty list --json  # Output in JSON format for automation
foltty rm docs      # Remove project or preset from list
foltty open docs    # Launch directly in your preferred terminal
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
├── install.sh               # 📦 1-line curl installer
├── Makefile                 # 🛠️ Central monorepo orchestrator
└── README.md                # Project documentation
```

---

## 🛠️ Central Makefile Commands

| Command | Action |
| :--- | :--- |
| `make help` | Display available targets |
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
