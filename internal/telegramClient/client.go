package telegramclient

import (
	"context"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

type ViewApp interface {
	GetViewApp() *tview.Application
	NewPasswordRequest(onSubmit func(password string), errorMsg string) *tview.Flex
	SetMainView()
	BackToMain()
}

func Run(qrView *tview.TextView, app ViewApp) {
	ctx := context.Background()

	godotenv.Load()

	sessionStorage := &telegram.FileSessionStorage{
		Path: "session.json",
	}

	d := tg.NewUpdateDispatcher()

	d.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		if msg, ok := update.Message.(*tg.Message); ok {
			app.GetViewApp().QueueUpdateDraw(func() {
				qrView.SetText("Новое сообщение: " + msg.Message)
			})
		}
		return nil
	})

	d.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		if msg, ok := update.Message.(*tg.Message); ok {
			app.GetViewApp().QueueUpdateDraw(func() {
				qrView.SetText("Новое сообщение в канале: " + msg.Message)
			})
		}

		return nil
	})

	gaps := updates.New(updates.Config{
		Handler: d,
		Logger:  nil,
	})

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		UpdateHandler:  gaps,
		SessionStorage: sessionStorage,
	})

	if err != nil {
		app.GetViewApp().Stop()
	}

	// Запускаем клиент
	if err := client.Run(ctx, func(ctx context.Context) error {

		status, err := client.Auth().Status(ctx)

		if err != nil {
			qrView.SetText(err.Error())
		}

		if !status.Authorized {
			QrAuth(ctx, client, app, qrView, &d)
		}

		user, err := client.Self(ctx)

		if err != nil {
			qrView.SetText(err.Error())
		}

		app.GetViewApp().QueueUpdateDraw(func() {
			app.SetMainView()
		})

		gaps.Run(ctx, client.API(), user.ID, updates.AuthOptions{})

		return nil
	}); err != nil {
		log.Fatal("Ошибка запуска клиента:", err)
	}

}
