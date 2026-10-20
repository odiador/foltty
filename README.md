# 👻 Ghostty Project Manager (Raycast Extension)

Selector visual, rápido y moderno para gestionar tus proyectos locales y abrirlos instantáneamente en [Ghostty](https://ghostty.org/) desde Raycast.

Está diseñado para integrarse al 100% con tu script CLI `pman` y compartir el mismo archivo de registro (`~/.pman_projects`).

---

## 🚀 Características

- ⚡ **Apertura Directa e Instantánea:** Lanza una ventana limpia en Ghostty posicionada en la raíz del proyecto usando `open -na Ghostty --args --working-directory="<ruta>"`.
- 📁 **Sincronización con `pman`:** Lee y actualiza en tiempo real `~/.pman_projects`. Todo lo que agregues o borres aquí o desde la terminal (`padd` / `pman add`) se sincroniza automáticamente.
- ➕ **Agregar Proyectos desde Raycast (`⌘ + N`):** Selector visual con diálogo nativo de macOS (`Form.FilePicker`) o campo para escribir rutas manuales (soporta `~/`).
- 🗑️ **Remover Proyectos (`⌘ + ⌫`):** Elimina cualquier proyecto de la lista y actualiza el archivo persistente al instante.
- 🔍 **Búsqueda Inteligente:** Filtra proyectos al vuelo por nombre de carpeta o por ruta completa.
- ⚠️ **Detección de Rutas Inaccesibles:** Si una carpeta fue movida o borrada del disco, se marca con un distintivo visual rojo (`Ruta inaccesible`).
- 🛠️ **Acciones Rápidas:**
  - `Enter`: Abrir en Ghostty.
  - `⌘ + Enter`: Mostrar en Finder.
  - `⌘ + ⇧ + O`: Abrir con otra aplicación (Cursor, VS Code, Zed, etc.).
  - `⌘ + C`: Copiar ruta al portapapeles.
  - `⌘ + N`: Agregar nuevo proyecto a la lista.
  - `⌘ + R`: Recargar lista de proyectos.
  - `⌘ + ⌫`: Remover de la lista.

---

## 📦 Estructura del Proyecto

```tree
tty-raycast-pr/
├── assets/
│   ├── icon.png           # Icono principal de la extensión (512x512)
│   └── command.png        # Icono del comando (512x512)
├── src/
│   └── index.tsx          # Componente React principal (List + AddProjectForm)
├── package.json           # Manifiesto de Raycast, comandos y dependencias
├── tsconfig.json          # Configuración de TypeScript ES2022
├── .eslintrc.json         # Reglas oficiales de linting de Raycast
└── README.md
```

---

## 🛠️ Instalación y Modo de Desarrollo (pnpm)

1. **Instalar dependencias:**
   ```bash
   pnpm install
   ```

2. **Iniciar en modo desarrollo:**
   ```bash
   pnpm run dev
   ```
   *Raycast detectará automáticamente la extensión local y la compilará con Hot Reload.*

3. **Compilar para producción:**
   ```bash
   pnpm run build
   ```

4. **Validar formato y linter oficial:**
   ```bash
   pnpm run lint
   # o para corregir automáticamente:
   pnpm run fix-lint
   ```

---

## 💡 Mejoras implementadas respecto a la propuesta inicial

1. **Corrección de bugs críticos de ejecución:**
   - En el código inicial, `homedir` se importaba de `@raycast/api` (donde no existe), lo que provocaba que fallara al compilar. Ahora se importa de `os` nativo de Node.
   - Se eliminó `"fs": "^0.0.1-security"` del `package.json`, ya que `fs` es un módulo interno de Node.js.
   - Se añadieron los scripts obligatorios en `package.json` (`dev`, `build`, `lint`, `fix-lint`).
2. **Ejecución segura de procesos:**
   - Se reemplazó la concatenación insegura de strings en `exec(...)` por `execFileAsync("open", ["-na", "Ghostty", ...])`, evitando fallos ante rutas con espacios, comillas o caracteres especiales.
3. **Cierre limpio de Raycast:**
   - Se añadió `await closeMainWindow()` al abrir un proyecto, lo que oculta Raycast de inmediato y pasa el foco limpiamente a la nueva ventana de Ghostty.
4. **Formulario Integrado de Creación (`AddProjectForm`):**
   - Ahora puedes añadir carpetas directamente con la interfaz visual de Raycast sin depender únicamente de la terminal.
5. **Estado Vacío Amigable (`List.EmptyView`):**
   - Si no tienes proyectos registrados en `~/.pman_projects`, Raycast te muestra una guía interactiva con un botón directo para añadir tu primer proyecto.
