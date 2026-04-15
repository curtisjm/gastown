# Window Mode: tmux windows-per-rig instead of sessions-per-agent

## Status: Shelved (2026-04-15)

This branch contains a complete but unmerged implementation of an alternative
tmux topology where each rig gets one tmux session with agents as named windows,
instead of the default one-session-per-agent model.

## Motivation

With many rigs and agents, the default model creates tmux session sprawl. A
gastown rig with a witness, refinery, and 8 polecats produces 10 separate tmux
sessions. Window mode groups them into a single session with named windows,
giving a cleaner `tmux ls` and tab-based navigation within a rig.

## Why it was shelved

The complexity cost outweighs the cosmetic benefit:

- **Two permanent code paths.** Every manager (witness, refinery, crew, polecat,
  daemon) has `if windowMode { ... } else { ... }` branching. Every future
  lifecycle change must consider both modes.
- **Fighting tmux's architecture.** tmux environment variables are
  session-scoped only — no window-scoped env vars. Per-agent identity vars
  (GT_ROLE, BD_ACTOR, etc.) must be baked into the command string via
  `exec env VAR=val ...` (the PrependEnv pattern). This is a workaround, not
  a natural fit.
- **Feature surface area.** GT_PANE_ID, pane-died hooks, status bars, health
  checks, ghost detection, cycle bindings, and agent listing all need separate
  window-mode implementations.
- **~2000+ lines of branching logic** across 15+ files to solve what is
  primarily a UI/cosmetic concern.

A lighter alternative (e.g., a `gt dashboard` aggregation view, or better
grouping in `gt agents`) could address the sprawl without the maintenance
burden.

## Design decisions

| # | Decision | Rationale |
|---|----------|-----------|
| D1 | Rig session lifecycle: create-on-first-agent, destroy-on-last | No empty sessions cluttering listings |
| D2 | `TmuxTarget` type encapsulates session or session:window | Avoids rewriting 53 function signatures |
| D3 | Town-level agents (mayor, deacon, boot, dogs) stay as sessions | Singletons, not rig-bound |
| D4 | Polecats are windows in the rig session | Always bound to one rig |
| D5 | Env vars via `exec env` command prefix | tmux has no window-scoped env vars |
| D6 | Town-level setting: `window_mode` in TownSettings | One topology active at a time across all rigs |

## Implementation phases (all complete on this branch)

| Phase | Commit | Description |
|-------|--------|-------------|
| 0 | e95125a | Foundation: TmuxTarget type, TownSettings.WindowMode, naming functions, IsWindowMode helper, gt config command |
| 1 | 787fdea | Tmux primitives: HasWindow, NewWindowWithCommand, KillWindow, ListWindows, WindowCount, RenameWindow |
| 2 | 7492e67 | Witness: Target(), branch Start/IsRunning/Stop, AgentEnvSplit, daemon ensureWitnessRunning |
| 3 | f6d8da4 | Refinery + crew: generalize lifecycle.go SessionConfig, branch both managers |
| 4 | 8d7535b | Polecat: session_manager.go Start/Stop/IsRunning/Attach/Capture/Inject/List branches |
| 5 | c6619e9 | Theming: ConfigureWindowModeSession, ApplyWindowTheme, per-window tab bar styling |
| 6 | 5663968 | Keybindings: GT_WINDOW_MODE env var, cycle bindings fall through to native next/prev-window |
| 7 | 64ed386 | Daemon: ghost killer, health checks, idle reaper, per-window pane-died/auto-respawn hooks |
| 8 | be5f7ef | Agent listing: window discovery in gt agents, select-window in menus |
| 9 | a14349b | Bulk ops: gt rig stop/start window mode, mode switch guard |

## Key files changed

- `internal/tmux/target.go` — TmuxTarget type (new)
- `internal/tmux/tmux.go` — Window lifecycle methods, theming, hooks, bindings
- `internal/config/types.go` — WindowMode field in TownSettings
- `internal/config/loader.go` — IsWindowMode() helper
- `internal/config/env.go` — AgentEnvSplit, SplitAgentEnv
- `internal/session/names.go` — RigSessionName(), WindowName()
- `internal/witness/manager.go` — Target(), window-mode branches
- `internal/refinery/manager.go` — Target(), window-mode branches
- `internal/crew/manager.go` — Target(), window-mode branches
- `internal/polecat/session_manager.go` — Target(), ensureRigSession(), all operation branches
- `internal/session/lifecycle.go` — WindowMode in SessionConfig
- `internal/daemon/daemon.go` — Ghost killer, health checks, idle reaper branches
- `internal/cmd/agents.go` — Window discovery, menu targeting
- `internal/cmd/rig.go` — Bulk stop/start, mode switch guard
- `internal/cmd/config.go` — gt config set window_mode

## How to activate (if resuming work)

```bash
gt config set window_mode true    # enable for all rigs
gt rig stop <rig> && gt rig start <rig>   # restart with new topology
```

## Testing notes

Each phase includes tests. Pre-existing test failures (TestCrossPlatformBuild,
BuiltProperly ldflags) are unrelated to window mode. Run:

```bash
cd rig
go test ./internal/tmux/... ./internal/polecat/... ./internal/witness/... \
  ./internal/config/... ./internal/session/... ./internal/cmd/... -count=1
```
