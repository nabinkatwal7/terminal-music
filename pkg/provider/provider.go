package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Track represents a single audio track or stream.
type Track struct {
	ID     string
	Title  string
	Artist string
	URL    string
}

// Playlist represents a collection of tracks.
type Playlist struct {
	ID          string
	Name        string
	Description string
	Tracks      []Track
}

// Provider defines the interface for fetching music.
type Provider interface {
	GetPlaylists() ([]Playlist, error)
	GetName() string
}

type youtubeQuery struct {
	ID          string
	Name        string
	Description string
	Query       string
}

var defaultQueries = []youtubeQuery{
	{
		ID:          "lofi-focus",
		Name:        "Lofi Coding",
		Description: "Lo-fi beats for focused coding sessions.",
		Query:       "lofi hip hop beats to relax/study to",
	},
	{
		ID:          "jazz-cafe",
		Name:        "Jazz Cafe",
		Description: "Instrumental jazz and cafe ambience for flow state.",
		Query:       "instrumental jazz cafe music",
	},
	{
		ID:          "ambient-focus",
		Name:        "Ambient Focus",
		Description: "Deep ambient textures for distraction-free work.",
		Query:       "ambient focus music no lyrics",
	},
}

// YouTubeProvider implements the Provider interface using yt-dlp search.
type YouTubeProvider struct {
	queries []youtubeQuery
}

// NewYouTubeProvider initializes a YouTube provider without API keys.
func NewYouTubeProvider() *YouTubeProvider {
	return &YouTubeProvider{queries: defaultQueries}
}

func (p *YouTubeProvider) GetName() string {
	return "YouTube (no API key)"
}

// GetPlaylists builds curated categories from YouTube search results.
func (p *YouTubeProvider) GetPlaylists() ([]Playlist, error) {
	playlists := make([]Playlist, 0, len(p.queries))

	for _, q := range p.queries {
		tracks, err := p.searchTracks(q.Query)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, Playlist{
			ID:          q.ID,
			Name:        q.Name,
			Description: q.Description,
			Tracks:      tracks,
		})
	}

	return playlists, nil
}

type ytDLPEntry struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Uploader string `json:"uploader"`
	Channel  string `json:"channel"`
}

func (p *YouTubeProvider) searchTracks(query string) ([]Track, error) {
	ytDLPPath, err := exec.LookPath("yt-dlp")
	if err != nil {
		return nil, fmt.Errorf("yt-dlp is required in PATH for YouTube search")
	}

	cmd := exec.Command(ytDLPPath, "--dump-json", "--flat-playlist", "ytsearch12:"+query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp search failed: %w", err)
	}

	lines := bytes.Split(out, []byte{'\n'})
	tracks := make([]Track, 0, len(lines))
	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		var entry ytDLPEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}
		if entry.ID == "" || entry.Title == "" {
			continue
		}

		artist := strings.TrimSpace(entry.Uploader)
		if artist == "" {
			artist = strings.TrimSpace(entry.Channel)
		}
		if artist == "" {
			artist = "YouTube"
		}

		tracks = append(tracks, Track{
			ID:     entry.ID,
			Title:  entry.Title,
			Artist: artist,
			URL:    "https://www.youtube.com/watch?v=" + entry.ID,
		})
	}

	if len(tracks) == 0 {
		return nil, fmt.Errorf("no YouTube tracks found for query %q", query)
	}

	return tracks, nil
}
