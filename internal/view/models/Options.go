package models

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Options struct {
	View *tview.Flex
}

func (a *App) CreateOptionsMenu() Options {
	list := tview.NewList().
		AddItem("фильтр", "", 0, nil).
		AddItem("аккаунт телеграмма", "", 0, nil).
		AddItem("Назад", "", 0, nil)

	a.setStyles(list)

	a.setCustomBorder(list.Box, "options")

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

	form := a.newOptionForm()

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			a.View.SetRoot(form, true)
		case 2:
			a.View.SetRoot(a.MainMenu.View, true)
			list.SetCurrentItem(0)
		}
	})

	return Options{
		View: flex,
	}
}

func (a *App) newOptionForm() *tview.Form {
	// Форма для редактирования
	form := tview.NewForm()
	a.setCustomBorder(form.Box, "test")
	form.SetBorder(true).SetTitle("Edit Option1").SetTitleAlign(tview.AlignLeft)

	// Загружаем данные из конфига
	name, filter := loadConfig()

	form.AddInputField("Name", name, 20, nil, nil).
		AddInputField("Filter", filter, 20, nil, nil).
		AddButton("Save", func() {

			name := form.GetFormItem(0).(*tview.InputField).GetText()
			filter := form.GetFormItem(1).(*tview.InputField).GetText()
			saveConfig(name, filter)

			a.View.SetRoot(a.Options.View, true)
		}).
		AddButton("Cancel", func() {
			a.View.SetRoot(a.Options.View, true)
		})
	return form
}

func loadConfig() (string, string) {
	return "Default Name", "Default Filter"
}

func saveConfig(name, filter string) {
	// fmt.Printf("Saved: Name=%s, Filter=%s\n", name, filter)
}
