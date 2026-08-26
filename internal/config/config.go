// Package config loads and persists user presets for qrcode-generator.
// A preset is a named prefix: NAME=VALUE per line. Names are case-insensitive.
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

// DefaultKey is the reserved config key naming which preset to use when no
// preset is requested. legacyDefault keeps old PREFIX= configs working.
const DefaultKey = "@default"

const legacyDefault = "prefix"

// Resolve returns the prefix for the given preset name. When name is empty it
// falls back to the configured default preset, then the legacy "prefix" preset.
// ok reports whether the resolved preset actually exists.
func Resolve(presets map[string]string, name string) (prefix string, ok bool) {
	if name == "" {
		name = presets[DefaultKey]
	}

	if name == "" {
		name = legacyDefault
	}

	prefix, ok = presets[strings.ToLower(name)]

	return prefix, ok
}

// Load reads presets (name -> prefix) from the user's XDG config directory.
// If the file does not exist, Load returns an empty map and nil error.
func Load() (map[string]string, error) {
	path := path()

	data, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	presets := map[string]string{}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		trimmed := strings.TrimSpace(scanner.Text())

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		equalIndex := strings.IndexByte(trimmed, '=')

		if equalIndex <= 0 {
			return nil, fmt.Errorf("config: line %d: missing '=' separator", lineNumber)
		}

		name := strings.ToLower(strings.TrimSpace(trimmed[:equalIndex]))
		presets[name] = trimmed[equalIndex+1:]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("config: scan %s: %w", path, err)
	}

	return presets, nil
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

// Set writes or updates a single preset in the config file.
// Comments and blank lines in the existing file are preserved.
func Set(name, value string) error {
	name = strings.ToLower(strings.TrimSpace(name))

	if name == "" {
		return errors.New("config: preset name must not be empty")
	}

	path := path()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("config: mkdir: %w", err)
	}

	existing, err := os.ReadFile(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config: read: %w", err)
	}

	var out bytes.Buffer

	replaced := false
	prefix := name + "="

	scanner := bufio.NewScanner(bytes.NewReader(existing))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.ToLower(strings.TrimSpace(line))

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
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, out.Bytes(), 0o644); err != nil {
		return fmt.Errorf("config: write tmp: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)

		return fmt.Errorf("config: rename: %w", err)
	}

	return nil
}
