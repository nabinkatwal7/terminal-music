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
	Focused     string // playlists or tracks
	ProgressBar progress.Model
	Width       int
	Height      int
	Error       error
	Quitting    bool
	Volume      float64 // 0 - 100
}

func NewModel(engine *audio.Engine, prov provider.Provider) Model {
	playlists, err := prov.GetPlaylists()

	pb := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(30),
	)

	return Model{
		Engine:      engine,
		Provider:    prov,
		Playlists:   playlists,
		Focused:     "playlists",
		ProgressBar: pb,
		Volume:      50,
		Error:       err,
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
				m.SelectedTr = 0
			} else if m.Focused == "tracks" && m.SelectedTr > 0 {
				m.SelectedTr--
			}
		case "down", "j":
			if m.Focused == "playlists" && m.SelectedPl < len(m.Playlists)-1 {
				m.SelectedPl++
				m.SelectedTr = 0
			} else if m.Focused == "tracks" && m.SelectedTr < len(m.currentTracks())-1 {
				m.SelectedTr++
			}
		case "tab", "right", "l":
			m.Focused = "tracks"
		case "shift+tab", "left", "h":
			m.Focused = "playlists"
		case "enter", " ":
			track, ok := m.currentTrack()
			if !ok {
				break
			}
			m.Error = m.Engine.Play(track.URL)
		case "s":
			m.Engine.Stop()
		case "+", "=":
			if m.Volume < 100 {
				m.Volume += 5
				m.Engine.SetVolume(m.Volume)
			}
		case "-", "_":
			if m.Volume > 0 {
				m.Volume -= 5
				m.Engine.SetVolume(m.Volume)
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if msg.Width > 20 {
			m.ProgressBar.Width = msg.Width - 24
		}

	case tickMsg:
		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return "\n  " + TitleStyle.Render("Terminal Music") + " - bye\n\n"
	}

	if len(m.Playlists) == 0 {
		errMsg := "No playlists loaded"
		if m.Error != nil {
			errMsg = m.Error.Error()
		}
		return RootStyle.Render(TitleStyle.Render(GetBanner()) + "\n\n" + ErrorStyle.Render(errMsg))
	}

	pl := m.currentPlaylist()
	track, hasTrack := m.currentTrack()

	leftColWidth := max(28, m.Width/3)
	rightColWidth := max(36, m.Width-leftColWidth-10)

	var playlists strings.Builder
	playlists.WriteString(TitleStyle.Render("Playlists") + "\n")
	playlists.WriteString(SubtitleStyle.Render("Source: "+m.Provider.GetName()) + "\n\n")
	for i, item := range m.Playlists {
		prefix := "  "
		style := MutedStyle
		if i == m.SelectedPl && m.Focused == "playlists" {
			prefix = "> "
			style = FocusedItemStyle
		} else if i == m.SelectedPl {
			prefix = "* "
			style = SelectedItemStyle
		}
		playlists.WriteString(style.Render(prefix+item.Name) + "\n")
	}

	var tracks strings.Builder
	tracks.WriteString(TitleStyle.Render("Tracks") + "\n")
	tracks.WriteString(SubtitleStyle.Render(pl.Description) + "\n\n")
	for i, t := range pl.Tracks {
		prefix := "  "
		style := MutedStyle
		if i == m.SelectedTr && m.Focused == "tracks" {
			prefix = "> "
			style = FocusedItemStyle
		} else if i == m.SelectedTr {
			prefix = "* "
			style = SelectedItemStyle
		}
		tracks.WriteString(style.Render(prefix+t.Title) + "\n")
		tracks.WriteString(MutedStyle.Render("   "+t.Artist) + "\n")
	}

	sidebar := PanelStyle.Width(leftColWidth).Render(playlists.String())
	main := PanelStyle.Width(rightColWidth).Render(tracks.String())
	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, "  ", main)

	status := "Stopped"
	if m.Engine.IsPlaying() {
		status = "Streaming"
	}

	duration, _ := m.Engine.GetProgress()
	progressValue := float64((duration / time.Second) % 180)
	progressValue = progressValue / 180.0
	progressBar := m.ProgressBar.ViewAs(progressValue)

	nowPlayingText := "Nothing playing"
	if hasTrack {
		nowPlayingText = fmt.Sprintf("%s - %s", track.Title, track.Artist)
	}

	nowPlaying := strings.Builder{}
	nowPlaying.WriteString(TitleStyle.Render(GetBanner()) + "\n")
	nowPlaying.WriteString(SubtitleStyle.Render("YouTube stream player") + "\n\n")
	nowPlaying.WriteString(fmt.Sprintf("Status: %s\n", lipgloss.NewStyle().Foreground(Success).Render(status)))
	nowPlaying.WriteString(fmt.Sprintf("Now: %s\n", nowPlayingText))
	nowPlaying.WriteString(fmt.Sprintf("Elapsed: %s\n", formatDuration(duration)))
	nowPlaying.WriteString(fmt.Sprintf("Volume: %.0f%%\n\n", m.Volume))
	nowPlaying.WriteString(progressBar)

	controls := "Enter/Space play  s stop  Tab switch panel  Arrows move  +/- volume  q quit"

	view := strings.Builder{}
	view.WriteString(content + "\n\n")
	view.WriteString(NowPlayingStyle.Width(m.Width - 6).Render(nowPlaying.String()))
	view.WriteString("\n")
	view.WriteString(ControlStyle.Render(controls))

	if m.Error != nil {
		view.WriteString("\n" + ErrorStyle.Render("Error: "+m.Error.Error()))
	}

	return RootStyle.Render(view.String())
}

func (m Model) currentPlaylist() provider.Playlist {
	if len(m.Playlists) == 0 {
		return provider.Playlist{}
	}
	if m.SelectedPl < 0 {
		m.SelectedPl = 0
	}
	if m.SelectedPl >= len(m.Playlists) {
		m.SelectedPl = len(m.Playlists) - 1
	}
	return m.Playlists[m.SelectedPl]
}

func (m Model) currentTracks() []provider.Track {
	pl := m.currentPlaylist()
	return pl.Tracks
}

func (m Model) currentTrack() (provider.Track, bool) {
	tracks := m.currentTracks()
	if len(tracks) == 0 {
		return provider.Track{}, false
	}
	if m.SelectedTr < 0 {
		m.SelectedTr = 0
	}
	if m.SelectedTr >= len(tracks) {
		m.SelectedTr = len(tracks) - 1
	}
	return tracks[m.SelectedTr], true
}

func tick() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		return "00:00"
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
