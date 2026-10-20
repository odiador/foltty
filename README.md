# 👻 Foltty

> **El gestor de proyectos ultra-rápido para Ghostty.**  
> Disponible como **Extensión Nativa para Raycast** y como **TUI / CLI Nativa en Go**.

---

## 💡 Arquitectura Híbrida: Raycast (React/TSX) + Go (TUI/CLI)

### ¿Es necesario que el plugin de Raycast sea en Node.js con TSX?
**Sí, para la capa de presentación de Raycast:**
- Raycast tiene un runtime interno basado en Node.js y React. Toda la UI de Raycast (`List`, `ActionPanel`, formularios, navegación y renderizado flotante de macOS) requiere obligatoriamente TypeScript/TSX y `@raycast/api`.

**Pero la lógica y la TUI pueden ser 100% en Go:**
- **Go para la TUI/CLI:** Arranque instantáneo (<5ms), binario nativo compilado, cero dependencias en tiempo de ejecución.
- **Raycast para la UI de macOS:** Maneja la ventana flotante y los atajos del sistema.
- **Single Source of Truth:** Ambos mundos leen y escriben en el mismo archivo (`~/.foltty_projects`), asegurando sincronización instantánea y sin fricción.

---

## 📂 Estructura del Proyecto

```tree
foltty/
├── cmd/
│   └── foltty/             # Entrypoint CLI/TUI en Go
│       └── main.go
├── internal/
│   └── project/            # Lógica de proyectos en Go (List, Add, Rm, Ghostty)
│       └── project.go
├── go.mod                  # Módulo de Go (github.com/odiador/foltty)
├── bin/
│   └── foltty.js          # Launcher JS/TSX para pnpm link
├── src/
│   ├── core/
│   │   └── projects.ts    # Lógica de persistencia en TypeScript
│   ├── cli.ts             # CLI interactivo en Node
│   └── index.tsx          # Componente React para Raycast
├── assets/
│   └── icon.png           # Icono oficial (512x512)
├── Makefile               # Automatización para Go y Raycast
└── package.json           # Manifiesto de Raycast
```

---

## ⚡ Comandos Rápidos

### 1. Extensión de Raycast
```bash
make dev
```
Inicia la extensión en modo desarrollo en Raycast con Hot Reload.

### 2. TUI en Go (Compilación y Ejecución Nativa)
```bash
make tui-go
```
O para instalar el binario `foltty` compilado en tu sistema con Go:
```bash
make install-go
```

### 3. Comandos de Terminal
```bash
foltty              # Abre la TUI interactiva con theming de Ghostty
foltty add          # Registra el directorio actual en Foltty
foltty add ~/repo   # Registra una ruta específica
foltty list         # Lista todos los proyectos
foltty list --json  # Salida en JSON para integraciones
foltty rm <nombre>  # Remueve un proyecto
foltty open <query> # Abre directamente por coincidencia de nombre
```

---

## 🛠️ Makefile

Ejecuta `make help` para ver todas las tareas disponibles:

| Comando | Descripción |
| :--- | :--- |
| `make dev` | Inicia el modo de desarrollo de Raycast con Hot Reload |
| `make build` | Compila la extensión de Raycast para producción (`dist/`) |
| `make build-go` | Compila el binario nativo en Go (`bin/foltty-go`) |
| `make install-go` | Instala el binario `foltty` compilado en `~/go/bin` |
| `make tui-go` | Ejecuta la TUI nativa en Go |
| `make lint` | Ejecuta las verificaciones de linter de Raycast |
| `make fix` | Aplica correcciones automáticas de linter |
