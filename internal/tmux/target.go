package tmux

import "fmt"

// TmuxTarget identifies a tmux session or a window within a session.
// In session mode (default), only Session is set.
// In window mode, both Session and Window are set.
type TmuxTarget struct {
	Session string
	Window  string
}

// String returns the tmux target string.
// For session-only targets: "session-name"
// For window targets: "session-name:window-name"
func (t TmuxTarget) String() string {
	if t.Window != "" {
		return fmt.Sprintf("%s:%s", t.Session, t.Window)
	}
	return t.Session
}

// IsWindow returns true if this target refers to a window within a session.
func (t TmuxTarget) IsWindow() bool {
	return t.Window != ""
}
