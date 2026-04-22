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
	Prefix string
}

// known is the whitelist of recognized config keys.
var known = map[string]struct{}{
	"PREFIX": {},
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
		case "PREFIX":
			cfg.Prefix = value
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

// Set writes or updates a single key in the config file.
// The key must be in the whitelist of known keys.
// Other keys, comments, and blank lines in the existing file are preserved.
func Set(key, value string) error {
	if _, ok := known[key]; !ok {
		return fmt.Errorf("config: unknown key %q", key)
	}

	p := path()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("config: mkdir: %w", err)
	}

	existing, err := os.ReadFile(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config: read: %w", err)
	}

	var out bytes.Buffer
	replaced := false
	prefix := key + "="
	scanner := bufio.NewScanner(bytes.NewReader(existing))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			out.WriteString(prefix + value + "\n")
			replaced = true
			continue
		}
		out.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("config: scan existing: %w", err)
	}
	if !replaced {
		out.WriteString(prefix + value + "\n")
	}

	// Atomic write: tmp file in same dir, then rename.
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("config: write tmp: %w", err)
	}
	if err := os.Rename(tmp, p); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("config: rename: %w", err)
	}
	return nil
}
