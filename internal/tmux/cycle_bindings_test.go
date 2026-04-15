package tmux

import (
	"strings"
	"testing"
)

// TestIsGTBindingCurrent_DetectsStalePattern verifies that isGTBindingCurrent
// returns false when the baked-in pattern doesn't match the current pattern.
// This is the core of the gt rig add fix: after adding a rig, the prefix
// pattern changes and existing bindings become stale.
func TestIsGTBindingCurrent_DetectsStalePattern(t *testing.T) {
	tm := newTestTmux(t)

	session := "gt-test-stale-" + t.Name()
	_ = tm.KillSession(session)
	defer func() { _ = tm.KillSession(session) }()

	if err := tm.NewSessionWithCommand(session, "", "sleep 30"); err != nil {
		t.Fatalf("session creation: %v", err)
	}

	// Install a binding with an OLD pattern (missing a hypothetical "qu" prefix)
	oldPattern := "^(gt|hq)-"
	oldIfShell := "echo '#{session_name}' | grep -Eq '" + oldPattern + "'"
	if _, err := tm.run("bind-key", "-T", "prefix", "n",
		"if-shell", oldIfShell,
		"run-shell 'gt cycle next --session #{session_name} --client #{client_tty}'",
		"next-window"); err != nil {
		t.Fatalf("installing old binding: %v", err)
	}

	// Verify the binding has --client (so isGTBindingWithClient returns true)
	if !tm.isGTBindingWithClient("prefix", "n") {
		t.Fatal("expected isGTBindingWithClient to return true for the installed binding")
	}

	// But the pattern is stale — a new pattern with "qu" should not match
	newPattern := "^(gt|hq|qu)-"
	if tm.isGTBindingCurrent("prefix", "n", newPattern) {
		t.Error("expected isGTBindingCurrent to return false for stale pattern")
	}

	// The old pattern should still match
	if !tm.isGTBindingCurrent("prefix", "n", oldPattern) {
		t.Error("expected isGTBindingCurrent to return true for matching pattern")
	}
}

// TestSetCycleBindings_RefreshesStalePattern verifies that SetCycleBindings
// re-binds when the existing binding has a stale prefix pattern, even though
// it already has --client support.
func TestSetCycleBindings_RefreshesStalePattern(t *testing.T) {
	tm := newTestTmux(t)

	session := "gt-test-refresh-" + t.Name()
	_ = tm.KillSession(session)
	defer func() { _ = tm.KillSession(session) }()

	if err := tm.NewSessionWithCommand(session, "", "sleep 30"); err != nil {
		t.Fatalf("session creation: %v", err)
	}

	// Install a binding with a STALE pattern (only gt|hq, missing other prefixes)
	stalePattern := "^(gt|hq)-"
	staleIfShell := "echo '#{session_name}' | grep -Eq '" + stalePattern + "'"
	if _, err := tm.run("bind-key", "-T", "prefix", "n",
		"if-shell", staleIfShell,
		"run-shell 'gt cycle next --session #{session_name} --client #{client_tty}'",
		"next-window"); err != nil {
		t.Fatalf("installing stale binding: %v", err)
	}
	if _, err := tm.run("bind-key", "-T", "prefix", "p",
		"if-shell", staleIfShell,
		"run-shell 'gt cycle prev --session #{session_name} --client #{client_tty}'",
		"previous-window"); err != nil {
		t.Fatalf("installing stale binding for p: %v", err)
	}

	// Call SetCycleBindings — it should detect the stale pattern and re-bind
	if err := tm.SetCycleBindings(session); err != nil {
		t.Fatalf("SetCycleBindings: %v", err)
	}

	// Verify the binding was updated with the current pattern
	currentPattern := sessionPrefixPattern()
	output, err := tm.run("list-keys", "-T", "prefix", "n")
	if err != nil {
		t.Fatalf("listing keys: %v", err)
	}
	if !strings.Contains(output, currentPattern) {
		t.Errorf("expected binding to contain current pattern %q, got: %s", currentPattern, output)
	}
}

// TestSetCycleBindings_WindowModeCheck verifies that SetCycleBindings installs
// bindings that check GT_WINDOW_MODE and fall through to native window cycling.
func TestSetCycleBindings_WindowModeCheck(t *testing.T) {
	tm := newTestTmux(t)

	session := "gt-test-winmode-" + t.Name()
	_ = tm.KillSession(session)
	defer func() { _ = tm.KillSession(session) }()

	if err := tm.NewSessionWithCommand(session, "", "sleep 30"); err != nil {
		t.Fatalf("session creation: %v", err)
	}

	// Call SetCycleBindings to install fresh bindings
	if err := tm.SetCycleBindings(session); err != nil {
		t.Fatalf("SetCycleBindings: %v", err)
	}

	// Verify C-b n binding contains the GT_WINDOW_MODE check
	nOutput, err := tm.run("list-keys", "-T", "prefix", "n")
	if err != nil {
		t.Fatalf("listing n key: %v", err)
	}
	if !strings.Contains(nOutput, "GT_WINDOW_MODE") {
		t.Errorf("C-b n binding missing GT_WINDOW_MODE check, got: %s", nOutput)
	}
	if !strings.Contains(nOutput, "next-window") {
		t.Errorf("C-b n binding missing next-window fallback for window mode, got: %s", nOutput)
	}
	if !strings.Contains(nOutput, "gt cycle next") {
		t.Errorf("C-b n binding missing gt cycle next for session mode, got: %s", nOutput)
	}

	// Verify C-b p binding contains the GT_WINDOW_MODE check
	pOutput, err := tm.run("list-keys", "-T", "prefix", "p")
	if err != nil {
		t.Fatalf("listing p key: %v", err)
	}
	if !strings.Contains(pOutput, "GT_WINDOW_MODE") {
		t.Errorf("C-b p binding missing GT_WINDOW_MODE check, got: %s", pOutput)
	}
	if !strings.Contains(pOutput, "previous-window") {
		t.Errorf("C-b p binding missing previous-window fallback for window mode, got: %s", pOutput)
	}
	if !strings.Contains(pOutput, "gt cycle prev") {
		t.Errorf("C-b p binding missing gt cycle prev for session mode, got: %s", pOutput)
	}
}

// TestWindowModeCycleCmd verifies the generated shell command structure.
func TestWindowModeCycleCmd(t *testing.T) {
	cmd := windowModeCycleCmd("next", "next-window")
	if !strings.Contains(cmd, "GT_WINDOW_MODE") {
		t.Error("command missing GT_WINDOW_MODE check")
	}
	if !strings.Contains(cmd, "tmux next-window") {
		t.Error("command missing tmux next-window for window mode")
	}
	if !strings.Contains(cmd, "gt cycle next") {
		t.Error("command missing gt cycle next for session mode")
	}
	if !strings.HasPrefix(cmd, "run-shell '") {
		t.Error("command should start with run-shell")
	}

	cmd = windowModeCycleCmd("prev", "previous-window")
	if !strings.Contains(cmd, "tmux previous-window") {
		t.Error("command missing tmux previous-window for window mode")
	}
	if !strings.Contains(cmd, "gt cycle prev") {
		t.Error("command missing gt cycle prev for session mode")
	}
}
