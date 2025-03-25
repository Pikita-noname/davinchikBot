package models

import (
	"image/jpeg"
	"os"

	telegramclient "github.com/Pikita-noname/davinchikTgApp/internal/telegramClient"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Auth struct {
	View   *tview.Flex
	QRView *tview.TextView
}

type Davinchi struct {
	View  *tview.Flex
	Image *tview.Image
	Text  *tview.TextView
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
	go telegramclient.Run(t.Auth.QRView, app)
}

func (a *App) newDavinchiView() Davinchi {

	file, err := os.Open("C:\\Users\\kuzne\\OneDrive\\Desktop\\projects\\davincikTgApp\\photo_2025-03-12_23-31-10.jpg")
	if err != nil {
		panic(err) // Обработайте ошибку
	}
	defer file.Close()

	// Декодируем изображение
	photo, err := jpeg.Decode(file)
	if err != nil {
		panic(err) // Обработайте ошибку
	}

	image := tview.NewImage().SetImage(photo) // Установите изображение позже
	textTop := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Верхний текст")

	textBottom := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("Нижний текст")

	textFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(textTop, 0, 1, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 1, 0, false).
		AddItem(textBottom, 0, 1, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)

	contentContainer := tview.NewFlex().
		AddItem(textFlex, 0, 3, true).
		AddItem(image, 0, 4, false)

	a.setCustomBorder(contentContainer.Box, "Davinchi")

	mainFlex := tview.NewFlex().
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 1, 0, false).
			AddItem(contentContainer, 0, 5, true).
			AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 0, false),
			0, 5, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorDefault), 0, 1, false)

	mainFlex.SetBackgroundColor(tcell.ColorDefault)

	return Davinchi{
		View:  mainFlex,
		Image: image,
		Text:  textTop,
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
