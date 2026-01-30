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

### Activity Tracking
- Prompt for activity name when starting work sessions
- Pre-set activity via CLI: `workpulse -a "task name"`
- Activities saved to `~/.workpulse/activities.json`

### Daily Summary
- Press `d` to view today's completed sessions
- Shows time spent per mode
- Displays total focused time

### Notifications
- Desktop notifications when timer completes
- Sound alert on macOS (Glass.aiff)
- Cross-platform beep via beeep library

### Visual
- Gradient progress bar (purple → blue)
- Color-coded modes
- Centered terminal UI with rounded borders
- Running/paused status indicators

### Keyboard Controls
| Key | Action |
|-----|--------|
| `s` | Start timer |
| `p` | Pause timer |
| `r` | Reset timer |
| `n` | Next mode |
| `d` | Daily summary |
| `q` | Quit |
| `1-6` | Quick switch modes |
| `w/b/l` | Work/Break/Long break |

---

## Roadmap

### In Progress
- [ ] Auto-advance to next mode after completion
- [ ] Session counter (Pomodoro X/4)
- [ ] Weekly and monthly stats view

### Planned
- [ ] Quiet mode (disable sound)
- [ ] Custom notification sounds
- [ ] Export sessions to CSV
- [ ] Themes (dark/light/custom)
- [ ] Streaks and goals
