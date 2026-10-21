package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/odiador/foltty/internal/config"
	"github.com/odiador/foltty/internal/project"
	"github.com/odiador/foltty/internal/terminal"
	"github.com/odiador/foltty/internal/ui"
)

func printHelp() {
	fmt.Print(`
👻 Foltty (Go Core) - Multi-Terminal Project Manager

USAGE:
  foltty                     Launch interactive Spotlight-like TUI
  foltty add <path> [name]   Add folder or preset with optional display name
  foltty rename <target> <n> Rename a project's display name
  foltty rm <path|preset>    Remove project or preset from list
  foltty list [--json]       List all registered projects
  foltty open <query> [term] Open project in preferred or specified terminal
  foltty sort [mode]         View or set sort order (recent, name, added)

TERMINAL EMULATOR CONFIGURATION:
  foltty terminal            Show available and current default terminal
  foltty terminal set <id>   Set default terminal (ghostty, iterm, wezterm, warp, alacritty, kitty, terminal)

PRESETS MANAGEMENT:
  foltty preset              List all available presets (Home, Docs, Dev, Customs)
  foltty preset add <id> <p> Add custom preset alias (e.g. foltty preset add work ~/Work)
  foltty preset rm <id>      Remove custom preset

DEFAULT PRESETS:
  ~ / home                   User home directory ($HOME)
  docs / documents           ~/Documents
  downloads                  ~/Downloads
  desktop                    ~/Desktop
  dev                        ~/Development, ~/Developer, or ~/Projects
  config                     ~/.config
`)
}

func handleTerminalCommands(args []string) {
	if len(args) == 0 || args[0] == "list" || args[0] == "ls" {
		current := config.GetPreferredTerminal()
		terminals := terminal.GetAvailableTerminals()

		fmt.Println("👻 Terminal Emulators on this Mac:")
		fmt.Println()
		for _, t := range terminals {
			availTag := "\x1b[32m[Installed]\x1b[0m"
			if !t.Available {
				availTag = "\x1b[90m[Not Found]\x1b[0m"
			}
			isDefault := ""
			if t.ID == current {
				isDefault = " \x1b[36m(Default ✦)\x1b[0m"
			}
			fmt.Printf("  • \x1b[1m%-12s\x1b[0m %-16s %s%s\n", t.ID, t.Name, availTag, isDefault)
		}
		fmt.Printf("\n👉 Change default: foltty terminal set <id>\n")
		return
	}

	switch args[0] {
	case "set":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Usage: foltty terminal set <ghostty|iterm|wezterm|warp|alacritty|kitty|terminal>")
			os.Exit(1)
		}
		termID := strings.ToLower(args[1])
		if err := config.SetPreferredTerminal(termID); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Default terminal set to: %s\n", termID)

	case "get":
		fmt.Println(config.GetPreferredTerminal())

	default:
		fmt.Fprintf(os.Stderr, "Unknown terminal command: %s\n", args[0])
		os.Exit(1)
	}
}

func handlePresetCommands(args []string) {
	if len(args) == 0 || args[0] == "list" || args[0] == "ls" {
		presets, err := project.GetAllPresets()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("👻 Available Presets (%d):\n\n", len(presets))
		for _, p := range presets {
			tag := "\x1b[34m[Built-in]\x1b[0m"
			if p.IsCustom {
				tag = "\x1b[32m[Custom]\x1b[0m"
			}
			fmt.Printf("  • \x1b[1m%-14s\x1b[0m %s -> \x1b[90m%s\x1b[0m\n", p.ID, tag, p.Path)
		}
		fmt.Println("\n👉 Add preset to project list: foltty add <preset>")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "❌ Usage: foltty preset add <name> <path>")
			os.Exit(1)
		}
		name := args[1]
		targetPath := args[2]
		if err := project.AddCustomPreset(name, targetPath); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		norm, _ := project.NormalizePath(targetPath)
		fmt.Printf("✅ Preset saved: %s -> %s\n", name, norm)

	case "rm", "remove", "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Usage: foltty preset rm <name>")
			os.Exit(1)
		}
		if err := project.RemoveCustomPreset(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🗑️ Preset removed: %s\n", args[1])

	default:
		target := args[0]
		if err := project.AddProject(target); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		resolved, _ := project.ResolvePresetOrPath(target)
		fmt.Printf("✅ Project added via preset: %s\n", resolved)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "tui" {
		if err := ui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error in TUI: %v\n", err)
			os.Exit(1)
		}
		return
	}

	switch args[0] {
	case "terminal", "term":
		handleTerminalCommands(args[1:])

	case "preset", "presets":
		handlePresetCommands(args[1:])

	case "add":
		target := "."
		customName := ""
		if len(args) > 1 {
			target = args[1]
		}
		if len(args) > 2 {
			customName = strings.Join(args[2:], " ")
		}
		if target == "." {
			cwd, err := os.Getwd()
			if err == nil {
				target = cwd
			}
		}
		if err := project.AddProjectWithName(target, customName); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		resolved, _ := project.ResolvePresetOrPath(target)
		nameMsg := ""
		if customName != "" {
			nameMsg = fmt.Sprintf(" ('%s')", customName)
		}
		fmt.Printf("✅ Project added%s: %s\n", nameMsg, resolved)

	case "rename", "label":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "❌ Usage: foltty rename <path|name> <new display name>")
			os.Exit(1)
		}
		target := args[1]
		newName := strings.Join(args[2:], " ")
		if err := project.SetProjectName(target, newName); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✏️ Renamed '%s' to '%s'\n", target, newName)

	case "rm", "remove", "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Please specify project path, name, or preset to remove.")
			os.Exit(1)
		}
		if err := project.RemoveProject(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🗑️ Project removed: %s\n", args[1])

	case "list", "ls":
		projects, err := project.GetProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		sortBy := config.GetSortBy()
		projects = project.SortProjects(projects, sortBy)

		if len(args) > 1 && args[1] == "--json" {
			data, _ := json.MarshalIndent(projects, "", "  ")
			fmt.Println(string(data))
			return
		}

		if len(projects) == 0 {
			fmt.Println("👻 No registered projects.")
			return
		}

		fmt.Printf("👻 Registered Projects (%d) [Sort: %s]:\n\n", len(projects), strings.ToUpper(sortBy))
		for _, p := range projects {
			mark := "📂"
			if !p.Exists {
				mark = "⚠️ "
			}
			name := p.Name
			timeStr := ""
			if p.LastOpenedAt > 0 {
				rel := project.FormatRelativeTime(p.LastOpenedAt)
				if rel != "" {
					timeStr = fmt.Sprintf(" \x1b[90m(%s)\x1b[0m", rel)
				}
			}
			fmt.Printf("  %s \x1b[1m%-24s\x1b[0m%s -> \x1b[90m%s\x1b[0m\n", mark, name, timeStr, p.Path)
		}

	case "sort":
		if len(args) < 2 {
			fmt.Printf("✦ Current default sort: %s\nUsage: foltty sort <recent|name|added>\n", config.GetSortBy())
			return
		}
		if err := config.SetSortBy(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Default sort order set to: %s\n", config.GetSortBy())

	case "open":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "❌ Please specify project name, path, or preset to open.")
			os.Exit(1)
		}

		targetTerm := config.GetPreferredTerminal()
		if len(args) > 2 {
			targetTerm = terminal.TerminalID(args[2])
		}

		query := strings.ToLower(args[1])
		projects, err := project.GetProjects()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var targetPath string
		for _, p := range projects {
			if strings.ToLower(p.Name) == query || strings.ToLower(p.CustomName) == query || strings.Contains(strings.ToLower(p.Path), query) {
				targetPath = p.Path
				break
			}
		}

		if targetPath == "" {
			if resolved, err := project.ResolvePresetOrPath(args[1]); err == nil {
				targetPath = resolved
			}
		}

		if targetPath == "" {
			fmt.Fprintf(os.Stderr, "❌ No project or preset found matching: %s\n", args[1])
			os.Exit(1)
		}

		fmt.Printf("🚀 Opening in %s: %s\n", targetTerm, targetPath)
		if err := project.OpenInTerminal(targetPath, targetTerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening terminal: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
		printHelp()
		os.Exit(1)
	}
}
