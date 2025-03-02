package models

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	telegramBlue     = lipgloss.Color("#0088CC")
	telegramWhite    = lipgloss.Color("#FFFFFF")
	telegramDarkGray = lipgloss.Color("#2A2A2A")
)

var container = lipgloss.NewStyle().Padding(2, 6).Border(lipgloss.RoundedBorder()).BorderForeground(telegramBlue)
