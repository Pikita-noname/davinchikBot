package models

import (
	telegramclient "github.com/Pikita-noname/davinchikTgApp/internal/telegramClient"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Auth struct {
	View   *tview.Flex
	QRView *tview.TextView
}

type Telegram struct {
	Auth Auth
}

func (a *App) NewTelegram() Telegram {
	qrView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Запускаемся...")

	a.setCustomBorder(qrView.Box, "QR Code")

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
			AddItem(qrView, 30, 5, true).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false),
			0, 5, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)

	flex.SetBackgroundColor(tcell.ColorDefault)

	return Telegram{
		Auth: Auth{
			View:   flex,
			QRView: qrView,
		},
	}
}

func (t *Telegram) Run(app *App) {
	go telegramclient.Run(t.Auth.QRView, app)
}
