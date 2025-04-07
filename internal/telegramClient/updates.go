package telegramclient

import (
	"context"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"golang.org/x/exp/slices"
)

type ViewController interface {
	SetTitle(string)
	SetDescription(string)
}

type Filter struct {
	Name        string
	Age         string
	Description string
}

type UpdateHandler struct {
	client        *telegram.Client
	sender        *message.Sender
	view          ViewController
	updateDrawer  func(f func())
	AdsIndex      int
	Buttons       []string
	filter        Filter
	AdsIterations int
	pause         bool
}

func (h *UpdateHandler) HandleUpdate(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	if msg, ok := update.Message.(*tg.Message); ok {
		if !h.isDavinchiBot(msg, e) {
			return nil
		}
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

func (h *UpdateHandler) handleMessage(ctx context.Context, msg *tg.Message, e tg.Entities) error {
	if h.pause {
		return nil
	}

	err := h.getButtons(msg)
	if err != nil {
		h.view.SetDescription(err.Error())
	}

	peer, err := h.getInputPeer(msg, e)
	if err != nil {
		h.view.SetDescription(err.Error())
		return nil
	}

	if !h.isProfile(msg) {
		h.tryToSkip(msg, e)
		return nil
	}

	h.AdsIndex = 0

	profile := h.getProfile(msg)

	h.updateView(profile)

	var message string
	if h.filterProfile(profile) {
		message = "❤️"
	} else {
		message = "👎"
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
	return len(strings.Split(msg.GetMessage(), ",")) > 1 || len(h.Buttons) < 4
}

func (h *UpdateHandler) updateView(profile Profile) error {
	h.view.SetTitle(fmt.Sprintf("%s,%s", profile.Name, string(profile.Age)))
	h.view.SetDescription(profile.Description)

	return nil
}

func (h *UpdateHandler) tryToSkip(msg *tg.Message, e tg.Entities) error {

	peer, err := h.getInputPeer(msg, e)
	if err != nil {
		h.view.SetDescription(err.Error())
		return nil
	}

	h.AdsIterations++

	if h.AdsIterations > 5 {
		h.pause = true
		go h.unPause(msg, e)
	}

	if slices.Contains(h.Buttons, "1 🚀") {
		h.sendAnswer(peer, "1 🚀")
		h.AdsIndex = 0
	}

	h.sendAnswer(peer, h.Buttons[h.AdsIndex])

	h.AdsIndex++
	if h.AdsIndex >= len(h.Buttons) {
		h.AdsIndex = 0
	}
	return nil
}

func (h *UpdateHandler) unPause(msg *tg.Message, e tg.Entities) {
	time.Sleep(time.Hour * 1)
	h.pause = false
	h.AdsIterations = 0
	h.handleMessage(context.Background(), msg, e)
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
	return h.filterDescription(profile.Description) && h.filterAge(profile.Age) && h.filterName(profile.Name)
}

func (h *UpdateHandler) filterName(name string) bool {
	return h.filter.Name == "" || h.filter.Name == name
}

func (h *UpdateHandler) filterAge(age string) bool {
	if h.filter.Age == "" {
		return true
	}

	operator := string(h.filter.Age[0])
	ageValue := h.filter.Age[1:]

	ageInt, err := strconv.Atoi(ageValue)
	if err != nil {
		return false
	}

	profileAge, err := strconv.Atoi(age)
	if err != nil {
		return false
	}

	switch operator {
	case "<":
		return profileAge < ageInt
	case ">":
		return profileAge > ageInt
	case "=":
		return profileAge == ageInt
	default:
		return false
	}
}

func (h *UpdateHandler) filterDescription(description string) bool {
	if h.filter.Description == "" {
		return true
	}

	return strings.Contains(description, h.filter.Description)
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
