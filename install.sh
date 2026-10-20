#!/usr/bin/env bash
set -e

# ==============================================================================
# 👻 Foltty Installer (folder + tty = foltty)
# Installs the Go CLI/TUI and sets up the Raycast extension.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash
# Or locally:
#   ./install.sh
# ==============================================================================

BOLD="\033[1m"
GREEN="\033[32m"
CYAN="\033[36m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

echo ""
echo -e "${CYAN}${BOLD}👻 Foltty Installer${RESET} (${YELLOW}folder + tty = foltty${RESET})"
echo -e "${BOLD}====================================================${RESET}"
echo ""

# 1. Verify macOS
OS="$(uname -s)"
if [ "$OS" != "Darwin" ]; then
    echo -e "${RED}❌ Foltty is designed exclusively for macOS (Ghostty, iTerm2, WezTerm, Terminal + Raycast).${RESET}"
    exit 1
fi

ARCH="$(uname -m)"
echo -e "🍏 System detected: ${GREEN}macOS ($ARCH)${RESET}"

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

REPO_DIR="$HOME/.foltty"

# 2. Check if inside local cloned repository or need to clone
if [ -f "./tui/go.mod" ] && [ -d "./raycast" ]; then
    SOURCE_DIR="$(pwd)"
    echo -e "📁 Using local source repository: ${CYAN}$SOURCE_DIR${RESET}"
else
    echo -e "📦 Fetching latest version from ${CYAN}https://github.com/odiador/foltty.git${RESET}..."
    if [ -d "$REPO_DIR/.git" ]; then
        cd "$REPO_DIR" && git pull --quiet
    else
        rm -rf "$REPO_DIR"
        git clone --quiet https://github.com/odiador/foltty.git "$REPO_DIR"
    fi
    SOURCE_DIR="$REPO_DIR"
fi

# 3. Build & Install Go CLI binary
echo ""
echo -e "${BOLD}🔨 Building Go binary (foltty)...${RESET}"
if command -v go >/dev/null 2>&1; then
    cd "$SOURCE_DIR/tui"
    go build -ldflags="-s -w" -o "$INSTALL_DIR/foltty" ./cmd/foltty
    chmod +x "$INSTALL_DIR/foltty"
    echo -e "✅ Binary successfully installed to: ${GREEN}$INSTALL_DIR/foltty${RESET}"
else
    echo -e "${YELLOW}⚠️  Go is not installed. If you have Homebrew, install it with: brew install go${RESET}"
fi

# 4. Ensure ~/.local/bin is in PATH
SHELL_CONFIG=""
if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -f "$HOME/.bash_profile" ]; then
    SHELL_CONFIG="$HOME/.bash_profile"
elif [ -f "$HOME/.bashrc" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    if [ -n "$SHELL_CONFIG" ]; then
        echo ""
        echo -e "⚙️  Adding ${CYAN}$INSTALL_DIR${RESET} to your ${CYAN}$SHELL_CONFIG${RESET}..."
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_CONFIG"
        export PATH="$HOME/.local/bin:$PATH"
    fi
fi

# 5. Build Raycast extension
echo ""
echo -e "${BOLD}🟣 Configuring Raycast extension...${RESET}"
if [ -d "/Applications/Raycast.app" ]; then
    if command -v pnpm >/dev/null 2>&1 || command -v npm >/dev/null 2>&1; then
        echo -e "📦 Building extension bundle..."
        cd "$SOURCE_DIR/raycast"
        if command -v pnpm >/dev/null 2>&1; then
            pnpm install --silent
            pnpm run build >/dev/null
        else
            npm ci --silent
            npm run build >/dev/null
        fi
        echo -e "✅ Extension built at: ${GREEN}$SOURCE_DIR/raycast/dist${RESET}"
        echo ""
        echo -e "${CYAN}💡 To link the extension in Raycast:${RESET}"
        echo -e "   1. Open Raycast and search for ${BOLD}'Manage Extensions'${RESET}"
        echo -e "   2. Press ${BOLD}⌘ + K${RESET} ➔ ${BOLD}'Add Development Extension'${RESET}"
        echo -e "   3. Select folder: ${CYAN}$SOURCE_DIR/raycast${RESET}"
    else
        echo -e "${YELLOW}⚠️  Node package manager not found. Install Node & pnpm to use the Raycast extension.${RESET}"
    fi
else
    echo -e "${YELLOW}⚠️  Raycast.app not found in /Applications.${RESET}"
fi

echo ""
echo -e "${GREEN}${BOLD}🎉 Foltty installation completed!${RESET}"
echo -e "----------------------------------------------------"
echo -e "🚀 Try running in your terminal:"
echo -e "   ${CYAN}foltty terminal${RESET}  ➔ Inspect & select terminal emulators"
echo -e "   ${CYAN}foltty preset${RESET}    ➔ View available presets"
echo -e "   ${CYAN}foltty add docs${RESET}  ➔ Add Documents folder"
echo -e "   ${CYAN}foltty${RESET}           ➔ Open interactive TUI"
echo ""
