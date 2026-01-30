package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Text       lipgloss.Color
	TextMuted  lipgloss.Color
	Background lipgloss.Color
	Border     lipgloss.Color
	InputBg    lipgloss.Color
}

var (
	DarkTheme = Theme{
		Name:       "dark",
		Primary:    lipgloss.Color("#8B5CF6"),
		Secondary:  lipgloss.Color("#3B82F6"),
		Accent:     lipgloss.Color("#06B6D4"),
		Success:    lipgloss.Color("#10B981"),
		Warning:    lipgloss.Color("#F59E0B"),
		Error:      lipgloss.Color("#EF4444"),
		Text:       lipgloss.Color("#F9FAFB"),
		TextMuted:  lipgloss.Color("#9CA3AF"),
		Background: lipgloss.Color("#1F2937"),
		Border:     lipgloss.Color("#8B5CF6"),
		InputBg:    lipgloss.Color("#374151"),
	}

	LightTheme = Theme{
		Name:       "light",
		Primary:    lipgloss.Color("#7C3AED"),
		Secondary:  lipgloss.Color("#2563EB"),
		Accent:     lipgloss.Color("#0891B2"),
		Success:    lipgloss.Color("#059669"),
		Warning:    lipgloss.Color("#D97706"),
		Error:      lipgloss.Color("#DC2626"),
		Text:       lipgloss.Color("#111827"),
		TextMuted:  lipgloss.Color("#6B7280"),
		Background: lipgloss.Color("#F9FAFB"),
		Border:     lipgloss.Color("#7C3AED"),
		InputBg:    lipgloss.Color("#E5E7EB"),
	}

	DraculaTheme = Theme{
		Name:       "dracula",
		Primary:    lipgloss.Color("#BD93F9"),
		Secondary:  lipgloss.Color("#8BE9FD"),
		Accent:     lipgloss.Color("#50FA7B"),
		Success:    lipgloss.Color("#50FA7B"),
		Warning:    lipgloss.Color("#FFB86C"),
		Error:      lipgloss.Color("#FF5555"),
		Text:       lipgloss.Color("#F8F8F2"),
		TextMuted:  lipgloss.Color("#6272A4"),
		Background: lipgloss.Color("#282A36"),
		Border:     lipgloss.Color("#BD93F9"),
		InputBg:    lipgloss.Color("#44475A"),
	}

	NordTheme = Theme{
		Name:       "nord",
		Primary:    lipgloss.Color("#88C0D0"),
		Secondary:  lipgloss.Color("#81A1C1"),
		Accent:     lipgloss.Color("#A3BE8C"),
		Success:    lipgloss.Color("#A3BE8C"),
		Warning:    lipgloss.Color("#EBCB8B"),
		Error:      lipgloss.Color("#BF616A"),
		Text:       lipgloss.Color("#ECEFF4"),
		TextMuted:  lipgloss.Color("#D8DEE9"),
		Background: lipgloss.Color("#2E3440"),
		Border:     lipgloss.Color("#88C0D0"),
		InputBg:    lipgloss.Color("#3B4252"),
	}

	currentTheme = DarkTheme
)

func GetTheme(name string) Theme {
	switch name {
	case "light":
		return LightTheme
	case "dracula":
		return DraculaTheme
	case "nord":
		return NordTheme
	default:
		return DarkTheme
	}
}

func SetTheme(theme Theme) {
	currentTheme = theme
	updateStyles()
}

func CurrentTheme() Theme {
	return currentTheme
}

func updateStyles() {
	t := currentTheme

	Purple = t.Primary
	Blue = t.Secondary
	Cyan = t.Accent
	Green = t.Success
	Yellow = t.Warning
	Red = t.Error
	White = t.Text
	LightGray = t.TextMuted
	Background = t.Background
	DarkGray = t.InputBg

	ContainerStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Text).
		MarginBottom(1)

	MessageStyle = lipgloss.NewStyle().
		Foreground(t.TextMuted).
		Italic(true).
		MarginBottom(1)

	TimerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Text).
		MarginBottom(1)

	ProgressFilledStyle = lipgloss.NewStyle().
		Foreground(t.Primary)

	ProgressEmptyStyle = lipgloss.NewStyle().
		Foreground(t.InputBg)

	PercentStyle = lipgloss.NewStyle().
		Foreground(t.TextMuted).
		MarginLeft(1)

	HelpStyle = lipgloss.NewStyle().
		Foreground(t.TextMuted).
		MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)

	RunningStyle = lipgloss.NewStyle().
		Foreground(t.Success).
		Bold(true)

	PausedStyle = lipgloss.NewStyle().
		Foreground(t.Warning).
		Bold(true)

	CompletedStyle = lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)

	ActivityStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		Italic(true)

	InputStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.InputBg).
		Padding(0, 1)

	SummaryHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		MarginBottom(1)

	SummaryItemStyle = lipgloss.NewStyle().
		Foreground(t.Text)

	SummaryValueStyle = lipgloss.NewStyle().
		Foreground(t.Accent).
		Bold(true)
}

func AvailableThemes() []string {
	return []string{"dark", "light", "dracula", "nord"}
}
