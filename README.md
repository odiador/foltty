# 👻 Foltty

> **`folder` + `tty` = `foltty`** ⚡  
> El selector ultra-rápido de proyectos para **Ghostty**, diseñado como un monorepo con **Extensión Nativa de Raycast** y **TUI / CLI en Go**.

---

## 💡 El Concepto: ¿Por qué Foltty?

El nombre nace de la combinación de:
- **`folder`** (tus carpetas y proyectos de desarrollo)
- **`tty`** (tu terminal, optimizada especialmente para **Ghostty**)

La misión de **Foltty** es eliminar toda la fricción entre encontrar un proyecto y empezar a trabajar en él:
1. **Fuera de la terminal (Navegador, Slack, Figma):** Invocas Raycast (`foltty`), filtras tu proyecto y con un `Enter` se abre una ventana limpia de Ghostty posicionada en la raíz del proyecto.
2. **Dentro de la terminal (Ghostty, tmux, zsh):** No necesitas cambiar de aplicación ni usar el ratón; ejecutas `foltty` (o `Ctrl+P`) y navegas por la TUI interactiva con previsualización en tiempo real.
3. **Estado 100% Sincronizado:** Ambos mundos leen y escriben en el mismo archivo (`~/.foltty_projects`), por lo que cualquier proyecto o preset agregado se refleja al instante en ambos entornos.

---

## 🏗️ Estructura del Monorepo

```tree
foltty/
├── raycast/                 # 🟣 App 1: Extensión Oficial de Raycast (React + TypeScript)
│   ├── assets/
│   │   └── icon.png         # Icono oficial (512x512)
│   ├── src/
│   │   ├── core/
│   │   │   └── projects.ts  # Motor de persistencia y llamada a Ghostty
│   │   └── index.tsx        # Interfaz de Raycast con ActionPanels y submenús
│   ├── package.json         # Dependencias (@raycast/api, React 19)
│   └── tsconfig.json
│
├── tui/                     # 🔵 App 2: CLI y TUI Nativa (Go)
│   ├── cmd/
│   │   └── foltty/
│   │       └── main.go      # Entrypoint de Go (TUI interactiva + CLI)
│   ├── internal/
│   │   └── project/
│   │       └── project.go   # Lógica compartida de rutas, presets y Ghostty
│   └── go.mod
│
├── Makefile                 # 🛠️ Orquestador raíz único para ambos proyectos
└── README.md                # Documentación central
```

---

## 🌟 Sistema de Presets y Alias

Foltty incluye atajos preconfigurados para las carpetas más comunes de macOS y te permite definir tus propios alias personalizados:

### Presets Nativos Incluidos:
- **`~` / `home`** ➔ Directorio personal (`$HOME`)
- **`docs` / `documents`** ➔ `~/Documents`
- **`downloads`** ➔ `~/Downloads`
- **`desktop`** ➔ `~/Desktop`
- **`dev`** ➔ Carpeta de desarrollo detectada (`~/Development`, `~/Developer`, `~/Projects`)
- **`config`** ➔ `~/.config`

---

## 🚀 Guía de Uso

### 1. Extensión para Raycast (GUI)

Inicia la extensión en modo desarrollo con Hot Reload:
```bash
make dev
```
1. Abre Raycast (`⌥ + Espacio` o `⌘ + Espacio`).
2. Escribe **`Foltty`** o **`Buscar Proyectos`** y presiona `Enter`.
3. **Atajos clave en Raycast:**
   - `Enter`: Abre la carpeta seleccionada en una nueva ventana de Ghostty.
   - `⌘ + Enter`: Muestra el proyecto en Finder.
   - `⌘ + ⇧ + A`: Abre el submenú de **Presets Rápidos** (Home, Docs, Dev, etc.).
   - `⌘ + N`: Formulario para agregar una carpeta manual o desde preset con dropdown.
   - `⌘ + ⌫`: Remueve el proyecto de la lista.
   - Menú de acciones (`⌘ + K`) ➔ *"Guardar como Preset Custom..."*: Asigna un alias a cualquier carpeta.

---

### 2. TUI y CLI en la Terminal (Go)

Compila y lanza la TUI interactiva nativa escrita en Go:
```bash
make tui
```

O instala el binario `foltty` en tu sistema (`~/go/bin`):
```bash
make install-tui
```

Una vez instalado en tu `PATH`:
```bash
# Navegación interactiva con previsualización
foltty

# Gestión de rutas y proyectos
foltty add           # Registra el directorio actual
foltty add docs      # Agrega ~/Documents mediante preset
foltty add ~         # Agrega tu Home
foltty list          # Lista los proyectos registrados
foltty list --json   # Salida en JSON para scripts
foltty rm docs       # Remueve un proyecto o preset
foltty open docs     # Abre directamente en Ghostty

# Gestión de presets personalizados (alias)
foltty preset                                   # Lista todos los presets disponibles
foltty preset add work ~/Development/mi-empresa  # Crea un preset personalizado
foltty add work                                 # Usa tu preset
foltty preset rm work                           # Elimina el preset
```

---

## 🛠️ Comandos del Makefile Central

El proyecto utiliza un **único Makefile orquestador** en la raíz para evitar la sobrecarga de entrar y salir de carpetas:

| Comando | Acción |
| :--- | :--- |
| `make help` | Muestra el menú de ayuda con todos los comandos disponibles |
| `make build` | Compila ambos proyectos (`bin/foltty` y `raycast/dist/`) |
| `make dev` | Inicia el modo de desarrollo de Raycast con Hot Reload |
| `make tui` | Compila y ejecuta la TUI nativa en Go |
| `make install-tui` | Instala el binario `foltty` de Go en tu sistema (`~/go/bin`) |
| `make preset` | Muestra todos los presets nativos y custom |
| `make list` | Lista los proyectos registrados en la terminal |
| `make add` | Agrega el directorio actual a Foltty |
| `make lint-raycast`| Ejecuta el linter oficial de Raycast |
| `make clean` | Limpia los binarios y carpetas de build temporales |
