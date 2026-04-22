package config

import (
	"os"
	"path/filepath"
	"strings"
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
	if cfg.Prefix != "" {
		t.Errorf("Prefix = %q, want empty", cfg.Prefix)
	}
}

func writeConfig(t *testing.T, contents string) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, "qrcode-generator")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(cfgDir, "config")
	if err := os.WriteFile(p, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidFile_ParsesPrefix(t *testing.T) {
	writeConfig(t, "PREFIX=numero_do_zap=\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Prefix != "numero_do_zap=" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "numero_do_zap=")
	}
}

func TestLoad_IgnoresBlanksAndComments(t *testing.T) {
	writeConfig(t, "# a comment\n\n   # indented comment\nPREFIX=foo\n\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Prefix != "foo" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "foo")
	}
}

func TestLoad_UnknownKey_SilentlySkipped(t *testing.T) {
	writeConfig(t, "PREFIX=x\nUNKNOWN_KEY=y\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Prefix != "x" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "x")
	}
}

func TestLoad_MalformedLine_ReturnsErrorWithLineNumber(t *testing.T) {
	writeConfig(t, "PREFIX=ok\nno_equals_sign_here\n")

	_, err := Load()
	if err == nil {
		t.Fatal("Load: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error %q does not mention line 2", err)
	}
}

func TestSet_CreatesFileAndDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	if err := Set("PREFIX", "numero_do_zap="); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load after Set error: %v", err)
	}
	if cfg.Prefix != "numero_do_zap=" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "numero_do_zap=")
	}
}

func TestSet_UpdatesExistingKeyWithoutDuplicating(t *testing.T) {
	writeConfig(t, "PREFIX=old\n")

	if err := Set("PREFIX", "new"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	data, err := os.ReadFile(path())
	if err != nil {
		t.Fatal(err)
	}
	count := strings.Count(string(data), "PREFIX=")
	if count != 1 {
		t.Errorf("PREFIX appears %d times, want 1. File:\n%s", count, data)
	}
	cfg, _ := Load()
	if cfg.Prefix != "new" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "new")
	}
}

func TestSet_PreservesCommentsAndBlanks(t *testing.T) {
	writeConfig(t, "# my config\n\nPREFIX=old\n# trailing note\n")

	if err := Set("PREFIX", "new"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	data, err := os.ReadFile(path())
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for _, want := range []string{"# my config", "# trailing note", "PREFIX=new"} {
		if !strings.Contains(s, want) {
			t.Errorf("file missing %q. Contents:\n%s", want, s)
		}
	}
}

func TestSet_RejectsUnknownKey(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	err := Set("BOGUS", "value")
	if err == nil {
		t.Fatal("Set: expected error for unknown key, got nil")
	}
}
