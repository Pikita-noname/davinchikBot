package models

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PasswordRequest struct {
	View *tview.Flex
}

func (a *App) NewPasswordRequest(onSubmit func(password string), errorMsg string) *tview.Flex {
	passwordField := tview.NewInputField().SetLabel("🔑 Введите пароль 2FA: ")
	passwordField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			onSubmit(passwordField.GetText())
		}
	})

	errorView := tview.NewTextView().
		SetText(errorMsg).
		SetTextColor(tcell.ColorRed).
		SetDynamicColors(true).
		SetWrap(true)

	a.setCustomBorder(passwordField.Box, "2FA")

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
			AddItem(passwordField, 0, 1, true).
			AddItem(errorView, 0, 1, false). // Добавляем поле ошибки
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false),
			0, 1, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)

	flex.SetBackgroundColor(tcell.ColorDefault)

	a.View.SetFocus(passwordField)

	return flex
}
