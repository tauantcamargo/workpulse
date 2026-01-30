# WorkPulse

A minimal, beautiful terminal-based Pomodoro timer built in Go.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-macOS%20|%20Linux%20|%20Windows-lightgrey)

```
┌───────────────────────────────────────────────┐
│                                               │
│              🎯 WORK                          │
│                                               │
│       "Stay focused, you got this!"           │
│                                               │
│            ⏱️  14:32 / 25:00 ▶                │
│                                               │
│  ████████████████░░░░░░░░░░░░░░░░░░  58%     │
│                                               │
│   [s]tart  [p]ause  [r]eset  [n]ext  [q]uit  │
│                                               │
└───────────────────────────────────────────────┘
```

## Features

- **6 Timer Modes** — Work, Short Break, Long Break, Walk, Water, Video
- **Session Planner** — Set a total time goal with work/break ratios
- **Auto-Update** — Get notified when new versions are available
- **Settings UI** — Configure everything from within the app
- **Gradient Progress Bar** — Beautiful purple-to-blue color gradient
- **Desktop Notifications** — Get notified when timers complete
- **Sound Alerts** — Audio feedback on timer completion
- **Activity Tracking** — Log what you're working on
- **Daily/Weekly/Monthly Summary** — View your productivity stats
- **Export to CSV** — Export your session data
- **Themes** — Dark, Light, Dracula, Nord
- **Persistent Storage** — Sessions saved locally

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/tauantcamargo/workpulse/main/install.sh | sh
```

### Homebrew (macOS/Linux)

```bash
brew install tauantcamargo/tap/workpulse
```

### Download Binary

Download the latest binary for your platform from the [Releases](https://github.com/tauantcamargo/workpulse/releases) page.

### Go Install

```bash
go install github.com/tauantcamargo/workpulse@latest
```

### From Source

```bash
git clone https://github.com/tauantcamargo/workpulse.git
cd workpulse
go build -o workpulse
./workpulse
```

## Usage

```bash
# Start the timer
workpulse

# Start with a pre-set activity name
workpulse --activity "Building new feature"
workpulse -a "Code review"

# Start a 2-hour session with 25/5 ratio
workpulse --session 2h --ratio 25/5

# Use preset ratios
workpulse --session 90m --preset deep-work
```

### CLI Commands

```bash
workpulse                    # Start the timer
workpulse config             # View/edit configuration
workpulse stats              # View productivity stats
workpulse export             # Export sessions to CSV
workpulse update             # Check for updates and install
workpulse version            # Show version information
```

### Session Planning Flags

```bash
--session <duration>    # Total session goal (e.g., 2h, 90m)
--ratio <work/break>    # Work/break ratio (e.g., 25/5)
--preset <name>         # Use preset: standard, short-burst, deep-work
```

## Keyboard Shortcuts

### Timer View

| Key | Action |
|-----|--------|
| `s` | Start timer |
| `p` | Pause timer |
| `r` | Reset timer |
| `n` | Next mode |
| `d` | Stats/Summary |
| `c` | Settings |
| `S` | Session Planner |
| `q` | Quit |

### Update Banner (when available)

| Key | Action |
|-----|--------|
| `u` | Install update & restart |
| `U` | Dismiss banner |

### Quick Mode Switch

| Key | Mode |
|-----|------|
| `1` or `w` | Work (25 min) |
| `2` or `b` | Short Break (5 min) |
| `3` or `l` | Long Break (15 min) |
| `4` | Walk (10 min) |
| `5` | Water (2 min) |
| `6` | Video (20 min) |

### Settings View

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Navigate |
| `h/l` or `←/→` | Adjust value |
| `Enter/Space` | Toggle |
| `s` | Save settings |
| `Esc` | Cancel |

### Summary View

| Key | Action |
|-----|--------|
| `d` | Cycle period (Daily/Weekly/Monthly) |
| `e` | Export to CSV |
| `Esc` | Return |

## Session Planner

Plan focused work sessions with automatic cycling between work and breaks:

**Presets:**
- **Standard** — 25m work / 5m break (classic Pomodoro)
- **Short Burst** — 15m work / 10m break (frequent breaks)
- **Deep Work** — 45m work / 15m break (extended focus)

Press `S` from the timer view or use CLI flags:

```bash
# 2-hour session with standard Pomodoro
workpulse --session 2h --preset standard

# Custom ratio
workpulse --session 90m --ratio 30/10
```

## Timer Modes

| Mode | Duration | Purpose |
|------|----------|---------|
| 🎯 Work | 25 min | Focused work session |
| ☕ Short Break | 5 min | Quick rest between sessions |
| 🌴 Long Break | 15 min | Extended rest after multiple sessions |
| 🚶 Walk | 10 min | Get moving, stretch |
| 💧 Water | 2 min | Hydration reminder |
| 🎬 Video | 20 min | Video break time |

## Configuration

Configure via CLI or the in-app settings (`c` key):

```bash
workpulse config --work 30m
workpulse config --short-break 10m
workpulse config --theme dracula
workpulse config --sound=false
workpulse config --auto=true
workpulse config --daily-goal 4h
```

## Data Storage

WorkPulse stores data in `~/.workpulse/`:

```
~/.workpulse/
├── activities.json      # Session history
├── config.json          # User preferences
└── update_cache.json    # Update check cache (24h TTL)
```

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style definitions
- [Beeep](https://github.com/gen2brain/beeep) — Cross-platform notifications

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the [Pomodoro Technique](https://francescocirillo.com/products/the-pomodoro-technique) by Francesco Cirillo
- Built with the excellent [Charm](https://charm.sh/) ecosystem
