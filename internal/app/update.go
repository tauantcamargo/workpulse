package app

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
	"github.com/tauantcamargo/workpulse/internal/config"
	"github.com/tauantcamargo/workpulse/internal/session"
	"github.com/tauantcamargo/workpulse/internal/storage"
	"github.com/tauantcamargo/workpulse/internal/timer"
	"github.com/tauantcamargo/workpulse/internal/ui"
	"github.com/tauantcamargo/workpulse/internal/update"
)

type NotifyMsg struct{}

type AutoAdvanceMsg struct{}

type UpdateCompleteMsg struct {
	Success bool
	Error   error
}

func (m Model) Init() tea.Cmd {
	if m.UpdateChecker != nil {
		return m.UpdateChecker.Check()
	}
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

	case update.UpdateCheckMsg:
		return m.handleUpdateCheck(msg)

	case UpdateCompleteMsg:
		return m.handleUpdateComplete(msg)
	}

	return m, nil
}

func (m Model) handleUpdateCheck(msg update.UpdateCheckMsg) (tea.Model, tea.Cmd) {
	if msg.Error != nil {
		return m, nil
	}
	if msg.Info.HasUpdate {
		return m.SetUpdateInfo(&msg.Info), nil
	}
	return m, nil
}

func (m Model) handleUpdateComplete(msg UpdateCompleteMsg) (tea.Model, tea.Cmd) {
	if msg.Success {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.ViewState {
	case ViewActivityInput:
		return m.handleActivityInput(msg)
	case ViewSummary:
		return m.handleSummaryKeys(msg)
	case ViewSessionSetup:
		return m.handleSessionSetupKeys(msg)
	case ViewSessionComplete:
		return m.handleSessionCompleteKeys(msg)
	case ViewSettings:
		return m.handleSettingsKeys(msg)
	case ViewSessionDurationInput:
		return m.handleSessionDurationInput(msg)
	case ViewSessionRatioInput:
		return m.handleSessionRatioInput(msg)
	default:
		return m.handleTimerKeys(msg)
	}
}

func (m Model) handleSummaryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "d", "tab", "right", "left":
		return m.NextSummaryPeriod(), nil
	case "esc", "enter":
		return m.SetViewState(ViewTimer).SetSummaryPeriod(PeriodDaily), nil
	case "e":
		return m.exportSessions()
	}
	return m, nil
}

func (m Model) exportSessions() (tea.Model, tea.Cmd) {
	if m.Storage == nil {
		return m, nil
	}

	var sessions []storage.Session
	var filename string
	switch m.SummaryPeriod {
	case PeriodWeekly:
		sessions = m.Storage.GetWeekSessions()
		filename = "workpulse_weekly.csv"
	case PeriodMonthly:
		sessions = m.Storage.GetMonthSessions()
		filename = "workpulse_monthly.csv"
	default:
		sessions = m.Storage.GetTodaySessions()
		filename = "workpulse_daily.csv"
	}

	if len(sessions) == 0 {
		return m, nil
	}

	csvData := m.Storage.ExportCSV(sessions)
	homeDir, _ := os.UserHomeDir()
	exportPath := homeDir + "/Desktop/" + filename
	os.WriteFile(exportPath, []byte(csvData), 0644)

	return m, nil
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

	case tea.KeySpace:
		return m.SetInputBuffer(m.InputBuffer + " "), nil

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

	case "u":
		if m.ShowUpdateBanner && m.UpdateInfo != nil {
			return m, runUpdate()
		}
		return m, nil

	case "U":
		if m.ShowUpdateBanner {
			return m.DismissUpdateBanner(), nil
		}
		return m, nil

	case "S":
		return m.SetViewState(ViewSessionSetup), nil

	case "x":
		if m.SessionPlan != nil && m.SessionPlan.Active {
			newPlan := m.SessionPlan.Cancel()
			return m.SetSessionPlan(&newPlan), nil
		}
		return m, nil

	case "c":
		return m.SetViewState(ViewSettings).SetSettingsField(SettingWorkDuration), nil

	case "e":
		return m.SetViewState(ViewSummary), nil
	}

	return m, nil
}

func runUpdate() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("sh", "-c", "curl -sSL https://raw.githubusercontent.com/tauantcamargo/workpulse/main/install.sh | sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return UpdateCompleteMsg{Success: false, Error: err}
		}

		binary, err := os.Executable()
		if err != nil {
			return UpdateCompleteMsg{Success: false, Error: err}
		}

		if err := syscall.Exec(binary, os.Args, os.Environ()); err != nil {
			return UpdateCompleteMsg{Success: false, Error: err}
		}

		return UpdateCompleteMsg{Success: true}
	}
}

func (m Model) startTimer() (tea.Model, tea.Cmd) {
	if m.Timer.Mode.Type == timer.ModeWork && m.ActivityName == "" {
		return m.SetViewState(ViewActivityInput), nil
	}

	newTimer := m.Timer.Start()
	return m.SetTimer(newTimer), timer.Tick()
}

func (m Model) nextMode() (tea.Model, tea.Cmd) {
	nextMode := m.Modes.NextWithPomodoro(m.Timer.Mode.Type, m.PomodoroCount, DefaultPomodorosBeforeLongBreak)
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

		if m.Timer.Mode.Type == timer.ModeWork {
			newModel = newModel.IncrementPomodoro()
		} else if m.Timer.Mode.Type == timer.ModeLongBreak {
			newModel = newModel.ResetPomodoroCount()
		}

		if newModel.SessionPlan != nil && newModel.SessionPlan.Active {
			updatedPlan := newModel.SessionPlan.AddElapsed(m.Timer.Mode.Duration)
			newModel = newModel.SetSessionPlan(&updatedPlan)
		}

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
	if m.SessionPlan != nil && m.SessionPlan.Active {
		if m.Timer.Mode.Type == timer.ModeShortBreak || m.Timer.Mode.Type == timer.ModeLongBreak {
			updatedPlan := m.SessionPlan.NextCycle()
			newModel := m.SetSessionPlan(&updatedPlan)

			if updatedPlan.IsComplete() {
				return newModel.SetViewState(ViewSessionComplete), nil
			}

			workMode := m.Modes.Get(timer.ModeWork)
			newTimer := m.Timer.SetMode(workMode).Start()
			return newModel.SetTimer(newTimer), timer.Tick()
		}

		breakMode := m.Modes.Get(timer.ModeShortBreak)
		newTimer := m.Timer.SetMode(breakMode).Start()
		return m.SetTimer(newTimer), timer.Tick()
	}

	nextMode := m.Modes.NextWithPomodoro(m.Timer.Mode.Type, m.PomodoroCount, DefaultPomodorosBeforeLongBreak)
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

func (m Model) handleSessionSetupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m.SetViewState(ViewTimer), nil
	case "1":
		return m.startSessionWithPreset("standard")
	case "2":
		return m.startSessionWithPreset("short-burst")
	case "3":
		return m.startSessionWithPreset("deep-work")
	case "4":
		return m.SetViewState(ViewSessionDurationInput).SetInputBuffer("").SetSessionSetupStep(0), nil
	}
	return m, nil
}

func (m Model) handleSessionCompleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc", "enter":
		return m.SetViewState(ViewTimer).SetSessionPlan(nil), nil
	case "r":
		if m.SessionPlan != nil {
			newPlan := m.SessionPlan.Start()
			newPlan.CurrentCycle = 1
			newPlan.ElapsedTotal = 0
			return m.SetSessionPlan(&newPlan).SetViewState(ViewTimer), nil
		}
		return m.SetViewState(ViewTimer), nil
	}
	return m, nil
}

func (m Model) startSessionWithPreset(presetName string) (tea.Model, tea.Cmd) {
	ratio, ok := session.GetPreset(presetName)
	if !ok {
		return m.SetViewState(ViewTimer), nil
	}

	plan := session.NewPlan(2*time.Hour, ratio)
	newModel := m.SetSessionPlan(plan).SetViewState(ViewTimer)

	newModel.Modes = timer.NewModes(timer.DurationConfig{
		Work:       ratio.Work,
		ShortBreak: ratio.Break,
		LongBreak:  15 * time.Minute,
		Walk:       10 * time.Minute,
		Water:      2 * time.Minute,
		Video:      20 * time.Minute,
	})

	workMode := newModel.Modes.Get(timer.ModeWork)
	newTimer := m.Timer.SetMode(workMode)
	newModel = newModel.SetTimer(newTimer)
	newModel.AutoAdvance = true

	return newModel, nil
}

func (m Model) handleSettingsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m.SetViewState(ViewTimer), nil
	case "up", "k":
		return m.PrevSettingsField(), nil
	case "down", "j":
		return m.NextSettingsField(), nil
	case "left", "h":
		return m.decrementSetting()
	case "right", "l":
		return m.incrementSetting()
	case "enter", " ":
		return m.toggleSetting()
	case "s":
		return m.saveSettings()
	}
	return m, nil
}

func (m Model) incrementSetting() (tea.Model, tea.Cmd) {
	switch m.SettingsField {
	case SettingWorkDuration:
		newDuration := m.Modes.Get(timer.ModeWork).Duration + 5*time.Minute
		if newDuration > 60*time.Minute {
			newDuration = 60 * time.Minute
		}
		return m.updateModeDuration(timer.ModeWork, newDuration), nil
	case SettingShortBreak:
		newDuration := m.Modes.Get(timer.ModeShortBreak).Duration + time.Minute
		if newDuration > 30*time.Minute {
			newDuration = 30 * time.Minute
		}
		return m.updateModeDuration(timer.ModeShortBreak, newDuration), nil
	case SettingLongBreak:
		newDuration := m.Modes.Get(timer.ModeLongBreak).Duration + 5*time.Minute
		if newDuration > 60*time.Minute {
			newDuration = 60 * time.Minute
		}
		return m.updateModeDuration(timer.ModeLongBreak, newDuration), nil
	case SettingTheme:
		themes := ui.AvailableThemes()
		newIndex := (m.ThemeIndex + 1) % len(themes)
		ui.SetTheme(ui.GetTheme(themes[newIndex]))
		return m.SetThemeIndex(newIndex), nil
	case SettingDailyGoal:
		newGoal := m.DailyGoal + 30*time.Minute
		if newGoal > 12*time.Hour {
			newGoal = 12 * time.Hour
		}
		return m.SetDailyGoal(newGoal), nil
	}
	return m, nil
}

func (m Model) decrementSetting() (tea.Model, tea.Cmd) {
	switch m.SettingsField {
	case SettingWorkDuration:
		newDuration := m.Modes.Get(timer.ModeWork).Duration - 5*time.Minute
		if newDuration < 5*time.Minute {
			newDuration = 5 * time.Minute
		}
		return m.updateModeDuration(timer.ModeWork, newDuration), nil
	case SettingShortBreak:
		newDuration := m.Modes.Get(timer.ModeShortBreak).Duration - time.Minute
		if newDuration < time.Minute {
			newDuration = time.Minute
		}
		return m.updateModeDuration(timer.ModeShortBreak, newDuration), nil
	case SettingLongBreak:
		newDuration := m.Modes.Get(timer.ModeLongBreak).Duration - 5*time.Minute
		if newDuration < 5*time.Minute {
			newDuration = 5 * time.Minute
		}
		return m.updateModeDuration(timer.ModeLongBreak, newDuration), nil
	case SettingTheme:
		themes := ui.AvailableThemes()
		newIndex := m.ThemeIndex - 1
		if newIndex < 0 {
			newIndex = len(themes) - 1
		}
		ui.SetTheme(ui.GetTheme(themes[newIndex]))
		return m.SetThemeIndex(newIndex), nil
	case SettingDailyGoal:
		newGoal := m.DailyGoal - 30*time.Minute
		if newGoal < 30*time.Minute {
			newGoal = 30 * time.Minute
		}
		return m.SetDailyGoal(newGoal), nil
	}
	return m, nil
}

func (m Model) toggleSetting() (tea.Model, tea.Cmd) {
	switch m.SettingsField {
	case SettingSound:
		return m.SetSoundEnabled(!m.SoundEnabled), nil
	case SettingNotify:
		return m.SetNotifyEnabled(!m.NotifyEnabled), nil
	case SettingAutoAdvance:
		return m.SetAutoAdvance(!m.AutoAdvance), nil
	case SettingTheme:
		return m.incrementSetting()
	}
	return m, nil
}

func (m Model) updateModeDuration(modeType timer.ModeType, duration time.Duration) Model {
	cfg := timer.DurationConfig{
		Work:       m.Modes.Get(timer.ModeWork).Duration,
		ShortBreak: m.Modes.Get(timer.ModeShortBreak).Duration,
		LongBreak:  m.Modes.Get(timer.ModeLongBreak).Duration,
		Walk:       m.Modes.Get(timer.ModeWalk).Duration,
		Water:      m.Modes.Get(timer.ModeWater).Duration,
		Video:      m.Modes.Get(timer.ModeVideo).Duration,
	}

	switch modeType {
	case timer.ModeWork:
		cfg.Work = duration
	case timer.ModeShortBreak:
		cfg.ShortBreak = duration
	case timer.ModeLongBreak:
		cfg.LongBreak = duration
	}

	newModes := timer.NewModes(cfg)
	newModel := m.SetModes(newModes)

	if m.Timer.Mode.Type == modeType {
		newMode := newModes.Get(modeType)
		newModel = newModel.SetTimer(m.Timer.SetMode(newMode))
	}

	return newModel
}

func (m Model) saveSettings() (tea.Model, tea.Cmd) {
	cfg := config.Config{
		Durations: config.ModeDuration{
			Work:       m.Modes.Get(timer.ModeWork).Duration,
			ShortBreak: m.Modes.Get(timer.ModeShortBreak).Duration,
			LongBreak:  m.Modes.Get(timer.ModeLongBreak).Duration,
			Walk:       m.Modes.Get(timer.ModeWalk).Duration,
			Water:      m.Modes.Get(timer.ModeWater).Duration,
			Video:      m.Modes.Get(timer.ModeVideo).Duration,
		},
		SoundEnabled:             m.SoundEnabled,
		NotifyEnabled:            m.NotifyEnabled,
		AutoAdvance:              m.AutoAdvance,
		PomodorosBeforeLongBreak: DefaultPomodorosBeforeLongBreak,
		DailyGoal:                m.DailyGoal,
		Theme:                    ui.AvailableThemes()[m.ThemeIndex],
	}

	config.Save(cfg)
	return m.SetViewState(ViewTimer), nil
}

func (m Model) handleSessionDurationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		duration, err := parseDurationInput(m.InputBuffer)
		if err != nil || duration <= 0 {
			return m.SetInputBuffer(""), nil
		}
		return m.SetInputBuffer("").
			SetPendingSessionDuration(duration).
			SetViewState(ViewSessionRatioInput), nil

	case tea.KeyBackspace:
		if len(m.InputBuffer) > 0 {
			return m.SetInputBuffer(m.InputBuffer[:len(m.InputBuffer)-1]), nil
		}
		return m, nil

	case tea.KeyEsc:
		return m.SetInputBuffer("").SetViewState(ViewSessionSetup), nil

	case tea.KeyRunes:
		return m.SetInputBuffer(m.InputBuffer + string(msg.Runes)), nil
	}
	return m, nil
}

func (m Model) handleSessionRatioInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		ratio, err := parseRatioInput(m.InputBuffer)
		if err != nil || ratio.Work <= 0 || ratio.Break <= 0 {
			return m.SetInputBuffer(""), nil
		}

		duration := m.PendingSessionDuration
		if duration <= 0 {
			duration = 2 * time.Hour
		}

		return m.startCustomSession(duration, ratio)

	case tea.KeyBackspace:
		if len(m.InputBuffer) > 0 {
			return m.SetInputBuffer(m.InputBuffer[:len(m.InputBuffer)-1]), nil
		}
		return m, nil

	case tea.KeyEsc:
		return m.SetInputBuffer("").SetViewState(ViewSessionSetup), nil

	case tea.KeyRunes:
		return m.SetInputBuffer(m.InputBuffer + string(msg.Runes)), nil
	}
	return m, nil
}

func (m Model) startCustomSession(totalDuration time.Duration, ratio session.Ratio) (tea.Model, tea.Cmd) {
	plan := session.NewPlan(totalDuration, ratio)
	newModel := m.SetSessionPlan(plan).SetViewState(ViewTimer).SetInputBuffer("")

	newModel.Modes = timer.NewModes(timer.DurationConfig{
		Work:       ratio.Work,
		ShortBreak: ratio.Break,
		LongBreak:  15 * time.Minute,
		Walk:       10 * time.Minute,
		Water:      2 * time.Minute,
		Video:      20 * time.Minute,
	})

	workMode := newModel.Modes.Get(timer.ModeWork)
	newTimer := m.Timer.SetMode(workMode)
	newModel = newModel.SetTimer(newTimer)
	newModel.AutoAdvance = true

	return newModel, nil
}

func parseDurationInput(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	if strings.HasSuffix(s, "h") {
		hours, err := strconv.Atoi(strings.TrimSuffix(s, "h"))
		if err == nil {
			return time.Duration(hours) * time.Hour, nil
		}
	}

	if strings.HasSuffix(s, "m") {
		mins, err := strconv.Atoi(strings.TrimSuffix(s, "m"))
		if err == nil {
			return time.Duration(mins) * time.Minute, nil
		}
	}

	mins, err := strconv.Atoi(s)
	if err == nil {
		return time.Duration(mins) * time.Minute, nil
	}

	return 0, err
}

func parseRatioInput(s string) (session.Ratio, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return session.Ratio{}, nil
	}

	work, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || work <= 0 {
		return session.Ratio{}, err
	}

	breakDur, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || breakDur <= 0 {
		return session.Ratio{}, err
	}

	return session.Ratio{
		Work:  time.Duration(work) * time.Minute,
		Break: time.Duration(breakDur) * time.Minute,
	}, nil
}
