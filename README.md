# 👻 Foltty (Monorepo)

> **El gestor de proyectos ultra-rápido para Ghostty.**  
> Arquitectura monorepo compuesta por la **Extensión Nativa para Raycast** y la **TUI / CLI Nativa en Go**.

---

## 🏗️ Estructura del Monorepo

```tree
foltty/
├── raycast/                 # 🟣 Extensión oficial para Raycast (React + TypeScript)
│   ├── assets/
│   │   └── icon.png         # Icono oficial (512x512)
│   ├── src/
│   │   ├── core/
│   │   │   └── projects.ts  # Lógica de lectura/escritura de proyectos
│   │   └── index.tsx        # Interfaz de lista y acciones en Raycast
│   ├── package.json         # Dependencias (@raycast/api, React 19)
│   └── tsconfig.json
├── tui/                     # 🔵 Aplicación CLI y TUI Nativa (Go)
│   ├── cmd/
│   │   └── foltty/
│   │       └── main.go      # Punto de entrada CLI/TUI
│   ├── internal/
│   │   └── project/
│   │       └── project.go   # Motor de proyectos y llamada a Ghostty
│   └── go.mod
├── Makefile                 # 🛠️ Orquestador central para ambos entornos
└── README.md
```

---

## ⚡ Estado Compartido (Single Source of Truth)

Ambas aplicaciones (`raycast/` y `tui/`) leen y escriben en el mismo archivo persistente:
- **`~/.foltty_projects`** (con sincronización y compatibilidad automática hacia `~/.pman_projects`).

Cualquier proyecto agregado o eliminado desde Raycast se refleja de inmediato en la terminal, y viceversa.

---

## 🚀 Inicio Rápido

### 1. Instalar dependencias
```bash
make install
```
*(Instala los paquetes de `raycast` vía pnpm y descarga los módulos de `tui` vía Go).*

### 2. Ejecutar la Extensión en Raycast
```bash
make dev
# o:
make dev-raycast
```
Abre Raycast y busca **"Buscar Proyectos"** o **"Foltty"**.

### 3. Ejecutar la TUI en la Terminal
```bash
make tui
```
O compilar e instalar el binario `foltty` de Go en tu sistema (`~/go/bin`):
```bash
make install-go
```

Una vez instalado en tu `PATH`, puedes usarlo desde cualquier directorio:
```bash
foltty              # Abre la TUI interactiva (con preview)
foltty add          # Registra el directorio actual
foltty add ~/repo   # Registra una ruta específica
foltty list         # Lista los proyectos registrados
foltty list --json  # Salida en JSON para scripts
foltty rm <nombre>  # Elimina un proyecto
foltty open <query> # Abre en Ghostty por búsqueda
```

---

## 🛠️ Comandos del Makefile Central

| Comando | Descripción |
| :--- | :--- |
| `make help` | Muestra el menú de ayuda con todos los comandos |
| `make install` | Instala dependencias de Raycast (pnpm) y Go |
| `make build` | Compila ambos proyectos (`bin/foltty` y `raycast/dist/`) |
| `make dev` | Inicia el modo desarrollo de Raycast con Hot Reload |
| `make tui` | Compila y ejecuta la TUI nativa en Go |
| `make install-tui` | Instala el binario `foltty` en tu sistema (`~/go/bin`) |
| `make list` | Lista los proyectos registrados en la terminal |
| `make add` | Agrega la carpeta actual a la lista |
| `make lint-raycast`| Ejecuta el linter oficial de Raycast |
| `make clean` | Limpia los binarios y artefactos generados |
