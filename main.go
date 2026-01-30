package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauantcamargo/workpulse/internal/app"
	"github.com/tauantcamargo/workpulse/internal/storage"
	"github.com/tauantcamargo/workpulse/internal/timer"
)

func main() {
	defaults := timer.DefaultDurations()

	activityFlag := flag.String("activity", "", "Activity name for work sessions")
	activityShort := flag.String("a", "", "Activity name for work sessions (shorthand)")

	workDuration := flag.Duration("work", defaults.Work, "Work session duration (e.g., 25m, 30m)")
	shortBreak := flag.Duration("short-break", defaults.ShortBreak, "Short break duration (e.g., 5m)")
	longBreak := flag.Duration("long-break", defaults.LongBreak, "Long break duration (e.g., 15m)")
	walkDuration := flag.Duration("walk", defaults.Walk, "Walk break duration (e.g., 10m)")
	waterDuration := flag.Duration("water", defaults.Water, "Water break duration (e.g., 2m)")
	videoDuration := flag.Duration("video", defaults.Video, "Video break duration (e.g., 20m)")

	flag.Parse()

	activityName := *activityFlag
	if activityName == "" {
		activityName = *activityShort
	}

	durations := timer.DurationConfig{
		Work:       parseDuration(*workDuration, defaults.Work),
		ShortBreak: parseDuration(*shortBreak, defaults.ShortBreak),
		LongBreak:  parseDuration(*longBreak, defaults.LongBreak),
		Walk:       parseDuration(*walkDuration, defaults.Walk),
		Water:      parseDuration(*waterDuration, defaults.Water),
		Video:      parseDuration(*videoDuration, defaults.Video),
	}

	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	model := app.NewModel(activityName, store, durations)

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running WorkPulse: %v\n", err)
		os.Exit(1)
	}
}

func parseDuration(d time.Duration, fallback time.Duration) time.Duration {
	if d <= 0 {
		return fallback
	}
	return d
}
