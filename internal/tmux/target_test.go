package tmux

import "testing"

func TestTmuxTarget_String(t *testing.T) {
	tests := []struct {
		name   string
		target TmuxTarget
		want   string
	}{
		{"session only", TmuxTarget{Session: "gt-witness"}, "gt-witness"},
		{"session and window", TmuxTarget{Session: "gt-rig", Window: "witness"}, "gt-rig:witness"},
		{"polecat window", TmuxTarget{Session: "gt-rig", Window: "Toast"}, "gt-rig:Toast"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.String()
			if got != tt.want {
				t.Errorf("TmuxTarget.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTmuxTarget_IsWindow(t *testing.T) {
	tests := []struct {
		name   string
		target TmuxTarget
		want   bool
	}{
		{"session only", TmuxTarget{Session: "gt-witness"}, false},
		{"with window", TmuxTarget{Session: "gt-rig", Window: "witness"}, true},
		{"empty window", TmuxTarget{Session: "gt-rig", Window: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.IsWindow()
			if got != tt.want {
				t.Errorf("TmuxTarget.IsWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
