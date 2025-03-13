package models

import "github.com/rivo/tview"

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
