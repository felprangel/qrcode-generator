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
	if len(cfg) != 0 {
		t.Errorf("presets = %v, want empty", cfg)
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

func TestLoad_ParsesPresets(t *testing.T) {
	writeConfig(t, "PREFIX=numero_do_zap=\nwork=https://work/\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg["prefix"] != "numero_do_zap=" {
		t.Errorf("prefix = %q, want %q", cfg["prefix"], "numero_do_zap=")
	}
	if cfg["work"] != "https://work/" {
		t.Errorf("work = %q, want %q", cfg["work"], "https://work/")
	}
}

// Backward compat: an old PREFIX= config resolves as the default preset.
func TestResolve_LegacyPrefixIsDefault(t *testing.T) {
	writeConfig(t, "PREFIX=legacy_\n")

	cfg, _ := Load()
	got, ok := Resolve(cfg, "")
	if !ok || got != "legacy_" {
		t.Errorf("Resolve(default) = %q, %v; want %q, true", got, ok, "legacy_")
	}
}

func TestResolve_ConfiguredDefaultWins(t *testing.T) {
	writeConfig(t, "PREFIX=legacy_\nwork=work_\n@default=work\n")

	cfg, _ := Load()
	got, ok := Resolve(cfg, "")
	if !ok || got != "work_" {
		t.Errorf("Resolve(default) = %q, %v; want %q, true", got, ok, "work_")
	}
}

func TestResolve_UnknownPreset(t *testing.T) {
	writeConfig(t, "work=work_\n")

	cfg, _ := Load()
	if _, ok := Resolve(cfg, "nope"); ok {
		t.Error("Resolve(nope) ok = true, want false")
	}
}

func TestLoad_IgnoresBlanksAndComments(t *testing.T) {
	writeConfig(t, "# a comment\n\n   # indented comment\nprefix=foo\n\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg["prefix"] != "foo" {
		t.Errorf("prefix = %q, want %q", cfg["prefix"], "foo")
	}
}

func TestLoad_MalformedLine_ReturnsErrorWithLineNumber(t *testing.T) {
	writeConfig(t, "prefix=ok\nno_equals_sign_here\n")

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

	if err := Set("work", "numero_do_zap="); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load after Set error: %v", err)
	}
	if cfg["work"] != "numero_do_zap=" {
		t.Errorf("work = %q, want %q", cfg["work"], "numero_do_zap=")
	}
}

func TestSet_UpdatesExistingKeyWithoutDuplicating(t *testing.T) {
	writeConfig(t, "prefix=old\n")

	if err := Set("prefix", "new"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	data, err := os.ReadFile(path())
	if err != nil {
		t.Fatal(err)
	}
	count := strings.Count(string(data), "prefix=")
	if count != 1 {
		t.Errorf("prefix appears %d times, want 1. File:\n%s", count, data)
	}
	cfg, _ := Load()
	if cfg["prefix"] != "new" {
		t.Errorf("prefix = %q, want %q", cfg["prefix"], "new")
	}
}

func TestSet_PreservesCommentsAndBlanks(t *testing.T) {
	writeConfig(t, "# my config\n\nprefix=old\n# trailing note\n")

	if err := Set("prefix", "new"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	data, err := os.ReadFile(path())
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for _, want := range []string{"# my config", "# trailing note", "prefix=new"} {
		if !strings.Contains(s, want) {
			t.Errorf("file missing %q. Contents:\n%s", want, s)
		}
	}
}

func TestSet_RejectsEmptyName(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	if err := Set("   ", "value"); err == nil {
		t.Fatal("Set: expected error for empty name, got nil")
	}
}
