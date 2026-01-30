package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Purple     = lipgloss.Color("#8B5CF6")
	Blue       = lipgloss.Color("#3B82F6")
	Cyan       = lipgloss.Color("#06B6D4")
	Green      = lipgloss.Color("#10B981")
	Yellow     = lipgloss.Color("#F59E0B")
	Red        = lipgloss.Color("#EF4444")
	Pink       = lipgloss.Color("#EC4899")
	Gray       = lipgloss.Color("#6B7280")
	DarkGray   = lipgloss.Color("#374151")
	LightGray  = lipgloss.Color("#9CA3AF")
	White      = lipgloss.Color("#F9FAFB")
	Background = lipgloss.Color("#1F2937")

	// Container style
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Purple)

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(White).
			MarginBottom(1)

	// Emoji style (larger display)
	EmojiStyle = lipgloss.NewStyle().
			MarginBottom(1)

	// Message style
	MessageStyle = lipgloss.NewStyle().
			Foreground(LightGray).
			Italic(true).
			MarginBottom(1)

	// Timer display style
	TimerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(White).
			MarginBottom(1)

	// Progress bar styles
	ProgressFilledStyle = lipgloss.NewStyle().
				Foreground(Purple)

	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(DarkGray)

	// Percentage style
	PercentStyle = lipgloss.NewStyle().
			Foreground(LightGray).
			MarginLeft(1)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(Gray).
			MarginTop(1)

	// Key style for shortcuts
	KeyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	// Status styles
	RunningStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	PausedStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	CompletedStyle = lipgloss.NewStyle().
			Foreground(Pink).
			Bold(true)

	// Activity name style
	ActivityStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Italic(true)

	// Input style
	InputStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(DarkGray).
			Padding(0, 1)

	// Summary styles
	SummaryHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Purple).
				MarginBottom(1)

	SummaryItemStyle = lipgloss.NewStyle().
				Foreground(White)

	SummaryValueStyle = lipgloss.NewStyle().
				Foreground(Cyan).
				Bold(true)
)

func ModeColor(mode string) lipgloss.Color {
	switch mode {
	case "WORK":
		return Purple
	case "SHORT BREAK":
		return Green
	case "LONG BREAK":
		return Cyan
	case "WALK":
		return Yellow
	case "WATER":
		return Blue
	case "VIDEO":
		return Pink
	default:
		return Purple
	}
}
