package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(emoji, name, message string) string {
	var b strings.Builder

	modeColor := ModeColor(name)
	emojiLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Render(emoji + " " + lipgloss.NewStyle().Bold(true).Foreground(modeColor).Render(name))

	messageLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(LightGray).
		Italic(true).
		Render("\"" + message + "\"")

	b.WriteString(emojiLine)
	b.WriteString("\n\n")
	b.WriteString(messageLine)

	return b.String()
}

func RenderTimer(elapsed, total string, running bool) string {
	timerIcon := "⏱️ "
	timerText := elapsed + " / " + total

	var status string
	if running {
		status = RunningStyle.Render(" ▶")
	} else {
		status = PausedStyle.Render(" ⏸")
	}

	timerLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Bold(true).
		Foreground(White).
		Render(timerIcon + timerText + status)

	return timerLine
}

func RenderActivity(name string) string {
	if name == "" {
		return ""
	}
	activityLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Cyan).
		Italic(true).
		Render("📝 " + name)

	return activityLine
}

func RenderHelp() string {
	helpItems := []struct {
		key  string
		desc string
	}{
		{"s", "start"},
		{"p", "pause"},
		{"r", "reset"},
		{"n", "next"},
		{"d", "stats"},
		{"c", "config"},
		{"S", "session"},
		{"q", "quit"},
	}

	var parts []string
	for _, item := range helpItems {
		key := KeyStyle.Render("[" + item.key + "]")
		parts = append(parts, key+item.desc)
	}

	helpLine := strings.Join(parts, " ")

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render(helpLine)
}

func RenderModeHelp() string {
	modes := []struct {
		key  string
		mode string
	}{
		{"1", "work"},
		{"2", "short"},
		{"3", "long"},
		{"4", "walk"},
		{"5", "water"},
		{"6", "video"},
	}

	var parts []string
	for _, m := range modes {
		key := KeyStyle.Render("[" + m.key + "]")
		parts = append(parts, key+m.mode)
	}

	modeLine := strings.Join(parts, " ")

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render(modeLine)
}

func RenderCompletedMessage(modeName string) string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Pink).
		Bold(true).
		Render("🎉 " + modeName + " completed! Press [n] for next")
}

func RenderInputPrompt(prompt, input string) string {
	promptStyle := lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true)

	inputStyle := lipgloss.NewStyle().
		Foreground(White)

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Render(promptStyle.Render(prompt) + inputStyle.Render(input+"_"))
}

func RenderDailySummary(stats map[string]int, totalMinutes int) string {
	return RenderSummary(stats, totalMinutes, "Daily")
}

func RenderSummary(stats map[string]int, totalMinutes int, period string) string {
	var b strings.Builder

	header := SummaryHeaderStyle.Render("📊 " + period + " Summary")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(header))
	b.WriteString("\n\n")

	modeOrder := []struct {
		key   string
		label string
		emoji string
	}{
		{"work", "Work", "🎯"},
		{"short_break", "Short Break", "☕"},
		{"long_break", "Long Break", "🌴"},
		{"walk", "Walk", "🚶"},
		{"water", "Water", "💧"},
		{"video", "Video", "🎬"},
	}

	hasStats := false
	for _, m := range modeOrder {
		minutes := stats[m.key]
		if minutes > 0 {
			hasStats = true
			line := m.emoji + " " + m.label + ": " + SummaryValueStyle.Render(formatMinutes(minutes))
			b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(line))
			b.WriteString("\n")
		}
	}

	if !hasStats {
		noDataLine := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(45).
			Foreground(Gray).
			Italic(true).
			Render("No sessions recorded")
		b.WriteString(noDataLine)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	totalLine := "⏰ Total: " + SummaryValueStyle.Render(formatMinutes(totalMinutes))
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(totalLine))
	b.WriteString("\n\n")

	helpLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render("[d] cycle period  [e] export  [esc] return")
	b.WriteString(helpLine)

	return b.String()
}

func formatMinutes(minutes int) string {
	if minutes < 60 {
		return itoa(minutes) + "m"
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return itoa(hours) + "h"
	}
	return itoa(hours) + "h " + itoa(mins) + "m"
}

// RenderUpdateBanner renders the update notification banner
func RenderUpdateBanner(latestVersion string) string {
	var b strings.Builder

	bannerStyle := lipgloss.NewStyle().
		Foreground(White).
		Background(Purple).
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Center).
		Width(43)

	updateLine := bannerStyle.Render("Update available: " + latestVersion)
	b.WriteString(updateLine)
	b.WriteString("\n")

	instructionStyle := lipgloss.NewStyle().
		Foreground(LightGray).
		Align(lipgloss.Center).
		Width(45)

	b.WriteString(instructionStyle.Render("[u] install  [U] dismiss"))

	return b.String()
}

// RenderSessionProgress renders the session progress indicator
func RenderSessionProgress(currentCycle, totalCycles int, elapsed, total string, progress float64) string {
	var b strings.Builder

	cycleText := "Cycle " + itoa(currentCycle) + "/" + itoa(totalCycles)
	cycleLineRendered := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Cyan).
		Bold(true).
		Render(cycleText)
	b.WriteString(cycleLineRendered)
	b.WriteString("\n")

	timeText := elapsed + " / " + total
	timeLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(LightGray).
		Render(timeText)
	b.WriteString(timeLine)
	b.WriteString("\n")

	progressBar := RenderProgressWithPercent(progress, 35)
	centeredProgress := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Render(progressBar)
	b.WriteString(centeredProgress)

	return b.String()
}

// RenderSessionSetup renders the session setup view
func RenderSessionSetup() string {
	var b strings.Builder

	header := SummaryHeaderStyle.Render("Session Planner")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(header))
	b.WriteString("\n\n")

	presets := []struct {
		key   string
		name  string
		ratio string
	}{
		{"1", "Standard", "25m work / 5m break"},
		{"2", "Short Burst", "15m work / 10m break"},
		{"3", "Deep Work", "45m work / 15m break"},
		{"4", "Custom", "enter duration & ratio"},
	}

	for _, p := range presets {
		keyStyle := KeyStyle.Render("[" + p.key + "]")
		line := keyStyle + " " + p.name + " - " + lipgloss.NewStyle().Foreground(LightGray).Render(p.ratio)
		b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render("[esc] cancel")
	b.WriteString(helpLine)

	return b.String()
}

// RenderSessionComplete renders the session completion view
func RenderSessionComplete(totalCycles int, totalTime string) string {
	var b strings.Builder

	header := CompletedStyle.Render("Session Complete!")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(header))
	b.WriteString("\n\n")

	statsLine := "Completed " + itoa(totalCycles) + " cycles in " + totalTime
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Foreground(White).Render(statsLine))
	b.WriteString("\n\n")

	helpLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render("[r] restart  [esc] close")
	b.WriteString(helpLine)

	return b.String()
}

// SettingsItem represents a single settings row
type SettingsItem struct {
	Label    string
	Value    string
	Selected bool
}

// RenderSettings renders the settings view
func RenderSettings(items []SettingsItem) string {
	var b strings.Builder

	header := SummaryHeaderStyle.Render("Settings")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(header))
	b.WriteString("\n\n")

	for _, item := range items {
		var line string
		labelStyle := lipgloss.NewStyle().Foreground(LightGray).Width(20)
		valueStyle := lipgloss.NewStyle().Foreground(Cyan).Bold(true)

		if item.Selected {
			labelStyle = labelStyle.Foreground(White).Bold(true)
			valueStyle = valueStyle.Foreground(Purple)
			line = "> " + labelStyle.Render(item.Label) + valueStyle.Render(item.Value)
		} else {
			line = "  " + labelStyle.Render(item.Label) + valueStyle.Render(item.Value)
		}

		b.WriteString(lipgloss.NewStyle().Width(45).Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render("[j/k] navigate  [h/l] adjust  [s] save  [esc] cancel")
	b.WriteString(helpLine)

	return b.String()
}

// RenderDurationInput renders a duration input prompt
func RenderDurationInput(prompt, input string) string {
	var b strings.Builder

	header := SummaryHeaderStyle.Render("Session Duration")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(header))
	b.WriteString("\n\n")

	b.WriteString(RenderInputPrompt(prompt, input))
	b.WriteString("\n\n")

	hint := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render("Examples: 2h, 90m, 1h30m")
	b.WriteString(hint)
	b.WriteString("\n")

	return b.String()
}

// RenderRatioInput renders a ratio input prompt
func RenderRatioInput(prompt, input string) string {
	var b strings.Builder

	header := SummaryHeaderStyle.Render("Work/Break Ratio")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(header))
	b.WriteString("\n\n")

	b.WriteString(RenderInputPrompt(prompt, input))
	b.WriteString("\n\n")

	hint := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(45).
		Foreground(Gray).
		Render("Format: work/break (e.g., 25/5)")
	b.WriteString(hint)
	b.WriteString("\n")

	return b.String()
}
