package telegram

import (
	"questionnaire-bot/internal/handlers"
	"questionnaire-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Bot *tgbotapi.BotAPI
	l   *logger.Logger
}

func New(token string, l *logger.Logger) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Telegram{Bot: bot, l: l}, nil
}

func (t *Telegram) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handlers.Update(update, t.l)
	}
}
