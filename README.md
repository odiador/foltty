# 👻 Foltty

> **El gestor de proyectos ultra-rápido para Ghostty.**  
> Disponible como **Extensión Nativa para Raycast** y como **TUI Interactiva para la Terminal**.

---

## 💡 ¿Por qué Raycast + TUI? (Dual-Mode)

1. **Cuando estás en el sistema (navegador, Slack, editor, etc.):**
   * Presionas el atajo global de **Raycast** (`foltty`), buscas tu proyecto y con un solo `Enter` se abre una ventana limpia de **Ghostty** en la ruta exacta.
2. **Cuando YA estás dentro de Ghostty o tmux:**
   * No necesitas romper tu flujo de teclado ni cambiar de aplicación: ejecutas `foltty` (o atajo `Ctrl+P`) y se abre la **TUI interactiva** ultra-rápida.
3. **Estado 100% compartido (Single Source of Truth):**
   * Ambos modos leen y escriben en el mismo archivo persistente (`~/.foltty_projects` y compatible con `~/.pman_projects`). Cualquier proyecto que agregues o borres en uno se refleja inmediatamente en el otro.

---

## ⚡ Modos de Uso

### 1. Extensión para Raycast
Inicia la extensión en modo desarrollo:
```bash
make dev
# o con pnpm:
pnpm run dev
```
- Busca **"Buscar Proyectos"** o **"Foltty"** en Raycast.
- `Enter`: Abre en Ghostty.
- `⌘ + Enter`: Abre en Finder.
- `⌘ + N`: Abre el formulario visual para agregar una carpeta (`Form.FilePicker`).
- `⌘ + ⌫`: Remueve el proyecto de la lista.

### 2. TUI en Terminal (CLI)
Ejecuta la TUI interactiva directamente:
```bash
make tui
# o directamente:
./bin/foltty.js
```

Para disponer del comando global `foltty` en cualquier directorio de tu Mac:
```bash
make link
# o:
pnpm link --global
```

Una vez enlazado, puedes usar:
```bash
foltty              # Abre la TUI interactiva (con preview de FZF o selector nativo)
foltty add          # Agrega el directorio actual a la lista
foltty add ~/repo   # Agrega una ruta específica
foltty list         # Lista todos los proyectos registrados
foltty rm <nombre>  # Remueve un proyecto
foltty open <query> # Abre directamente el proyecto que coincida
```

---

## 📂 Arquitectura del Código

```tree
foltty/
├── bin/
│   └── foltty.js          # Ejecutable CLI / TUI para la terminal
├── src/
│   ├── core/
│   │   └── projects.ts    # Lógica unificada de almacenamiento, rutas y Ghostty
│   ├── cli.ts             # TUI interactiva (FZF theming + fallback nativo en Node)
│   └── index.tsx          # Componente React para la extensión de Raycast
├── assets/
│   ├── icon.png           # Icono principal de Raycast (512x512)
│   └── command.png        # Icono del comando
├── Makefile               # Automatización de tareas
├── package.json           # Manifiesto unificado de Raycast y binario CLI
└── tsconfig.json
```

---

## 🛠️ Comandos del Makefile

| Comando | Descripción |
| :--- | :--- |
| `make tui` | Abre la TUI interactiva en tu terminal |
| `make dev` | Inicia el modo desarrollo de Raycast con Hot Reload |
| `make link` | Registra el binario `foltty` de forma global en tu Mac |
| `make add` | Agrega la carpeta actual a la lista de Foltty |
| `make list` | Muestra la lista de proyectos registrados |
| `make build` | Compila la extensión para producción (`dist/`) |
| `make lint` | Ejecuta las verificaciones de linter de Raycast |
| `make fix` | Aplica correcciones automáticas de estilo y linter |
