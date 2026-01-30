# Manual Testing Guide

## Quick Start

```bash
# Build the binary
go build -o workpulse

# Run with default settings
./workpulse

# Run with custom settings for faster testing
./workpulse --work 10s --short-break 5s --long-break 10s --auto
```

## Test Scenarios

### 1. Basic Timer Controls

```bash
./workpulse
```

| Action | Expected Result |
|--------|-----------------|
| Press `s` | Prompts for activity name |
| Type name + Enter | Timer starts counting |
| Press `p` | Timer pauses |
| Press `s` | Timer resumes |
| Press `r` | Timer resets to 00:00 |
| Press `q` | App exits cleanly |

### 2. Mode Switching

```bash
./workpulse
```

| Action | Expected Result |
|--------|-----------------|
| Press `1` or `w` | Switches to Work mode |
| Press `2` or `b` | Switches to Short Break |
| Press `3` or `l` | Switches to Long Break |
| Press `4` | Switches to Walk |
| Press `5` | Switches to Water |
| Press `6` | Switches to Video |
| Press `n` | Goes to next mode in sequence |

### 3. Custom Durations

```bash
./workpulse --work 30m --short-break 10m
```

- Verify Work mode shows 30:00
- Verify Short Break shows 10:00

### 4. Timer Completion (Fast Test)

```bash
./workpulse --work 5s --short-break 3s
```

| Action | Expected Result |
|--------|-----------------|
| Press `s`, enter activity | Timer starts |
| Wait 5 seconds | Timer completes |
| | Desktop notification appears |
| | Sound plays (Glass.aiff on macOS) |
| | "WORK completed!" message shows |

### 5. Auto-Advance

```bash
./workpulse --work 5s --short-break 3s --auto
```

| Action | Expected Result |
|--------|-----------------|
| Press `s`, enter activity | Work timer starts |
| Wait 5 seconds | Work completes |
| Wait 2 seconds | Short Break auto-starts |
| Wait 3 seconds | Short Break completes |
| Wait 2 seconds | Work auto-starts |

### 6. Pomodoro Counter

```bash
./workpulse --work 3s --short-break 2s --long-break 5s --auto
```

| Action | Expected Result |
|--------|-----------------|
| Start timer | Shows "Pomodoro 0/4" |
| Complete 1st work | Shows "Pomodoro 1/4" |
| Complete 2nd work | Shows "Pomodoro 2/4" |
| Complete 3rd work | Shows "Pomodoro 3/4" |
| Complete 4th work | Shows "Pomodoro 4/4" |
| Press `n` after 4th | Goes to Long Break (not Short) |
| Complete Long Break | Counter resets to 0 |

### 7. Activity Tracking

```bash
./workpulse -a "Testing WorkPulse"
```

- Timer should start without prompting for activity name
- Activity name "Testing WorkPulse" displays in the UI

### 8. Daily Summary

```bash
./workpulse --work 5s --short-break 3s
```

1. Complete a few timers
2. Press `d` to view daily summary
3. Verify time is tracked per mode
4. Press any key to return to timer

### 9. Data Persistence

```bash
# Complete some sessions
./workpulse --work 5s

# Check stored data
cat ~/.workpulse/activities.json
```

Verify JSON contains session entries with:
- `mode`: "work", "short_break", etc.
- `activity`: Activity name
- `started`: ISO timestamp
- `duration`: Duration in nanoseconds
- `completed`: true

## Keyboard Shortcut Reference

| Key | Action |
|-----|--------|
| `s` | Start |
| `p` | Pause |
| `r` | Reset |
| `n` | Next mode |
| `d` | Daily summary |
| `q` | Quit |
| `1` | Work |
| `2` | Short Break |
| `3` | Long Break |
| `4` | Walk |
| `5` | Water |
| `6` | Video |
| `w` | Work |
| `b` | Short Break |
| `l` | Long Break |

## CLI Flags

```bash
./workpulse --help

# Activity name
-a, --activity "name"

# Durations (use Go duration format: 25m, 5m, 30s)
--work 25m
--short-break 5m
--long-break 15m
--walk 10m
--water 2m
--video 20m

# Auto-advance
--auto
```
