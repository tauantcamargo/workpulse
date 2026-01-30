package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	FilledBlock = "█"
	EmptyBlock  = "░"
	BarWidth    = 40
)

type GradientColor struct {
	R, G, B int
}

var (
	GradientStart = GradientColor{139, 92, 246}  // Purple #8B5CF6
	GradientEnd   = GradientColor{59, 130, 246}  // Blue #3B82F6
)

func interpolateColor(start, end GradientColor, t float64) lipgloss.Color {
	r := int(float64(start.R) + t*float64(end.R-start.R))
	g := int(float64(start.G) + t*float64(end.G-start.G))
	b := int(float64(start.B) + t*float64(end.B-start.B))

	return lipgloss.Color(rgbToHex(r, g, b))
}

func rgbToHex(r, g, b int) string {
	return "#" + toHex(r) + toHex(g) + toHex(b)
}

func toHex(n int) string {
	if n < 0 {
		n = 0
	}
	if n > 255 {
		n = 255
	}
	hex := "0123456789ABCDEF"
	return string(hex[n/16]) + string(hex[n%16])
}

func RenderProgressBar(progress float64, width int) string {
	if width <= 0 {
		width = BarWidth
	}

	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}

	var result strings.Builder

	for i := 0; i < filled; i++ {
		t := float64(i) / float64(width)
		color := interpolateColor(GradientStart, GradientEnd, t)
		style := lipgloss.NewStyle().Foreground(color)
		result.WriteString(style.Render(FilledBlock))
	}

	emptyStyle := lipgloss.NewStyle().Foreground(DarkGray)
	for i := filled; i < width; i++ {
		result.WriteString(emptyStyle.Render(EmptyBlock))
	}

	return result.String()
}

func RenderProgressWithPercent(progress float64, width int) string {
	bar := RenderProgressBar(progress, width)
	percent := int(progress * 100)
	percentStr := PercentStyle.Render(padPercent(percent) + "%")
	return bar + " " + percentStr
}

func padPercent(n int) string {
	if n < 10 {
		return "  " + itoa(n)
	}
	if n < 100 {
		return " " + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
