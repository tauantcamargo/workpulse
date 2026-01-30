package app

import (
	"github.com/tauantcamargo/workpulse/internal/storage"
	"github.com/tauantcamargo/workpulse/internal/timer"
)

type ViewState int

const (
	ViewTimer ViewState = iota
	ViewSummary
	ViewActivityInput
)

type Model struct {
	Timer        timer.Timer
	Modes        *timer.Modes
	ActivityName string
	ViewState    ViewState
	InputBuffer  string
	Storage      *storage.Storage
	Width        int
	Height       int
}

func NewModel(activityName string, store *storage.Storage, durations timer.DurationConfig) Model {
	modes := timer.NewModes(durations)
	mode := modes.Get(timer.ModeWork)
	return Model{
		Timer:        timer.New(mode),
		Modes:        modes,
		ActivityName: activityName,
		ViewState:    ViewTimer,
		InputBuffer:  "",
		Storage:      store,
		Width:        50,
		Height:       20,
	}
}

func (m Model) SetTimer(t timer.Timer) Model {
	return Model{
		Timer:        t,
		Modes:        m.Modes,
		ActivityName: m.ActivityName,
		ViewState:    m.ViewState,
		InputBuffer:  m.InputBuffer,
		Storage:      m.Storage,
		Width:        m.Width,
		Height:       m.Height,
	}
}

func (m Model) SetActivityName(name string) Model {
	return Model{
		Timer:        m.Timer,
		Modes:        m.Modes,
		ActivityName: name,
		ViewState:    m.ViewState,
		InputBuffer:  m.InputBuffer,
		Storage:      m.Storage,
		Width:        m.Width,
		Height:       m.Height,
	}
}

func (m Model) SetViewState(state ViewState) Model {
	return Model{
		Timer:        m.Timer,
		Modes:        m.Modes,
		ActivityName: m.ActivityName,
		ViewState:    state,
		InputBuffer:  m.InputBuffer,
		Storage:      m.Storage,
		Width:        m.Width,
		Height:       m.Height,
	}
}

func (m Model) SetInputBuffer(input string) Model {
	return Model{
		Timer:        m.Timer,
		Modes:        m.Modes,
		ActivityName: m.ActivityName,
		ViewState:    m.ViewState,
		InputBuffer:  input,
		Storage:      m.Storage,
		Width:        m.Width,
		Height:       m.Height,
	}
}

func (m Model) SetDimensions(width, height int) Model {
	return Model{
		Timer:        m.Timer,
		Modes:        m.Modes,
		ActivityName: m.ActivityName,
		ViewState:    m.ViewState,
		InputBuffer:  m.InputBuffer,
		Storage:      m.Storage,
		Width:        width,
		Height:       height,
	}
}
