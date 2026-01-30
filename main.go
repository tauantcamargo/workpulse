package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauantcamargo/workpulse/internal/app"
	"github.com/tauantcamargo/workpulse/internal/config"
	"github.com/tauantcamargo/workpulse/internal/storage"
	"github.com/tauantcamargo/workpulse/internal/timer"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			runConfigCommand(os.Args[2:])
			return
		case "export":
			runExportCommand(os.Args[2:])
			return
		}
	}

	runTimerApp()
}

func runExportCommand(args []string) {
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)

	outputFlag := exportCmd.String("o", "", "Output file path (default: stdout)")
	periodFlag := exportCmd.String("period", "all", "Period to export: all, today, week, month")

	exportCmd.Parse(args)

	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	var sessions []storage.Session
	switch *periodFlag {
	case "today":
		sessions = store.GetTodaySessions()
	case "week":
		sessions = store.GetWeekSessions()
	case "month":
		sessions = store.GetMonthSessions()
	default:
		sessions = store.GetAllSessions()
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found for the specified period.")
		return
	}

	csvData := store.ExportCSV(sessions)

	if *outputFlag != "" {
		if err := os.WriteFile(*outputFlag, []byte(csvData), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Exported %d sessions to %s\n", len(sessions), *outputFlag)
	} else {
		fmt.Print(csvData)
	}
}

func runConfigCommand(args []string) {
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)

	listFlag := configCmd.Bool("list", false, "List current configuration")
	workFlag := configCmd.Duration("work", 0, "Set work duration (e.g., 25m, 30m)")
	shortBreakFlag := configCmd.Duration("short-break", 0, "Set short break duration (e.g., 5m)")
	longBreakFlag := configCmd.Duration("long-break", 0, "Set long break duration (e.g., 15m)")
	walkFlag := configCmd.Duration("walk", 0, "Set walk duration (e.g., 10m)")
	waterFlag := configCmd.Duration("water", 0, "Set water duration (e.g., 2m)")
	videoFlag := configCmd.Duration("video", 0, "Set video duration (e.g., 20m)")
	soundFlag := configCmd.String("sound", "", "Enable/disable sound (true/false)")
	notifyFlag := configCmd.String("notify", "", "Enable/disable notifications (true/false)")
	autoFlag := configCmd.String("auto", "", "Enable/disable auto-advance (true/false)")
	pomodorosFlag := configCmd.Int("pomodoros", 0, "Set pomodoros before long break (e.g., 4)")
	resetFlag := configCmd.Bool("reset", false, "Reset to default configuration")

	configCmd.Parse(args)

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *resetFlag {
		cfg = config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration reset to defaults.")
		printConfig(cfg)
		return
	}

	modified := false

	if *workFlag > 0 {
		cfg.Durations.Work = *workFlag
		modified = true
	}
	if *shortBreakFlag > 0 {
		cfg.Durations.ShortBreak = *shortBreakFlag
		modified = true
	}
	if *longBreakFlag > 0 {
		cfg.Durations.LongBreak = *longBreakFlag
		modified = true
	}
	if *walkFlag > 0 {
		cfg.Durations.Walk = *walkFlag
		modified = true
	}
	if *waterFlag > 0 {
		cfg.Durations.Water = *waterFlag
		modified = true
	}
	if *videoFlag > 0 {
		cfg.Durations.Video = *videoFlag
		modified = true
	}
	if *soundFlag != "" {
		cfg.SoundEnabled = parseBool(*soundFlag)
		modified = true
	}
	if *notifyFlag != "" {
		cfg.NotifyEnabled = parseBool(*notifyFlag)
		modified = true
	}
	if *autoFlag != "" {
		cfg.AutoAdvance = parseBool(*autoFlag)
		modified = true
	}
	if *pomodorosFlag > 0 {
		cfg.PomodorosBeforeLongBreak = *pomodorosFlag
		modified = true
	}

	if modified {
		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration updated.")
	}

	if *listFlag || !modified {
		printConfig(cfg)
	}
}

func printConfig(cfg config.Config) {
	fmt.Println()
	fmt.Println("WorkPulse Configuration")
	fmt.Println("========================")
	fmt.Printf("Config file: %s\n", config.Path())
	fmt.Println()
	fmt.Println("Durations:")
	fmt.Printf("  work:        %s\n", config.FormatDuration(cfg.Durations.Work))
	fmt.Printf("  short-break: %s\n", config.FormatDuration(cfg.Durations.ShortBreak))
	fmt.Printf("  long-break:  %s\n", config.FormatDuration(cfg.Durations.LongBreak))
	fmt.Printf("  walk:        %s\n", config.FormatDuration(cfg.Durations.Walk))
	fmt.Printf("  water:       %s\n", config.FormatDuration(cfg.Durations.Water))
	fmt.Printf("  video:       %s\n", config.FormatDuration(cfg.Durations.Video))
	fmt.Println()
	fmt.Println("Settings:")
	fmt.Printf("  sound:      %s\n", formatBool(cfg.SoundEnabled))
	fmt.Printf("  notify:     %s\n", formatBool(cfg.NotifyEnabled))
	fmt.Printf("  auto:       %s\n", formatBool(cfg.AutoAdvance))
	fmt.Printf("  pomodoros:  %d (before long break)\n", cfg.PomodorosBeforeLongBreak)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  workpulse config --work 30m")
	fmt.Println("  workpulse config --sound=false")
	fmt.Println("  workpulse config --reset")
}

func parseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func formatBool(b bool) string {
	if b {
		return "enabled"
	}
	return "disabled"
}

func runTimerApp() {
	cfg, _ := config.Load()

	defaults := timer.DurationConfig{
		Work:       cfg.Durations.Work,
		ShortBreak: cfg.Durations.ShortBreak,
		LongBreak:  cfg.Durations.LongBreak,
		Walk:       cfg.Durations.Walk,
		Water:      cfg.Durations.Water,
		Video:      cfg.Durations.Video,
	}

	activityFlag := flag.String("activity", "", "Activity name for work sessions")
	activityShort := flag.String("a", "", "Activity name for work sessions (shorthand)")

	workDuration := flag.Duration("work", defaults.Work, "Work session duration (e.g., 25m, 30m)")
	shortBreak := flag.Duration("short-break", defaults.ShortBreak, "Short break duration (e.g., 5m)")
	longBreak := flag.Duration("long-break", defaults.LongBreak, "Long break duration (e.g., 15m)")
	walkDuration := flag.Duration("walk", defaults.Walk, "Walk break duration (e.g., 10m)")
	waterDuration := flag.Duration("water", defaults.Water, "Water break duration (e.g., 2m)")
	videoDuration := flag.Duration("video", defaults.Video, "Video break duration (e.g., 20m)")

	autoAdvance := flag.Bool("auto", cfg.AutoAdvance, "Auto-advance to next mode after completion")

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

	opts := app.Options{
		ActivityName: activityName,
		Durations:    durations,
		AutoAdvance:  *autoAdvance,
	}

	model := app.NewModel(opts, store)

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
