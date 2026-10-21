package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/odiador/foltty/internal/terminal"
)

type Project struct {
	Name         string `json:"name"`
	CustomName   string `json:"customName,omitempty"`
	Path         string `json:"path"`
	Exists       bool   `json:"exists"`
	LastOpenedAt int64  `json:"lastOpenedAt,omitempty"`
}

type Preset struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsCustom bool   `json:"isCustom"`
}

func GetStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	folttyFile := filepath.Join(home, ".foltty_projects")
	pmanFile := filepath.Join(home, ".pman_projects")

	if _, err := os.Stat(folttyFile); err == nil {
		return folttyFile, nil
	}

	if _, err := os.Stat(pmanFile); err == nil {
		content, err := os.ReadFile(pmanFile)
		if err == nil {
			_ = os.WriteFile(folttyFile, content, 0644)
			return folttyFile, nil
		}
		return pmanFile, nil
	}

	return folttyFile, nil
}

func GetNamesFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".foltty_names.json"), nil
}

func GetCustomNames() (map[string]string, error) {
	filePath, err := GetNamesFilePath()
	if err != nil {
		return make(map[string]string), err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return make(map[string]string), nil
	}

	var names map[string]string
	if err := json.Unmarshal(data, &names); err != nil {
		return make(map[string]string), nil
	}
	return names, nil
}

func GetRecentFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".foltty_recent.json"), nil
}

func GetRecentHistory() (map[string]int64, error) {
	filePath, err := GetRecentFilePath()
	if err != nil {
		return make(map[string]int64), err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return make(map[string]int64), nil
	}

	var history map[string]int64
	if err := json.Unmarshal(data, &history); err != nil {
		return make(map[string]int64), nil
	}
	return history, nil
}

func RecordProjectOpened(targetPath string) error {
	resolved, err := ResolvePresetOrPath(targetPath)
	if err != nil {
		resolved = targetPath
	}

	norm, err := NormalizePath(resolved)
	if err == nil {
		resolved = norm
	}

	history, _ := GetRecentHistory()
	if history == nil {
		history = make(map[string]int64)
	}
	history[resolved] = time.Now().UnixMilli()

	filePath, err := GetRecentFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func SetProjectName(target, customName string) error {
	projects, _ := GetProjects()
	var matchedPath string
	for _, p := range projects {
		if strings.EqualFold(p.Name, target) || strings.EqualFold(p.CustomName, target) || p.Path == target {
			matchedPath = p.Path
			break
		}
	}
	if matchedPath == "" {
		norm, err := NormalizePath(target)
		if err == nil {
			matchedPath = norm
		} else {
			matchedPath = target
		}
	}

	names, _ := GetCustomNames()
	trimmed := strings.TrimSpace(customName)

	if trimmed == "" {
		delete(names, matchedPath)
	} else {
		names[matchedPath] = trimmed
	}

	filePath, err := GetNamesFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(names, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func GetPresetsFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".foltty_presets.json"), nil
}

func GetDefaultPresets() []Preset {
	home, _ := os.UserHomeDir()
	presets := []Preset{
		{ID: "home", Name: "Home (~)", Path: home, IsCustom: false},
		{ID: "documents", Name: "Documents", Path: filepath.Join(home, "Documents"), IsCustom: false},
		{ID: "downloads", Name: "Downloads", Path: filepath.Join(home, "Downloads"), IsCustom: false},
		{ID: "desktop", Name: "Desktop", Path: filepath.Join(home, "Desktop"), IsCustom: false},
		{ID: "config", Name: "Config (.config)", Path: filepath.Join(home, ".config"), IsCustom: false},
	}

	devDirs := []string{"Development", "Developer", "Projects", "Code", "workspace"}
	for _, dir := range devDirs {
		p := filepath.Join(home, dir)
		if stat, err := os.Stat(p); err == nil && stat.IsDir() {
			presets = append(presets, Preset{
				ID:       "dev",
				Name:     fmt.Sprintf("Dev (%s)", dir),
				Path:     p,
				IsCustom: false,
			})
			break
		}
	}

	return presets
}

func GetCustomPresets() (map[string]string, error) {
	presetsPath, err := GetPresetsFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(presetsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	var custom map[string]string
	if err := json.Unmarshal(data, &custom); err != nil {
		return make(map[string]string), nil
	}
	return custom, nil
}

func GetAllPresets() ([]Preset, error) {
	presets := GetDefaultPresets()
	custom, err := GetCustomPresets()
	if err != nil {
		return presets, nil
	}

	for id, pathStr := range custom {
		norm, _ := NormalizePath(pathStr)
		presets = append(presets, Preset{
			ID:       id,
			Name:     fmt.Sprintf("%s (Custom)", id),
			Path:     norm,
			IsCustom: true,
		})
	}

	return presets, nil
}

func AddCustomPreset(name, targetPath string) error {
	trimmedName := strings.ToLower(strings.TrimSpace(name))
	if trimmedName == "" {
		return errors.New("preset name cannot be empty")
	}

	norm, err := NormalizePath(targetPath)
	if err != nil {
		return err
	}

	custom, err := GetCustomPresets()
	if err != nil {
		custom = make(map[string]string)
	}

	custom[trimmedName] = norm

	presetsPath, err := GetPresetsFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(custom, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(presetsPath, data, 0644)
}

func RemoveCustomPreset(name string) error {
	trimmedName := strings.ToLower(strings.TrimSpace(name))
	custom, err := GetCustomPresets()
	if err != nil {
		return err
	}

	if _, ok := custom[trimmedName]; !ok {
		return fmt.Errorf("custom preset not found: %s", name)
	}

	delete(custom, trimmedName)

	presetsPath, err := GetPresetsFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(custom, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(presetsPath, data, 0644)
}

func ResolvePresetOrPath(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("path or preset cannot be empty")
	}

	lower := strings.ToLower(trimmed)
	presets, _ := GetAllPresets()
	for _, p := range presets {
		if strings.ToLower(p.ID) == lower || (lower == "docs" && p.ID == "documents") || (lower == "dl" && p.ID == "downloads") {
			return p.Path, nil
		}
	}

	return NormalizePath(trimmed)
}

func NormalizePath(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("path cannot be empty")
	}

	if strings.HasPrefix(trimmed, "~/") || trimmed == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if trimmed == "~" {
			return home, nil
		}
		return filepath.Join(home, trimmed[2:]), nil
	}

	return filepath.Abs(trimmed)
}

func GetProjects() ([]Project, error) {
	storagePath, err := GetStoragePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Project{}, nil
		}
		return nil, err
	}

	customNames, _ := GetCustomNames()
	recentHistory, _ := GetRecentHistory()

	lines := strings.Split(string(data), "\n")
	seen := make(map[string]bool)
	var projects []Project

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		expanded, err := NormalizePath(trimmed)
		if err != nil {
			continue
		}

		if seen[expanded] {
			continue
		}
		seen[expanded] = true

		_, statErr := os.Stat(expanded)
		exists := statErr == nil

		name := filepath.Base(expanded)
		if name == "" || name == "/" || name == "." {
			name = expanded
		}

		customLabel := ""
		if cName, ok := customNames[expanded]; ok && cName != "" {
			name = cName
			customLabel = cName
		}

		var lastOpened int64
		if ts, ok := recentHistory[expanded]; ok {
			lastOpened = ts
		}

		projects = append(projects, Project{
			Name:         name,
			CustomName:   customLabel,
			Path:         expanded,
			Exists:       exists,
			LastOpenedAt: lastOpened,
		})
	}

	return projects, nil
}

func SortProjects(projects []Project, sortBy string) []Project {
	res := make([]Project, len(projects))
	copy(res, projects)

	switch strings.ToLower(sortBy) {
	case "recent":
		sort.SliceStable(res, func(i, j int) bool {
			if res[i].LastOpenedAt != res[j].LastOpenedAt {
				return res[i].LastOpenedAt > res[j].LastOpenedAt
			}
			return strings.ToLower(res[i].Name) < strings.ToLower(res[j].Name)
		})
	case "name", "name-asc", "az":
		sort.SliceStable(res, func(i, j int) bool {
			return strings.ToLower(res[i].Name) < strings.ToLower(res[j].Name)
		})
	case "name-desc", "za":
		sort.SliceStable(res, func(i, j int) bool {
			return strings.ToLower(res[i].Name) > strings.ToLower(res[j].Name)
		})
	case "added", "file":
		// Preserves file order
	default:
		// Default to recent
		sort.SliceStable(res, func(i, j int) bool {
			if res[i].LastOpenedAt != res[j].LastOpenedAt {
				return res[i].LastOpenedAt > res[j].LastOpenedAt
			}
			return strings.ToLower(res[i].Name) < strings.ToLower(res[j].Name)
		})
	}
	return res
}

func FormatRelativeTime(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	now := time.Now().UnixMilli()
	diffMs := now - timestamp
	if diffMs < 0 {
		return "just now"
	}

	diffSec := diffMs / 1000
	diffMin := diffSec / 60
	diffHour := diffMin / 60
	diffDay := diffHour / 24

	if diffMin < 1 {
		return "just now"
	}
	if diffMin < 60 {
		return fmt.Sprintf("%dm ago", diffMin)
	}
	if diffHour < 24 {
		return fmt.Sprintf("%dh ago", diffHour)
	}
	if diffDay == 1 {
		return "yesterday"
	}
	if diffDay < 7 {
		return fmt.Sprintf("%dd ago", diffDay)
	}

	t := time.UnixMilli(timestamp)
	return t.Format("Jan 2")
}

func AddProject(targetPath string) error {
	return AddProjectWithName(targetPath, "")
}

func AddProjectWithName(targetPath, customName string) error {
	resolved, err := ResolvePresetOrPath(targetPath)
	if err != nil {
		return err
	}

	stat, err := os.Stat(resolved)
	if err != nil {
		return fmt.Errorf("path does not exist: %s", resolved)
	}
	if !stat.IsDir() {
		return fmt.Errorf("path is not a directory: %s", resolved)
	}

	projects, err := GetProjects()
	if err != nil {
		return err
	}

	alreadyRegistered := false
	for _, p := range projects {
		if p.Path == resolved {
			alreadyRegistered = true
			break
		}
	}

	if !alreadyRegistered {
		storagePath, err := GetStoragePath()
		if err != nil {
			return err
		}

		f, err := os.OpenFile(storagePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString(resolved + "\n"); err != nil {
			return err
		}
	}

	if strings.TrimSpace(customName) != "" {
		_ = SetProjectName(resolved, customName)
	}

	return nil
}

func RemoveProject(targetPath string) error {
	resolved, err := ResolvePresetOrPath(targetPath)
	if err != nil {
		resolved = targetPath
	}

	projects, err := GetProjects()
	if err != nil {
		return err
	}

	var updated []string
	found := false

	for _, p := range projects {
		if p.Path == resolved || p.Name == targetPath || p.CustomName == targetPath || strings.EqualFold(p.Name, targetPath) {
			found = true
			_ = SetProjectName(p.Path, "") // remove custom name
			recentHistory, _ := GetRecentHistory()
			if _, ok := recentHistory[p.Path]; ok {
				delete(recentHistory, p.Path)
				if recentPath, err := GetRecentFilePath(); err == nil {
					if data, err := json.MarshalIndent(recentHistory, "", "  "); err == nil {
						_ = os.WriteFile(recentPath, data, 0644)
					}
				}
			}
			continue
		}
		updated = append(updated, p.Path)
	}

	if !found {
		return fmt.Errorf("project not found in list: %s", targetPath)
	}

	storagePath, err := GetStoragePath()
	if err != nil {
		return err
	}

	content := strings.Join(updated, "\n")
	if len(updated) > 0 {
		content += "\n"
	}

	return os.WriteFile(storagePath, []byte(content), 0644)
}

func OpenInTerminalWithMode(targetPath string, termID terminal.TerminalID, ghosttyMode string, title string) error {
	resolved, err := ResolvePresetOrPath(targetPath)
	if err != nil {
		return err
	}

	_ = RecordProjectOpened(resolved)

	if title == "" {
		names, _ := GetCustomNames()
		if n, ok := names[resolved]; ok && strings.TrimSpace(n) != "" {
			title = strings.TrimSpace(n)
		} else {
			title = filepath.Base(resolved)
		}
	}

	return terminal.LaunchWithGhosttyMode(termID, resolved, ghosttyMode, title)
}

func OpenInTerminal(targetPath string, termID terminal.TerminalID) error {
	return OpenInTerminalWithMode(targetPath, termID, "window", "")
}
