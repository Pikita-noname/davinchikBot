package models

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	View *tview.Application
	MainMenu
	Options
	AuthFlow
}

func NewApp() App {
	tviewApp := tview.NewApplication()

	return App{
		View: tviewApp,
	}
}

func (a App) setStyles(list *tview.List) {
	list.SetBackgroundColor(tcell.ColorDefault)
	list.SetMainTextColor(tcell.ColorGray)
	list.SetSecondaryTextColor(tcell.ColorDefault)
	list.SetSelectedTextColor(tcell.ColorWhite)
	list.SetSelectedBackgroundColor(tcell.ColorDefault)
	list.SetSelectedStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorDefault).
		Bold(true))
}

func (a App) setCustomBorder(box *tview.Box, title string) {
	box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

		for i := x; i < x+width; i++ {
			screen.SetContent(i, y, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
			screen.SetContent(i, y+height-1, '─', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
		}
		for j := y; j < y+height; j++ {
			screen.SetContent(x, j, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
			screen.SetContent(x+width-1, j, '│', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
		}
		screen.SetContent(x, y, '┌', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
		screen.SetContent(x+width-1, y, '┐', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
		screen.SetContent(x, y+height-1, '└', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
		screen.SetContent(x+width-1, y+height-1, '┘', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))

		title := title

		startX := x + (width-len(title))/2

		for i, r := range title {
			screen.SetContent(startX+i, y-1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}

		return x + 2, y + 2, width - 4, height - 4
	})
}
