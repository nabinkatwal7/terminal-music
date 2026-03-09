package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nabinkatwal7/terminal-music/pkg/audio"
	"github.com/nabinkatwal7/terminal-music/pkg/provider"
	"github.com/nabinkatwal7/terminal-music/pkg/ui"
)

func main() {
	engine, err := audio.NewEngine()
	if err != nil {
		fmt.Printf("Error initializing audio engine: %v\n", err)
		os.Exit(1)
	}
	defer engine.Stop()

	prov := provider.NewYouTubeProvider()
	m := ui.NewModel(engine, prov)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
