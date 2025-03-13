package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
)

func main() {
	ctx := context.Background()

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
			qrString, err := qrcode.New(qrData, qrcode.Medium)
			if err != nil {
				return fmt.Errorf("ошибка генерации QR-кода: %w", err)
			}

			fmt.Println("\n🔹 Отсканируйте этот QR-код в Telegram:\n")
			fmt.Println(qrString.ToSmallString(false))
			fmt.Printf("\n🔗 Или откройте ссылку вручную: %s\n", qrData)

			return nil
		})

		// Обработка ошибки SESSION_PASSWORD_NEEDED
		if err != nil {
			if tgerr.Is(err, "SESSION_PASSWORD_NEEDED") {
				// Запрашиваем пароль у пользователя
				password, inputErr := promptPassword()
				if inputErr != nil {
					return fmt.Errorf("ошибка ввода пароля: %w", inputErr)
				}

				// Используем метод Password из пакета auth и обновляем authorization
				authorization, err = authClient.Password(ctx, password)
				if err != nil {
					if errors.Is(err, auth.ErrPasswordInvalid) {
						return fmt.Errorf("введен неверный пароль")
					}
					return fmt.Errorf("ошибка проверки пароля: %w", err)
				}
			} else {
				return fmt.Errorf("ошибка аутентификации: %w", err)
			}
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

// promptPassword запрашивает пароль 2FA
func promptPassword() (string, error) {
	fmt.Print("🔑 Введите пароль двухфакторной аутентификации: ")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}
