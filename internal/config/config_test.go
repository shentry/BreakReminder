package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.IntervalMinutes != 30 {
		t.Errorf("expected IntervalMinutes=30, got %d", cfg.IntervalMinutes)
	}
	if cfg.NotificationStyle != NotifyBoth {
		t.Errorf("expected NotificationStyle=both, got %s", cfg.NotificationStyle)
	}
	if cfg.BreakDurationSec != 300 {
		t.Errorf("expected BreakDurationSec=300, got %d", cfg.BreakDurationSec)
	}
	if cfg.LaunchAtLogin != false {
		t.Error("expected LaunchAtLogin=false")
	}
	if cfg.SoundEnabled != true {
		t.Error("expected SoundEnabled=true")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Use a temp dir to avoid modifying real config
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := Config{
		IntervalMinutes:   15,
		NotificationStyle: NotifySystem,
		LaunchAtLogin:     true,
		SoundEnabled:      false,
		BreakDurationSec:  120,
	}

	err := Save(cfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	path := filepath.Join(tmpDir, ".breakreminder", "config.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file not created")
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.IntervalMinutes != 15 {
		t.Errorf("expected IntervalMinutes=15, got %d", loaded.IntervalMinutes)
	}
	if loaded.NotificationStyle != NotifySystem {
		t.Errorf("expected NotificationStyle=system, got %s", loaded.NotificationStyle)
	}
	if loaded.LaunchAtLogin != true {
		t.Error("expected LaunchAtLogin=true")
	}
	if loaded.SoundEnabled != false {
		t.Error("expected SoundEnabled=false")
	}
	if loaded.BreakDurationSec != 120 {
		t.Errorf("expected BreakDurationSec=120, got %d", loaded.BreakDurationSec)
	}
}

func TestLoadDefaultsOnMissing(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.IntervalMinutes != 30 {
		t.Errorf("expected default IntervalMinutes=30, got %d", cfg.IntervalMinutes)
	}
}
