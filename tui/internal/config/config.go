package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/odiador/foltty/internal/terminal"
)

type Config struct {
	DefaultTerminal string `json:"defaultTerminal"`
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
