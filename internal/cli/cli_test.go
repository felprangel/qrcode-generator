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

func TestRun_NonNumericArg_Returns1AndNamesBadArg(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"1", "abc", "3"}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "abc") {
		t.Errorf("stderr = %q, expected to name bad arg \"abc\"", stderr.String())
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

func TestRun_ConfigSet_UnknownKey_Returns1(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout, stderr bytes.Buffer
	code := Run([]string{"config", "set", "BOGUS=x"}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
}
