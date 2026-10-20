package project

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Project struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
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

func NormalizePath(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New("la ruta no puede estar vacía")
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
	normalized, err := NormalizePath(targetPath)
	if err != nil {
		return err
	}

	stat, err := os.Stat(normalized)
	if err != nil {
		return fmt.Errorf("la ruta no existe: %s", normalized)
	}
	if !stat.IsDir() {
		return fmt.Errorf("la ruta no es un directorio: %s", normalized)
	}

	projects, err := GetProjects()
	if err != nil {
		return err
	}

	for _, p := range projects {
		if p.Path == normalized {
			return fmt.Errorf("el proyecto ya está registrado: %s", normalized)
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

	if _, err := f.WriteString(normalized + "\n"); err != nil {
		return err
	}

	return nil
}

func RemoveProject(targetPath string) error {
	normalized, err := NormalizePath(targetPath)
	if err != nil {
		normalized = targetPath
	}

	projects, err := GetProjects()
	if err != nil {
		return err
	}

	var updated []string
	found := false

	for _, p := range projects {
		if p.Path == normalized || p.Name == targetPath {
			found = true
			continue
		}
		updated = append(updated, p.Path)
	}

	if !found {
		return fmt.Errorf("proyecto no encontrado: %s", targetPath)
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

func OpenInGhostty(targetPath string) error {
	normalized, err := NormalizePath(targetPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(normalized); err != nil {
		return fmt.Errorf("el directorio no existe: %s", normalized)
	}

	cmd := exec.Command("open", "-na", "Ghostty", "--args", fmt.Sprintf("--working-directory=%s", normalized))
	return cmd.Run()
}
