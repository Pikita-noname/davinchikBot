package models

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Options struct {
	View *tview.Flex
}

// Функция для создания подменю "Options"
func (a *App) CreateOptionsMenu() Options {
	list := tview.NewList().
		AddItem("Option1", "", 0, nil).
		AddItem("Option2", "", 0, nil).
		AddItem("Back", "", 0, nil)

	a.setStyles(list)

	a.setCustomBorder(list, "options")

	list.SetBorderPadding(1, 1, 2, 2)

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
			AddItem(list, 0, 5, true).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false),
			0, 1, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)
	// Устанавливаем фон для подменю Flex
	flex.SetBackgroundColor(tcell.ColorDefault)

	// Устанавливаем обработчик для возврата в главное меню
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Back" {
			a.View.SetRoot(a.MainMenu.View, true)
			list.SetCurrentItem(0)
		}
	})

	return Options{
		View: flex,
	}
}
