#!/usr/bin/env bash
set -e

# ==============================================================================
# 👻 Foltty Installer (folder + tty = foltty)
# Instala la CLI/TUI en Go y configura la extensión de Raycast.
#
# Uso:
#   curl -fsSL https://raw.githubusercontent.com/odiador/foltty/main/install.sh | bash
# O localmente:
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

# 1. Verificar macOS
OS="$(uname -s)"
if [ "$OS" != "Darwin" ]; then
    echo -e "${RED}❌ Foltty está diseñado exclusivamente para macOS (Ghostty + Raycast).${RESET}"
    exit 1
fi

ARCH="$(uname -m)"
echo -e "🍏 Sistema detectado: ${GREEN}macOS ($ARCH)${RESET}"

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

REPO_DIR="$HOME/.foltty"

# 2. Si estamos dentro del repositorio clonado, usamos la carpeta actual; si no, clonamos
if [ -f "./tui/go.mod" ] && [ -d "./raycast" ]; then
    SOURCE_DIR="$(pwd)"
    echo -e "📁 Usando código fuente local: ${CYAN}$SOURCE_DIR${RESET}"
else
    echo -e "📦 Obteniendo la última versión de ${CYAN}https://github.com/odiador/foltty.git${RESET}..."
    if [ -d "$REPO_DIR/.git" ]; then
        cd "$REPO_DIR" && git pull --quiet
    else
        rm -rf "$REPO_DIR"
        git clone --quiet https://github.com/odiador/foltty.git "$REPO_DIR"
    fi
    SOURCE_DIR="$REPO_DIR"
fi

# 3. Compilar e Instalar el binario CLI en Go
echo ""
echo -e "${BOLD}🔨 Compilando binario de Go (foltty)...${RESET}"
if command -v go >/dev/null 2>&1; then
    cd "$SOURCE_DIR/tui"
    go build -ldflags="-s -w" -o "$INSTALL_DIR/foltty" ./cmd/foltty
    chmod +x "$INSTALL_DIR/foltty"
    echo -e "✅ Binario instalado con éxito en: ${GREEN}$INSTALL_DIR/foltty${RESET}"
else
    echo -e "${YELLOW}⚠️  Go no está instalado. Si tienes Homebrew, puedes instalarlo con: brew install go${RESET}"
fi

# 4. Asegurar que ~/.local/bin esté en el PATH
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
        echo -e "⚙️  Agregando ${CYAN}$INSTALL_DIR${RESET} a tu ${CYAN}$SHELL_CONFIG${RESET}..."
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_CONFIG"
        export PATH="$HOME/.local/bin:$PATH"
    fi
fi

# 5. Compilar la Extensión de Raycast si Node/pnpm están disponibles
echo ""
echo -e "${BOLD}🟣 Configurando extensión de Raycast...${RESET}"
if [ -d "/Applications/Raycast.app" ]; then
    if command -v pnpm >/dev/null 2>&1; then
        echo -e "📦 Instalando dependencias y compilando extensión..."
        cd "$SOURCE_DIR/raycast"
        pnpm install --silent
        pnpm run build >/dev/null
        echo -e "✅ Extensión compilada en: ${GREEN}$SOURCE_DIR/raycast/dist${RESET}"
        echo ""
        echo -e "${CYAN}💡 Para vincularla en Raycast:${RESET}"
        echo -e "   1. Abre Raycast y busca ${BOLD}'Manage Extensions'${RESET}"
        echo -e "   2. Presiona ${BOLD}⌘ + K${RESET} ➔ ${BOLD}'Add Development Extension'${RESET}"
        echo -e "   3. Selecciona la carpeta: ${CYAN}$SOURCE_DIR/raycast${RESET}"
    else
        echo -e "${YELLOW}⚠️  pnpm no está instalado. Para compilar la extensión ejecuta: brew install pnpm && make dev${RESET}"
    fi
else
    echo -e "${YELLOW}⚠️  Raycast.app no encontrada en /Applications.${RESET}"
fi

echo ""
echo -e "${GREEN}${BOLD}🎉 ¡Instalación de Foltty completada con éxito!${RESET}"
echo -e "----------------------------------------------------"
echo -e "🚀 Prueba ejecutando en tu terminal:"
echo -e "   ${CYAN}foltty preset${RESET}   ➔ Ver presets disponibles"
echo -e "   ${CYAN}foltty add docs${RESET} ➔ Agregar tu carpeta Documentos"
echo -e "   ${CYAN}foltty${RESET}          ➔ Abrir la TUI interactiva para Ghostty"
echo ""
