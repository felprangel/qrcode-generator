// Package config loads and persists user settings for qrcode-generator.
package config

import (
	"errors"
	"os"
	"path/filepath"
)

// Config holds user-configurable settings.
type Config struct {
	URLPrefix string
}

// Load reads the config file from the user's XDG config directory.
// If the file does not exist, Load returns a zero-value Config and nil error.
func Load() (*Config, error) {
	p := path()
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	_ = data // parsing comes in Task 4
	return &Config{}, nil
}

// path returns the absolute path to the config file, honoring XDG_CONFIG_HOME.
func path() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "qrcode-generator", "config")
}
