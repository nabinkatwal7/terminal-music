package audio

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// Engine manages the audio playback state.
type Engine struct {
	mu           sync.Mutex
	streamer     beep.StreamSeekCloser
	format       beep.Format
	ctrl         *beep.Ctrl
	volume       *effects.Volume
	sampleRate   beep.SampleRate
	playing      bool
	currentTrack string
	errorChan    chan error
}

// NewEngine initializes the audio speaker and returns a new Engine.
func NewEngine() (*Engine, error) {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return nil, fmt.Errorf("failed to init speaker: %w", err)
	}

	return &Engine{
		sampleRate: sr,
		errorChan:  make(chan error, 1),
	}, nil
}

// Play streams audio from a URL.
func (e *Engine) Play(url string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.playing {
		e.stopInternal()
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Terminally/1.0 (Terminal Music Player; Go)")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// mp3.Decode requires a ReadCloser
	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("failed to decode mp3: %w", err)
	}

	e.streamer = streamer
	e.format = format
	e.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	e.volume = &effects.Volume{Streamer: e.ctrl, Base: 2, Volume: 0, Silent: false}
	e.playing = true
	e.currentTrack = url

	speaker.Play(beep.Seq(e.volume, beep.Callback(func() {
		e.mu.Lock()
		e.playing = false
		e.mu.Unlock()
	})))

	return nil
}

// Pause toggles the pause state.
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.ctrl != nil {
		e.ctrl.Paused = !e.ctrl.Paused
	}
}

// Stop stops the playback and closes the streamer.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stopInternal()
}

func (e *Engine) stopInternal() {
	speaker.Clear()
	if e.streamer != nil {
		e.streamer.Close()
		e.streamer = nil
	}
	e.playing = false
	e.currentTrack = ""
}

// SetVolume sets the volume level (-10.0 to 0.0 usually, where 0 is max).
// We'll abstract this to 0-100 for the UI.
func (e *Engine) SetVolume(vol float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.volume != nil {
		// beep.Volume.Volume is logarithmic. 0 is original, -1 is half amplitude, etc.
		// We'll map UI 0-100 to volume range.
		e.volume.Volume = vol
	}
}

// IsPlaying returns true if audio is currently playing.
func (e *Engine) IsPlaying() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.playing && e.ctrl != nil && !e.ctrl.Paused
}

// GetProgress returns the current position and total duration (if available).
func (e *Engine) GetProgress() (time.Duration, time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.streamer == nil {
		return 0, 0
	}
	pos := e.format.SampleRate.D(e.streamer.Position())
	len := e.format.SampleRate.D(e.streamer.Len())
	return pos, len
}
