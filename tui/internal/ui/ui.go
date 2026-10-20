package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/odiador/foltty/internal/config"
	"github.com/odiador/foltty/internal/project"
)

func HasFzf() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

func RunFzfTui() error {
	projects, err := project.GetProjects()
	if err != nil {
		return err
	}

	storagePath, err := project.GetStoragePath()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("👻 No projects registered in Foltty.")
		fmt.Println("👉 Add one with: foltty add, or use a preset: foltty add docs")
		return nil
	}

	var paths []string
	for _, p := range projects {
		paths = append(paths, p.Path)
	}

	preferredTerm := config.GetPreferredTerminal()

	fzfArgs := []string{
		"--prompt=👻 Foltty > ",
		"--height=75%",
		"--layout=reverse",
		"--border=rounded",
		fmt.Sprintf("--border-label= 📂 FOLTTY (%s) ", strings.ToUpper(string(preferredTerm))),
		"--color=border:#8b5cf6,label:#ec4899,prompt:#38bdf8,pointer:#f43f5e",
		"--header=[ENTER] Open Terminal  |  [CTRL-X] Remove  |  [ESC] Exit",
		"--header-first",
		fmt.Sprintf("--bind=ctrl-x:execute(foltty rm '{}')+reload(cat %s)", storagePath),
		"--preview=echo '📁 Path: {}'; echo ''; if [ -d '{}' ]; then ls -la --color=always '{}' | head -n 15; else echo '❌ Path not found on disk'; fi",
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
		fmt.Printf("🚀 Opening in %s: %s\n", preferredTerm, selected)
		return project.OpenInTerminal(selected, preferredTerm)
	}

	return nil
}

func RunNativeTui() error {
	projects, err := project.GetProjects()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("👻 No projects registered in Foltty.")
		fmt.Println("👉 Add one with: foltty add [path|preset]")
		return nil
	}

	preferredTerm := config.GetPreferredTerminal()
	fmt.Printf("\x1b[35m👻 FOLTTY - Project Selector (Terminal: %s)\x1b[0m\n", preferredTerm)
	for i, p := range projects {
		status := ""
		if !p.Exists {
			status = " \x1b[31m[Missing]\x1b[0m"
		}
		fmt.Printf("[%d] %s (%s)%s\n", i+1, p.Name, p.Path, status)
	}

	fmt.Print("\nSelect a number to launch (or press Enter to exit): ")
	var input string
	_, _ = fmt.Scanln(&input)

	var choice int
	_, err = fmt.Sscanf(input, "%d", &choice)
	if err == nil && choice >= 1 && choice <= len(projects) {
		selected := projects[choice-1]
		fmt.Printf("🚀 Opening in %s: %s\n", preferredTerm, selected.Path)
		return project.OpenInTerminal(selected.Path, preferredTerm)
	}

	return nil
}
