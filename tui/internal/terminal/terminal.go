package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TerminalID string

const (
	Ghostty TerminalID = "ghostty"
	ITerm   TerminalID = "iterm"
	WezTerm TerminalID = "wezterm"
	Alacritty TerminalID = "alacritty"
	Kitty   TerminalID = "kitty"
	Warp    TerminalID = "warp"
	AppleTerminal TerminalID = "terminal"
)

type TerminalInfo struct {
	ID        TerminalID
	Name      string
	AppPath   string
	Available bool
}

func GetAvailableTerminals() []TerminalInfo {
	home, _ := os.UserHomeDir()

	candidates := []TerminalInfo{
		{ID: Ghostty, Name: "Ghostty", AppPath: "/Applications/Ghostty.app"},
		{ID: ITerm, Name: "iTerm2", AppPath: "/Applications/iTerm.app"},
		{ID: WezTerm, Name: "WezTerm", AppPath: "/Applications/WezTerm.app"},
		{ID: Warp, Name: "Warp", AppPath: "/Applications/Warp.app"},
		{ID: Alacritty, Name: "Alacritty", AppPath: "/Applications/Alacritty.app"},
		{ID: Kitty, Name: "Kitty", AppPath: "/Applications/kitty.app"},
		{ID: AppleTerminal, Name: "Apple Terminal", AppPath: "/System/Applications/Utilities/Terminal.app"},
	}

	var result []TerminalInfo
	for _, term := range candidates {
		_, errMain := os.Stat(term.AppPath)
		userPath := filepath.Join(home, "Applications", filepath.Base(term.AppPath))
		_, errUser := os.Stat(userPath)

		available := errMain == nil || errUser == nil
		if term.ID == AppleTerminal {
			available = true
		}

		term.Available = available
		result = append(result, term)
	}

	return result
}

func DetectDefaultTerminal() TerminalID {
	terminals := GetAvailableTerminals()
	for _, t := range terminals {
		if t.Available {
			return t.ID
		}
	}
	return AppleTerminal
}

func Launch(terminalID TerminalID, projectPath string) error {
	resolvedPath, err := filepath.Abs(projectPath)
	if err != nil {
		resolvedPath = projectPath
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		return fmt.Errorf("directory not found: %s", resolvedPath)
	}

	switch strings.ToLower(string(terminalID)) {
	case "ghostty":
		cmd := exec.Command("open", "-na", "Ghostty", "--args", fmt.Sprintf("--working-directory=%s", resolvedPath))
		return cmd.Run()

	case "iterm", "iterm2":
		cmd := exec.Command("open", "-a", "iTerm", resolvedPath)
		return cmd.Run()

	case "wezterm":
		cmd := exec.Command("open", "-na", "WezTerm", "--args", "start", "--cwd", resolvedPath)
		return cmd.Run()

	case "alacritty":
		cmd := exec.Command("open", "-na", "Alacritty", "--args", "--working-directory", resolvedPath)
		return cmd.Run()

	case "kitty":
		cmd := exec.Command("open", "-na", "kitty", "--args", "--directory", resolvedPath)
		return cmd.Run()

	case "warp":
		cmd := exec.Command("open", "-na", "Warp", resolvedPath)
		return cmd.Run()

	case "terminal", "apple":
		cmd := exec.Command("open", "-a", "Terminal", resolvedPath)
		return cmd.Run()

	default:
		// Fallback to auto-detected terminal
		def := DetectDefaultTerminal()
		if def != terminalID {
			return Launch(def, resolvedPath)
		}
		cmd := exec.Command("open", "-a", "Terminal", resolvedPath)
		return cmd.Run()
	}
}
