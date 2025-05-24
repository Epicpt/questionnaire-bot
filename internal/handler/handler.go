package handler

import (
	"errors"
	"fmt"
	"questionnaire-bot/internal/entity"
	"questionnaire-bot/internal/usecase"
	"questionnaire-bot/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Handler interface {
	SendToAdmin(string)
	Send(int64, string, any)
	AdvanceStep(*entity.User)
	SendNextQuestion(*entity.User)
	FinishSurvey(*entity.User)
	ProcessMessage(*entity.User, string)
	ProcessContact(*entity.User, string)
}

type BotHandler struct {
	Bot     telegram.Telegram
	l       zerolog.Logger
	u       usecase.Usecase
	adminID int64
}

func New(bot telegram.Telegram, l zerolog.Logger, u usecase.Usecase, adminID int64) *BotHandler {
	return &BotHandler{Bot: bot, l: l, u: u, adminID: adminID}
}

func (t *BotHandler) Start() {
	u := t.Bot.NewUpdate(0)
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

func (t *BotHandler) Update(update tgbotapi.Update) {
	user, err := t.u.GetUser(update.Message.From.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			user = entity.New(update.Message.From.ID, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName)
			t.l.Info().Int64("id", user.TgID).Msgf("Новый пользователь %s!", user.Username)
		} else {
			t.l.Warn().Err(err).Int64("id", update.Message.From.ID).Str("user", update.Message.From.FirstName).Msg("failed get user from DB")
		}

	}

	if spamCheck(user.EmailSentCnt) {
		t.Send(user.ChatID, fmt.Sprintf("Менеджер свяжется с Вами в ближайшее время!"), nil)
		return
	}

	if update.Message.Contact != nil {
		t.ProcessContact(user, update.Message.Contact.PhoneNumber)
	} else {
		t.ProcessMessage(user, update.Message.Text)
	}

	user.MaxStepReached = max(user.MaxStepReached, user.CurrentStep)

	if err = t.u.SaveUser(user); err != nil {
		t.l.Error().Err(err).Int64("id", user.TgID).Msg("Ошибка при сохранении пользователя")
	}

}

func (t *BotHandler) Send(chatID int64, text string, keyboard any) {
	msg := t.Bot.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard
	_, err := t.Bot.Send(msg)
	if err != nil {
		t.l.Error().Err(err).Int64("chatID", chatID).Str("msg", text).Msg("Ошибка отправки сообщения")
	}
}

func (t *BotHandler) SendToAdmin(msg string) {
	t.Send(t.adminID, msg, nil)
}

func spamCheck(cnt int) bool {
	if cnt > 3 {
		return true
	}
	return false
}
