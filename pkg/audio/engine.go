package audio

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Engine manages audio playback by streaming YouTube audio through yt-dlp + ffplay.
type Engine struct {
	mu           sync.Mutex
	ytdlpCmd     *exec.Cmd
	ffplayCmd    *exec.Cmd
	playing      bool
	currentTrack string
	startedAt    time.Time
	volume       int
}

// NewEngine initializes a new Engine.
func NewEngine() (*Engine, error) {
	return &Engine{
		volume: 50,
	}, nil
}

// Play streams audio from a YouTube URL.
func (e *Engine) Play(url string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.playing {
		e.stopInternal()
	}

	ytdlpPath, err := exec.LookPath("yt-dlp")
	if err != nil {
		return fmt.Errorf("yt-dlp is required in PATH for YouTube streaming")
	}
	ffplayPath, err := exec.LookPath("ffplay")
	if err != nil {
		return fmt.Errorf("ffplay is required in PATH for YouTube streaming")
	}

	ytdlp := exec.Command(ytdlpPath,
		"-f", "bestaudio",
		"-o", "-",
		url,
	)

	ffplay := exec.Command(ffplayPath,
		"-nodisp",
		"-autoexit",
		"-loglevel", "quiet",
		"-volume", fmt.Sprintf("%d", e.volume),
		"-i", "pipe:0",
	)

	streamPipe, err := ytdlp.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create yt-dlp stream pipe: %w", err)
	}
	ffplay.Stdin = streamPipe

	var ytdlpErr strings.Builder
	var ffplayErr strings.Builder
	ytdlp.Stderr = &ytdlpErr
	ffplay.Stderr = &ffplayErr

	if err := ffplay.Start(); err != nil {
		return fmt.Errorf("failed to start ffplay: %w", err)
	}

	if err := ytdlp.Start(); err != nil {
		_ = ffplay.Process.Kill()
		return fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	e.ytdlpCmd = ytdlp
	e.ffplayCmd = ffplay
	e.playing = true
	e.currentTrack = url
	e.startedAt = time.Now()

	go e.waitForCompletion(&ytdlpErr, &ffplayErr)

	return nil
}

func (e *Engine) waitForCompletion(ytdlpErr, ffplayErr *strings.Builder) {
	e.mu.Lock()
	ytdlp := e.ytdlpCmd
	ffplay := e.ffplayCmd
	e.mu.Unlock()

	if ytdlp != nil {
		_ = ytdlp.Wait()
	}
	if ffplay != nil {
		_ = ffplay.Wait()
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.playing = false
	e.ytdlpCmd = nil
	e.ffplayCmd = nil
	if e.currentTrack != "" && (ytdlpErr.Len() > 0 || ffplayErr.Len() > 0) {
		// Keep currentTrack for UI context until user starts/stops next playback.
	}
}

// Pause is not supported with ffplay pipe mode.
func (e *Engine) Pause() {
	// Intentionally left as no-op.
}

// Stop stops the playback and kills active processes.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stopInternal()
}

func (e *Engine) stopInternal() {
	if e.ytdlpCmd != nil && e.ytdlpCmd.Process != nil {
		_ = e.ytdlpCmd.Process.Kill()
		_, _ = e.ytdlpCmd.Process.Wait()
		e.ytdlpCmd = nil
	}
	if e.ffplayCmd != nil && e.ffplayCmd.Process != nil {
		_ = e.ffplayCmd.Process.Kill()
		_, _ = e.ffplayCmd.Process.Wait()
		e.ffplayCmd = nil
	}
	e.playing = false
	e.currentTrack = ""
	e.startedAt = time.Time{}
}

// SetVolume sets playback volume from 0-100 for the next playback start.
func (e *Engine) SetVolume(vol float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if vol < 0 {
		vol = 0
	}
	if vol > 100 {
		vol = 100
	}
	e.volume = int(vol)
}

// IsPlaying returns true if audio is currently playing.
func (e *Engine) IsPlaying() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.playing
}

// GetProgress returns elapsed duration for the active stream.
func (e *Engine) GetProgress() (time.Duration, time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.playing || e.startedAt.IsZero() {
		return 0, 0
	}
	return time.Since(e.startedAt), 0
}
