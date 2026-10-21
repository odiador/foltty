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

func LaunchWithGhosttyMode(terminalID TerminalID, projectPath string, ghosttyMode string, title string) error {
	resolvedPath, err := filepath.Abs(projectPath)
	if err != nil {
		resolvedPath = projectPath
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		return fmt.Errorf("directory not found: %s", resolvedPath)
	}

	switch strings.ToLower(string(terminalID)) {
	case "ghostty":
		if ghosttyMode == "" {
			ghosttyMode = "smart"
		}
		script := `
on run argv
    set targetDir to item 1 of argv
    set openMode to item 2 of argv
    set targetTitle to ""
    if (count of argv) >= 3 then
        set targetTitle to item 3 of argv
    end if

    tell application "Ghostty"
        activate

        if openMode is "smart" or openMode is "" then
            repeat with w in windows
                repeat with t in tabs of w
                    set isMatch to false
                    set tTitle to ""
                    try
                        set tTitle to name of t
                    end try
                    set tDir to ""
                    try
                        set tDir to working directory of (focused terminal of t)
                    end try

                    if targetTitle is not "" and (tTitle is targetTitle or tTitle is ("🔔 " & targetTitle) or tTitle contains targetTitle) then
                        set isMatch to true
                    else if targetDir is not "" and (tDir is targetDir or tDir is (targetDir & "/") or tDir starts with (targetDir & "/")) then
                        set isMatch to true
                    end if

                    if isMatch then
                        try
                            activate window w
                        end try
                        select tab t
                        try
                            focus (focused terminal of t)
                        end try
                        return
                    end if
                end repeat
            end repeat
        end if

        set cfg to new surface configuration
        set initial working directory of cfg to targetDir
        if openMode is "window" then
            set myWin to new window with configuration cfg
            set myTab to selected tab of myWin
            set myTerm to focused terminal of myTab
        else if (count of windows) > 0 then
            set myTab to new tab in front window with configuration cfg
            set myTerm to focused terminal of myTab
        else
            set myWin to new window with configuration cfg
            set myTab to selected tab of myWin
            set myTerm to focused terminal of myTab
        end if
        if targetTitle is not "" then
            perform action ("set_tab_title:" & targetTitle) on myTerm
        end if
    end tell
end run`
		cmd := exec.Command("osascript", "-e", script, resolvedPath, ghosttyMode, title)
		if err := cmd.Run(); err != nil {
			// Fallback to macOS open without -n to preserve instance
			return exec.Command("open", "-a", "Ghostty", resolvedPath).Run()
		}
		return nil

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
			return LaunchWithGhosttyMode(def, resolvedPath, ghosttyMode, title)
		}
		cmd := exec.Command("open", "-a", "Terminal", resolvedPath)
		return cmd.Run()
	}
}

func Launch(terminalID TerminalID, projectPath string) error {
	return LaunchWithGhosttyMode(terminalID, projectPath, "window", filepath.Base(projectPath))
}
