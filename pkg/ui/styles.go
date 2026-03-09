package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	MidnightBlue = lipgloss.Color("#1e1e2e")
	NeonPurple   = lipgloss.Color("#c678dd")
	Cyan         = lipgloss.Color("#61afef")
	Gray         = lipgloss.Color("#5c6370")
	Green        = lipgloss.Color("#98c379")
	White        = lipgloss.Color("#abb2bf")
	Amber        = lipgloss.Color("#d19a66")

	// Dynamic Banner Colors
	BannerColors = []lipgloss.Color{
		lipgloss.Color("#c678dd"), // Purple
		lipgloss.Color("#61afef"), // Blue
		lipgloss.Color("#98c379"), // Green
		lipgloss.Color("#d19a66"), // Amber
		lipgloss.Color("#e06c75"), // Red
		lipgloss.Color("#56b6c2"), // Cyan
	}

	// Styles
	MainStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(NeonPurple)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(NeonPurple).
			Bold(true)

	StatusStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	ArtistStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	ActiveStyle = lipgloss.NewStyle().
			Foreground(Amber).
			Bold(true).
			Background(MidnightBlue)

	InactiveStyle = lipgloss.NewStyle().
			Foreground(White)

	ControlStyle = lipgloss.NewStyle().
			Foreground(Gray).
			MarginTop(1).
			Italic(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			BorderForeground(Gray)
)

func DrawBox(width int, content string) string {
	return MainStyle.Width(width).Render(content)
}

func GetBanner() string {
	return `
  ▀█▀ █▀▀ █▀▀█ █▀▄▀█ ░▀░ █▀▀▄ █▀▀█ █── █── █──█
  ─█─ █▀▀ █▄▄▀ █─▀─█ ▀█▀ █──█ █▄▄█ █── █── █▄▄█
  ─▀─ ▀▀▀ ▀─▀▀ ▀───▀ ▀▀▀ ▀──▀ ▀──▀ ▀▀▀ ▀▀▀ ▄▄▄█
       ~ Minimalist Focus Music for Hackers ~
`
}
