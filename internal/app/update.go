package app

import (
	"os/exec"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
	"github.com/tauantcamargo/workpulse/internal/timer"
)

type NotifyMsg struct{}

type AutoAdvanceMsg struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		return m.SetDimensions(msg.Width, msg.Height), nil

	case timer.TickMsg:
		return m.handleTick()

	case NotifyMsg:
		return m, nil

	case AutoAdvanceMsg:
		return m.handleAutoAdvance()
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.ViewState {
	case ViewActivityInput:
		return m.handleActivityInput(msg)
	case ViewSummary:
		return m.SetViewState(ViewTimer), nil
	default:
		return m.handleTimerKeys(msg)
	}
}

func (m Model) handleActivityInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		newModel := m.SetActivityName(m.InputBuffer).
			SetInputBuffer("").
			SetViewState(ViewTimer).
			SetTimer(m.Timer.Start())
		return newModel, timer.Tick()

	case tea.KeyBackspace:
		if len(m.InputBuffer) > 0 {
			return m.SetInputBuffer(m.InputBuffer[:len(m.InputBuffer)-1]), nil
		}
		return m, nil

	case tea.KeyEsc:
		return m.SetInputBuffer("").SetViewState(ViewTimer), nil

	case tea.KeyRunes:
		return m.SetInputBuffer(m.InputBuffer + string(msg.Runes)), nil
	}

	return m, nil
}

func (m Model) handleTimerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "s":
		return m.startTimer()

	case "p":
		return m.SetTimer(m.Timer.Pause()), nil

	case "r":
		return m.SetTimer(m.Timer.Reset()), nil

	case "n":
		return m.nextMode()

	case "w":
		return m.switchMode(timer.ModeWork)

	case "b":
		return m.switchMode(timer.ModeShortBreak)

	case "l":
		return m.switchMode(timer.ModeLongBreak)

	case "1":
		return m.switchMode(timer.ModeWork)

	case "2":
		return m.switchMode(timer.ModeShortBreak)

	case "3":
		return m.switchMode(timer.ModeLongBreak)

	case "4":
		return m.switchMode(timer.ModeWalk)

	case "5":
		return m.switchMode(timer.ModeWater)

	case "6":
		return m.switchMode(timer.ModeVideo)

	case "d":
		return m.SetViewState(ViewSummary), nil
	}

	return m, nil
}

func (m Model) startTimer() (tea.Model, tea.Cmd) {
	if m.Timer.Mode.Type == timer.ModeWork && m.ActivityName == "" {
		return m.SetViewState(ViewActivityInput), nil
	}

	newTimer := m.Timer.Start()
	return m.SetTimer(newTimer), timer.Tick()
}

func (m Model) nextMode() (tea.Model, tea.Cmd) {
	nextMode := m.Modes.Next(m.Timer.Mode.Type)
	newTimer := m.Timer.SetMode(nextMode)
	return m.SetTimer(newTimer), nil
}

func (m Model) switchMode(modeType timer.ModeType) (tea.Model, tea.Cmd) {
	mode := m.Modes.Get(modeType)
	newTimer := m.Timer.SetMode(mode)
	return m.SetTimer(newTimer), nil
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if !m.Timer.Running {
		return m, nil
	}

	newTimer := m.Timer.Tick()

	if newTimer.Completed && !m.Timer.Completed {
		newModel := m.SetTimer(newTimer)
		if m.Storage != nil {
			m.Storage.SaveSession(
				string(m.Timer.Mode.Type),
				m.ActivityName,
				m.Timer.Mode.Duration,
				time.Now().Add(-m.Timer.Mode.Duration),
			)
		}

		cmds := []tea.Cmd{
			sendNotification(m.Timer.Mode.Name),
			playSound(),
		}

		if m.AutoAdvance {
			cmds = append(cmds, scheduleAutoAdvance())
		}

		return newModel, tea.Batch(cmds...)
	}

	return m.SetTimer(newTimer), timer.Tick()
}

func (m Model) handleAutoAdvance() (tea.Model, tea.Cmd) {
	nextMode := m.Modes.Next(m.Timer.Mode.Type)
	newTimer := m.Timer.SetMode(nextMode).Start()
	return m.SetTimer(newTimer), timer.Tick()
}

func scheduleAutoAdvance() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return AutoAdvanceMsg{}
	})
}

func sendNotification(modeName string) tea.Cmd {
	return func() tea.Msg {
		beeep.Notify(
			"WorkPulse",
			modeName+" completed! Time for a change.",
			"",
		)
		return NotifyMsg{}
	}
}

func playSound() tea.Cmd {
	return func() tea.Msg {
		beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)

		if runtime.GOOS == "darwin" {
			exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Start()
		}

		return NotifyMsg{}
	}
}
