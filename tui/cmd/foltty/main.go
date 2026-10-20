package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/odiador/foltty/internal/project"
)

func printHelp() {
	fmt.Println(`
👻 Foltty (Go Core) - Project Manager para Ghostty & Raycast

USO BÁSICO:
  foltty                     Abre la TUI interactiva
  foltty add [ruta|preset]   Agrega una carpeta o preset a la lista (ej: foltty add docs)
  foltty rm <ruta|preset>    Elimina un proyecto de la lista (ej: foltty rm docs)
  foltty list [--json]       Muestra todos los proyectos registrados
  foltty open <nombre>       Abre directamente un proyecto por nombre o búsqueda

GESTIÓN DE PRESETS:
  foltty preset              Lista todos los presets disponibles (Home, Docs, Dev, Customs)
  foltty preset add <id> <p> Agrega un preset personalizado (ej: foltty preset add work ~/Work)
  foltty preset rm <id>      Elimina un preset personalizado

PRESETS POR DEFECTO:
  ~ / home                   Directorio personal ($HOME)
  docs / documents           ~/Documents
  downloads                  ~/Downloads
  desktop                    ~/Desktop
  dev                        ~/Development o ~/Projects
  config                     ~/.config
`)
}

func hasFzf() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

func runFzfTui() error {
	projects, err := project.GetProjects()
	if err != nil {
		return err
	}

	storagePath, err := project.GetStoragePath()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("👻 No hay proyectos registrados en Foltty.")
		fmt.Println("👉 Agrega uno con: foltty add, o usa un preset: foltty add docs")
		return nil
	}

	var paths []string
	for _, p := range projects {
		paths = append(paths, p.Path)
	}

	fzfArgs := []string{
		"--prompt=👻 Foltty > ",
		"--height=75%",
		"--layout=reverse",
		"--border=rounded",
		"--border-label= 📂 FOLTTY - GHOSTTY SELECTOR ",
		"--color=border:#8b5cf6,label:#ec4899,prompt:#38bdf8,pointer:#f43f5e",
		"--header=[ENTER] Abrir en Ghostty  |  [CTRL-X] Eliminar  |  [ESC] Salir",
		"--header-first",
		fmt.Sprintf("--bind=ctrl-x:execute(foltty rm '{}')+reload(cat %s)", storagePath),
		"--preview=echo '📁 Ruta: {}'; echo ''; if [ -d '{}' ]; then ls -la --color=always '{}' | head -n 15; else echo '❌ Ruta no encontrada'; fi",
		"--preview-window=right:50%:border-left",
	}

	cmd := exec.Command("fzf", fzfArgs...)
	cmd.Stdin = strings.NewReader(strings.Join(paths, "\n"))
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	selected := strings.TrimSpace(string(out))
	if selected != "" {
		fmt.Printf("🚀 Abriendo en Ghostty: %s\n", selected)
		return project.OpenInGhostty(selected)
	}

	return nil
}

func runSimpleTui() error {
	projects, err := project.GetProjects()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("👻 No hay proyectos registrados en Foltty.")
		fmt.Println("👉 Agrega uno con: foltty add [ruta|preset]")
		return nil
	}

	fmt.Println("\x1b[35m👻 FOLTTY - Selector de Proyectos\x1b[0m")
	for i, p := range projects {
		fmt.Printf("[%d] %s (%s)\n", i+1, p.Name, p.Path)
	}

	fmt.Print("\nSelecciona un número para abrir en Ghostty (o Enter para salir): ")
	var input string
	_, _ = fmt.Scanln(&input)

	var choice int
	_, err = fmt.Sscanf(input, "%d", &choice)
	if err == nil && choice >= 1 && choice <= len(projects) {
		selected := projects[choice-1]
		fmt.Printf("🚀 Abriendo en Ghostty: %s\n", selected.Path)
		return project.OpenInGhostty(selected.Path)
	}

	return nil
}

func handlePresetCommands(args []string) {
	if len(args) == 0 || args[0] == "list" || args[0] == "ls" {
		presets, err := project.GetAllPresets()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("👻 Presets disponibles (%d):\n\n", len(presets))
		for _, p := range presets {
			tag := "\x1b[34m[Built-in]\x1b[0m"
			if p.IsCustom {
				tag = "\x1b[32m[Custom]\x1b[0m"
			}
			fmt.Printf("  • \x1b[1m%-14s\x1b[0m %s -> \x1b[90m%s\x1b[0m\n", p.ID, tag, p.Path)
		}
		fmt.Println("\n👉 Para agregarlo a tu lista: foltty add <preset>")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "❌ Uso: foltty preset add <nombre> <ruta>")
			os.Exit(1)
		}
		name := args[1]
		targetPath := args[2]
		if err := project.AddCustomPreset(name, targetPath); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		norm, _ := project.NormalizePath(targetPath)
		fmt.Printf("✅ Preset guardado: %s -> %s\n", name, norm)

	case "rm", "remove", "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Uso: foltty preset rm <nombre>")
			os.Exit(1)
		}
		if err := project.RemoveCustomPreset(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🗑️ Preset removido: %s\n", args[1])

	default:
		// Si escribieron `foltty preset docs`, resolvemos y agregamos
		target := args[0]
		if err := project.AddProject(target); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		resolved, _ := project.ResolvePresetOrPath(target)
		fmt.Printf("✅ Proyecto agregado vía preset: %s\n", resolved)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "tui" {
		if hasFzf() {
			if err := runFzfTui(); err != nil {
				fmt.Fprintf(os.Stderr, "Error en TUI: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := runSimpleTui(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
		return
	}

	switch args[0] {
	case "preset", "presets":
		handlePresetCommands(args[1:])

	case "add":
		target := "."
		if len(args) > 1 {
			target = args[1]
		}
		if target == "." {
			cwd, err := os.Getwd()
			if err == nil {
				target = cwd
			}
		}
		if err := project.AddProject(target); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		resolved, _ := project.ResolvePresetOrPath(target)
		fmt.Printf("✅ Proyecto agregado: %s\n", resolved)

	case "rm", "remove", "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Especifica la ruta, nombre o preset a eliminar.")
			os.Exit(1)
		}
		if err := project.RemoveProject(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🗑️ Proyecto removido: %s\n", args[1])

	case "list", "ls":
		projects, err := project.GetProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(args) > 1 && args[1] == "--json" {
			data, _ := json.MarshalIndent(projects, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(projects) == 0 {
			fmt.Println("👻 No hay proyectos registrados.")
			return
		}

		fmt.Printf("👻 Proyectos registrados en Foltty (%d):\n\n", len(projects))
		for _, p := range projects {
			mark := "📂"
			if !p.Exists {
				mark = "⚠️ "
			}
			fmt.Printf("  %s \x1b[1m%-18s\x1b[0m -> \x1b[90m%s\x1b[0m\n", mark, p.Name, p.Path)
		}

	case "open":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Especifica el nombre, ruta o preset a abrir.")
			os.Exit(1)
		}
		query := strings.ToLower(args[1])
		projects, err := project.GetProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var found *project.Project
		for _, p := range projects {
			if strings.ToLower(p.Name) == query || strings.Contains(strings.ToLower(p.Path), query) {
				found = &p
				break
			}
		}

		if found == nil {
			// Probar si es un preset directo
			if resolved, err := project.ResolvePresetOrPath(args[1]); err == nil {
				fmt.Printf("🚀 Abriendo preset en Ghostty: %s\n", resolved)
				if err := project.OpenInGhostty(resolved); err != nil {
					fmt.Fprintf(os.Stderr, "Error al abrir Ghostty: %v\n", err)
					os.Exit(1)
				}
				return
			}
			fmt.Fprintf(os.Stderr, "❌ No se encontró ningún proyecto con el criterio: %s\n", args[1])
			os.Exit(1)
		}

		fmt.Printf("🚀 Abriendo en Ghostty: %s\n", found.Path)
		if err := project.OpenInGhostty(found.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error al abrir Ghostty: %v\n", err)
			os.Exit(1)
		}

	case "file":
		p, err := project.GetStoragePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(p)

	case "--help", "-h", "help":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "Comando desconocido: %s\n", args[0])
		printHelp()
		os.Exit(1)
	}
}
