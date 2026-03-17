package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type NotificationStyle string

const (
	NotifySystem NotificationStyle = "system"
	NotifyPopup  NotificationStyle = "popup"
	NotifyBoth   NotificationStyle = "both"
)

type Config struct {
	IntervalMinutes   int               `json:"interval_minutes"`
	NotificationStyle NotificationStyle `json:"notification_style"`
	LaunchAtLogin     bool              `json:"launch_at_login"`
	SoundEnabled      bool              `json:"sound_enabled"`
	BreakDurationSec  int               `json:"break_duration_sec"`
}

func DefaultConfig() Config {
	return Config{
		IntervalMinutes:   30,
		NotificationStyle: NotifyBoth,
		LaunchAtLogin:     false,
		SoundEnabled:      true,
		BreakDurationSec:  300,
	}
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".breakreminder")
	return dir, os.MkdirAll(dir, 0755)
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			_ = Save(cfg)
			return cfg, nil
		}
		return DefaultConfig(), err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	// Validate and apply defaults for zero values
	if cfg.IntervalMinutes <= 0 {
		cfg.IntervalMinutes = 30
	}
	if cfg.BreakDurationSec <= 0 {
		cfg.BreakDurationSec = 300
	}
	if cfg.NotificationStyle == "" {
		cfg.NotificationStyle = NotifyBoth
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
