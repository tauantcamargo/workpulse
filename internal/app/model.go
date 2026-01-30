package app

import (
	"time"

	"github.com/tauantcamargo/workpulse/internal/session"
	"github.com/tauantcamargo/workpulse/internal/storage"
	"github.com/tauantcamargo/workpulse/internal/timer"
	"github.com/tauantcamargo/workpulse/internal/update"
)

type ViewState int

const (
	ViewTimer ViewState = iota
	ViewSummary
	ViewActivityInput
	ViewSessionSetup
	ViewSessionComplete
	ViewSettings
	ViewSessionDurationInput
	ViewSessionRatioInput
)

type SummaryPeriod int

const (
	PeriodDaily SummaryPeriod = iota
	PeriodWeekly
	PeriodMonthly
)

type SettingsField int

const (
	SettingWorkDuration SettingsField = iota
	SettingShortBreak
	SettingLongBreak
	SettingSound
	SettingNotify
	SettingAutoAdvance
	SettingTheme
	SettingDailyGoal
	SettingsFieldCount
)

const DefaultPomodorosBeforeLongBreak = 4

type Options struct {
	ActivityName   string
	Durations      timer.DurationConfig
	AutoAdvance    bool
	SessionPlan    *session.Plan
	CurrentVersion string
	SoundEnabled   bool
	NotifyEnabled  bool
	ThemeIndex     int
	DailyGoal      time.Duration
}

type Model struct {
	Timer            timer.Timer
	Modes            *timer.Modes
	ActivityName     string
	AutoAdvance      bool
	PomodoroCount    int
	ViewState        ViewState
	SummaryPeriod    SummaryPeriod
	InputBuffer      string
	Storage          *storage.Storage
	Width            int
	Height           int
	UpdateInfo       *update.UpdateInfo
	UpdateChecker    *update.Checker
	ShowUpdateBanner bool
	SessionPlan      *session.Plan
	SettingsField    SettingsField
	SoundEnabled     bool
	NotifyEnabled    bool
	ThemeIndex       int
	DailyGoal            time.Duration
	SessionSetupStep     int
	PendingSessionDuration time.Duration
}

func NewModel(opts Options, store *storage.Storage) Model {
	modes := timer.NewModes(opts.Durations)
	mode := modes.Get(timer.ModeWork)

	var checker *update.Checker
	if opts.CurrentVersion != "" {
		checker = update.NewChecker(opts.CurrentVersion)
	}

	return Model{
		Timer:            timer.New(mode),
		Modes:            modes,
		ActivityName:     opts.ActivityName,
		AutoAdvance:      opts.AutoAdvance,
		PomodoroCount:    0,
		ViewState:        ViewTimer,
		SummaryPeriod:    PeriodDaily,
		InputBuffer:      "",
		Storage:          store,
		Width:            50,
		Height:           20,
		UpdateInfo:       nil,
		UpdateChecker:    checker,
		ShowUpdateBanner: false,
		SessionPlan:      opts.SessionPlan,
		SettingsField:    SettingWorkDuration,
		SoundEnabled:     opts.SoundEnabled,
		NotifyEnabled:    opts.NotifyEnabled,
		ThemeIndex:       opts.ThemeIndex,
		DailyGoal:        opts.DailyGoal,
		SessionSetupStep: 0,
	}
}

// copy creates a shallow copy of the model
func (m Model) copy() Model {
	return Model{
		Timer:                  m.Timer,
		Modes:                  m.Modes,
		ActivityName:           m.ActivityName,
		AutoAdvance:            m.AutoAdvance,
		PomodoroCount:          m.PomodoroCount,
		ViewState:              m.ViewState,
		SummaryPeriod:          m.SummaryPeriod,
		InputBuffer:            m.InputBuffer,
		Storage:                m.Storage,
		Width:                  m.Width,
		Height:                 m.Height,
		UpdateInfo:             m.UpdateInfo,
		UpdateChecker:          m.UpdateChecker,
		ShowUpdateBanner:       m.ShowUpdateBanner,
		SessionPlan:            m.SessionPlan,
		SettingsField:          m.SettingsField,
		SoundEnabled:           m.SoundEnabled,
		NotifyEnabled:          m.NotifyEnabled,
		ThemeIndex:             m.ThemeIndex,
		DailyGoal:              m.DailyGoal,
		SessionSetupStep:       m.SessionSetupStep,
		PendingSessionDuration: m.PendingSessionDuration,
	}
}

func (m Model) SetTimer(t timer.Timer) Model {
	c := m.copy()
	c.Timer = t
	return c
}

func (m Model) SetActivityName(name string) Model {
	c := m.copy()
	c.ActivityName = name
	return c
}

func (m Model) SetViewState(state ViewState) Model {
	c := m.copy()
	c.ViewState = state
	return c
}

func (m Model) SetInputBuffer(input string) Model {
	c := m.copy()
	c.InputBuffer = input
	return c
}

func (m Model) SetDimensions(width, height int) Model {
	c := m.copy()
	c.Width = width
	c.Height = height
	return c
}

func (m Model) SetPomodoroCount(count int) Model {
	c := m.copy()
	c.PomodoroCount = count
	return c
}

func (m Model) SetSummaryPeriod(period SummaryPeriod) Model {
	c := m.copy()
	c.SummaryPeriod = period
	return c
}

func (m Model) NextSummaryPeriod() Model {
	next := (m.SummaryPeriod + 1) % 3
	return m.SetSummaryPeriod(next)
}

func (m Model) IncrementPomodoro() Model {
	newCount := m.PomodoroCount + 1
	if newCount > DefaultPomodorosBeforeLongBreak {
		newCount = 1
	}
	return m.SetPomodoroCount(newCount)
}

func (m Model) ResetPomodoroCount() Model {
	return m.SetPomodoroCount(0)
}

func (m Model) SetUpdateInfo(info *update.UpdateInfo) Model {
	c := m.copy()
	c.UpdateInfo = info
	c.ShowUpdateBanner = info != nil && info.HasUpdate
	return c
}

func (m Model) DismissUpdateBanner() Model {
	c := m.copy()
	c.ShowUpdateBanner = false
	return c
}

func (m Model) SetSessionPlan(plan *session.Plan) Model {
	c := m.copy()
	c.SessionPlan = plan
	return c
}

func (m Model) SetSettingsField(field SettingsField) Model {
	c := m.copy()
	c.SettingsField = field
	return c
}

func (m Model) NextSettingsField() Model {
	next := (m.SettingsField + 1) % SettingsFieldCount
	return m.SetSettingsField(next)
}

func (m Model) PrevSettingsField() Model {
	prev := m.SettingsField - 1
	if prev < 0 {
		prev = SettingsFieldCount - 1
	}
	return m.SetSettingsField(prev)
}

func (m Model) SetSoundEnabled(enabled bool) Model {
	c := m.copy()
	c.SoundEnabled = enabled
	return c
}

func (m Model) SetNotifyEnabled(enabled bool) Model {
	c := m.copy()
	c.NotifyEnabled = enabled
	return c
}

func (m Model) SetAutoAdvance(enabled bool) Model {
	c := m.copy()
	c.AutoAdvance = enabled
	return c
}

func (m Model) SetThemeIndex(index int) Model {
	c := m.copy()
	c.ThemeIndex = index
	return c
}

func (m Model) SetDailyGoal(goal time.Duration) Model {
	c := m.copy()
	c.DailyGoal = goal
	return c
}

func (m Model) SetSessionSetupStep(step int) Model {
	c := m.copy()
	c.SessionSetupStep = step
	return c
}

func (m Model) SetModes(modes *timer.Modes) Model {
	c := m.copy()
	c.Modes = modes
	return c
}

func (m Model) SetPendingSessionDuration(d time.Duration) Model {
	c := m.copy()
	c.PendingSessionDuration = d
	return c
}
