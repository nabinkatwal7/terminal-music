package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nabinkatwal7/terminal-music/pkg/audio"
	"github.com/nabinkatwal7/terminal-music/pkg/provider"
)

type tickMsg time.Time

type Model struct {
	Engine      *audio.Engine
	Provider    provider.Provider
	Playlists   []provider.Playlist
	SelectedPl  int
	SelectedTr  int
	Focused     string // "playlists", "tracks"
	ProgressBar progress.Model
	Width       int
	Height      int
	Error       error
	Quitting    bool
	Volume      float64 // 0 - 100
	ColorIndex  int
}

func NewModel(engine *audio.Engine, prov provider.Provider) Model {
	playlists, _ := prov.GetPlaylists()
	return Model{
		Engine:      engine,
		Provider:    prov,
		Playlists:   playlists,
		Focused:     "playlists",
		ProgressBar: progress.New(progress.WithDefaultGradient()),
		Volume:      50,
		ColorIndex:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			m.Engine.Stop()
			return m, tea.Quit
		case "up", "k":
			if m.Focused == "playlists" && m.SelectedPl > 0 {
				m.SelectedPl--
			} else if m.Focused == "tracks" && m.SelectedTr > 0 {
				m.SelectedTr--
			}
		case "down", "j":
			if m.Focused == "playlists" && m.SelectedPl < len(m.Playlists)-1 {
				m.SelectedPl++
			} else if m.Focused == "tracks" && m.SelectedTr < len(m.Playlists[m.SelectedPl].Tracks)-1 {
				m.SelectedTr++
			}
		case "tab", "right", "l":
			m.Focused = "tracks"
		case "shift+tab", "left", "h":
			m.Focused = "playlists"
		case "enter", " ":
			if m.Focused == "tracks" {
				track := m.Playlists[m.SelectedPl].Tracks[m.SelectedTr]
				m.Error = m.Engine.Play(track.URL)
			}
		case "p":
			m.Engine.Pause()
		case "s":
			m.Engine.Stop()
		case "+", "=":
			if m.Volume < 100 {
				m.Volume += 5
				m.Engine.SetVolume(m.Volume/50 - 1) // Map 0-100 to -1 to 1 (roughly)
			}
		case "-", "_":
			if m.Volume > 0 {
				m.Volume -= 5
				m.Engine.SetVolume(m.Volume/50 - 1)
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.ProgressBar.Width = msg.Width - 10

	case tickMsg:
		if m.Engine.IsPlaying() {
			m.ColorIndex++
		}
		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return "\n  " + TitleStyle.Render("Terminally") + " - Stay focused. See you soon.\n\n"
	}

	var s strings.Builder

	// Dynamic Banner Color
	bannerStyle := TitleStyle
	if m.Engine.IsPlaying() {
		bannerStyle = bannerStyle.Copy().Foreground(BannerColors[m.ColorIndex%len(BannerColors)])
	}

	// Banner
	s.WriteString(bannerStyle.Render(GetBanner()))
	s.WriteString("\n")

	// Main Layout
	leftColWidth := 25
	rightColWidth := m.Width - leftColWidth - 8

	// Playlists Column
	var playlists strings.Builder
	playlists.WriteString(HighlightStyle.Render("   Playlists") + "\n\n")
	for i, pl := range m.Playlists {
		style := InactiveStyle
		prefix := "  "
		if i == m.SelectedPl && m.Focused == "playlists" {
			style = ActiveStyle
			prefix = "-> "
		} else if i == m.SelectedPl {
			style = TitleStyle
			prefix = "  "
		}
		playlists.WriteString(style.Render(prefix+pl.Name) + "\n")
	}

	// Tracks Column
	var tracks strings.Builder
	tracks.WriteString(HighlightStyle.Render("   Tracks: "+m.Playlists[m.SelectedPl].Name) + "\n\n")
	for i, tr := range m.Playlists[m.SelectedPl].Tracks {
		style := InactiveStyle
		prefix := "  "
		if i == m.SelectedTr && m.Focused == "tracks" {
			style = ActiveStyle
			prefix = "-> "
		} else if i == m.SelectedTr {
			style = TitleStyle
			prefix = "  "
		}
		tracks.WriteString(style.Render(prefix+tr.Title) + "\n      " + ArtistStyle.Render(tr.Artist) + "\n")
	}

	// Sidebar + Main View
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(leftColWidth).Render(playlists.String()),
		lipgloss.NewStyle().Width(rightColWidth).PaddingLeft(4).Render(tracks.String()),
	)

	s.WriteString(BoxStyle.Width(m.Width - 4).Render(mainContent) + "\n\n")

	// Now Playing Section
	pos, _ := m.Engine.GetProgress()
	status := "⏸  Paused"
	if m.Engine.IsPlaying() {
		status = "▶  Streaming"
	}

	// Simple animated wave/progress for streams
	wave := "~~~~"
	if m.Engine.IsPlaying() {
		ticks := int(time.Now().Unix() % 4)
		wave = strings.Repeat("~", ticks) + " " + strings.Repeat("~", 4-ticks)
	}

	nowPlaying := fmt.Sprintf("%s | %s | %s [%s]",
		StatusStyle.Render(status),
		m.Playlists[m.SelectedPl].Tracks[m.SelectedTr].Title,
		ArtistStyle.Render(m.Playlists[m.SelectedPl].Tracks[m.SelectedTr].Artist),
		HighlightStyle.Render(formatDuration(pos)),
	)

	s.WriteString(MainStyle.Width(m.Width - 6).Render(nowPlaying + " " + wave) + "\n")

	// Controls
	controls := " q: quit | space: play/pause | s: stop | tab: switch | arrows: move | +/-: volume "
	s.WriteString("\n" + ControlStyle.Width(m.Width).Align(lipgloss.Center).Render(controls))

	if m.Error != nil {
		s.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75")).Bold(true).Render(" ⚡ Error: "+m.Error.Error()))
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(s.String())
}

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func formatDuration(d time.Duration) string {
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}
