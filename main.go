package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauantcamargo/workpulse/internal/app"
	"github.com/tauantcamargo/workpulse/internal/storage"
)

func main() {
	activityFlag := flag.String("activity", "", "Activity name for work sessions")
	activityShort := flag.String("a", "", "Activity name for work sessions (shorthand)")
	flag.Parse()

	activityName := *activityFlag
	if activityName == "" {
		activityName = *activityShort
	}

	store, err := storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	model := app.NewModel(activityName, store)

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running WorkPulse: %v\n", err)
		os.Exit(1)
	}
}
