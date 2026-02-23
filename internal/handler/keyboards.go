package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const telegramPhone = "📱 Отправить Telegram-номер"

func KeyboardFromOptions(q Question, showBack bool) any {
	var keyboard [][]tgbotapi.KeyboardButton

	for _, opt := range q.Options {
		row := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(opt.Text))
		keyboard = append(keyboard, row)
	}

	switch q.SpecialKeyboard {
	case requestPhone:
		btn := tgbotapi.NewKeyboardButton(telegramPhone)
		btn.RequestContact = true
		keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(btn))
	}

	if showBack {
		keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BackButton)))
	}

	if len(keyboard) == 0 {
		return tgbotapi.NewRemoveKeyboard(true)
	}
	return tgbotapi.NewReplyKeyboard(keyboard...)
}
