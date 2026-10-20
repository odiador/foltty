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

USO:
  foltty                     Abre la TUI interactiva
  foltty add [ruta]          Agrega la carpeta actual o una ruta a la lista
  foltty rm <ruta|nombre>    Elimina un proyecto de la lista
  foltty list [--json]       Muestra todos los proyectos registrados
  foltty open <nombre>       Abre directamente un proyecto por nombre o búsqueda
  foltty file                Muestra la ruta del archivo de proyectos
  foltty --help              Muestra esta ayuda
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

	if len(projects) == 0 {
		fmt.Println("👻 No hay proyectos registrados en Foltty.")
		fmt.Println("👉 Agrega uno con: foltty add")
		return nil
	}

	storagePath, err := project.GetStoragePath()
	if err != nil {
		return err
	}

	var paths []string
	for _, p := range projects {
		paths = append(paths, p.Path)
	}

	fzfArgs := []string{
		"--prompt=👻 Foltty > ",
		"--height=70%",
		"--layout=reverse",
		"--border=rounded",
		"--border-label= 📂 FOLTTY - GHOSTTY SELECTOR ",
		"--color=border:#8b5cf6,label:#ec4899,prompt:#38bdf8,pointer:#f43f5e",
		"--header=[ENTER] Abrir en Ghostty  |  [CTRL-X] Eliminar",
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
		// Exit code 130 is user canceling (Esc / Ctrl-C)
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
		fmt.Println("👉 Agrega uno con: foltty add")
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
		normalized, _ := project.NormalizePath(target)
		fmt.Printf("✅ Proyecto agregado: %s\n", normalized)

	case "rm", "remove", "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Especifica la ruta o nombre del proyecto a eliminar.")
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
			fmt.Printf("  %s \x1b[1m%s\x1b[0m -> \x1b[90m%s\x1b[0m\n", mark, p.Name, p.Path)
		}

	case "open":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Especifica el nombre o ruta del proyecto a abrir.")
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
