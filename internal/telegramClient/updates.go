package telegramclient

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	// Регистрация декодера JPEG
	// Регистрация декодера PNG
)

type ViewController interface {
	SetTitle(string)
	SetDescription(string)
	SetImage(image.Image)
}

type UpdateHandler struct {
	client       *telegram.Client
	view         ViewController
	updateDrawer func(f func())
}

func (h *UpdateHandler) HandleUpdate(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	if msg, ok := update.Message.(*tg.Message); ok {
		go h.handlePhoto(ctx, msg)
		h.updateDrawer(func() {
			h.view.SetTitle("Новое сообщение: " + h.getUserName(msg, e))
		})

	}
	return nil
}

func (h *UpdateHandler) isDavinchiBot(msg *tg.Message, e tg.Entities) bool {
	userName := h.getUserName(msg, e)
	return userName == "@leomatchbot"
}

func (h *UpdateHandler) getUserName(msg *tg.Message, e tg.Entities) string {
	var userID int64

	switch peer := msg.FromID.(type) {
	case *tg.PeerUser:
		userID = peer.UserID
	case nil:
		switch peer := msg.PeerID.(type) {
		case *tg.PeerUser:
			userID = peer.UserID
		default:
			return "Неизвестный отправитель"
		}
	default:
		return "Неизвестный отправитель"
	}

	if user, exists := e.Users[userID]; exists {
		if user.Username != "" {
			return user.Username
		}
		return user.FirstName
	}

	return fmt.Sprintf("UserID:%d", userID)
}

func (h *UpdateHandler) handlePhoto(ctx context.Context, msg *tg.Message) error {
	switch media := msg.Media.(type) {
	case *tg.MessageMediaPhoto:
		if photo, ok := media.Photo.(*tg.Photo); ok {
			var largestSize tg.PhotoSizeClass
			var maxSize int
			for _, size := range photo.Sizes {
				switch s := size.(type) {
				case *tg.PhotoSize:
					if s.Size > maxSize {
						maxSize = s.Size
						largestSize = s
					}
				case *tg.PhotoSizeProgressive:
					if len(s.Sizes) > 0 && s.Sizes[len(s.Sizes)-1] > maxSize {
						maxSize = s.Sizes[len(s.Sizes)-1]
						largestSize = s
					}
				}
			}

			if largestSize == nil {
				return fmt.Errorf("no valid photo size found")
			}

			var thumbSize string
			switch s := largestSize.(type) {
			case *tg.PhotoSize:
				thumbSize = s.Type
			case *tg.PhotoSizeProgressive:
				thumbSize = s.Type
			}

			location := &tg.InputPhotoFileLocation{
				ID:            photo.ID,
				AccessHash:    photo.AccessHash,
				FileReference: photo.FileReference,
				ThumbSize:     thumbSize,
			}

			img, err := h.getImageFromLocation(ctx, location, maxSize)
			if err != nil {
				return fmt.Errorf("get image from photo: %v", err)
			}

			h.updateDrawer(func() {
				h.view.SetImage(img)
			})
		}
	default:
		h.updateDrawer(func() {
			h.view.SetDescription("Медиа не является фото")
		})
	}
	return nil
}

func (h *UpdateHandler) getImageFromLocation(ctx context.Context, location *tg.InputPhotoFileLocation, totalSize int) (image.Image, error) {
	var buf bytes.Buffer

	req := &tg.UploadGetFileRequest{
		Location: location,
		Offset:   0,
		Limit:    1024 * 1024, // Запрашиваем по 1MB за раз
	}

	for {
		resp, err := h.client.API().UploadGetFile(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("upload get file: %v", err)
		}

		file, ok := resp.(*tg.UploadFile)
		if !ok {
			return nil, fmt.Errorf("unexpected response type")
		}

		if _, err := buf.Write(file.Bytes); err != nil {
			return nil, fmt.Errorf("write to buffer: %v", err)
		}

		if len(file.Bytes) < req.Limit || buf.Len() >= totalSize {
			break
		}

		req.Offset += int64(len(file.Bytes))
	}

	img, _, err := image.Decode(&buf)
	if err != nil {
		return nil, fmt.Errorf("decode image: %v", err)
	}

	return img, nil
}
