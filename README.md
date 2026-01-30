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
- **Gradient Progress Bar** — Beautiful purple-to-blue color gradient
- **Desktop Notifications** — Get notified when timers complete
- **Sound Alerts** — Audio feedback on timer completion
- **Activity Tracking** — Log what you're working on
- **Daily Summary** — View your productivity stats
- **Persistent Storage** — Sessions saved locally

## Installation

### From Source

```bash
git clone https://github.com/tauantcamargo/workpulse.git
cd workpulse
go build -o workpulse
./workpulse
```

### Go Install

```bash
go install github.com/tauantcamargo/workpulse@latest
```

## Usage

```bash
# Start the timer
workpulse

# Start with a pre-set activity name
workpulse --activity "Building new feature"
workpulse -a "Code review"
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `s` | Start timer |
| `p` | Pause timer |
| `r` | Reset timer |
| `n` | Next mode |
| `d` | Daily summary |
| `q` | Quit |

### Quick Mode Switch

| Key | Mode |
|-----|------|
| `1` or `w` | Work (25 min) |
| `2` or `b` | Short Break (5 min) |
| `3` or `l` | Long Break (15 min) |
| `4` | Walk (10 min) |
| `5` | Water (2 min) |
| `6` | Video (20 min) |

## Timer Modes

| Mode | Duration | Purpose |
|------|----------|---------|
| 🎯 Work | 25 min | Focused work session |
| ☕ Short Break | 5 min | Quick rest between sessions |
| 🌴 Long Break | 15 min | Extended rest after multiple sessions |
| 🚶 Walk | 10 min | Get moving, stretch |
| 💧 Water | 2 min | Hydration reminder |
| 🎬 Video | 20 min | Video break time |

## Data Storage

WorkPulse stores data in `~/.workpulse/`:

```
~/.workpulse/
├── activities.json    # Session history
└── config.json        # User preferences
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
