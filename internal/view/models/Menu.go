package models

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainMenu struct {
	View *tview.Flex
}

func (a *App) CreateMainMenu() MainMenu {
	list := tview.NewList().
		AddItem("Начать!", "", 0, nil).
		AddItem("Настройки", "", 0, nil).
		AddItem("Выход", "", 0, nil)

	a.setStyles(list)

	a.setCustomBorder(list.Box, "Davinity")

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

	flex.SetBackgroundColor(tcell.ColorDefault)

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			fmt.Println("Start selected!")
		case 1:
			a.View.SetRoot(a.Options.View, true)
		case 2:
			a.View.Stop()
		}
	})

	return MainMenu{
		View: flex,
	}
}
