package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	list list.Model
}

type item struct {
	title       string
	description string
	onEnter     func(tea.Msg, Menu) (tea.Model, tea.Cmd)
}

func (i item) OnEnter(msg tea.Msg, m Menu) (tea.Model, tea.Cmd) {
	if i.onEnter == nil {
		return m, nil
	}
	return i.onEnter(msg, m)
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func NewMenu() tea.Model {

	items := []list.Item{
		item{title: "start"},
		item{title: "options", onEnter: func(msg tea.Msg, m Menu) (tea.Model, tea.Cmd) {
			optionItems := []list.Item{
				item{title: "option1"},
				item{title: "option2"},
			}

			return m, m.list.SetItems(optionItems)
		}},
		item{title: "exit", onEnter: func(msg tea.Msg, m Menu) (tea.Model, tea.Cmd) {
			return NewExitPage(), NewExitPage().Init()
		}},
	}

	return Menu{list: NewCustomList(items, "Davinchi cheat")}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC, tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyEnter:
			return m.EnterHandler(msg)
		}
	case tea.WindowSizeMsg:
		h, v := container.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Menu) EnterHandler(msg tea.Msg) (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		return i.OnEnter(msg, m)
	}
	return m, nil
}

func (m Menu) View() string {
	return container.Render(m.list.View())
}

func NewCustomDelegate() list.ItemDelegate {
	delegate := list.NewDefaultDelegate()

	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(telegramDarkGray).
		Bold(false).
		Padding(0, 5)

	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(telegramWhite).
		Bold(true).
		Padding(0, 5)

	delegate.Styles.NormalDesc = lipgloss.NewStyle()
	delegate.Styles.SelectedDesc = lipgloss.NewStyle()

	delegate.SetSpacing(0)

	return delegate
}

func NewCustomList(items []list.Item, title string) list.Model {
	l := list.New(items, NewCustomDelegate(), 0, 0)

	l.Title = title
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(telegramBlue).
		Bold(true).
		Padding(0, 1)

	l.SetShowPagination(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	return l
}
