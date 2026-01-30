package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/tauantcamargo/workpulse/internal/config"
	"github.com/tauantcamargo/workpulse/internal/storage"
	"github.com/tauantcamargo/workpulse/internal/timer"
	"github.com/tauantcamargo/workpulse/internal/ui"
)

func (m Model) View() string {
	switch m.ViewState {
	case ViewSummary:
		return m.renderSummary()
	case ViewActivityInput:
		return m.renderActivityInput()
	case ViewSessionSetup:
		return m.renderSessionSetup()
	case ViewSessionComplete:
		return m.renderSessionComplete()
	case ViewSettings:
		return m.renderSettings()
	case ViewSessionDurationInput:
		return m.renderSessionDurationInput()
	case ViewSessionRatioInput:
		return m.renderSessionRatioInput()
	default:
		return m.renderTimer()
	}
}

func (m Model) renderTimer() string {
	var content strings.Builder

	if m.ShowUpdateBanner && m.UpdateInfo != nil {
		content.WriteString(ui.RenderUpdateBanner(m.UpdateInfo.LatestVersion.String()))
		content.WriteString("\n\n")
	} else {
		content.WriteString("\n")
	}
	content.WriteString(ui.RenderHeader(
		m.Timer.Mode.Emoji,
		m.Timer.Mode.Name,
		m.Timer.Mode.Message,
	))
	content.WriteString("\n\n")

	if m.SessionPlan != nil && m.SessionPlan.Active {
		cycleText := "Cycle " + strconv.Itoa(m.SessionPlan.CurrentCycle) + "/" + strconv.Itoa(m.SessionPlan.TotalCycles)
		cycleLine := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(45).
			Foreground(ui.Cyan).
			Bold(true).
			Render(cycleText)
		content.WriteString(cycleLine)
		content.WriteString("\n")

		elapsedTotal := config.FormatDuration(m.SessionPlan.ElapsedTotal)
		totalGoal := config.FormatDuration(m.SessionPlan.TotalGoal)
		sessionTimeText := elapsedTotal + " / " + totalGoal
		sessionTimeLine := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(45).
			Foreground(ui.LightGray).
			Render(sessionTimeText)
		content.WriteString(sessionTimeLine)
		content.WriteString("\n")

		sessionProgress := ui.RenderProgressWithPercent(m.SessionPlan.Progress(), 35)
		centeredSessionProgress := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(45).
			Render(sessionProgress)
		content.WriteString(centeredSessionProgress)
		content.WriteString("\n\n")
	} else if m.Timer.Mode.Type == timer.ModeWork || m.Timer.Mode.Type == timer.ModeShortBreak || m.Timer.Mode.Type == timer.ModeLongBreak {
		pomodoroText := "Pomodoro " + strconv.Itoa(m.PomodoroCount) + "/" + strconv.Itoa(DefaultPomodorosBeforeLongBreak)
		pomodoroLine := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(45).
			Foreground(ui.Purple).
			Bold(true).
			Render(pomodoroText)
		content.WriteString(pomodoroLine)
		content.WriteString("\n\n")
	}

	if m.ActivityName != "" && m.Timer.Mode.Type == timer.ModeWork {
		content.WriteString(ui.RenderActivity(m.ActivityName))
		content.WriteString("\n\n")
	}

	elapsed := timer.FormatDuration(m.Timer.Elapsed)
	total := timer.FormatDuration(m.Timer.Mode.Duration)
	content.WriteString(ui.RenderTimer(elapsed, total, m.Timer.Running))
	content.WriteString("\n\n")

	progressBar := ui.RenderProgressWithPercent(m.Timer.Progress(), 35)
	centeredProgress := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Render(progressBar)
	content.WriteString(centeredProgress)
	content.WriteString("\n\n")

	if m.Timer.Completed {
		content.WriteString(ui.RenderCompletedMessage(m.Timer.Mode.Name))
		content.WriteString("\n\n")
	}

	content.WriteString(ui.RenderHelp())
	content.WriteString("\n")
	content.WriteString(ui.RenderModeHelp())
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderActivityInput() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString(ui.RenderHeader(
		m.Timer.Mode.Emoji,
		m.Timer.Mode.Name,
		m.Timer.Mode.Message,
	))
	content.WriteString("\n\n")

	content.WriteString(ui.RenderInputPrompt("Activity name: ", m.InputBuffer))
	content.WriteString("\n\n")

	hint := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(ui.Gray).
		Render("Press Enter to start, Esc to cancel")
	content.WriteString(hint)
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderSummary() string {
	stats := make(map[string]int)
	totalMinutes := 0

	var periodTitle string
	if m.Storage != nil {
		var sessions []storage.Session
		switch m.SummaryPeriod {
		case PeriodWeekly:
			sessions = m.Storage.GetWeekSessions()
			periodTitle = "Weekly"
		case PeriodMonthly:
			sessions = m.Storage.GetMonthSessions()
			periodTitle = "Monthly"
		default:
			sessions = m.Storage.GetTodaySessions()
			periodTitle = "Daily"
		}

		for _, session := range sessions {
			minutes := int(session.Duration.Minutes())
			stats[session.Mode] += minutes
			totalMinutes += minutes
		}
	}

	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(ui.RenderSummary(stats, totalMinutes, periodTitle))
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderSessionSetup() string {
	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(ui.RenderSessionSetup())
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderSessionComplete() string {
	var content strings.Builder
	content.WriteString("\n")

	totalCycles := 0
	var totalTime time.Duration
	if m.SessionPlan != nil {
		totalCycles = m.SessionPlan.TotalCycles
		totalTime = m.SessionPlan.ElapsedTotal
	}

	content.WriteString(ui.RenderSessionComplete(totalCycles, config.FormatDuration(totalTime)))
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderSettings() string {
	var content strings.Builder
	content.WriteString("\n")

	themes := ui.AvailableThemes()
	themeName := "dark"
	if m.ThemeIndex < len(themes) {
		themeName = themes[m.ThemeIndex]
	}

	items := []ui.SettingsItem{
		{Label: "Work Duration:", Value: config.FormatDuration(m.Modes.Get(timer.ModeWork).Duration), Selected: m.SettingsField == SettingWorkDuration},
		{Label: "Short Break:", Value: config.FormatDuration(m.Modes.Get(timer.ModeShortBreak).Duration), Selected: m.SettingsField == SettingShortBreak},
		{Label: "Long Break:", Value: config.FormatDuration(m.Modes.Get(timer.ModeLongBreak).Duration), Selected: m.SettingsField == SettingLongBreak},
		{Label: "Sound:", Value: formatEnabled(m.SoundEnabled), Selected: m.SettingsField == SettingSound},
		{Label: "Notifications:", Value: formatEnabled(m.NotifyEnabled), Selected: m.SettingsField == SettingNotify},
		{Label: "Auto-advance:", Value: formatEnabled(m.AutoAdvance), Selected: m.SettingsField == SettingAutoAdvance},
		{Label: "Theme:", Value: themeName, Selected: m.SettingsField == SettingTheme},
		{Label: "Daily Goal:", Value: config.FormatDuration(m.DailyGoal), Selected: m.SettingsField == SettingDailyGoal},
	}

	content.WriteString(ui.RenderSettings(items))
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderSessionDurationInput() string {
	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(ui.RenderDurationInput("Duration: ", m.InputBuffer))
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func (m Model) renderSessionRatioInput() string {
	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(ui.RenderRatioInput("Ratio: ", m.InputBuffer))
	content.WriteString("\n")

	return ui.ContainerStyle.Render(content.String())
}

func formatEnabled(enabled bool) string {
	if enabled {
		return "ON"
	}
	return "OFF"
}
