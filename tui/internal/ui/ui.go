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
	launchTitle    string
	sortBy         string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type to search projects, paths or presets..."
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 62

	ci := textinput.New()
	ci.Prompt = ":"
	ci.Placeholder = "help, add <path>, rename <name>, rm, term <id>, q"
	ci.CharLimit = 150
	ci.Width = 62

	api := textinput.New()
	api.Placeholder = "Folder path or preset (e.g. ~/Dev/app, docs, .)"
	api.CharLimit = 150
	api.Width = 62

	ani := textinput.New()
	ani.Placeholder = "Friendly name (optional, press Enter to skip)"
	ani.CharLimit = 50
	ani.Width = 62

	ri := textinput.New()
	ri.Placeholder = "New display name"
	ri.CharLimit = 50
	ri.Width = 62

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
		sortBy:         config.GetSortBy(),
	}
	m.applyFilter()
	return m
}

func (m *model) applyFilter() {
	query := strings.TrimSpace(strings.ToLower(m.searchInput.Value()))
	var res []project.Project
	if query == "" {
		res = m.projects
	} else {
		for _, p := range m.projects {
			matchName := strings.Contains(strings.ToLower(p.Name), query)
			matchPath := strings.Contains(strings.ToLower(p.Path), query)
			matchCustom := strings.Contains(strings.ToLower(p.CustomName), query)
			if matchName || matchPath || matchCustom {
				res = append(res, p)
			}
		}
	}

	m.filtered = project.SortProjects(res, m.sortBy)

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

			case "backspace", "delete":
				if m.addStep == 0 && m.addPathInput.Value() == "" {
					m.mode = modeSearch
					m.addPathInput.Reset()
					m.addNameInput.Reset()
					m.searchInput.Focus()
					return m, nil
				}

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

			case "backspace", "delete":
				if m.renameInput.Value() == "" {
					m.mode = modeSearch
					m.renameInput.Reset()
					m.searchInput.Focus()
					return m, nil
				}

			case "enter":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					cur := m.filtered[m.selectedIndex]
					newName := m.renameInput.Value()
					_ = project.SetProjectName(cur.Path, newName)
					m.statusMsg = fmt.Sprintf("Renamed to '%s'", newName)
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

			case "backspace", "delete":
				if m.cmdInput.Value() == "" {
					m.mode = modeSearch
					m.cmdInput.Reset()
					m.searchInput.Focus()
					return m, nil
				}

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
					if m.selectedIndex > 0 {
						m.selectedIndex--
					}
				}
				return m, nil

			case "down", "ctrl+j", "ctrl+n":
				if len(m.filtered) > 0 {
					if m.selectedIndex < len(m.filtered)-1 {
						m.selectedIndex++
					}
				}
				return m, nil

			case "enter":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					m.launchTarget = target.Path
					m.launchTitle = target.Name
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil

			case "ctrl+t":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					m.launchTarget = target.Path
					m.launchTitle = target.Name
					m.ghosttyMode = "tab"
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil

			case "ctrl+w":
				if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
					target := m.filtered[m.selectedIndex]
					m.launchTarget = target.Path
					m.launchTitle = target.Name
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

			case "ctrl+s":
				switch m.sortBy {
				case "recent":
					m.sortBy = "name"
				case "name":
					m.sortBy = "name-desc"
				case "name-desc":
					m.sortBy = "added"
				default:
					m.sortBy = "recent"
				}
				_ = config.SetSortBy(m.sortBy)
				m.statusMsg = fmt.Sprintf("✦ Sort order: %s", strings.ToUpper(m.sortBy))
				m.applyFilter()
				return m, nil

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
		m.statusMsg = fmt.Sprintf("✓ Renamed to '%s'", newName)
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
			m.launchTitle = cur.Name
			m.ghosttyMode = "tab"
			m.quitting = true
		} else {
			m.statusMsg = "❌ No project selected"
		}

	case "win", "window":
		if len(m.filtered) > 0 && m.selectedIndex < len(m.filtered) {
			cur := m.filtered[m.selectedIndex]
			m.launchTarget = cur.Path
			m.launchTitle = cur.Name
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

	case "sort", "s":
		if len(parts) < 2 {
			m.statusMsg = fmt.Sprintf("Sort: %s. Use :sort <recent|name|added>", m.sortBy)
			return
		}
		targetSort := strings.ToLower(parts[1])
		_ = config.SetSortBy(targetSort)
		m.sortBy = config.GetSortBy()
		m.statusMsg = fmt.Sprintf("✦ Sort order set to: %s", strings.ToUpper(m.sortBy))
		m.applyFilter()

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

	contentWidth := 68

	// 1. Sleek Compact Header (No giant ASCII banner in main view)
	titlePart := lipgloss.NewStyle().Bold(true).Foreground(cyanColor).Render("👻 FOLTTY") + " " +
		lipgloss.NewStyle().Foreground(grayColor).Render("• Spotlight Project Launcher")
	termBadge := badgeStyle.Render("TTY: " + strings.ToUpper(string(m.preferredTerm)))
	sortBadge := badgeStyle.Render("SORT: " + strings.ToUpper(m.sortBy))
	badges := termBadge + " " + sortBadge

	gapWidth := max(2, contentWidth-lipgloss.Width(titlePart)-lipgloss.Width(badges))
	headerBar := lipgloss.JoinHorizontal(lipgloss.Center, titlePart, strings.Repeat(" ", gapWidth), badges)

	// 2. Search / Command Bar (Full width)
	var inputArea string
	switch m.mode {
	case modeCommand:
		inputArea = cmdBoxStyle.Width(contentWidth).Render(m.cmdInput.View())
	case modeAdd:
		if m.addStep == 0 {
			inputArea = searchBoxStyle.Width(contentWidth).Render("📁 Path: " + m.addPathInput.View())
		} else {
			inputArea = searchBoxStyle.Width(contentWidth).Render("Display Name: " + m.addNameInput.View())
		}
	case modeRename:
		inputArea = searchBoxStyle.Width(contentWidth).Render("Rename: " + m.renameInput.View())
	default:
		inputArea = searchBoxStyle.Width(contentWidth).Render("🔍 " + m.searchInput.View())
	}

	// 3. Projects List with Midpoint Centered Window Scrolling
	var listItems []string
	homeDir, _ := os.UserHomeDir()

	visibleCount := 7
	total := len(m.filtered)

	if total == 0 {
		listItems = append(listItems, lipgloss.NewStyle().Foreground(grayColor).Italic(true).Render("\n  No matching projects found.\n  Press Ctrl+A or type :a <path> to add one!\n"))
	} else {
		start := 0
		if total > visibleCount {
			mid := visibleCount / 2 // 3 -> cursor stays centered around the 4th item (middle)
			start = max(0, min(m.selectedIndex-mid, total-visibleCount))
		}
		end := min(total, start+visibleCount)

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

			timeBadge := ""
			if p.LastOpenedAt > 0 {
				relTime := project.FormatRelativeTime(p.LastOpenedAt)
				if relTime != "" {
					timeBadge = lipgloss.NewStyle().Foreground(grayColor).Render(" • " + relTime)
				}
			}

			line1 := pointer + titleStyle.Render(displayName) + missingBadge + timeBadge
			line2 := "     " + itemPathStyle.Render(displayPath)
			listItems = append(listItems, line1+"\n"+line2)
		}
	}
	listBox := lipgloss.NewStyle().Width(contentWidth).Render(lipgloss.JoinVertical(lipgloss.Left, listItems...))

	// 4. Compact Horizontal Controls Bar
	keyStyle := lipgloss.NewStyle().Foreground(cyanColor).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(grayColor)
	bullet := lipgloss.NewStyle().Foreground(subtleColor).Render(" • ")

	controlsBar := lipgloss.JoinHorizontal(lipgloss.Center,
		keyStyle.Render("↵ ") + descStyle.Render("open"),
		bullet,
		keyStyle.Render("^w ") + descStyle.Render("win"),
		bullet,
		keyStyle.Render("^t ") + descStyle.Render("tab"),
		bullet,
		keyStyle.Render("^a ") + descStyle.Render("add"),
		bullet,
		keyStyle.Render("^r ") + descStyle.Render("rename"),
		bullet,
		keyStyle.Render("^s ") + descStyle.Render("sort"),
		bullet,
		keyStyle.Render(": ") + descStyle.Render("vim"),
		bullet,
		keyStyle.Render("esc ") + descStyle.Render("quit"),
	)

	// 5. Status / Info Footer
	var statusLine string
	if m.statusMsg != "" {
		statusLine = statusStyle.Render(m.statusMsg)
	} else {
		if total == 0 {
			statusLine = lipgloss.NewStyle().Foreground(grayColor).Render("0 projects  •  Type :h for help")
		} else {
			statusLine = lipgloss.NewStyle().Foreground(grayColor).Render(fmt.Sprintf("%d/%d projects  •  Type :h for help", m.selectedIndex+1, total))
		}
	}

	divider := lipgloss.NewStyle().Foreground(subtleColor).Render(strings.Repeat("─", contentWidth))

	// 6. Combine Layout
	content := lipgloss.JoinVertical(lipgloss.Left,
		headerBar,
		"\n",
		inputArea,
		"\n",
		listBox,
		"\n",
		divider,
		controlsBar,
		statusLine,
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
  Ctrl + S            Cycle sort order (recent, name, added)
  Ctrl + D / Ctrl + X Remove selected project from list
  :                   Enter Vim command mode
  Esc                 Back to search / Quit

VIM COMMANDS (:):
  :h, :help           Show this help manual
  :tab                Open selected project in Ghostty (New Tab)
  :win, :window       Open selected project in Ghostty (New Window)
  :mode <tab|window>  Set Ghostty default launch mode
  :s, :sort <mode>    Set sort order (recent, name, added)
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
		title := finalModel.launchTitle
		if title == "" {
			title = filepath.Base(finalModel.launchTarget)
		}
		fmt.Printf("🚀 Launching %s in %s (%s)...\n", title, finalModel.preferredTerm, finalModel.ghosttyMode)
		return project.OpenInTerminalWithMode(finalModel.launchTarget, finalModel.preferredTerm, finalModel.ghosttyMode, title)
	}

	return nil
}
