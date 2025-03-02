package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type OptionsPage struct {
	list list.Model
}

func (m OptionsPage) Init() tea.Cmd {
	return nil
}

func (m OptionsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC, tea.KeyCtrlQ:
			return m, tea.Quit

		case tea.KeyEnter:
			return m.EnterHandler(msg)
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m OptionsPage) EnterHandler(msg tea.Msg) (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		m.list.Title = string(i.Title())
	}
	return m, nil
}

func (m OptionsPage) View() string {
	return container.Render(m.list.View())
}

func NewOptionsPage() tea.Model {

	items := []list.Item{
		item{title: "option1"},
		item{title: "option2"},
		item{title: "option3"},
	}

	return OptionsPage{list: NewCustomList(items, "options")}
}
