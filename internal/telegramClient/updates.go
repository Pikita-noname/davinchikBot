package telegramclient

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type ViewController interface {
	SetTitle(string)
	SetDescription(string)
	SetImage(image.Image)
}

type UpdateHandler struct {
	client       *telegram.Client
	sender       *message.Sender
	view         ViewController
	updateDrawer func(f func())
	AdsIndex     int
	Buttons      []string
}

func (h *UpdateHandler) HandleUpdate(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	if msg, ok := update.Message.(*tg.Message); ok {
		if !h.isDavinchiBot(msg, e) {
			return nil
		}
		// go h.handlePhoto(ctx, msg)
		h.updateDrawer(func() {
			if len([]rune(msg.Message)) > 4 {
				h.handleMessage(ctx, msg, e)
			}
		})

	}
	time.Sleep(time.Second)
	return nil
}

func (h *UpdateHandler) isDavinchiBot(msg *tg.Message, e tg.Entities) bool {
	userName := h.getUserName(msg, e)
	return userName == "leomatchbot"
}

func (h *UpdateHandler) getUserName(msg *tg.Message, e tg.Entities) string {

	userID := h.getUserId(msg)

	if userID == 0 {
		return "unknown"
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
		Limit:    1024 * 1024,
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

func (h *UpdateHandler) handleMessage(ctx context.Context, msg *tg.Message, e tg.Entities) error {
	err := h.getButtons(msg)
	if err != nil {
		h.view.SetDescription(err.Error())
	}

	if !h.isProfile(msg) {
		h.tryToSkip()
		return nil
	}

	h.AdsIndex = 0

	profile := h.getProfile(msg)

	h.updateView(profile)

	var message string
	if h.filterProfile(profile) {
		message = "👎"
	} else {
		message = "👎"
	}

	peer, err := h.getInputPeer(msg, e)
	if err != nil {
		h.view.SetDescription(err.Error())
		return nil
	}

	h.sendAnswer(peer, message)
	return nil
}

func (h *UpdateHandler) getButtons(msg *tg.Message) error {
	if msg.ReplyMarkup == nil {
		return fmt.Errorf("message dont have buttons")
	}

	keyboard, ok := msg.ReplyMarkup.(*tg.ReplyKeyboardMarkup)
	if !ok {
		return fmt.Errorf("cant convert to type *tg.ReplyKeyboardMarkup")
	}

	var buttons []string

	for _, row := range keyboard.Rows {
		for _, button := range row.Buttons {
			buttons = append(buttons, button.GetText())
		}
	}

	h.Buttons = buttons

	return nil
}

func (h *UpdateHandler) isProfile(msg *tg.Message) bool {
	return len(strings.Split(msg.GetMessage(), ",")) > 1
}

func (h *UpdateHandler) updateView(profile Profile) error {
	h.view.SetTitle(fmt.Sprintf("%s,%s", profile.Name, string(profile.Age)))
	h.view.SetDescription(profile.Description)

	return nil
}

func (h *UpdateHandler) tryToSkip() error {
	// h.sendAnswer(h.Buttons[h.AdsIndex])
	h.AdsIndex++
	return nil
}

type Profile struct {
	Name        string
	Age         string
	Description string
}

func (h *UpdateHandler) getProfile(msg *tg.Message) Profile {
	if msg == nil || msg.Message == "" {
		return Profile{}
	}

	text := strings.TrimSpace(msg.Message)

	parameters := strings.SplitN(text, ",", 3)

	return Profile{
		Name:        parameters[0],
		Age:         parameters[1],
		Description: parameters[2],
	}
}

func (h *UpdateHandler) filterProfile(profile Profile) bool {
	return true
}

func (h *UpdateHandler) sendAnswer(peer tg.InputPeerClass, message string) error {
	if h.sender == nil || message == "" || len(h.Buttons) == 0 {
		return nil
	}

	h.sender.To(peer).Text(context.Background(), message)

	return nil
}

func (h *UpdateHandler) getInputPeer(msg *tg.Message, e tg.Entities) (tg.InputPeerClass, error) {
	switch peer := msg.PeerID.(type) {
	case *tg.PeerUser:
		if user, exists := e.Users[peer.UserID]; exists {
			return &tg.InputPeerUser{
				UserID:     peer.UserID,
				AccessHash: user.AccessHash,
			}, nil
		}
		return nil, fmt.Errorf("нет данных о пользователе %d в Entities", peer.UserID)
	case *tg.PeerChat:
		return &tg.InputPeerChat{
			ChatID: peer.ChatID,
		}, nil
	case *tg.PeerChannel:
		if _, exists := e.Chats[peer.ChannelID]; exists {
			return &tg.InputPeerChannel{
				ChannelID:  peer.ChannelID,
				AccessHash: 0,
			}, nil
		}
		return nil, fmt.Errorf("нет данных о канале %d в Entities", peer.ChannelID)
	default:
		return nil, fmt.Errorf("неизвестный тип пира: %T", peer)
	}
}

func (h *UpdateHandler) getUserId(msg *tg.Message) int64 {
	var userID int64

	switch peer := msg.FromID.(type) {
	case *tg.PeerUser:
		userID = peer.UserID
	case nil:
		switch peer := msg.PeerID.(type) {
		case *tg.PeerUser:
			userID = peer.UserID
		default:
			return 0
		}
	default:
		return 0
	}
	return userID
}
