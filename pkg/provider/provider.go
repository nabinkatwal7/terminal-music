package provider

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

// StaticProvider implements the Provider interface with curated lofi streams.
type StaticProvider struct{}

// NewStaticProvider returns a new StaticProvider.
func NewStaticProvider() *StaticProvider {
	return &StaticProvider{}
}

func (p *StaticProvider) GetName() string {
	return "Curated Lofi"
}

// GetPlaylists returns a list of predefined lofi playlists.
func (p *StaticProvider) GetPlaylists() ([]Playlist, error) {
	return []Playlist{
		{
			ID:          "lofi-focus",
			Name:        "Lofi & Study Beats",
			Description: "Smooth beats for deep coding sessions.",
			Tracks: []Track{
				{
					ID:     "lofi-radio-ru",
					Title:  "Lofi Radio",
					Artist: "lofiradio.ru",
					URL:    "https://lofiradio.ru/stream",
				},
				{
					ID:     "chillhop-flux",
					Title:  "Chillhop Radio",
					Artist: "FluxFM",
					URL:    "https://streams.fluxfm.de/Chillhop/mp3-128/streams.fluxfm.de/",
				},
				{
					ID:     "lofi-hiphop",
					Title:  "Lofi HipHop",
					Artist: "hearme.fm",
					URL:    "https://hearme.fm/radio/lofi-hiphop",
				},
			},
		},
		{
			ID:          "classical-piano",
			Name:        "Classical & Piano",
			Description: "Timeless masterpieces for deep work.",
			Tracks: []Track{
				{
					ID:     "klassik-piano",
					Title:  "Klassik Radio Piano",
					Artist: "Klassik Radio",
					URL:    "https://stream.klassikradio.de/piano/mp3-128",
				},
				{
					ID:     "radioparadise-mellow",
					Title:  "Mellow Mix",
					Artist: "Radio Paradise",
					URL:    "http://stream.radioparadise.com/mellow-128",
				},
			},
		},
		{
			ID:          "ambient-focus",
			Name:        "Ambient & Nature",
			Description: "Natural soundscapes for focus.",
			Tracks: []Track{
				{
					ID:     "nature-ambient",
					Title:  "Nature Ambient",
					Artist: "Ambience Radio",
					URL:    "https://stream.zeno.fm/03n82msm068uv",
				},
				{
					ID:     "deep-focus-ambient",
					Title:  "Deep Focus",
					Artist: "Focus FM",
					URL:    "https://stream.zeno.fm/f3pv487220hvv",
				},
			},
		},
	}, nil
}
