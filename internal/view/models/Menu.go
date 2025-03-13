package models

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainMenu struct {
	View *tview.Flex
}

// Функция для создания главного меню
func (a *App) CreateMainMenu() MainMenu {
	list := tview.NewList().
		AddItem("Start", "", 0, nil).
		AddItem("Options", "", 0, nil).
		AddItem("Exit", "", 0, nil)

	a.setStyles(list)

	a.setCustomBorder(list, "Davinity")

	// Устанавливаем отступы и размеры
	list.SetBorderPadding(1, 1, 2, 2)

	// Создаем Flex для центрирования списка
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
			AddItem(list, 0, 5, true).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false),
			0, 1, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)

	// Устанавливаем фон для всего Flex
	flex.SetBackgroundColor(tcell.ColorDefault)

	// Обработчик выбора пунктов меню
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch mainText {
		case "Start":
			fmt.Println("Start selected!")
		case "Options":
			a.View.SetRoot(a.Options.View, true)
		case "Exit":
			a.View.Stop()
		}
	})

	return MainMenu{
		View: flex,
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

func (a App) setCustomBorder(list *tview.List, title string) {
	list.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Внешняя граница

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

		// Рассчитываем позицию для центрального выравнивания
		startX := x + (width-len(title))/2

		// Отображаем заголовок в центре
		for i, r := range title {
			screen.SetContent(startX+i, y-1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}

		return x + 2, y + 2, width - 4, height - 4
	})
}
