package telegramclient

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"

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

	viewApp := app.GetViewApp()

	// Создаем обработчик обновлений
	d := tg.NewUpdateDispatcher()
	loggedIn := qrlogin.OnLoginToken(d)

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Errorf("ошибка чтения .env: %w", err)
	}

	sessionStorage := &telegram.FileSessionStorage{
		Path: "session.json",
	}

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		UpdateHandler:  d,
		SessionStorage: sessionStorage,
	})

	if err != nil {
		fmt.Errorf("ошибка создания клиента из .env: %w", err)
	}

	// Создаем экземпляр auth.Client с добавлением rand.Reader
	authClient := auth.NewClient(client.API(), rand.Reader, 29708230, "428ef65a36ade933259c4c832cd65bfd")

	// Запускаем клиент
	if err := client.Run(ctx, func(ctx context.Context) error {
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

		// Обработка ошибки SESSION_PASSWORD_NEEDED
		if err != nil {
			if tgerr.Is(err, "SESSION_PASSWORD_NEEDED") {

				passwordChan := make(chan string)

				viewApp.QueueUpdateDraw(func() {
					log.Println("📢 Отображение окна ввода пароля")

					viewApp.SetRoot(app.NewPasswordRequest(func(password string) {
						log.Println("📝 Введен пароль, отправка запроса...")
						app.BackToMain()
						passwordChan <- password
					}), true)

				})

				password := <-passwordChan

				authorization, err = authClient.Password(ctx, password)
				if err != nil {
					if errors.Is(err, auth.ErrPasswordInvalid) {
						log.Println("❌ Введен неверный пароль")
					} else {
						log.Printf("❌ Ошибка при аутентификации паролем: %v", err)
					}
				} else {
					log.Println("✅ Пароль успешно принят!")
				}
			}
		}

		if authorization == nil {
			return nil
		}

		// Проверяем данные пользователя
		u, ok := authorization.User.AsNotEmpty()
		if !ok {
			return fmt.Errorf("не удалось получить данные пользователя: %T", authorization.User)
		}

		fmt.Println("\n✅ Успешный вход!")
		fmt.Printf("ID: %d | Username: %s | Бот: %t\n", u.ID, u.Username, u.Bot)
		return nil
	}); err != nil {
		log.Fatal("Ошибка запуска клиента:", err)
	}
}
