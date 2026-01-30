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
		{"d", "daily"},
		{"q", "quit"},
	}

	var parts []string
	for _, item := range helpItems {
		key := KeyStyle.Render("[" + item.key + "]")
		parts = append(parts, key+item.desc)
	}

	helpLine := strings.Join(parts, "  ")

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
	var b strings.Builder

	header := SummaryHeaderStyle.Render("📊 Daily Summary")
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

	for _, m := range modeOrder {
		minutes := stats[m.key]
		if minutes > 0 {
			line := m.emoji + " " + m.label + ": " + SummaryValueStyle.Render(formatMinutes(minutes))
			b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	totalLine := "⏰ Total: " + SummaryValueStyle.Render(formatMinutes(totalMinutes))
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Render(totalLine))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(45).Foreground(Gray).Render("Press any key to return"))

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
