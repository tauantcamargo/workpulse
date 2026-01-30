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

type Options struct {
	ActivityName string
	Durations    timer.DurationConfig
	AutoAdvance  bool
}

type Model struct {
	Timer        timer.Timer
	Modes        *timer.Modes
	ActivityName string
	AutoAdvance  bool
	ViewState    ViewState
	InputBuffer  string
	Storage      *storage.Storage
	Width        int
	Height       int
}

func NewModel(opts Options, store *storage.Storage) Model {
	modes := timer.NewModes(opts.Durations)
	mode := modes.Get(timer.ModeWork)
	return Model{
		Timer:        timer.New(mode),
		Modes:        modes,
		ActivityName: opts.ActivityName,
		AutoAdvance:  opts.AutoAdvance,
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
		AutoAdvance:  m.AutoAdvance,
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
		AutoAdvance:  m.AutoAdvance,
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
		AutoAdvance:  m.AutoAdvance,
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
		AutoAdvance:  m.AutoAdvance,
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
		AutoAdvance:  m.AutoAdvance,
		ViewState:    m.ViewState,
		InputBuffer:  m.InputBuffer,
		Storage:      m.Storage,
		Width:        width,
		Height:       height,
	}
}
