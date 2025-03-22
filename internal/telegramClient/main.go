package telegramclient

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
	"github.com/skip2/go-qrcode"
)

type ViewApp interface {
	GetViewApp() *tview.Application
	NewPasswordRequest(onSubmit func(password string)) *tview.Flex
	BackToMain()
}

func Run(qrView *tview.TextView, app ViewApp) {
	ctx := context.Background()

	godotenv.Load()

	sessionStorage := &telegram.FileSessionStorage{
		Path: "session.json",
	}

	d := tg.NewUpdateDispatcher()

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		UpdateHandler:  d,
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
			QrAuth(ctx, client, app, qrView)
		}

		user, err := client.Self(ctx)

		if err != nil {
			qrView.SetText(err.Error())
		}

		app.GetViewApp().QueueUpdateDraw(func() {
			qrView.SetText(user.Username)
		})

		return nil
	}); err != nil {
		log.Fatal("Ошибка запуска клиента:", err)
	}
}

func QrAuth(ctx context.Context, client *telegram.Client, app ViewApp, qrView *tview.TextView) (*tg.AuthAuthorization, error) {

	viewApp := app.GetViewApp()

	d := tg.NewUpdateDispatcher()
	loggedIn := qrlogin.OnLoginToken(d)

	authClient := auth.NewClient(client.API(), rand.Reader, 29708230, "428ef65a36ade933259c4c832cd65bfd")

	qr := client.QR()

	// Объявляем authorization один раз в начале
	var authorization *tg.AuthAuthorization
	var err error

	// Первая попытка QR-аутентификации
	authorization, err = qr.Auth(ctx, loggedIn, func(ctx context.Context, token qrlogin.Token) error {
		qrData := token.URL()

		// Генерируем и выводим QR-код
		qr, err := qrcode.New(qrData, qrcode.Medium)
		if err != nil {
			return fmt.Errorf("ошибка генерации QR-кода: %w", err)
		}

		viewApp.QueueUpdateDraw(func() {
			qrView.SetText(qr.ToSmallString(false))
		})

		return nil
	})

	if err != nil {
		if tgerr.Is(err, "SESSION_PASSWORD_NEEDED") {

			passwordChan := make(chan string)

			viewApp.QueueUpdateDraw(func() {
				qrView.SetText("📢 Отображение окна ввода пароля")

				viewApp.SetRoot(app.NewPasswordRequest(func(password string) {
					app.BackToMain()
					passwordChan <- password
				}), true)

			})

			password := <-passwordChan

			authorization, err = authClient.Password(ctx, password)

			if err != nil {
				if errors.Is(err, auth.ErrPasswordInvalid) {
					qrView.SetText("❌ Введен неверный пароль")
					viewApp.Stop()
				} else {
					qrView.SetText("❌ Ошибка при аутентификации паролем: " + err.Error())
					viewApp.Stop()
				}
			} else {
				qrView.SetText("✅ Пароль успешно принят!")
			}
		}
	}

	return authorization, nil
}
