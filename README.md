# 👻 Foltty (Monorepo)

> **El gestor de proyectos ultra-rápido para Ghostty.**  
> Arquitectura monorepo con **Extensión Nativa para Raycast** y **TUI / CLI Nativa en Go**, con soporte completo para **Presets y Alias Personalizados**.

---

## 🎯 Gestión de Rutas y Presets

Foltty permite agregar o eliminar rutas rápidamente usando nombres estándar o tus propios presets personalizados.

### 🌟 Presets Nativos Incluidos:
- `~` o `home` ➔ Tu directorio personal (`$HOME`)
- `docs` o `documents` ➔ `~/Documents`
- `downloads` ➔ `~/Downloads`
- `desktop` ➔ `~/Desktop`
- `dev` ➔ Carpeta de desarrollo detectada (`~/Development`, `~/Developer` o `~/Projects`)
- `config` ➔ `~/.config`

---

## ⚡ Cómo Usar Presets

### 1. Desde la Terminal (CLI en Go)

```bash
# Ver todos los presets disponibles (built-in + customs)
foltty preset

# Agregar un preset directamente a tu lista activa
foltty add docs
foltty add ~
foltty add dev

# Crear tu propio preset personalizado (alias)
foltty preset add work ~/Development/empresa/proyecto-principal
foltty preset add blog ~/Sites/mi-blog

# Usar tu preset personalizado
foltty add work

# Abrir directamente en Ghostty mediante un preset
foltty open docs

# Eliminar un proyecto o preset de tu lista
foltty rm docs
foltty preset rm work
```

### 2. Desde Raycast (Interfaz Gráfica)

- **Submenú de Presets (`⌘ + ⇧ + A`):** Dentro del menú de acciones de cualquier proyecto, entra en **"Agregar desde Preset..."** y selecciona `Documentos`, `Home`, `Descargas` o tus presets custom para añadirlos en 1 clic.
- **Formulario de Adición (`⌘ + N`):** Cuenta con un selector dropdown superior donde puedes elegir un preset o seleccionar una carpeta en Finder.
- **Guardar como Preset Custom:** Puedes convertir cualquier proyecto de tu lista en un preset personalizado con su propio alias para reutilizarlo en la terminal y en Raycast.
- **Vista Vacía:** Si la lista está vacía, Raycast te sugiere botones directos para agregar tus carpetas principales al instante.

---

## 🏗️ Estructura del Monorepo

```tree
foltty/
├── raycast/                 # 🟣 App 1: Extensión Oficial de Raycast (React + TS)
│   ├── assets/
│   │   └── icon.png         # Icono oficial (512x512)
│   ├── src/
│   │   ├── core/
│   │   │   └── projects.ts  # Soporte de proyectos y presets compartidos
│   │   └── index.tsx        # Interfaz de Raycast con Submenús de Presets
│   ├── package.json         # Dependencias (@raycast/api, React 19)
│   └── tsconfig.json
│
├── tui/                     # 🔵 App 2: CLI y TUI Nativa (Go)
│   ├── cmd/
│   │   └── foltty/
│   │       └── main.go      # Entrypoint de Go (CLI + Presets + TUI)
│   ├── internal/
│   │   └── project/
│   │       └── project.go   # Motor de rutas, presets y Ghostty
│   └── go.mod
│
├── Makefile                 # 🛠️ Orquestador central del monorepo
└── README.md
```

---

## 🛠️ Comandos del Makefile Central

| Comando | Descripción |
| :--- | :--- |
| `make help` | Muestra el menú de ayuda con todos los comandos |
| `make preset` | Muestra todos los presets disponibles (Home, Docs, Dev, Customs) |
| `make add` | Agrega el directorio actual a la lista de Foltty |
| `make list` | Lista los proyectos registrados en la terminal |
| `make build` | Compila ambos proyectos (`bin/foltty` y `raycast/dist/`) |
| `make dev` | Inicia el modo desarrollo de Raycast con Hot Reload |
| `make tui` | Compila y ejecuta la TUI nativa en Go |
| `make install-tui` | Instala el binario `foltty` de Go en tu sistema (`~/go/bin`) |
| `make clean` | Limpia los binarios y carpetas de build temporales |
