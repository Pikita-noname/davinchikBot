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

type Davinchi struct {
	View        *tview.Flex
	Title       *tview.TextView
	Description *tview.TextView
}

type Telegram struct {
	Auth Auth
	Main Davinchi
}

func (a *App) NewTelegram() Telegram {

	return Telegram{
		Auth: a.newAuthView(),
		Main: a.newDavinchiView(),
	}
}

func (t *Telegram) Run(app *App) {
	go telegramclient.Run(app)
}

func (a *App) newDavinchiView() Davinchi {

	title := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Имя, возраст")

	description := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Описание")

	// Контейнер для текста и изображения (вертикальный)
	contentContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).                                                 // Устанавливаем вертикальное направление
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false). // Отступ сверху
		AddItem(title, 0, 1, false).                                                 // Заголовок
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 1, 0, false).
		AddItem(description, 0, 1, false)

	// Устанавливаем рамку
	a.setCustomBorder(contentContainer.Box, "Davinchi")

	// Основной контейнер
	mainFlex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false). // Левый отступ
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 1, 0, false). // Верхний отступ
			AddItem(contentContainer, 0, 5, true).                                       // Основной контент
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 0, false), // Нижний отступ
												0, 5, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false) // Правый отступ

	mainFlex.SetBackgroundColor(tcell.ColorDefault)

	return Davinchi{
		View:        mainFlex,
		Title:       title,
		Description: description,
	}
}

func (a *App) newAuthView() Auth {
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

	return Auth{
		View:   flex,
		QRView: qrView,
	}
}

func (t Telegram) SetTitle(text string) {
	t.Main.Title.SetText(text)
}

func (t Telegram) SetDescription(text string) {
	t.Main.Description.SetText(text)
}
