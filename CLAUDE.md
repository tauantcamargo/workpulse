# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Build
go build -o workpulse

# Run
./workpulse
./workpulse --activity "task name"   # Pre-set activity
./workpulse -a "task name"           # Shorthand

# Dependencies
go mod tidy
```

## Architecture

WorkPulse is a terminal Pomodoro timer built with the Elm architecture via Bubble Tea.

### Core Packages

- **internal/app/** - Bubble Tea application (Model/Update/View pattern)
  - `model.go`: Immutable state struct with setter methods that return new copies
  - `update.go`: Message handlers, keyboard input, timer ticks, notifications
  - `view.go`: Renders UI based on ViewState (ViewTimer, ViewSummary, ViewActivityInput)

- **internal/timer/** - Timer logic
  - `timer.go`: Immutable Timer struct with Start/Pause/Reset/Tick methods
  - `modes.go`: Mode definitions (work, short_break, long_break, walk, water, video)

- **internal/ui/** - Lip Gloss styling
  - `styles.go`: Color palette and reusable styles
  - `progress.go`: Gradient progress bar (purple→blue interpolation)
  - `components.go`: Header, timer display, help text renderers

- **internal/storage/** - JSON persistence to `~/.workpulse/activities.json`

- **internal/config/** - User preferences at `~/.workpulse/config.json`

### Key Patterns

**Immutability**: All state updates return new structs rather than mutating:
```go
func (m Model) SetTimer(t timer.Timer) Model {
    return Model{Timer: t, ...}  // Returns new Model
}
```

**Bubble Tea Messages**: `timer.TickMsg` drives the countdown; `NotifyMsg` signals notification completion.

**Sound on Completion**: Uses `beeep.Beep()` plus macOS `afplay` for the Glass.aiff system sound.
