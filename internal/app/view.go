package app

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	default:
		return m.renderTimer()
	}
}

func (m Model) renderTimer() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString(ui.RenderHeader(
		m.Timer.Mode.Emoji,
		m.Timer.Mode.Name,
		m.Timer.Mode.Message,
	))
	content.WriteString("\n\n")

	if m.Timer.Mode.Type == timer.ModeWork || m.Timer.Mode.Type == timer.ModeShortBreak || m.Timer.Mode.Type == timer.ModeLongBreak {
		pomodoroText := "Pomodoro " + strconv.Itoa(m.PomodoroCount) + "/" + strconv.Itoa(DefaultPomodorosBeforeLongBreak)
		pomodoroLine := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(45).
			Foreground(ui.Purple).
			Bold(true).
			Render("🍅 " + pomodoroText)
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
