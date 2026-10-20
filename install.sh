#!/usr/bin/env bash
set -e

# ==============================================================================
# 👻 Foltty Installer (folder + tty = foltty)
# Flexible installer for Raycast Extension and Go CLI/TUI.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash
# Non-interactive flags:
#   curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash -s -- --cli
#   curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash -s -- --raycast
#   curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash -s -- --all
# ==============================================================================

BOLD="\033[1m"
GREEN="\033[32m"
CYAN="\033[36m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

INSTALL_MODE=""

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --all|-a)
            INSTALL_MODE="all"
            shift
            ;;
        --cli|--tui|--go|-c)
            INSTALL_MODE="cli"
            shift
            ;;
        --raycast|-r)
            INSTALL_MODE="raycast"
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [--all | --cli | --raycast]"
            echo "  --all      Install both Go CLI/TUI and Raycast extension (default)"
            echo "  --cli      Install only the standalone Go CLI/TUI binary"
            echo "  --raycast  Install only the Raycast extension"
            exit 0
            ;;
        *)
            shift
            ;;
    esac
done

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

# 2. Interactive Selection Menu if no flag provided
if [ -z "$INSTALL_MODE" ]; then
    if [ -t 0 ] || [ -e /dev/tty ]; then
        echo ""
        echo -e "${BOLD}What would you like to install?${RESET}"
        echo -e "  ${CYAN}1)${RESET} Both: Go CLI/TUI + Raycast Extension ${GREEN}(Recommended)${RESET}"
        echo -e "  ${CYAN}2)${RESET} Go CLI/TUI only ${YELLOW}(Standalone binary, zero Node.js needed)${RESET}"
        echo -e "  ${CYAN}3)${RESET} Raycast Extension only ${YELLOW}(Visual launcher for Raycast)${RESET}"
        echo ""
        printf "Enter choice [1-3] (Default: 1): "
        
        # Read from /dev/tty to support 'curl | bash'
        if [ -e /dev/tty ]; then
            read -r USER_CHOICE < /dev/tty || USER_CHOICE="1"
        else
            read -r USER_CHOICE || USER_CHOICE="1"
        fi

        case "$USER_CHOICE" in
            2) INSTALL_MODE="cli" ;;
            3) INSTALL_MODE="raycast" ;;
            *) INSTALL_MODE="all" ;;
        esac
    else
        INSTALL_MODE="all"
    fi
fi

echo -e "📦 Selected installation target: ${BOLD}${CYAN}$INSTALL_MODE${RESET}"

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

REPO_DIR="$HOME/.foltty"

# 3. Locate or clone source repository
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

# 4. Install Go CLI / TUI if requested
if [ "$INSTALL_MODE" = "all" ] || [ "$INSTALL_MODE" = "cli" ]; then
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

    # Ensure ~/.local/bin is in PATH
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
fi

# 5. Install Raycast extension if requested
if [ "$INSTALL_MODE" = "all" ] || [ "$INSTALL_MODE" = "raycast" ]; then
    echo ""
    echo -e "${BOLD}🟣 Configuring Raycast extension...${RESET}"
    if [ -d "/Applications/Raycast.app" ]; then
        if command -v npm >/dev/null 2>&1; then
            echo -e "📦 Building extension bundle..."
            cd "$SOURCE_DIR/raycast"
            npm ci --silent
            npm run build >/dev/null
            echo -e "✅ Extension built at: ${GREEN}$SOURCE_DIR/raycast/dist${RESET}"
            echo ""
            echo -e "${CYAN}💡 To link the extension in Raycast:${RESET}"
            echo -e "   1. Open Raycast and search for ${BOLD}'Manage Extensions'${RESET}"
            echo -e "   2. Press ${BOLD}⌘ + K${RESET} ➔ ${BOLD}'Add Development Extension'${RESET}"
            echo -e "   3. Select folder: ${CYAN}$SOURCE_DIR/raycast${RESET}"
        else
            echo -e "${YELLOW}⚠️  npm is not installed. Install Node.js to use the Raycast extension.${RESET}"
        fi
    else
        echo -e "${YELLOW}⚠️  Raycast.app not found in /Applications.${RESET}"
    fi
fi

echo ""
echo -e "${GREEN}${BOLD}🎉 Foltty installation completed!${RESET}"
echo -e "----------------------------------------------------"
if [ "$INSTALL_MODE" = "all" ] || [ "$INSTALL_MODE" = "cli" ]; then
    echo -e "🚀 Terminal CLI commands:"
    echo -e "   ${CYAN}foltty terminal${RESET}  ➔ Inspect & select terminal emulators"
    echo -e "   ${CYAN}foltty preset${RESET}    ➔ View available presets"
    echo -e "   ${CYAN}foltty add docs${RESET}  ➔ Add Documents folder"
    echo -e "   ${CYAN}foltty${RESET}           ➔ Open interactive TUI"
fi
if [ "$INSTALL_MODE" = "all" ] || [ "$INSTALL_MODE" = "raycast" ]; then
    echo -e "🟣 Raycast: Search for 'Foltty' or 'Search Projects' in Raycast"
fi
echo ""
