package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui"
)

func main() {
	gameEngine := engine.NewGame("Sang Petualang", "Oakhaven")
	gameEngine.CurrentState = engine.StateTitleScreen
	appModel := ui.NewAppModel(gameEngine)

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Terjadi kesalahan saat menjalankan game: %v\n", err)
		os.Exit(1)
	}
}
