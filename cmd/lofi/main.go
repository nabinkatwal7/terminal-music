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
	// Initialize Engine
	engine, err := audio.NewEngine()
	if err != nil {
		fmt.Printf("Error initializing audio engine: %v\n", err)
		os.Exit(1)
	}
	defer engine.Stop()

	// Initialize Provider
	prov := provider.NewStaticProvider()

	// Initialize UI Model
	m := ui.NewModel(engine, prov)

	// Run Bubble Tea
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
