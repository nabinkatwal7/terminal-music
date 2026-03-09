package ui

import "github.com/charmbracelet/lipgloss"

var (
	Background = lipgloss.Color("#0f1419")
	PanelBg    = lipgloss.Color("#162028")
	Border     = lipgloss.Color("#2f4b5a")
	Primary    = lipgloss.Color("#8fd3ff")
	Accent     = lipgloss.Color("#ffcc80")
	Success    = lipgloss.Color("#a5d6a7")
	Muted      = lipgloss.Color("#90a4ae")
	Text       = lipgloss.Color("#e3edf5")
	Danger     = lipgloss.Color("#ef9a9a")

	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted)

	FocusedItemStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Danger).
			Bold(true)

	RootStyle = lipgloss.NewStyle().
			Background(Background).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Background(PanelBg).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	NowPlayingStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Accent).
			Padding(1, 2)

	ControlStyle = lipgloss.NewStyle().
			Foreground(Muted).
			PaddingTop(1)
)

func GetBanner() string {
	return "TERMINAL MUSIC"
}
