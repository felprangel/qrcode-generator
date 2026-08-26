package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoArgs_PrintsUsageAndReturns1(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(nil, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "usage") {
		t.Errorf("stderr = %q, want substring \"usage\"", stderr.String())
	}
}

func TestRun_ArbitraryText_IsAccepted(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"https://example.com", "hello world"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "--- https://example.com ---") || !strings.Contains(out, "--- hello world ---") {
		t.Errorf("stdout missing per-content headers. got:\n%s", out)
	}
}

func TestRun_AllNumericArgs_ReturnsZeroAndPrintsHeaders(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"1", "2"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "--- 1 ---") || !strings.Contains(out, "--- 2 ---") {
		t.Errorf("stdout missing per-number headers. got:\n%s", out)
	}
}

func TestRun_ClearShortFlag_EmitsClearSequence(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"-c", "1"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), clearScreen) {
		t.Errorf("stdout should start with clear sequence")
	}
}

func TestRun_ClearLongFlag_EmitsClearSequence(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"1", "--clear"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), clearScreen) {
		t.Errorf("stdout should start with clear sequence")
	}
}

func TestRun_NoClearFlag_DoesNotEmitClearSequence(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"1"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
	if strings.Contains(stdout.String(), clearScreen) {
		t.Errorf("stdout should not contain clear sequence")
	}
}

func TestRun_ClearFlagOnlyNoNumbers_Returns1(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"-c"}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func TestRun_ConfigSet_PersistsValue(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"config", "set", "PREFIX=numero_do_zap="}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("config set exit code = %d, want 0. stderr=%q", code, stderr.String())
	}

	// Follow-up generate run should use the persisted prefix.
	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"42"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("generate exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "--- 42 ---") {
		t.Errorf("expected header for 42 in stdout:\n%s", stdout.String())
	}
}

func TestRun_ConfigSet_Malformed_Returns1(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"config", "set", "no_equals_here"}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}

func TestRun_Preset_SelectsPrefix(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	if code := Run([]string{"config", "set", "work=W-"}, &stdout, &stderr); code != 0 {
		t.Fatalf("config set exit = %d, stderr=%q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code := Run([]string{"-p", "work", "42"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
}

func TestRun_UnknownPreset_Returns1(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"-p", "nope", "42"}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "nope") {
		t.Errorf("stderr = %q, want it to name the bad preset", stderr.String())
	}
}

func TestRun_NoPreset_BypassesConfiguredDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	for _, args := range [][]string{
		{"config", "set", "work=W-"},
		{"config", "default", "work"},
	} {
		if code := Run(args, &stdout, &stderr); code != 0 {
			t.Fatalf("Run(%v) exit = %d, stderr=%q", args, code, stderr.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code := Run([]string{"--no-preset", "raw"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0. stderr=%q", code, stderr.String())
	}
	// --no-preset must not error even with a default configured, and no prefix
	// is applied (the QR encodes "raw" verbatim).
	if !strings.Contains(stdout.String(), "--- raw ---") {
		t.Errorf("stdout missing header. got:\n%s", stdout.String())
	}
}

func TestRun_ConfigDefault_SetsDefaultPreset(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	for _, args := range [][]string{
		{"config", "set", "work=W-"},
		{"config", "default", "work"},
	} {
		if code := Run(args, &stdout, &stderr); code != 0 {
			t.Fatalf("Run(%v) exit = %d, stderr=%q", args, code, stderr.String())
		}
	}

	// No -p flag should now resolve to the "work" preset without error.
	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"42"}, &stdout, &stderr); code != 0 {
		t.Errorf("generate with configured default exit = %d, stderr=%q", code, stderr.String())
	}
}

func TestRun_ConfigList_SortsPresetsAndMarksDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	for _, args := range [][]string{
		{"config", "set", "work=W-"},
		{"config", "set", "home=H-"},
		{"config", "default", "work"},
	} {
		if code := Run(args, &stdout, &stderr); code != 0 {
			t.Fatalf("Run(%v) exit = %d, stderr=%q", args, code, stderr.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"config", "list"}, &stdout, &stderr); code != 0 {
		t.Fatalf("config list exit = %d, stderr=%q", code, stderr.String())
	}

	// Sorted, default marked, and the reserved @default key itself hidden.
	want := "home=H-\nwork=W-  (default)\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}
