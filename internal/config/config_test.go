package config

import (
	"testing"
)

func TestLoad_NoFile_ReturnsEmptyConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load returned nil config")
	}
	if cfg.URLPrefix != "" {
		t.Errorf("URLPrefix = %q, want empty", cfg.URLPrefix)
	}
}
