package telegramclient

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/rivo/tview"
	"github.com/skip2/go-qrcode"
)

func QrAuth(ctx context.Context, client *telegram.Client, app ViewApp, qrView *tview.TextView, d *tg.UpdateDispatcher) (*tg.AuthAuthorization, error) {

	viewApp := app.GetViewApp()

	loggedIn := qrlogin.OnLoginToken(d)

	authClient := auth.NewClient(client.API(), rand.Reader, 29708230, "428ef65a36ade933259c4c832cd65bfd")

	qr := client.QR()

	authorization, err := qr.Auth(ctx, loggedIn, func(ctx context.Context, token qrlogin.Token) error {
		qrData := token.URL()

		qr, err := qrcode.New(qrData, qrcode.Medium)
		if err != nil {
			return fmt.Errorf("ошибка генерации QR-кода: %w", err)
		}

		viewApp.QueueUpdateDraw(func() {
			qrView.SetText(qr.ToSmallString(false) + "\n📲 Отсканируйте QR-код в Telegram")
		})

		return nil
	})

	if err != nil {
		if tgerr.Is(err, "SESSION_PASSWORD_NEEDED") {
			authorization = PasswordRequest(ctx, authClient, app)
		}
	}

	return authorization, nil
}

func PasswordRequest(ctx context.Context, authClient *auth.Client, app ViewApp) *tg.AuthAuthorization {
	viewApp := app.GetViewApp()
	var errorMsg string

	for {
		passwordChan := make(chan string)

		viewApp.QueueUpdateDraw(func() {
			viewApp.SetRoot(app.NewPasswordRequest(func(password string) {
				passwordChan <- password
			}, errorMsg), true)
		})

		password := <-passwordChan

		authorization, err := authClient.Password(ctx, password)
		if err != nil {
			if errors.Is(err, auth.ErrPasswordInvalid) {
				errorMsg = "❌ Неверный пароль, попробуйте еще раз"
			} else {
				errorMsg = "❌ Ошибка аутентификации: " + err.Error()
			}
		} else {
			return authorization
		}
	}
}
