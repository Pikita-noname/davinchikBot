package view

import (
	"github.com/Pikita-noname/davinchikTgApp/internal/view/models"
	tea "github.com/charmbracelet/bubbletea"
)

func NewApp() (m tea.Model) {
	return models.NewMenu()
}
