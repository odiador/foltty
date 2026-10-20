package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/odiador/foltty/internal/terminal"
)

type Project struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
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

		projects = append(projects, Project{
			Name:   name,
			Path:   expanded,
			Exists: exists,
		})
	}

	return projects, nil
}

func AddProject(targetPath string) error {
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

	for _, p := range projects {
		if p.Path == resolved {
			return fmt.Errorf("project already registered: %s", resolved)
		}
	}

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
		if p.Path == resolved || p.Name == targetPath || strings.EqualFold(p.Name, targetPath) {
			found = true
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

func OpenInTerminal(targetPath string, termID terminal.TerminalID) error {
	resolved, err := ResolvePresetOrPath(targetPath)
	if err != nil {
		return err
	}

	return terminal.Launch(termID, resolved)
}
