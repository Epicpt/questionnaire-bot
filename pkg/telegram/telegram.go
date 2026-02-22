package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	NewMessage(int64, string) tgbotapi.MessageConfig
	NewUpdate(int) tgbotapi.UpdateConfig
	NewDocument(chatID int64, file tgbotapi.RequestFileData) tgbotapi.DocumentConfig
}

type Bot struct {
	Bot *tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{Bot: bot}, nil
}

func (b *Bot) NewUpdate(offset int) tgbotapi.UpdateConfig {
	return tgbotapi.NewUpdate(offset)
}

func (b *Bot) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return b.Bot.GetUpdatesChan(config)
}
func (b *Bot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return b.Bot.Send(c)
}
func (b *Bot) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, text)
}
func (b *Bot) NewDocument(chatID int64, file tgbotapi.RequestFileData) tgbotapi.DocumentConfig {
	return tgbotapi.NewDocument(chatID, file)
}
