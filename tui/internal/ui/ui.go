package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/odiador/foltty/internal/config"
	"github.com/odiador/foltty/internal/project"
	"github.com/odiador/foltty/internal/terminal"
)

const asciiBanner = ` ______   ______     __         ______   ______   __  __    
/\  ___\ /\  __ \   /\ \       /\__  _\ /\__  _\ /\ \_\ \   
\ \  __\ \ \ \/\ \  \ \ \____  \/_/\ \/ \/_/\ \/ \ \____ \  
 \ \_\    \ \_____\  \ \_____\    \ \_\    \ \_\  \/\_____\ 
  \/_/     \/_____/   \/_____/     \/_/     \/_/   \/_____/ `

type mode int

const (
	modeSearch mode = iota
	modeCommand
	modeHelp
	modeAdd
	modeRename
)

type model struct {
	width          int
	height         int
	mode           mode
	projects       []project.Project
	filtered       []project.Project
	selectedIndex  int
	searchInput    textinput.Model
	cmdInput       textinput.Model
	addPathInput   textinput.Model
	addNameInput   textinput.Model
	addStep        int // 0 = path, 1 = name
	renameInput    textinput.Model
	statusMsg      string
	preferredTerm  terminal.TerminalID
	terminalActive string
	ghosttyMode    string
	quitting       bool
	launchTarget   string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type to search projects, paths or presets..."
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 46

	ci := textinput.New()
	ci.Prompt = ":"
	ci.Placeholder = "help, add <path>, rename <name>, rm, term <id>, q"
	ci.CharLimit = 150
	ci.Width = 50

	api := textinput.New()
	api.Placeholder = "Folder path or preset (e.g. ~/Dev/app, docs, .)"
	api.CharLimit = 150
	api.Width = 46

	ani := textinput.New()
	ani.Placeholder = "Friendly name (optional, press Enter to skip)"
	ani.CharLimit = 50
	ani.Width = 46

	ri := textinput.New()
	ri.Placeholder = "New display name"
	ri.CharLimit = 50
	ri.Width = 46

	prefTerm := config.GetPreferredTerminal()
	projs, _ := project.GetProjects()

	m := model{
		mode:           modeSearch,
		projects:       projs,
		filtered:       projs,
		searchInput:    ti,
		cmdInput:       ci,
		addPathInput:   api,
		addNameInput:   ani,
		renameInput:    ri,
		preferredTerm:  prefTerm,
		terminalActive: string(prefTerm),
		ghosttyMode:    config.GetGhosttyMode(),
	}
	m.applyFilter()
	return m
}

func (m *model) applyFilter() {
	query := strings.TrimSpace(strings.ToLower(m.searchInput.Value()))
	if query == "" {
		m.filtered = m.projects
	} else {
		var res []project.Project
		for _, p := range m.projects {
			matchName := strings.Contains(strings.ToLower(p.Name), query)
			matchPath := strings.Contains(strings.ToLower(p.Path), query)
			matchCustom := strings.Contains(strings.ToLower(p.CustomName), query)
			if matchName || matchPath || matchCustom {
				res = append(res, p)
			}
		}
		m.filtered = res
	}

	if m.selectedIndex >= len(m.filtered) {
		m.selectedIndex = max(0, len(m.filtered)-1)
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		switch m.mode {
		case modeHelp:
			if msg.String() == "esc" || msg.String() == "q" || msg.String() == "enter" {
				m.mode = modeSearch
				m.searchInput.Focus()
				return m, nil
			}
			return m, nil

		case modeAdd:
			switch msg.String() {
			case "esc":
				m.mode = modeSearch
				m.addStep = 0
				m.addPathInput.Reset()
				m.addNameInput.Reset()
				m.searchInput.Focus()
				return m, nil

			case "enter":
				if m.addStep == 0 {
					if strings.TrimSpace(m.addPathInput.Value()) == "" {
						m.statusMsg = "❌ Path cannot be empty"
						return m, nil
					}
					m.addStep = 1
					m.addNameInput.Focus()
					return m, nil
				} else {
					targetPath := m.addPathInput.Value()
					customName := m.addNameInput.Value()
					err := project.AddProjectWithName(targetPath, customName)
					if err != nil {
						m.statusMsg = fmt.Sprintf("❌ %v", err)
					} else {
						m.statusMsg = fmt.Sprintf("✅ Added: %s", targetPath)
						m.projects, _ = project.GetProjects()
						m.applyFilter()
					}
					m.mode = modeSearch
					m.addStep = 0
					m.addPathInput.Reset()
					m.addNameInput.Reset()
					m.searchInput.Focus()
					return m, nil
				}
			}

			if m.addStep == 0 {
				m.addPathInput, cmd = m.addPathInput.Update(msg)
			} else {
				m.addNameInput, cmd = m.addNameInput.Update(msg)
			}
			return m, cmd

		case modeRename:
			switch msg.String() {
			case "esc":
				m.mode = modeSearch
				m.renameInput.Reset()
				m.searchInput.Focus()
				return m, nil

			case "enter":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					cur := m.filtered[m.selectedIndex]
					newName := m.renameInput.Value()
					_ = project.SetProjectName(cur.Path, newName)
					m.statusMsg = fmt.Sprintf("🏷️ Renamed to '%s'", newName)
					m.projects, _ = project.GetProjects()
					m.applyFilter()
				}
				m.mode = modeSearch
				m.renameInput.Reset()
				m.searchInput.Focus()
				return m, nil
			}
			m.renameInput, cmd = m.renameInput.Update(msg)
			return m, cmd

		case modeCommand:
			switch msg.String() {
			case "esc":
				m.mode = modeSearch
				m.cmdInput.Reset()
				m.searchInput.Focus()
				return m, nil

			case "enter":
				rawCmd := strings.TrimSpace(m.cmdInput.Value())
				m.executeVimCommand(rawCmd)
				m.cmdInput.Reset()
				if m.mode == modeCommand {
					m.mode = modeSearch
					m.searchInput.Focus()
				}
				return m, nil
			}
			m.cmdInput, cmd = m.cmdInput.Update(msg)
			return m, cmd

		case modeSearch:
			// Normal search and navigation
			switch msg.String() {
			case ":":
				m.mode = modeCommand
				m.cmdInput.Reset()
				m.cmdInput.Focus()
				return m, nil

			case "up", "ctrl+k", "ctrl+p":
				if len(m.filtered) > 0 {
					m.selectedIndex = (m.selectedIndex - 1 + len(m.filtered)) % len(m.filtered)
				}
				return m, nil

			case "down", "ctrl+j", "ctrl+n":
				if len(m.filtered) > 0 {
					m.selectedIndex = (m.selectedIndex + 1) % len(m.filtered)
				}
				return m, nil

			case "enter":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					m.launchTarget = target.Path
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil

			case "ctrl+t":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					m.launchTarget = target.Path
					m.ghosttyMode = "tab"
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil

			case "ctrl+w":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					m.launchTarget = target.Path
					m.ghosttyMode = "window"
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil

			case "ctrl+a":
				m.mode = modeAdd
				m.addStep = 0
				m.addPathInput.Reset()
				m.addNameInput.Reset()
				m.addPathInput.Focus()
				return m, nil

			case "ctrl+r":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					m.mode = modeRename
					cur := m.filtered[m.selectedIndex]
					m.renameInput.SetValue(cur.Name)
					m.renameInput.Focus()
					return m, nil
				}

			case "ctrl+d", "ctrl+x":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					_ = project.RemoveProject(target.Path)
					m.statusMsg = fmt.Sprintf("🗑️ Removed '%s'", target.Name)
					m.projects, _ = project.GetProjects()
					m.applyFilter()
				}
				return m, nil

			case "esc":
				if m.searchInput.Value() != "" {
					m.searchInput.Reset()
					m.applyFilter()
					return m, nil
				}
				m.quitting = true
				return m, tea.Quit
			}

			var tiCmd tea.Cmd
			m.searchInput, tiCmd = m.searchInput.Update(msg)
			m.applyFilter()
			return m, tiCmd
		}
	}

	return m, cmd
}

func (m *model) executeVimCommand(raw string) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return
	}

	verb := strings.ToLower(parts[0])

	switch verb {
	case "h", "help":
		m.mode = modeHelp

	case "q", "quit", "exit":
		m.quitting = true

	case "a", "add":
		if len(parts) < 2 {
			m.statusMsg = "❌ Usage: :a <path> [name]"
			return
		}
		targetPath := parts[1]
		customName := ""
		if len(parts) > 2 {
			customName = strings.Join(parts[2:], " ")
		}
		if err := project.AddProjectWithName(targetPath, customName); err != nil {
			m.statusMsg = fmt.Sprintf("❌ %v", err)
		} else {
			m.statusMsg = fmt.Sprintf("✅ Added: %s", targetPath)
			m.projects, _ = project.GetProjects()
			m.applyFilter()
		}

	case "r", "rename":
		if len(m.filtered) == 0 || m.selectedIndex >= len(m.filtered) {
			m.statusMsg = "❌ No project selected"
			return
		}
		if len(parts) < 2 {
			m.statusMsg = "❌ Usage: :r <new name>"
			return
		}
		newName := strings.Join(parts[1:], " ")
		cur := m.filtered[m.selectedIndex]
		_ = project.SetProjectName(cur.Path, newName)
		m.statusMsg = fmt.Sprintf("🏷️ Renamed to '%s'", newName)
		m.projects, _ = project.GetProjects()
		m.applyFilter()

	case "d", "delete", "rm":
		if len(m.filtered) == 0 || m.selectedIndex >= len(m.filtered) {
			m.statusMsg = "❌ No project selected"
			return
		}
		cur := m.filtered[m.selectedIndex]
		_ = project.RemoveProject(cur.Path)
		m.statusMsg = fmt.Sprintf("🗑️ Removed '%s'", cur.Name)
		m.projects, _ = project.GetProjects()
		m.applyFilter()

	case "term", "terminal":
		if len(parts) < 2 {
			current := config.GetPreferredTerminal()
			m.statusMsg = fmt.Sprintf("Active: %s. Use :term <ghostty|iterm|wezterm|warp|terminal>", current)
			return
		}
		termID := strings.ToLower(parts[1])
		_ = config.SetPreferredTerminal(termID)
		m.preferredTerm = terminal.TerminalID(termID)
		m.terminalActive = termID
		m.statusMsg = fmt.Sprintf("⚡ Terminal set to %s", termID)

	case "tab":
		if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
			cur := m.filtered[m.selectedIndex]
			m.launchTarget = cur.Path
			m.ghosttyMode = "tab"
			m.quitting = true
		} else {
			m.statusMsg = "❌ No project selected"
		}

	case "win", "window":
		if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
			cur := m.filtered[m.selectedIndex]
			m.launchTarget = cur.Path
			m.ghosttyMode = "window"
			m.quitting = true
		} else {
			m.statusMsg = "❌ No project selected"
		}

	case "mode":
		if len(parts) < 2 {
			m.statusMsg = fmt.Sprintf("Ghostty Mode: %s. Use :mode <tab|window>", m.ghosttyMode)
			return
		}
		targetMode := strings.ToLower(parts[1])
		if targetMode != "tab" && targetMode != "window" {
			m.statusMsg = "❌ Mode must be 'tab' or 'window'"
			return
		}
		_ = config.SetGhosttyMode(targetMode)
		m.ghosttyMode = targetMode
		m.statusMsg = fmt.Sprintf("⚡ Ghostty default mode set to %s", targetMode)

	default:
		m.statusMsg = fmt.Sprintf("❓ Unknown command: %s (Type :h for help)", verb)
	}
}

// -----------------------------------------------------------------------------
// Lipgloss Styles & Rendering
// -----------------------------------------------------------------------------

var (
	purpleColor = lipgloss.Color("#a855f7")
	pinkColor   = lipgloss.Color("#ec4899")
	cyanColor   = lipgloss.Color("#38bdf8")
	grayColor   = lipgloss.Color("#64748b")
	redColor    = lipgloss.Color("#f87171")
	subtleColor = lipgloss.Color("#334155")

	bannerStyle = lipgloss.NewStyle().
			Foreground(purpleColor).
			Bold(true)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purpleColor).
			Padding(1, 2)

	searchBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cyanColor).
			Padding(0, 1)

	cmdBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(pinkColor).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(cyanColor).
				Bold(true)

	itemPathStyle = lipgloss.NewStyle().
			Foreground(grayColor)

	badgeStyle = lipgloss.NewStyle().
			Background(purpleColor).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1).
			Bold(true)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(subtleColor).
			PaddingLeft(2)

	statusStyle = lipgloss.NewStyle().
			Foreground(pinkColor).
			Italic(true)
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	if m.mode == modeHelp {
		return m.renderHelpView()
	}

	// 1. Banner
	banner := bannerStyle.Render(asciiBanner)
	headerText := lipgloss.NewStyle().Foreground(grayColor).Render("folder + tty = foltty  •  Spotlight Project Launcher")
	header := lipgloss.JoinVertical(lipgloss.Center, banner, headerText)

	// 2. Search / Command Bar
	var inputArea string
	switch m.mode {
	case modeCommand:
		inputArea = cmdBoxStyle.Render(m.cmdInput.View())
	case modeAdd:
		if m.addStep == 0 {
			inputArea = searchBoxStyle.Render("📁 Add Path: " + m.addPathInput.View())
		} else {
			inputArea = searchBoxStyle.Render("🏷️ Display Name: " + m.addNameInput.View())
		}
	case modeRename:
		inputArea = searchBoxStyle.Render("✏️ Rename: " + m.renameInput.View())
	default:
		inputArea = searchBoxStyle.Render("🔍 " + m.searchInput.View())
	}

	// 3. Projects List (Left Column)
	var listItems []string
	homeDir, _ := os.UserHomeDir()

	if len(m.filtered) == 0 {
		listItems = append(listItems, lipgloss.NewStyle().Foreground(grayColor).Italic(true).Render("\n  No matching projects found.\n  Press Ctrl+A or type :a <path> to add one!"))
	} else {
		maxVisible := 8
		start := 0
		if m.selectedIndex >= maxVisible {
			start = m.selectedIndex - maxVisible + 1
		}
		end := min(len(m.filtered), start+maxVisible)

		for i := start; i < end; i++ {
			p := m.filtered[i]
			isSelected := i == m.selectedIndex

			displayPath := p.Path
			if strings.HasPrefix(displayPath, homeDir) {
				displayPath = "~" + displayPath[len(homeDir):]
			}

			pointer := "   "
			titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#e2e8f0"))
			if isSelected {
				pointer = " ❯ "
				titleStyle = selectedItemStyle
			}

			missingBadge := ""
			if !p.Exists {
				missingBadge = lipgloss.NewStyle().Foreground(redColor).Render(" [Missing]")
			}

			displayName := p.Name
			if p.CustomName != "" {
				displayName = fmt.Sprintf("%s 🏷️", p.CustomName)
			}

			line1 := pointer + titleStyle.Render(displayName) + missingBadge
			line2 := "     " + itemPathStyle.Render(displayPath)
			listItems = append(listItems, line1+"\n"+line2)
		}
	}
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, listItems...)
	leftBox := lipgloss.NewStyle().Width(52).Render(leftColumn)

	// 4. Sidebar / Keybindings & Context (Right Column)
	var sidebarItems []string
	termBadge := badgeStyle.Render("TTY: " + strings.ToUpper(string(m.preferredTerm)))
	sidebarItems = append(sidebarItems, termBadge+"\n")

	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Bold(true).Foreground(cyanColor).Render("Shortcuts"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render("↵ Enter   Open Terminal"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render("^A        Add project"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render("^R        Rename project"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render("^D        Delete project"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render(":         Vim commands"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render("Esc       Quit / Back"))

	sidebarItems = append(sidebarItems, "\n"+lipgloss.NewStyle().Bold(true).Foreground(pinkColor).Render("Vim Mode (:)"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render(":h        Open Help"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render(":a <path> Add project"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render(":r <name> Rename project"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render(":term <id>Switch terminal"))
	sidebarItems = append(sidebarItems, lipgloss.NewStyle().Foreground(grayColor).Render(":q        Quit"))

	rightBox := sidebarStyle.Width(30).Render(lipgloss.JoinVertical(lipgloss.Left, sidebarItems...))

	// 5. Main Body Grid (Left List + Right Sidebar)
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	// 6. Footer Status
	footer := ""
	if m.statusMsg != "" {
		footer = statusStyle.Render(m.statusMsg)
	} else {
		footer = lipgloss.NewStyle().Foreground(subtleColor).Render(fmt.Sprintf("%d projects registered  •  Press ':' for vim commands", len(m.projects)))
	}

	// 7. Combine Card Layout
	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"\n",
		inputArea,
		"\n",
		body,
		"\n",
		footer,
	)

	card := cardStyle.Render(content)

	// Center on screen
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
	}
	return card
}

func (m model) renderHelpView() string {
	title := bannerStyle.Render(asciiBanner)
	heading := lipgloss.NewStyle().Bold(true).Foreground(cyanColor).Render("\n📖 FOLTTY CHEATSHEET & VIM COMMANDS\n")

	helpContent := `
Foltty is your Spotlight-like project launcher for macOS terminal emulators.

NAVIGATION & ACTIONS:
  ↑ / ↓ / ^J / ^K     Navigate through project list
  Enter               Launch selected project in active terminal
  Ctrl + T            Quick launch selected in Ghostty (New Tab)
  Ctrl + W            Quick launch selected in Ghostty (New Window)
  Ctrl + A            Add new project or path interactively
  Ctrl + R            Rename / assign custom display name to selected project
  Ctrl + D / Ctrl + X Remove selected project from list
  :                   Enter Vim command mode
  Esc                 Back to search / Quit

VIM COMMANDS (:):
  :h, :help           Show this help manual
  :tab                Open selected project in Ghostty (New Tab)
  :win, :window       Open selected project in Ghostty (New Window)
  :mode <tab|window>  Set Ghostty default launch mode
  :a <path> [name]    Add project path (e.g. :a ~/Dev/my-api "Backend API")
  :r <name>           Rename selected project's display name
  :d, :rm             Delete selected project
  :term <id>          Set default terminal (ghostty, iterm, wezterm, warp, terminal)
  :q, :quit           Quit Foltty

PRESETS AVAILABLE:
  home (~)            User home folder
  docs                ~/Documents
  downloads           ~/Downloads
  dev                 ~/Development or ~/Projects
  config              ~/.config

Press [Esc], [q], or [Enter] to return to search.`

	body := lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1")).Render(helpContent)
	box := cardStyle.Render(lipgloss.JoinVertical(lipgloss.Center, title, heading, body))

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return box
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Run() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := m.(model)
	if finalModel.launchTarget != "" {
		fmt.Printf("🚀 Launching %s in %s (%s)...\n", filepath.Base(finalModel.launchTarget), finalModel.preferredTerm, finalModel.ghosttyMode)
		return project.OpenInTerminalWithMode(finalModel.launchTarget, finalModel.preferredTerm, finalModel.ghosttyMode)
	}

	return nil
}
