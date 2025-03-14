package view

import (
	"github.com/Pikita-noname/davinchikTgApp/internal/view/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewApp() models.App {

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:   tcell.ColorDefault,
		BorderColor:                tcell.ColorBlue,
		TitleColor:                 tcell.ColorWhite,
		GraphicsColor:              tcell.ColorDefault,
		PrimaryTextColor:           tcell.ColorWhite,
		SecondaryTextColor:         tcell.ColorDefault,
		TertiaryTextColor:          tcell.ColorDefault,
		InverseTextColor:           tcell.ColorWhite,
		ContrastBackgroundColor:    tcell.ColorDefault,
		ContrastSecondaryTextColor: tcell.ColorDefault,
	}

	app := models.NewApp()

	app.MainMenu = app.CreateMainMenu()
	app.Options = app.CreateOptionsMenu()
	app.Telegram = app.NewTelegram()

	return app
}
