package view

import (
	"github.com/Pikita-noname/davinchikTgApp/internal/view/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewApp() models.App {

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:   tcell.ColorDefault, // Фон всех примитивов — прозрачный
		BorderColor:                tcell.ColorBlue,    // Цвет границы — синий (как у вас)
		TitleColor:                 tcell.ColorWhite,   // Цвет заголовков — белый
		GraphicsColor:              tcell.ColorDefault, // Цвет графики — прозрачный
		PrimaryTextColor:           tcell.ColorWhite,   // Основной цвет текста — белый
		SecondaryTextColor:         tcell.ColorDefault, // Второстепенный цвет текста — прозрачный
		TertiaryTextColor:          tcell.ColorDefault, // Третичный цвет текста — прозрачный
		InverseTextColor:           tcell.ColorWhite,   // Цвет текста при инверсии — белый
		ContrastBackgroundColor:    tcell.ColorDefault, // Контрастный фон — прозрачный
		ContrastSecondaryTextColor: tcell.ColorDefault, // Контрастный второстепенный текст — прозрачный
	}

	app := models.NewApp()

	app.MainMenu = app.CreateMainMenu()
	app.Options = app.CreateOptionsMenu()

	return app
}
