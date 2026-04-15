// Package session provides polecat session lifecycle management.
package session

import (
	"fmt"
)

// DefaultPrefix is the default beads prefix used when no rig-specific prefix is known.
const DefaultPrefix = "gt"

// HQPrefix is the prefix for town-level services (Mayor, Deacon).
const HQPrefix = "hq-"

// MayorSessionName returns the session name for the Mayor agent.
// One mayor per machine - multi-town requires containers/VMs for isolation.
func MayorSessionName() string {
	return HQPrefix + "mayor"
}

// DeaconSessionName returns the session name for the Deacon agent.
// One deacon per machine - multi-town requires containers/VMs for isolation.
func DeaconSessionName() string {
	return HQPrefix + "deacon"
}

// WitnessSessionName returns the session name for a rig's Witness agent.
// rigPrefix is the rig's beads prefix (e.g., "gt" for gastown, "bd" for beads).
func WitnessSessionName(rigPrefix string) string {
	return fmt.Sprintf("%s-witness", rigPrefix)
}

// RefinerySessionName returns the session name for a rig's Refinery agent.
// rigPrefix is the rig's beads prefix (e.g., "gt" for gastown, "bd" for beads).
func RefinerySessionName(rigPrefix string) string {
	return fmt.Sprintf("%s-refinery", rigPrefix)
}

// CrewSessionName returns the session name for a crew worker in a rig.
// rigPrefix is the rig's beads prefix (e.g., "gt" for gastown, "bd" for beads).
func CrewSessionName(rigPrefix, name string) string {
	return fmt.Sprintf("%s-crew-%s", rigPrefix, name)
}

// PolecatSessionName returns the session name for a polecat in a rig.
// rigPrefix is the rig's beads prefix (e.g., "gt" for gastown, "bd" for beads).
func PolecatSessionName(rigPrefix, name string) string {
	return fmt.Sprintf("%s-%s", rigPrefix, name)
}

// OverseerSessionName returns the session name for the human operator.
// The overseer is the human who controls Gas Town, not an AI agent.
func OverseerSessionName() string {
	return HQPrefix + "overseer"
}

// BootSessionName returns the session name for the Boot watchdog.
// Boot is town-level (launched by deacon), so it uses the hq- prefix.
// "hq-boot" avoids tmux prefix-matching collisions with "hq-deacon".
func BootSessionName() string {
	return HQPrefix + "boot"
}

// DogSessionName returns the session name for a named dog agent.
// Dogs are town-level (managed by deacon), so they use the hq- prefix.
// Pattern: hq-dog-<name> (e.g., hq-dog-alpha).
func DogSessionName(name string) string {
	return fmt.Sprintf("%sdog-%s", HQPrefix, name)
}

// RigSessionName returns the tmux session name for a rig in window mode.
// All rig agents share this single session, each in a named window.
// rigPrefix is the rig's beads prefix (e.g., "gt" for gastown).
// Pattern: <prefix>-rig (e.g., "gt-rig").
func RigSessionName(rigPrefix string) string {
	return fmt.Sprintf("%s-rig", rigPrefix)
}

// WindowName returns the tmux window name for an agent in window mode.
// For named agents (polecats, crew): returns the agent name (e.g., "Toast", "jack").
// For singleton agents (witness, refinery): returns the role (e.g., "witness").
func WindowName(role, name string) string {
	if name != "" {
		return name
	}
	return role
}
