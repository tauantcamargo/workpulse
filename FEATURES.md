# Features

## Current Features

### Timer Modes
- **Work** (25 min default) - Focused work session
- **Short Break** (5 min default) - Quick rest between sessions
- **Long Break** (15 min default) - Extended rest after multiple sessions
- **Walk** (10 min default) - Movement break
- **Water** (2 min default) - Hydration reminder
- **Video** (20 min default) - Video break time

### Custom Durations
Set custom timer lengths via CLI flags:
```bash
workpulse --work 30m --short-break 10m --long-break 20m
workpulse --walk 15m --water 3m --video 25m
```

### Auto-Advance
Automatically start the next timer mode after completion:
```bash
workpulse --auto
```
- Work completes → Short Break starts automatically
- Short Break completes → Work starts automatically
- 2-second delay before advancing

### Session Counter (Pomodoro Technique)
- Displays "Pomodoro X/4" during work and break sessions
- After 4 work sessions, suggests Long Break instead of Short Break
- Counter resets after completing a Long Break
- Traditional Pomodoro workflow built-in

### Activity Tracking
- Prompt for activity name when starting work sessions
- Pre-set activity via CLI: `workpulse -a "task name"`
- Activities saved to `~/.workpulse/activities.json`

### Statistics Summary
- Press `d` to view summary (cycles through Daily → Weekly → Monthly)
- Shows time spent per mode
- Displays total focused time
- Press `d` or `tab` while in summary to switch periods
- Press `esc` to return to timer

### Configuration
Persistent configuration stored at `~/.workpulse/config.json`:
```bash
# View current configuration
workpulse config --list

# Set durations
workpulse config --work 30m --short-break 10m
workpulse config --long-break 20m --walk 15m

# Toggle settings
workpulse config --sound=false    # Disable sound
workpulse config --notify=true    # Enable notifications
workpulse config --auto=true      # Enable auto-advance

# Set pomodoros before long break
workpulse config --pomodoros 6

# Reset to defaults
workpulse config --reset
```

### Export Sessions
Export session history to CSV for external analysis:
```bash
# Export all sessions to stdout
workpulse export

# Export to file
workpulse export -o sessions.csv

# Filter by period
workpulse export --period today
workpulse export --period week
workpulse export --period month
workpulse export --period all    # default
```

CSV columns: `mode`, `activity`, `started`, `duration_minutes`, `completed`

### Streaks and Goals
Track daily focus time goals and consecutive day streaks:
```bash
# View current stats
workpulse stats

# Set daily focus time goal
workpulse config --daily-goal 2h
workpulse config --daily-goal 90m
```

Stats include:
- Today's focus time and goal progress
- Current streak (consecutive days with work sessions)
- Longest streak achieved

### Notifications
- Desktop notifications when timer completes
- Sound alert on macOS (Glass.aiff)
- Cross-platform beep via beeep library

### Visual
- Gradient progress bar (purple → blue)
- Color-coded modes
- Centered terminal UI with rounded borders
- Running/paused status indicators

### Themes
Customize the color scheme with built-in themes:
```bash
# Set theme
workpulse config --theme dark      # Default dark theme
workpulse config --theme light     # Light theme
workpulse config --theme dracula   # Dracula color scheme
workpulse config --theme nord      # Nord color scheme
```

### Keyboard Controls
| Key | Action |
|-----|--------|
| `s` | Start timer |
| `p` | Pause timer |
| `r` | Reset timer |
| `n` | Next mode |
| `d` | Stats summary (cycles periods) |
| `q` | Quit |
| `1-6` | Quick switch modes |
| `w/b/l` | Work/Break/Long break |

---

## Roadmap

### Planned
- [x] Quiet mode (disable sound) - `workpulse config --sound=false`
- [x] Export sessions to CSV - `workpulse export`
- [x] Streaks and goals - `workpulse stats`
- [x] Themes (dark/light/custom) - `workpulse config --theme <name>`
- [ ] Custom notification sounds
