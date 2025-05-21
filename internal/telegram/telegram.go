package telegram

import (
	"errors"
	"fmt"
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Notifier interface {
	SendToAdmin(msg string)
	Send(int64, string, any)
}

type Telegram struct {
	Bot     *tgbotapi.BotAPI
	l       zerolog.Logger
	u       usecase.Usecase
	adminID int64
}

func New(token string, l zerolog.Logger, u usecase.Usecase, adminID int64) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Telegram{Bot: bot, l: l, u: u, adminID: adminID}, nil
}

func (t *Telegram) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.l.Warn().Interface("panic", r).
						Str("text", update.Message.Text).
						Int64("user_id", update.Message.From.ID).
						Msg("Паника в обработке сообщения")

					t.Send(t.adminID, fmt.Sprintf("Произошла ошибка в обработке сообщения от пользователя %d:\n%v", update.Message.From.ID, r), nil)
				}
			}()
			t.Update(update)
		}()

	}
}

func (t *Telegram) Update(update tgbotapi.Update) {
	user, err := t.u.GetUser(update.Message.From.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			user = entity.New(update.Message.From.ID, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
			t.l.Info().Int64("id", user.TgID).Msgf("Новый пользователь %s!", user.Username)
		} else {
			t.l.Warn().Err(err).Int64("id", update.Message.From.ID).Str("user", update.Message.From.FirstName).Msg("failed get user from DB")
		}

	}

	if update.Message.Contact != nil {
		t.processContact(user, update.Message.Contact.PhoneNumber)
	} else {
		t.processMessage(user, update.Message.Text)
	}

	user.MaxStepReached = max(user.MaxStepReached, user.CurrentStep) // todo: add metric

	if err = t.u.SaveUser(user); err != nil {
		t.l.Error().Err(err).Int64("id", user.TgID).Msg("Ошибка при сохранении пользователя")
	}

}

func (t *Telegram) Send(chatID int64, text string, keyboard any) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	_, err := t.Bot.Send(msg)
	if err != nil {
		t.l.Error().Err(err).Int64("chatID", chatID).Str("msg", text).Msg("Ошибка отправки сообщения")
	}
}

func (t *Telegram) SendToAdmin(msg string) {
	adminID := t.adminID
	t.Send(adminID, msg, nil)
}
