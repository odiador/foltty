package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/odiador/foltty/internal/terminal"
)

type Config struct {
	DefaultTerminal string `json:"defaultTerminal"`
	GhosttyMode     string `json:"ghosttyMode,omitempty"` // "tab" or "window"
	SortBy          string `json:"sortBy,omitempty"`      // "recent", "name", "name-desc", "added"
}

func GetConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".foltty_config.json"), nil
}

func LoadConfig() (*Config, error) {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return &Config{DefaultTerminal: "auto"}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return &Config{DefaultTerminal: "auto"}, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &Config{DefaultTerminal: "auto"}, nil
	}

	if cfg.DefaultTerminal == "" {
		cfg.DefaultTerminal = "auto"
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func GetPreferredTerminal() terminal.TerminalID {
	cfg, err := LoadConfig()
	if err != nil || cfg.DefaultTerminal == "auto" || cfg.DefaultTerminal == "" {
		return terminal.DetectDefaultTerminal()
	}
	return terminal.TerminalID(cfg.DefaultTerminal)
}

func SetPreferredTerminal(termID string) error {
	cfg, _ := LoadConfig()
	cfg.DefaultTerminal = termID
	return SaveConfig(cfg)
}

func GetGhosttyMode() string {
	cfg, err := LoadConfig()
	if err != nil || cfg.GhosttyMode == "" {
		return "smart"
	}
	return cfg.GhosttyMode
}

func SetGhosttyMode(mode string) error {
	cfg, _ := LoadConfig()
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode != "window" && mode != "tab" {
		mode = "smart"
	}
	cfg.GhosttyMode = mode
	return SaveConfig(cfg)
}

func GetSortBy() string {
	cfg, err := LoadConfig()
	if err != nil || cfg.SortBy == "" {
		return "recent"
	}
	return cfg.SortBy
}

func SetSortBy(sortBy string) error {
	cfg, _ := LoadConfig()
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "name", "name-asc", "az":
		cfg.SortBy = "name"
	case "name-desc", "za":
		cfg.SortBy = "name-desc"
	case "added", "file":
		cfg.SortBy = "added"
	default:
		cfg.SortBy = "recent"
	}
	return SaveConfig(cfg)
}
