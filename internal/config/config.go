// Package config loads and persists user settings for qrcode-generator.
package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds user-configurable settings.
type Config struct {
	URLPrefix string
}

// known is the whitelist of recognized config keys.
var known = map[string]struct{}{
	"URL_PREFIX": {},
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
		return nil, fmt.Errorf("config: read %s: %w", p, err)
	}

	cfg := &Config{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		eq := strings.IndexByte(trimmed, '=')
		if eq <= 0 {
			return nil, fmt.Errorf("config: line %d: missing '=' separator", lineNo)
		}
		key := strings.TrimSpace(trimmed[:eq])
		value := trimmed[eq+1:]
		if _, ok := known[key]; !ok {
			continue
		}
		switch key {
		case "URL_PREFIX":
			cfg.URLPrefix = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("config: scan %s: %w", p, err)
	}
	return cfg, nil
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
