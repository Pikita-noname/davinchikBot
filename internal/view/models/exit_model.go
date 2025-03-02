package models

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ExitPage struct {
	count int
}

type tickMsg time.Time

func (m ExitPage) Init() tea.Cmd {
	return tick()
}

func (m ExitPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC, tea.KeyCtrlQ:
			return m, tea.Quit
		}

	case tickMsg:
		m.count = m.count - 1
		if m.count <= 0 {
			return m, tea.Quit
		}
		return m, tick()
	}

	return m, nil
}

func (m ExitPage) View() string {
	return fmt.Sprintf("goodbuy! %d", m.count)
}

func NewExitPage() tea.Model {
	return ExitPage{count: 2}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
