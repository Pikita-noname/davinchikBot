package main

import (
	"log"
	"os"

	"github.com/Pikita-noname/davinchikTgApp/internal/view"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := view.NewApp()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
